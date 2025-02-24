package collection

import "sync"

// Tx ...
type Tx struct {
	mode int8
	mx   sync.Mutex
	tx   *Tx
}

func (t *Tx) on(x *Tx) bool {
	switch t {
	case nil:
		return false
	case x:
		return true
	}
	t.mx.Lock()
	defer t.mx.Unlock()
	return t.tx.on(x)
}

func (t *Tx) off(x *Tx) *Tx {
	switch t {
	case nil:
		panic(x)
	case x:
		return x.tx
	}
	t.mx.Lock()
	defer t.mx.Unlock()
	t.tx = t.tx.off(x)
	return t
}

// Row ...
type Row[T any] struct {
	T  T
	N  uint64
	rw sync.RWMutex
	mx sync.Mutex
	tx *Tx
}

func (r *Row[T]) acquire(t *Tx) bool {
	r.mx.Lock()
	defer r.mx.Unlock()
	if r.tx.on(t) {
		return false
	}
	t.mx.Lock()
	defer t.mx.Unlock()
	t.tx = r.tx
	r.tx = t
	return true
}

func (r *Row[T]) release(t *Tx) {
	r.mx.Lock()
	defer r.mx.Unlock()
	t.mx.Lock()
	defer t.mx.Unlock()
	r.tx = r.tx.off(t)
	t.tx = nil
}

func (r *Row[T]) lock(t *Tx) {
	if !r.acquire(t) {
		panic(t)
	}
	r.rw.Lock()
}

func (r *Row[T]) unlock(t *Tx) {
	r.release(t)
	r.rw.Unlock()
}

func (r *Row[T]) read(t *Tx) bool {
	ok := r.acquire(t)
	if ok {
		r.rw.RLock()
	}
	return ok
}

func (r *Row[T]) unread(t *Tx) {
	r.release(t)
	r.rw.RUnlock()
}

func (r *Row[T]) committed(tx *Tx) bool {
	if r.read(tx) {
		defer r.unread(tx)
		return r.N > 0
	}
	return true
}

func (r *Row[T]) get(tx *Tx) (T, uint64, bool) {
	var z T
	if r.read(tx) {
		defer r.unread(tx)
		if r.N > 0 {
			return r.T, r.N, true
		}
	}
	return z, 0, false
}
