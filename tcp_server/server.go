package tcp_server

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
)

type Message struct {
	Value     string `json:"value"`
	IsWorking bool   `json:"is_working,string"`
}

func Handle(conn net.Conn) {

	var m Message

	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(conn)

	connHTTP := NewConnHTTP(conn)
	httpData := connHTTP.GetHTTPData()
	body := httpData.GetBody()

	if len(body) == 0 {
		return
	}

	jsonBody := httpData.GetJsonBody()

	if err := json.Unmarshal(jsonBody, &m); err != nil {
		log.Fatal(err)
	}

	fmt.Println(m)

}
