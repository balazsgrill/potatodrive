package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/balazsgrill/potatodrive/filesystem"
	"github.com/spf13/afero"
)

func main() {
	fs := afero.NewBasePathFs(afero.NewOsFs(), "C:\\work\\vfsbase")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	closer, err := filesystem.StartProjecting("C:\\work\\vfs", fs)
	if err != nil {
		log.Panic(err)
	}

	<-c
	closer.Close()
	os.Exit(1)
}
