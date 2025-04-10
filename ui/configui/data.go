package configui

import (
	"github.com/balazsgrill/potatodrive/bindings"
	"github.com/balazsgrill/potatodrive/bindings/gphotos"
	"github.com/balazsgrill/potatodrive/bindings/s3"
	"github.com/balazsgrill/potatodrive/bindings/sftp"
)

type ConfigValues struct {
	ID            string
	Base          bindings.BaseConfig
	HasValue      bool
	HasS3         bool
	S3Config      s3.Config
	HasSFTP       bool
	SFTPConfig    sftp.Config
	HasGPhotos    bool
	GPhotosConfig gphotos.Config

	NotHasValue bool
}

func ReadFrom(data *bindings.Config) *ConfigValues {
	if data == nil {
		return &ConfigValues{
			HasValue:    false,
			NotHasValue: true,
		}
	}
	result := &ConfigValues{
		ID:          data.ID,
		Base:        data.BaseConfig,
		HasValue:    true,
		NotHasValue: false,
	}
	if s3, ok := data.BindingConfig.(*s3.Config); ok {
		result.HasS3 = true
		result.S3Config = *s3
	}
	if sftp, ok := data.BindingConfig.(*sftp.Config); ok {
		result.HasSFTP = true
		result.SFTPConfig = *sftp
	}
	if gphotos, ok := data.BindingConfig.(*gphotos.Config); ok {
		result.HasGPhotos = true
		result.GPhotosConfig = *gphotos
	}
	return result
}

func WriteTo(data *ConfigValues) *bindings.Config {
	if data == nil {
		return nil
	}
	if !data.HasValue {
		return nil
	}
	result := &bindings.Config{
		ID:            data.ID,
		BaseConfig:    data.Base,
		BindingConfig: nil,
	}
	if data.HasS3 {
		result.BindingConfig = &data.S3Config
		result.Type = bindings.TYPE_S3
	}
	if data.HasSFTP {
		result.BindingConfig = &data.SFTPConfig
		result.Type = bindings.TYPE_SFTP
	}
	if data.HasGPhotos {
		result.BindingConfig = &data.GPhotosConfig
		result.Type = bindings.TYPE_GPHOTOS
	}
	return result
}
