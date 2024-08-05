package bindings

import (
	"flag"
	"io"
	"log"
	"os"
	"os/signal"
	"reflect"
	"strings"
	"syscall"
	"time"

	"github.com/balazsgrill/potatodrive/win/projfs/filesystem"
	"github.com/spf13/afero"
	"golang.org/x/sys/windows/registry"
)

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

func BindVirtualizationInstance(localpath string, remotefs afero.Fs) (io.Closer, error) {
	closer, err := filesystem.StartProjecting(localpath, remotefs)
	if err != nil {
		return nil, err
	}

	t := time.NewTicker(30 * time.Second)
	go func() {
		for range t.C {
			err = closer.PerformSynchronization()
			if err != nil {
				log.Println(err)
			}
		}
	}()

	return (closerFunc)(func() error {
		t.Stop()
		return closer.Close()
	}), nil
}
