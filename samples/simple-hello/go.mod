module hostlog

go 1.20

require github.com/tetratelabs/wazero v1.1.0

require github.com/bots-garden/capsule-host-sdk v0.0.2

require (
	github.com/go-resty/resty/v2 v2.7.0 // indirect
	golang.org/x/net v0.0.0-20211029224645-99673261e6eb // indirect
)

replace github.com/bots-garden/capsule-host-sdk => ../..
