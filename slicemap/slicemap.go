package slicemap

type kv struct {
	k []byte
	v []byte
}

type sliceMap []kv

// zero alloc
func (sm *sliceMap) Add(k, v []byte) {
	kvs := *sm
	if cap(kvs) > len(kvs) {
		kvs = kvs[:len(kvs)+1]
	} else {
		kvs = append(kvs, kv{})
	}
	kv := &kvs[len(kvs)-1]
	kv.k = append(kv.k[:0], k...)
	kv.v = append(kv.v[:0], v...)
	*sm = kvs
}
