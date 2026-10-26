.PHONY: all build test

all: build

build:
	mkdir -p bin
	go build -o bin/crux-server ./cmd/server

test:
	mkdir -p tests/bin
	@set -e; \
	for source in tests/*.go; do \
		[ -f "$$source" ] || continue; \
		case "$$source" in *_test.go) continue ;; esac; \
		name=$$(basename "$$source" .go); \
		go build -o "tests/bin/$$name" "$$source"; \
	done
