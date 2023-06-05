package agent

import "sync"

type pollCounter struct {
	value int64
	mutex *sync.Mutex
}

func newPollCounter() *pollCounter {
	return &pollCounter{
		value: 0,
		mutex: &sync.Mutex{},
	}
}

func (p *pollCounter) add() {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.value += 1
}

func (p *pollCounter) reset() {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.value = 0
}

func (p *pollCounter) get() int64 {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	return p.value
}
