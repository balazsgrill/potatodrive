package client

import (
	"time"

	"github.com/apache/thrift/lib/go/thrift"
	"github.com/balazsgrill/potatodrive/bindings/proxy"
	"github.com/spf13/afero"
)

func Connect(url string) (afero.Fs, error) {
	conf := &thrift.TConfiguration{
		ConnectTimeout: time.Second,
		SocketTimeout:  time.Second,

		MaxFrameSize: 1024 * 1024 * 256,

		TBinaryStrictRead:  thrift.BoolPtr(true),
		TBinaryStrictWrite: thrift.BoolPtr(true),
	}
	protocol := thrift.NewTCompactProtocolFactoryConf(conf)
	clientfactory := thrift.NewTHttpClientTransportFactory(url)
	transport, err := clientfactory.GetTransport(nil)
	if err != nil {
		return nil, err
	}
	client := proxy.NewFilesystemClientFactory(transport, protocol)
	return New(client), nil
}
