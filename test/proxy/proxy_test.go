package test

import (
	"net/http"
	"testing"

	"github.com/balazsgrill/potatodrive/bindings/proxy/client"
	"github.com/balazsgrill/potatodrive/bindings/proxy/server"
	"github.com/spf13/afero"
)

func TestProxyConnection(t *testing.T) {
	fs := afero.NewMemMapFs()
	mux := http.NewServeMux()
	mux.HandleFunc("/", server.Handler(fs))
	httpserver := http.Server{
		Addr:    "localhost:18080",
		Handler: mux,
	}
	go httpserver.ListenAndServe()
	defer httpserver.Close()

	clientconifg := &client.Config{
		URL:       "http://localhost:18080",
		KeyId:     "",
		KeySecret: "",
	}
	fs2, err := clientconifg.ToFileSystem()
	if err != nil {
		t.Fatal(err)
	}

	data := []byte("something")
	filename := "test.txt"

	err = afero.WriteFile(fs2, filename, data, 0644)
	if err != nil {
		t.Fatal(err)
	}

	data2, err := afero.ReadFile(fs, filename)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != string(data2) {
		t.Fatal("data mismatch")
	}
}
