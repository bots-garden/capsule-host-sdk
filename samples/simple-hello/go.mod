module hostlog

go 1.20

require github.com/tetratelabs/wazero v1.2.0

require github.com/bots-garden/capsule-host-sdk v0.0.2

require (
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/go-resty/resty/v2 v2.7.0 // indirect
	github.com/redis/go-redis/v9 v9.0.4 // indirect
	golang.org/x/net v0.0.0-20211029224645-99673261e6eb // indirect
)

replace github.com/bots-garden/capsule-host-sdk => ../..
