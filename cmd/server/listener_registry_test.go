package main

import (
	"net"
	"net/http"
	"reflect"
	"testing"
)

func TestListenerRegistryBindsAndTracksStartedListener(t *testing.T) {
	registry := newListenerRegistry()
	addr, err := registry.start("127.0.0.1:0", http.NotFoundHandler())
	if err != nil {
		t.Fatalf("start() error = %v", err)
	}

	registry.mu.RLock()
	entry := registry.listeners[addr]
	registry.mu.RUnlock()
	t.Cleanup(func() { _ = entry.server.Close() })

	if got, want := registry.addresses(), []string{addr}; !reflect.DeepEqual(got, want) {
		t.Fatalf("addresses() = %#v, want %#v", got, want)
	}
	if _, err := registry.start(addr, http.NotFoundHandler()); err == nil {
		t.Fatal("starting a second listener on the same address should fail")
	}
}

func TestListenerRegistryDoesNotTrackUnavailablePort(t *testing.T) {
	occupied, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("reserve test port: %v", err)
	}
	defer occupied.Close()

	registry := newListenerRegistry()
	if _, err := registry.start(occupied.Addr().String(), http.NotFoundHandler()); err == nil {
		t.Fatal("start() succeeded on an already occupied port")
	}
	if got := registry.addresses(); len(got) != 0 {
		t.Fatalf("addresses() = %#v, want no tracked listeners", got)
	}
}
