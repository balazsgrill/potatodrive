module github.com/balazsgrill/projfero

go 1.21.0

require (
	github.com/aws/aws-sdk-go v1.54.20
	github.com/balazsgrill/projfs v0.0.2
	github.com/fclairamb/afero-s3 v0.3.1
	github.com/google/uuid v1.6.0
	github.com/spf13/afero v1.11.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	golang.org/x/text v0.14.0 // indirect
)

//replace github.com/balazsgrill/projfs => ../projfs
