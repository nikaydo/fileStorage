package client

import (
	"fmt"
	"log"
	"sync"
)

type Msg struct {
	sync.WaitGroup
	Clt Client
}

func (m *Msg) RunGorutine() {
	m.Add(1)
	go func() {
		text, err := m.Clt.Read()
		if err != nil {
			log.Println(err)
			return
		}
		fmt.Println(string(text))
		m.Done()
	}()
}
