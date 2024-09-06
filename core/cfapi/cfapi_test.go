package cfapi_test

import (
	"testing"

	"github.com/balazsgrill/potatodrive/core/cfapi"
)

func Test_PlatformVersion(t *testing.T) {
	var platforminfo cfapi.CF_PLATFORM_INFO
	hr := cfapi.CfGetPlatformInfo(&platforminfo)
	if hr != 0 {
		t.Fatal(hr)
	}
	t.Logf("platform version: b%d i%d, r%d", platforminfo.BuildNumber, platforminfo.IntegrationNumber, platforminfo.RevisionNumber)
}
