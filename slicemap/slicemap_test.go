package slicemap

import "testing"

func TestSliceMap(t *testing.T) {
	sm := sliceMap{}
	sm.Add([]byte("key"), []byte("value"))
}
