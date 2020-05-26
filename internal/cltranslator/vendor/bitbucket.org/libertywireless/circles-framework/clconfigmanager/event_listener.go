package clconfigmanager

import "sync"

type EventListeners interface {
	Get(string) (KeyWatcher, bool)
	Store(string, KeyWatcher)
}

type eventListeners struct {
	watchers sync.Map
}

func NewEventListeners() EventListeners {
	return &eventListeners{
		watchers: sync.Map{},
	}
}

func (eventListeners *eventListeners) Get(key string) (KeyWatcher, bool) {
	val, _ := eventListeners.watchers.Load(key)
	keyWatcher, ok := val.(KeyWatcher)
	return keyWatcher, ok
}

func (eventListeners *eventListeners) Store(key string, keyWatcher KeyWatcher) {
	eventListeners.watchers.Store(key, keyWatcher)
}
