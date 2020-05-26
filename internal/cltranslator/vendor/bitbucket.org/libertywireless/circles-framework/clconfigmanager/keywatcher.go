package clconfigmanager

import (
	"context"
)

type KeyWatcher interface {
	OnChange(func(context.Context, string, Value))
	Execute(context.Context, string, Value)
}

type keyWatcher struct {
	key              string
	onChangeCallback func(context.Context, string, Value)
}

func NewKeyWatcher(key string) KeyWatcher {
	return &keyWatcher{
		key: key,
	}
}

func (keyWatcher *keyWatcher) OnChange(fn func(context.Context, string, Value)) {
	keyWatcher.onChangeCallback = fn
}

func (keyWatcher *keyWatcher) Execute(ctx context.Context, key string, val Value) {
	if keyWatcher.onChangeCallback != nil {
		keyWatcher.onChangeCallback(ctx, key, val)
	}
}
