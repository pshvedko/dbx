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

type Item interface{}

// Row ...
type Row struct {
	Item
	cas uint64
	rw  sync.RWMutex
	mx  sync.Mutex
	tx  *Tx
}

func (r *Row) acquire(t *Tx) bool {
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

func (r *Row) release(t *Tx) {
	r.mx.Lock()
	defer r.mx.Unlock()
	t.mx.Lock()
	defer t.mx.Unlock()
	r.tx = r.tx.off(t)
	t.tx = nil
}

func (r *Row) lock(t *Tx) {
	if !r.acquire(t) {
		panic(t)
	}
	r.rw.Lock()
}

func (r *Row) unlock(t *Tx) {
	r.release(t)
	r.rw.Unlock()
}

func (r *Row) read(t *Tx) bool {
	ok := r.acquire(t)
	if ok {
		r.rw.RLock()
	}
	return ok
}

func (r *Row) unread(t *Tx) {
	r.release(t)
	r.rw.RUnlock()
}

func (r *Row) committed(tx *Tx) bool {
	if r.read(tx) {
		defer r.unread(tx)
		return r.cas > 0
	}
	return true
}

func (r *Row) get(tx *Tx) (Item, uint64, bool) {
	if r.read(tx) {
		defer r.unread(tx)
		if r.Item != nil && r.cas > 0 {
			return r.Item, r.cas, true
		}
	}
	return nil, 0, false
}
