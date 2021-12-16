package tcp_server

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"sync"
)

type Message struct {
	Value     string `json:"value"`
	IsWorking bool   `json:"is_working,string"`
}

func Handle(conn net.Conn) {

	var m Message
	byteDataCH, stringBodyCH, bodyJsonData := make(chan []byte, 3072), make(chan string, 2), make(chan []byte, 3072)

	wg := sync.WaitGroup{}
	wg.Add(3)
	defer wg.Wait()

	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(conn)

	config := NewConfig(conn)

	go config.GetHTTPPostBody(byteDataCH, &wg)
	go config.FormatBody(byteDataCH, stringBodyCH, &wg)
	go config.GetJsonData(stringBodyCH, bodyJsonData, &wg)

	select {
	case returnedBody := <-bodyJsonData:
		err := json.Unmarshal(returnedBody, &m)
		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Printf("value -> %v | type -> %T\n", m.IsWorking, m.IsWorking)
}
