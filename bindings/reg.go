package bindings

import (
	"os"
	"reflect"
	"strings"

	"github.com/rs/zerolog"
	"golang.org/x/sys/windows/registry"
)

type registryConfigProvider struct {
	logger  zerolog.Logger
	basekey string
}

// DeleteConfig implements ConfigWriter.
func (r *registryConfigProvider) DeleteConfig(key string) error {
	parentkey, err := registry.OpenKey(registry.LOCAL_MACHINE, r.basekey, registry.QUERY_VALUE|registry.ENUMERATE_SUB_KEYS|registry.WRITE)
	if err != nil {
		r.logger.Err(err).Msgf("Open key: %s", r.basekey)
		return err
	}
	defer parentkey.Close()

	err = registry.DeleteKey(parentkey, key)
	if err != nil {
		r.logger.Err(err).Msgf("Delete key: %s", key)
		return err
	}
	return nil
}

// WriteConfig implements ConfigWriter.
func (r *registryConfigProvider) WriteConfig(config Config) error {
	keyname := config.ID
	parentkey, err := registry.OpenKey(registry.LOCAL_MACHINE, r.basekey, registry.QUERY_VALUE|registry.READ)
	if err != nil {
		r.logger.Err(err).Msgf("Open key: %s", r.basekey)
		return err
	}
	defer parentkey.Close()
	key, _, err := registry.CreateKey(parentkey, keyname, registry.SET_VALUE)
	if err != nil {
		r.logger.Err(err).Msgf("Create key: %s", keyname)
		return err
	}
	defer key.Close()

	err = writeConfigToRegistry(key, config.BaseConfig)
	if err != nil {
		return err
	}
	return writeConfigToRegistry(key, config.BindingConfig)
}

// Keys implements ConfigProvider.
func (r *registryConfigProvider) Keys() []string {
	parentkey, err := registry.OpenKey(registry.LOCAL_MACHINE, r.basekey, registry.QUERY_VALUE|registry.READ)
	if err != nil {
		r.logger.Err(err).Msgf("Open key: %s", r.basekey)
		return nil
	}
	defer parentkey.Close()

	keylist, err := parentkey.ReadSubKeyNames(0)
	if err != nil {
		r.logger.Err(err).Msgf("Read sub keys: %s", r.basekey)
		return nil
	}
	return keylist
}

// ReadConfig implements ConfigProvider.
func (r *registryConfigProvider) ReadConfig(keyname string) (Config, error) {
	var result Config
	result.ID = keyname
	parentkey, err := registry.OpenKey(registry.LOCAL_MACHINE, r.basekey, registry.QUERY_VALUE|registry.READ)
	if err != nil {
		r.logger.Err(err).Msgf("Open key: %s", r.basekey)
		return result, err
	}
	defer parentkey.Close()

	key, err := registry.OpenKey(parentkey, keyname, registry.QUERY_VALUE)
	if err != nil {
		r.logger.Err(err).Msgf("Open key: %s", keyname)
		return result, err
	}
	defer key.Close()

	err = ReadConfigFromRegistry(key, &result.BaseConfig)
	if err != nil {
		r.logger.Err(err).Msgf("Get base config: %v", err)
		return result, err
	}
	config := CreateConfigByType(result.Type)
	err = ReadConfigFromRegistry(key, &result.BaseConfig)
	if err != nil {
		r.logger.Err(err).Msgf("Read config: %v", err)
		return result, err
	}
	err = config.Validate()
	if err != nil {
		r.logger.Err(err).Msgf("Validate config: %v", err)
		return result, err
	}
	return result, nil
}

func NewRegistryConfigProvider(logger zerolog.Logger, basekey string) ConfigProvider {
	return &registryConfigProvider{logger: logger, basekey: basekey}
}

func NewRegistryConfigWriter(logger zerolog.Logger, basekey string) ConfigWriter {
	return &registryConfigProvider{logger: logger, basekey: basekey}
}

func writeConfigToRegistry(key registry.Key, config any) error {
	structPtrValue := reflect.ValueOf(config)
	structValue := structPtrValue.Elem()
	structType := structValue.Type()
	for i := range structType.NumField() {
		field := structType.Field(i)
		fieldValue := structValue.Field(i)
		tag := field.Tag.Get("reg")
		if tag != "" {
			switch field.Type.Kind() {
			case reflect.String:
				err := key.SetStringValue(tag, fieldValue.String())
				if err != nil {
					return err
				}
			case reflect.Bool:
				var value uint64
				if fieldValue.Bool() {
					value = 1
				} else {
					value = 0
				}
				err := key.SetQWordValue(tag, value)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func ReadConfigFromRegistry(key registry.Key, config any) error {
	structPtrValue := reflect.ValueOf(config)
	structValue := structPtrValue.Elem()
	structType := structValue.Type()
	for i := range structType.NumField() {
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
				if err == registry.ErrUnexpectedType {
					// attempt to read as multi-string
					values, _, err := key.GetStringsValue(tag)
					if err != nil {
						return err
					}
					value = strings.Join(values, "\n")
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
