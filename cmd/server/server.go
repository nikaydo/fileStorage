package main

import (
	"fmt"
	"log"

	"github.com/nikaydo/fileStorage/internal/server"
)

func main() {
	srv, err := server.ServerInit("9000", "localhost")
	if err != nil {
		log.Fatalln(err)
	}
	for {
		err := srv.Connect()
		if err != nil {
			log.Println(err.Error())
			continue
		}
		go hand(&srv)
	}

}

func hand(srv *server.Server) {
	b, err := srv.Read()
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println(string(b))
}
