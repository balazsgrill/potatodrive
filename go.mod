module github.com/balazsgrill/projfero

go 1.21.0

require (
	github.com/balazsgrill/projfs v0.0.0
	github.com/google/uuid v1.6.0
	github.com/spf13/afero v1.11.0
)

require golang.org/x/text v0.14.0 // indirect

replace github.com/balazsgrill/projfs => ../projfs
