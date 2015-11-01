package contextutil

import (
	"golang.org/x/net/context"
)

// implementation pulled from https://go-review.googlesource.com/#/c/10910/
func NewKey(str string) interface{} {
	return (*key)(&str)
}

type key string

func (k *key) String() string { return string(*k) }

type ValueStore interface {
	Get(context.Context) interface{}
	Set(context.Context, interface{}) context.Context
}

func NewValueStore(key interface{}) ValueStore {
	return &valueStore{key}
}

type valueStore struct {
	key interface{}
}

func (c *valueStore) Get(ctxt context.Context) interface{} {
	return ctxt.Value(c.key)
}

func (c *valueStore) Set(ctxt context.Context, val interface{}) context.Context {
	return context.WithValue(ctxt, c.key, val)
}
