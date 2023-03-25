package main

import (
	"fmt"
	"os"
)

type Mutex struct {
	mut chan struct{}
} 

func (m *Mutex) lock() {
	fmt.Fprintf(os.Stderr, "trying lock\n")
	<- m.mut
	fmt.Fprintf(os.Stderr, "finished lock\n")
}

func (m *Mutex) unlock() {
	fmt.Fprintf(os.Stderr, "trying unlock\n")
	m.mut <- struct{}{}
	fmt.Fprintf(os.Stderr, "finished unlock\n")
}

func initMutex() *Mutex {
	mut :=  &Mutex{mut: make(chan struct{}, 1)}
	mut.mut <- struct{}{}
	return mut 
}