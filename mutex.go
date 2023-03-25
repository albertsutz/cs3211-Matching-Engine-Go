package main

type Mutex struct {
	mut chan struct{}
} 

func (m *Mutex) lock() {
	<- m.mut
}

func (m *Mutex) unlock() {
	m.mut <- struct{}{}
}

func initMutex() *Mutex {
	mut :=  &Mutex{mut: make(chan struct{}, 1)}
	mut.unlock()
	return mut 
}