package main

import (
	"redis-lite/core"
	"redis-lite/storage"
)

func main() {
	inMemoryStore := storage.NewInMemoryStorage()
	processor := &core.CommandProcessor{Cache: inMemoryStore}
	core.NewServer(core.ServerOpts{Port: 8001}, processor).Start()
}
