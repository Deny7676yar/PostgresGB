package main

import "sync"

func NewInMemoryGet() *InMemoryGet{
	return &InMemoryGet{
		map[string]string{},
		sync.RWMutex{},
	}
}

type InMemoryGet struct {
	page map[string]string
	lock sync.RWMutex
}

func (i *InMemoryGet) RecordWin(link string){
	i.lock.Lock()
	defer i.lock.Unlock()
	i.page[link]++
}

func (i *InMemoryGet) GetResult(link string) string {
	i.lock.RLock()
	defer i.lock.RUnlock()
	return i.page[link]
}
