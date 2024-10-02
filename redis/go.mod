module github.com/EndlessUpHill/goakka/redis

go 1.22.5

replace github.com/EndlessUpHill/goakka/core v0.0.0 => ../core

require (
	github.com/EndlessUpHill/goakka/core v0.0.0
	github.com/go-redis/redis/v8 v8.11.5
)

require (
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect

)
