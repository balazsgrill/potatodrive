package bindings

import (
	"flag"
	"io"
	"os"
	"os/signal"
	"reflect"
	"strings"
	"syscall"
	"time"

	"github.com/rs/zerolog"

	"github.com/balazsgrill/potatodrive/bindings/s3"
	"github.com/balazsgrill/potatodrive/bindings/sftp"
	"github.com/balazsgrill/potatodrive/core"
	cfapi "github.com/balazsgrill/potatodrive/core/cfapi/filesystem"
	prjfs "github.com/balazsgrill/potatodrive/core/projfs/filesystem"
	"github.com/spf13/afero"
	"golang.org/x/sys/windows/registry"
)

const UseCFAPI bool = true

type BindingConfig interface {
	Validate() error
	ToFileSystem() (afero.Fs, error)
}

type BaseConfig struct {
	LocalPath string `flag:"localpath,Local folder" reg:"LocalPath"`
	Type      string `flag:"type,Type of binding" reg:"Type"`
}

func ConfigToFlags(config any) {
	structPtrValue := reflect.ValueOf(config)
	structValue := structPtrValue.Elem()
	structType := structValue.Type()
	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)

		tagstr := field.Tag.Get("flag")
		if tagstr != "" {
			tagdata := strings.Split(tagstr, ",")
			tag := tagdata[0]
			msg := tagdata[1]
			switch field.Type.Kind() {
			case reflect.String:
				flag.StringVar((*string)(structValue.Field(i).Addr().UnsafePointer()), tag, "", msg)
			case reflect.Bool:
				flag.BoolVar((*bool)(structValue.Field(i).Addr().UnsafePointer()), tag, false, msg)
			}
		}
	}
}

func CreateConfigByType(typestr string) BindingConfig {
	switch typestr {
	case "afero-s3":
		return &s3.Config{}
	case "afero-sftp":
		return &sftp.Config{}
	}
	return nil
}

func ReadConfigFromRegistry(key registry.Key, config any) error {
	structPtrValue := reflect.ValueOf(config)
	structValue := structPtrValue.Elem()
	structType := structValue.Type()
	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		fieldvalue := structValue.Field(i)
		tag := field.Tag.Get("reg")
		if tag != "" {
			switch field.Type.Kind() {
			case reflect.String:
				value, _, err := key.GetStringValue(tag)
				if os.IsNotExist(err) {
					continue
				}
				if err != nil {
					return err
				}
				fieldvalue.SetString(value)
			case reflect.Bool:
				value, _, err := key.GetIntegerValue(tag)
				if os.IsNotExist(err) {
					continue
				}
				if err != nil {
					return err
				}
				fieldvalue.SetBool(value != 0)
			}
		}
	}
	return nil
}

func CloseOnSigTerm(closers ...io.Closer) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	<-c
	for _, closer := range closers {
		closer.Close()
	}
	os.Exit(1)
}

type closerFunc func() error

func (f closerFunc) Close() error {
	return f()
}

type InstanceContext struct {
	Logger            zerolog.Logger
	StateCallback     func(core.ConnectionState)
	FileStateCallback func(core.FileSyncState)
}

func (context InstanceContext) ConnectionStateChanged(id string, syninprogress bool, err error) {
	if context.StateCallback == nil {
		return
	}
	state := core.ConnectionState{
		ID:             id,
		SyncInProgress: syninprogress,
		LastSyncError:  err,
	}
	context.StateCallback(state)
}

func BindVirtualizationInstance(id string, localpath string, remotefs afero.Fs, context InstanceContext) (io.Closer, error) {
	var closer core.Virtualization
	var err error
	if UseCFAPI {
		err = cfapi.RegisterRootPath(id, localpath)
		if err != nil {
			return nil, err
		}
		closer, err = cfapi.StartProjecting(localpath, remotefs, context.Logger)
	} else {
		closer, err = prjfs.StartProjecting(localpath, remotefs, context.Logger)
	}
	if err != nil {
		return nil, err
	}
	closer.SetFileStateHandler(context.FileStateCallback)

	internalSynchronize := func() {
		context.ConnectionStateChanged(id, true, nil)
		err = closer.PerformSynchronization()
		if err != nil {
			context.Logger.Err(err).Send()
			context.ConnectionStateChanged(id, false, err)
		} else {
			context.ConnectionStateChanged(id, false, nil)
		}
	}

	// initial sync
	internalSynchronize()

	t := time.NewTicker(30 * time.Second)
	go func() {
		for range t.C {
			internalSynchronize()
		}
	}()

	return (closerFunc)(func() error {
		t.Stop()
		return closer.Close()
	}), nil
}
