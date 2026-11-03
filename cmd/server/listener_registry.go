package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"sort"
	"sync"
)

type runningListener struct {
	listener net.Listener
	server   *http.Server
}

type listenerRegistry struct {
	mu        sync.RWMutex
	listeners map[string]runningListener
}

func newListenerRegistry() *listenerRegistry {
	return &listenerRegistry{listeners: make(map[string]runningListener)}
}

// start binds the requested address before serving, so an unavailable port
// fails synchronously and the bound listener is the one used by the server.
func (r *listenerRegistry) start(addr string, handler http.Handler) (string, error) {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return "", fmt.Errorf("bind %s: %w", addr, err)
	}

	actualAddr := ln.Addr().String()
	server := &http.Server{Addr: actualAddr, Handler: handler}

	r.mu.Lock()
	if _, exists := r.listeners[actualAddr]; exists {
		r.mu.Unlock()
		_ = ln.Close()
		return "", fmt.Errorf("listener on %s is already tracked", actualAddr)
	}
	r.listeners[actualAddr] = runningListener{listener: ln, server: server}
	r.mu.Unlock()

	go func() {
		if err := server.Serve(ln); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("HTTP listener on %s stopped: %v", actualAddr, err)
		}
		r.mu.Lock()
		delete(r.listeners, actualAddr)
		r.mu.Unlock()
	}()

	return actualAddr, nil
}

func (r *listenerRegistry) addresses() []string {
	r.mu.RLock()
	addresses := make([]string, 0, len(r.listeners))
	for addr := range r.listeners {
		addresses = append(addresses, addr)
	}
	r.mu.RUnlock()
	sort.Strings(addresses)
	return addresses
}
