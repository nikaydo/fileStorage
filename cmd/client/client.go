package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/nikaydo/fileStorage/internal/client"
)

func main() {
	file := flag.String("file", "", "send file to server")
	act := flag.String("mode", "", "action on server")
	port := flag.String("port", "9000", "port to connect the server")
	addr := flag.String("addr", "localhost", "the port through which the connection is made")
	flag.Parse()
	if len(*file) == 0 && len(*act) == 0 {
		log.Fatalln("u need chose file to send or actions")
		return
	}
	c, err := client.Connect(*port, *addr)
	if err != nil {
		log.Fatalln(err)
	}
	Message := client.Msg{WaitGroup: sync.WaitGroup{}, Clt: c}
	switch *act {
	case "send":
		if err := c.Write([]byte("send")); err != nil {
			log.Println(err)
			return
		}
		if err = c.WriteFile(*file); err != nil {
			log.Println(err)
			return
		}
	case "get":
		if err := c.Write([]byte("get")); err != nil {
			log.Println(err)
			return
		}
		if err := c.Write([]byte(*file)); err != nil {
			log.Println(err)
			return
		}
		f, err := c.Read()
		if err != nil {
			log.Println(err)
			return
		}
		err = os.WriteFile(filepath.Base(*file), f, 0755)
		if err != nil {
			log.Println(err)
			return
		}
		fmt.Println("file download")
	case "list":
		Message.RunGorutine()
		if err := c.Write([]byte("list")); err != nil {
			log.Println(err)
			return
		}
		Message.Wait()
	case "delete":
		Message.RunGorutine()
		if err := c.Write([]byte("delete")); err != nil {
			log.Println(err)
			return
		}
		if err := c.Write([]byte(filepath.Base(*file))); err != nil {
			log.Println(err)
			return
		}
		Message.Wait()
	default:
		log.Fatalln("chose actions in act flag is not valid")
	}
}
