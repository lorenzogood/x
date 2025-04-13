QUOTESTACK_CONFIG_LOCATION ?= config.hjson

export

.PHONY: dev
dev: config.hjson
	go run cmd/quotestack/main.go

.PHONY: test
test:
	go test ./... 

config.hjson:
	cp config.default.hjson config.hjson 
