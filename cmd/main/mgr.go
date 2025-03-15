package main

import (
	"errors"
	"io"
	"os"
	"path/filepath"

	"github.com/balazsgrill/potatodrive/bindings"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Manager struct {
	zerolog.Logger
	logf        io.Closer
	logfilepath string

	configProvider bindings.ConfigProvider
	instances      map[string]io.Closer
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

func startInstance(config bindings.Config, context bindings.InstanceContext) (io.Closer, error) {
	fs, err := config.ToFileSystem(context.Logger)
	if err != nil {
		context.Logger.Error().Msgf("Create file system: %v", err)
		return nil, err
	}

	context.Logger.Info().Msgf("Starting %s on %s", config.ID, config.LocalPath)
	innercontext := context
	innercontext.Logger = context.Logger.With().Str("instance", config.ID).Logger()
	c, err := bindings.BindVirtualizationInstance(config.ID, &config.BaseConfig, fs, innercontext)
	if err != nil {
		return nil, err
	}
	context.Logger.Info().Msgf("%s started", config.ID)
	return c, nil
}

func New() (*Manager, error) {
	m := &Manager{
		instances: make(map[string]io.Closer),
	}
	m.logfilepath, m.Logger, m.logf = initLogger()
	m.configProvider = bindings.NewRegistryConfigProvider(m.Logger, "SOFTWARE\\PotatoDrive")
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
	return m.configProvider.Keys(), nil
}

func (m *Manager) StartInstance(id string, context bindings.InstanceContext) error {
	config, err := m.configProvider.ReadConfig(id)
	if err != nil {
		return err
	}
	instance, err := startInstance(config, context)
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
