module github.com/balazsgrill/potatodrive

go 1.23.5

toolchain go1.23.6

require (
	github.com/apache/thrift v0.21.0
	github.com/aws/aws-sdk-go v1.54.20
	github.com/fclairamb/afero-s3 v0.3.1
	github.com/fsnotify/fsnotify v1.7.0
	github.com/go-ole/go-ole v1.2.6
	github.com/google/uuid v1.6.0
	github.com/lxn/walk v0.0.0-20210112085537-c389da54e794
	github.com/lxn/win v0.0.0-20210218163916-a377121e959e
	github.com/pkg/sftp v1.13.6
	github.com/rs/zerolog v1.33.0
	github.com/saltosystems/winrt-go v0.0.0-20240510082706-db61b37f5877
	github.com/spf13/afero v1.11.0
	golang.org/x/crypto v0.23.0
	golang.org/x/sys v0.20.0
)

require (
	github.com/integrii/flaggy v1.5.2 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/kr/fs v0.1.0 // indirect
	github.com/leonelquinteros/gotext v1.7.1 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.19 // indirect
	github.com/stretchr/testify v1.8.4 // indirect
	golang.org/x/text v0.16.0 // indirect
	gopkg.in/Knetic/govaluate.v3 v3.0.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

replace github.com/saltosystems/winrt-go => ./winrt-go
