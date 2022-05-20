package kvctx

import (
	"context"
	"sync"
)

type kindValueCtx string

const (
	ctxKey kindValueCtx = "dataValueCtx"
)

type dataValueCtxEntry struct {
	data map[string]interface{}
	mu   sync.Mutex
}

func WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, ctxKey, &dataValueCtxEntry{
		data: make(map[string]interface{}),
	})
}

func Add(ctx context.Context, key string, value interface{}) {
	v := ctx.Value(ctxKey)
	if v == nil {
		return
	}
	dataEntry, ok := v.(*dataValueCtxEntry)
	if !ok {
		return
	}
	dataEntry.mu.Lock()
	defer dataEntry.mu.Unlock()
	dataEntry.data[key] = value
}

func Get(ctx context.Context, key string) interface{} {
	v := ctx.Value(ctxKey)
	if v == nil {
		return ctx
	}
	dataEntry, ok := v.(*dataValueCtxEntry)
	if !ok {
		return ctx
	}
	dataEntry.mu.Lock()
	defer dataEntry.mu.Unlock()
	value, ok := dataEntry.data[key]
	if !ok {
		return ctx
	}
	return value
}
