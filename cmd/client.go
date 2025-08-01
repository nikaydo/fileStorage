package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/nikaydo/fileStorage/internal/client"
)

func main() {
	file := flag.String("file", "", "send file to server")
	mode := flag.String("mode", "", "action on server")
	flag.Parse()
	if len(*mode) == 0 {
		log.Fatalln("u need chose mode")
		return
	}
	c, err := client.Connect(os.Getenv("PORT"), os.Getenv("ADDR"))
	if err != nil {
		log.Fatalln(err)
	}
	Message := client.Msg{WaitGroup: sync.WaitGroup{}, Clt: c}
	switch *mode {
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
