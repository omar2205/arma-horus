package filewatcher

import (
	"log"
	"os"

	"github.com/fsnotify/fsnotify"
	"oskr.nl/arma-horus.go/internal/utils"
)

type FileWatcher struct {
	file        *os.File
	fileName    string
	subscribers []chan string
	watcher     *fsnotify.Watcher
}

func NewFileWatcher(fileName string) (*FileWatcher, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	err = watcher.Add(fileName)
	if err != nil {
		return nil, err
	}

	fw := &FileWatcher{
		file:        file,
		fileName:    fileName,
		watcher:     watcher,
		subscribers: make([]chan string, 0),
	}

	go fw.watch()
	return fw, nil
}

// Scans the file and notifies subscribers
func (fw *FileWatcher) watch() {
	for {
		select {
		case event := <-fw.watcher.Events:
			if event.Op == fsnotify.Write {

				lastLine, err := utils.ReadlastLine(fw.fileName)
				if err != nil {
					log.Println("Error reading file", err)
				}

				for _, sub := range fw.subscribers {
					sub <- lastLine
				}
			}

		case err := <-fw.watcher.Errors:
			log.Println(err)
		}
	}
}

func (fw *FileWatcher) ReadFile() ([]byte, error) {
	content, err := os.ReadFile(fw.fileName)
	if err != nil {
		return nil, err
	}

	return content, nil
}

func (fw *FileWatcher) Subscribe() chan string {
	sub := make(chan string)
	fw.subscribers = append(fw.subscribers, sub)

	return sub
}

func (fw *FileWatcher) Unsubscribe(sub chan string) {
	for i, s := range fw.subscribers {
		if s == sub {
			fw.subscribers = append(fw.subscribers[:i], fw.subscribers[i+1:]...)
			break
		}
	}
	close(sub)
}
