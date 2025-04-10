package ui

import (
	"io"
	"os"
	"path/filepath"

	"github.com/leonelquinteros/gotext"
	"github.com/lxn/walk"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type UIContext struct {
	*gotext.Po
	Logger  zerolog.Logger
	LogFile string
	Version string
	logf    io.Closer

	*walk.MainWindow
}

func NewUIContext(version string) *UIContext {
	logfilepath, logger, logf := initLogger()
	return &UIContext{
		Po:      GetLocalization(),
		Logger:  logger,
		LogFile: logfilepath,
		logf:    logf,
		Version: version,
	}
}

func (context *UIContext) Close() error {
	return context.logf.Close()
}

func initLogger() (string, zerolog.Logger, io.Closer) {
	cachedir, err := os.UserCacheDir()
	if err != nil {
		panic(err)
	}
	logfolder := filepath.Join(cachedir, "PotatoDrive")
	err = os.MkdirAll(logfolder, 0777)
	if err != nil {
		panic(err)
	}

	logfile := "potatodrive.log"
	logfilepath := filepath.Join(logfolder, logfile)
	logf, err := os.OpenFile(logfilepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}
	return logfilepath, log.Output(zerolog.MultiLevelWriter(logf, zerolog.NewConsoleWriter())).With().Timestamp().Logger(), logf
}
