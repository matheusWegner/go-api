package main

import (
	"log"
)

func main() {
	store , err := newPostgressStore()
	if err != nil {
		log.Fatal(err)
	}
	if err := store.init(); err != nil {
		log.Fatal(err)
	}
	server := newAPIServer(":8080",store)
	server.run()
}
