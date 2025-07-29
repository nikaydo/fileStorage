package main

import (
	"log"

	"github.com/nikaydo/fileStorage/internal/client"
)

func main() {
	c, err := client.Connect("9000", "localhost")
	if err != nil {
		log.Fatalln(err)
	}
	c.Write("hello")
}
