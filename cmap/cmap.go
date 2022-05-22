package cmap

import (
	"fmt"
	"reflect"
	"sync"
)

type ConcurrentMap struct {
	m         sync.Map
	keyType   reflect.Type
	valueType reflect.Type
}

func NewConcurrentMap(keyTyp, valueTyp reflect.Type) *ConcurrentMap {
	return &ConcurrentMap{
		keyType:   keyTyp,
		valueType: valueTyp,
	}
}

func (cm *ConcurrentMap) Load(key interface{}) (value interface{}, ok bool) {
	if reflect.TypeOf(key) != cm.keyType {
		return
	}
	return cm.m.Load(key)
}

func (cm *ConcurrentMap) Store(key, value interface{}) {
	if reflect.TypeOf(key) != cm.keyType {
		panic(fmt.Errorf("wrong key type: %v", reflect.TypeOf(key)))
	}
	if reflect.TypeOf(value) != cm.valueType {
		panic(fmt.Errorf("wrong value type: %v", reflect.TypeOf(value)))
	}
	cm.m.Store(key, value)
}
