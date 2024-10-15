package server

import (
	"net/http"
	"time"

	"github.com/apache/thrift/lib/go/thrift"
	"github.com/balazsgrill/potatodrive/bindings/proxy"
	"github.com/spf13/afero"
)

func Handler(fs afero.Fs) func(w http.ResponseWriter, r *http.Request) {
	fsserver := New(fs)
	processor := proxy.NewFilesystemProcessor(fsserver)
	conf := &thrift.TConfiguration{
		ConnectTimeout: time.Second,
		SocketTimeout:  time.Second,

		MaxFrameSize: 1024 * 1024 * 256,

		TBinaryStrictRead:  thrift.BoolPtr(true),
		TBinaryStrictWrite: thrift.BoolPtr(true),
	}
	protocol := thrift.NewTCompactProtocolFactoryConf(conf)
	return thrift.NewThriftHandlerFunc(processor, protocol, protocol)
}
