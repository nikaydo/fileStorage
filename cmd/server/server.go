package main

import (
	"bufio"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/nikaydo/fileStorage/internal/server"
)

func main() {
	err := loadenv(".env")
	if err != nil {
		log.Fatalln(err)
	}
	srv, err := server.ServerInit(os.Getenv("PORT"), os.Getenv("ADDR"))
	if err != nil {
		log.Fatalln(err)
	}
	os.Mkdir(os.Getenv("STORAGEPATH"), 0755)
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
		_, err = os.Stat(filepath.Join(os.Getenv("STORAGEPATH"), string(file)))
		if os.IsNotExist(err) {
			srv.Write([]byte("error: " + err.Error()))
		}
		srv.WriteFile(filepath.Join(os.Getenv("STORAGEPATH"), string(file)))
	case "list":
		n, err := os.ReadDir(os.Getenv("STORAGEPATH"))
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
		_, err = os.Stat(filepath.Join(os.Getenv("STORAGEPATH"), string(file)))
		if os.IsNotExist(err) {
			srv.Write([]byte("error: " + err.Error()))
		}
		if err := os.Remove(filepath.Join(os.Getenv("STORAGEPATH"), string(file))); err != nil {
			srv.Write([]byte("error: " + err.Error()))
		}
		srv.Write(append([]byte("deleted file: "), file...))
	case "send":
		file, err := srv.Read()
		if err != nil {
			log.Println(err)
			return
		}
		path := filepath.Join(os.Getenv("STORAGEPATH"), string(file))
		err = os.WriteFile(path, file, 0755)
		if err != nil {
			log.Fatal(err)
		}
		srv.Write(append([]byte("successful sended: "), file...))
	}
}

func loadenv(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		if os.Getenv(key) != "" {
			continue
		}

		os.Setenv(key, value)
	}

	return scanner.Err()
}
