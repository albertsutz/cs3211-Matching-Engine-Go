package main

import (
	// "fmt"
	// "os"
)

type Mutex struct {
	mut chan struct{}
} 

func (m *Mutex) lock() {
	// fmt.Fprintf(os.Stderr, "trying lock %v\n", m.name)
	<- m.mut
	// fmt.Fprintf(os.Stderr, "finished lock %v\n",m.name)
}

func (m *Mutex) unlock() {
	// fmt.Fprintf(os.Stderr, "trying unlock %v\n",m.name)
	m.mut <- struct{}{}
	// fmt.Fprintf(os.Stderr, "finished unlock %v\n",m.name)
}

func initMutex() *Mutex {
	mut :=  &Mutex{mut: make(chan struct{}, 1)}
	mut.mut <- struct{}{}
	return mut 
}