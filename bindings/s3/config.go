package s3

import (
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/balazsgrill/potatodrive/bindings"
	s3 "github.com/fclairamb/afero-s3"
	"github.com/spf13/afero"
)

type Config struct {
	LocalPath string `flag:"localpath,Local folder" reg:"LocalPath"`
	Endpoint  string `flag:"endpoint,S3 endpoint" reg:"Endpoint"`
	Region    string `flag:"region,Region" reg:"Region"`
	Bucket    string `flag:"bucket,Bucket" reg:"Bucket"`
	KeyId     string `flag:"keyid,Access Key ID" reg:"KeyID"`
	KeySecret string `flag:"secret,Access Key Secret" reg:"KeySecret"`
	UseSSL    bool   `flag:"useSSL,Use SSL encryption for S3 connection" reg:"UseSSL"`
}

func (c *Config) Validate() error {
	if c.Endpoint == "" {
		return errors.New("endpoint is mandatory")
	}
	if c.LocalPath == "" {
		return errors.New("localpath is mandatory")
	}
	if c.Region == "" {
		return errors.New("region is mandatory")
	}
	if c.Bucket == "" {
		return errors.New("bucket is mandatory")
	}
	if c.KeyId == "" {
		return errors.New("keyid is mandatory")
	}
	if c.KeySecret == "" {
		return errors.New("secret is mandatory")
	}
	return nil
}

func (c *Config) ToFileSystem() (afero.Fs, error) {
	sess, err := session.NewSession(&aws.Config{
		Region:           aws.String(c.Region),
		Endpoint:         aws.String(c.Endpoint),
		DisableSSL:       aws.Bool(!c.UseSSL),
		S3ForcePathStyle: aws.Bool(true),
		Credentials:      credentials.NewStaticCredentials(c.KeyId, c.KeySecret, ""),
	})
	if err != nil {
		return nil, err
	}

	fs := s3.NewFs(c.Bucket, sess)
	fs.MkdirAll("root", 0777)
	rootfs := bindings.NewBasePathFs(fs, "root")
	return rootfs, nil
}
