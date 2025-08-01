package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/nikaydo/fileStorage/internal/server"
)

const (
	storagePath = "storage"
)

func main() {
	srv, err := server.ServerInit("9000", "localhost")
	if err != nil {
		log.Fatalln(err)
	}
	os.Mkdir(storagePath, 0755)
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
	mode, err := srv.Read()
	if err != nil {
		log.Println(err)
		return
	}
	switch string(mode) {
	case "get":
		file, err := srv.Read()
		if err != nil {
			log.Println(err)
			return
		}
		_, err = os.Stat(filepath.Join(storagePath, string(file)))
		if os.IsNotExist(err) {
			srv.Write([]byte("error: " + err.Error()))
		}
		srv.WriteFile(filepath.Join(storagePath, string(file)))
	case "list":
		n, err := os.ReadDir(storagePath)
		if err != nil {
			log.Println(err)
			return
		}
		var list string
		for i, file := range n {
			if i == len(n)-1 {
				list += file.Name()
				continue
			}
			list += file.Name() + "\n"
		}
		srv.Write([]byte(list))
	case "delete":
		file, err := srv.Read()
		if err != nil {
			log.Println(err)
			return
		}
		_, err = os.Stat(filepath.Join(storagePath, string(file)))
		if os.IsNotExist(err) {
			srv.Write([]byte("error: " + err.Error()))
		}
		if err := os.Remove(filepath.Join(storagePath, string(file))); err != nil {
			srv.Write([]byte("error: " + err.Error()))
		}
		srv.Write(append([]byte("deleted file: "), file...))
	case "send":
		file, err := srv.Read()
		if err != nil {
			log.Println(err)
			return
		}
		path := filepath.Join(storagePath, string(file))
		err = os.WriteFile(path, file, 0755)
		if err != nil {
			log.Fatal(err)
		}
		srv.Write(append([]byte("successful sended: "), file...))
	}
}
