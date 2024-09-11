package main

import (
	"errors"
	"io"
	"os"
	"path/filepath"

	"github.com/balazsgrill/potatodrive/bindings"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/sys/windows/registry"
)

type Manager struct {
	zerolog.Logger
	logf        io.Closer
	logfilepath string

	parentkey registry.Key

	keylist   []string
	instances map[string]io.Closer
}

func initLogger() (string, zerolog.Logger, io.Closer) {
	cachedir, err := os.UserCacheDir()
	if err != nil {
		panic(err)
	}
	logfolder := filepath.Join(cachedir, "PotatoDrive")
	err = os.MkdirAll(logfolder, 0777)
	if err != nil {
		panic(err)
	}

	logfile := "potatodrive.log"
	logfilepath := filepath.Join(logfolder, logfile)
	logf, err := os.OpenFile(logfilepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}
	return logfilepath, log.Output(zerolog.MultiLevelWriter(logf, zerolog.NewConsoleWriter())).With().Timestamp().Logger(), logf
}

func startInstance(parentkey registry.Key, keyname string, context bindings.InstanceContext) (io.Closer, error) {
	key, err := registry.OpenKey(parentkey, keyname, registry.QUERY_VALUE)
	if err != nil {
		context.Logger.Printf("Open key: %v", err)
		return nil, err
	}

	var basec bindings.BaseConfig
	err = bindings.ReadConfigFromRegistry(key, &basec)
	if err != nil {
		context.Logger.Printf("Get base config: %v", err)
		return nil, err
	}
	config := bindings.CreateConfigByType(basec.Type)
	bindings.ReadConfigFromRegistry(key, config)
	err = config.Validate()
	if err != nil {
		context.Logger.Printf("Validate config: %v", err)
		return nil, err
	}
	fs, err := config.ToFileSystem()
	if err != nil {
		context.Logger.Printf("Create file system: %v", err)
		return nil, err
	}

	context.Logger.Info().Msgf("Starting %s on %s", keyname, basec.LocalPath)
	innercontext := context
	innercontext.Logger = context.Logger.With().Str("instance", keyname).Logger()
	c, err := bindings.BindVirtualizationInstance(keyname, basec.LocalPath, fs, innercontext)
	if err != nil {
		return nil, err
	}
	context.Logger.Info().Msgf("%s started", keyname)
	return c, nil
}

func New() (*Manager, error) {
	m := &Manager{
		instances: make(map[string]io.Closer),
	}
	m.logfilepath, m.Logger, m.logf = initLogger()
	var err error
	m.parentkey, err = registry.OpenKey(registry.LOCAL_MACHINE, "SOFTWARE\\PotatoDrive", registry.QUERY_VALUE|registry.READ)
	if err != nil {
		return nil, err
	}

	m.keylist, err = m.parentkey.ReadSubKeyNames(0)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (m *Manager) Close() error {
	for key, instance := range m.instances {
		m.Logger.Info().Msgf("Closing %s", key)
		err := instance.Close()
		if err != nil {
			m.Logger.Err(err).Msg("Failed to close instance")
		}
	}
	return m.logf.Close()
}

func (m *Manager) InstanceList() ([]string, error) {
	return m.keylist, nil
}

func (m *Manager) StartInstance(id string, context bindings.InstanceContext) error {
	instance, err := startInstance(m.parentkey, id, context)
	if err != nil {
		return err
	}
	if instance == nil {
		return errors.New("instance is nil")
	}
	m.instances[id] = instance
	return nil
}

func (m *Manager) StopInstance(id string) error {
	instance := m.instances[id]
	err := instance.Close()
	if err != nil {
		return err
	}
	delete(m.instances, id)
	return nil
}
