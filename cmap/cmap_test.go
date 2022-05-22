package cmap

import (
	"testing"
)

func Test_cmap_store(t *testing.T) {
	cm := &ConcurrentMap{}

	tests := []struct {
		name  string
		key   interface{}
		value interface{}
	}{
		{
			"test-string",
			"key",
			"value",
		},
		{
			"test-int",
			1,
			2,
		},
		{
			"test-p",
			"pp",
			&ConcurrentMap{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cm.m.Store(tt.key, tt.value)
			_, ok := cm.m.Load(tt.key)
			if !ok {
				t.Errorf("%s load error", tt.name)
			}
		})
	}
}
