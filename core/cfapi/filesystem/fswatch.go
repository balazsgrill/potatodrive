package filesystem

import (
	"github.com/fsnotify/fsnotify"
)

func (instance *VirtualizationInstance) watch() {
	instance.Logger.Debug().Msg("Watching for changes")
	go func() {
		for err := range instance.watcher.Errors {
			instance.Logger.Debug().Msgf("Received error: %s", err)
		}
	}()
	for event := range instance.watcher.Events {
		instance.Logger.Debug().Msgf("Received event: %s", event)
		if event.Op&fsnotify.Remove == fsnotify.Remove {
			instance.handleDeletion(event.Name)
		}
	}
}
