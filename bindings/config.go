package bindings

type Config struct {
	ID string
	BaseConfig
	BindingConfig
}

type ConfigProvider interface {
	Keys() []string
	ReadConfig(key string) (Config, error)
}

type ConfigWriter interface {
	ConfigProvider
	WriteConfig(config Config) error
	DeleteConfig(key string) error
}
