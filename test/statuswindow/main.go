package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/balazsgrill/potatodrive/core/tasks"
	"github.com/balazsgrill/potatodrive/ui"
	"github.com/rs/zerolog"
)

func main() {
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	output.FormatLevel = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("| %-6s|", i))
	}
	output.FormatMessage = func(i interface{}) string {
		return fmt.Sprintf("***%s****", i)
	}
	output.FormatFieldName = func(i interface{}) string {
		return fmt.Sprintf("%s:", i)
	}
	output.FormatFieldValue = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("%s", i))
	}

	log := zerolog.New(output).With().Timestamp().Logger()
	list := ui.NewTaskListModel()
	list.TaskStateListener(tasks.TaskState{
		ID: 1, Progress: 0,
		Name: "name1"})
	uicontext := &ui.UIContext{
		Logger:  log,
		Version: "Version",
	}
	ui.CreateTaskListWindow(uicontext, list)
	uicontext.MainWindow.Show()
	uicontext.MainWindow.Run()
}
