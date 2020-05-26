package option

import (
	"context"
	"reflect"
)

type Options struct {
	// for alternative data
	Context        context.Context
	optionalParams []KV
}

type KV struct {
	Key   interface{}
	Value interface{}
}

// pass option to config struct
type Option func(o *Options)

func NewOptions(opts ...Option) Options {
	options := Options{
		Context: context.Background(),
	}

	for _, o := range opts {
		o(&options)
	}

	return options
}

type OptionalKV func() (interface{}, interface{})

func Opt(key, value interface{}) OptionalKV {
	if key == nil || !reflect.TypeOf(key).Comparable() {
		return nil
	}
	return func() (interface{}, interface{}) {
		return key, value
	}

}
func NewOptionsFromOptional(opts ...OptionalKV) Options {
	options := Options{
		Context: context.Background(),
	}

	for _, o := range opts {
		if o == nil {
			continue
		}
		key, val := o()

		options.optionalParams = append(options.optionalParams, KV{key, val})
		options.Context = context.WithValue(options.Context, key, val)
	}

	return options
}

func (o Options) GetValue(key interface{}) interface{} {
	for _, kv := range o.optionalParams {
		if kv.Key == key {
			return kv.Value
		}
	}
	return nil
}

func (o Options) GetValues() []KV {
	return o.optionalParams
}
