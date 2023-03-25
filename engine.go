package main

import "C"
import (
	"context"
	"fmt"
	"io"
	"net"
	"os"
	// "time"
)

type Request struct {
	in input
	clientChan chan struct{}
}

type Engine struct{
}

func (e *Engine) processRequest(ctx context.Context) {
	done := make(chan struct{})
	go func() {
		<-ctx.Done()
		close(done)
	}()
	ob := newOrderBook(done) 
	reqChan := make(chan *Request, 10000)

	for {
		req := <-reqChan
		switch req.in.

		
	}
}

func (e *Engine) accept(ctx context.Context, conn net.Conn) {
	go func() {
		<-ctx.Done()
		conn.Close()
	}()
	go e.handleConn(conn)
}

func (e *Engine) handleConn(conn net.Conn, reqChan chan *Request) {
	defer conn.Close()
	clientChan := make(chan struct{})
	for {
		in, err := readInput(conn)
		if err != nil {
			if err != io.EOF {
				_, _ = fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
			}
			return
		}
		request := &Request{in: input, clientChan: clientChan}
		reqChan <- request
		<- clientChan 
	}
}
