module examples

go 1.22.5

replace github.com/EndlessUpHill/goakka/core v0.0.0 => ../core

replace github.com/EndlessUpHill/goakka/redis v0.0.0 => ../redis

replace github.com/EndlessUpHill/goakka/nats v0.0.0 => ../nats

require (
	github.com/EndlessUpHill/goakka/core v0.0.3
	github.com/EndlessUpHill/goakka/nats v0.0.0
	github.com/EndlessUpHill/goakka/redis v0.0.0
)

require (
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/go-redis/redis/v8 v8.11.5 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/klauspost/compress v1.17.2 // indirect
	github.com/nats-io/nats.go v1.37.0 // indirect
	github.com/nats-io/nkeys v0.4.7 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	golang.org/x/crypto v0.18.0 // indirect
	golang.org/x/sys v0.21.0 // indirect
)
