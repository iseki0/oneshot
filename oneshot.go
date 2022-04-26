package oneshot

import (
	"context"
	"sync"
)

type Source interface {
	Wait(ctx context.Context, keys ...string) (timeout bool)
}

type Subject interface {
	Emit(key string)
}

func New() (source Source, subject Subject) {
	s := &_S{
		L: &sync.Mutex{},
		M: map[string]map[chan struct{}]struct{}{},
	}
	return s, s
}

type _S struct {
	L sync.Locker
	M map[string]map[chan struct{}]struct{}
}

type _Sub *_S

func (this *_S) Emit(key string) {
	var s []chan struct{}
	this.L.Lock()
	for it := range this.M[key] {
		s = append(s, it)
	}
	this.L.Unlock()
	for _, it := range s {
		close(it)
	}
}
func (this *_S) Wait(ctx context.Context, keys ...string) (timeout bool) {
	if len(keys) == 0 {
		return
	}
	if ctx == nil {
		ctx = context.TODO()
	}
	ch := make(chan struct{})
	this.L.Lock()
	for _, it := range keys {
		if this.M[it] == nil {
			this.M[it] = map[chan struct{}]struct{}{}
		}
		this.M[it][ch] = struct{}{}
	}
	this.L.Unlock()

	select {
	case <-ctx.Done():
		timeout = true
	case <-ch:
		timeout = false
	}

	this.L.Lock()
	for _, it := range keys {
		if this.M[it] == nil {
			continue
		}
		delete(this.M[it], ch)
		if len(this.M[it]) == 0 {
			delete(this.M, it)
		}
	}
	this.L.Unlock()
	return
}
