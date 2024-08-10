package filesystem

import (
	"log"

	"github.com/fsnotify/fsnotify"
)

func (instance *VirtualizationInstance) watch() {
	log.Println("Watching for changes")
	go func() {
		for err := range instance.watcher.Errors {
			log.Printf("Received error: %s", err)
		}
	}()
	for event := range instance.watcher.Events {
		log.Printf("Received event: %s", event)
		if event.Op&fsnotify.Remove == fsnotify.Remove {
			instance.handleDeletion(event.Name)
		}
	}
}
