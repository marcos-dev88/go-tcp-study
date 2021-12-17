package tcp_server

import (
	"bufio"
	"bytes"
	"encoding/json"
	"log"
	"net"
	"regexp"
	"strings"
	"sync"
)

type (
	HTTPPost interface {
		GetBody(newWg *sync.WaitGroup)
		GetJsonBody(newWg *sync.WaitGroup)
	}

	HTTPData interface {
		GetURL() string
		GetHeaders() map[string]string
	}
)

type ConnHTTP struct {
	Conn       net.Conn
	Body       *chan []byte
	JsonBody   *chan []byte
	StringBody *chan string
}

func NewConnHTTP(conn net.Conn, body *chan []byte, jsonBody *chan []byte, stringBody *chan string) *ConnHTTP {
	return &ConnHTTP{Conn: conn, Body: body, JsonBody: jsonBody, StringBody: stringBody}
}

// GetBody - It gets all scanned data sent by HTTP POST request and sends to channel what accepts []byte
func (c *ConnHTTP) GetBody(newWg *sync.WaitGroup) {
	defer newWg.Done()

	sc := bufio.NewScanner(c.Conn)

	var dataBytes = make([]byte, 3072)
	for sc.Scan() {
		dataBytes = append(dataBytes, sc.Bytes()...)
		if bytes.Contains(dataBytes, []byte("}")) {
			break
		}
	}

	if err := sc.Err(); err != nil {
		log.Fatal(err)
	}

	*c.Body <- dataBytes
	close(*c.Body)
}

// GetJsonBody - It gets data from body already formatted and sends data to channel what accepts []byte
func (c *ConnHTTP) GetJsonBody(newWg *sync.WaitGroup) {
	defer newWg.Done()

	dataSlice := strings.Split(<-*c.StringBody, ",")

	var bodyMap = make(map[string]interface{})
	var rgex = regexp.MustCompile(`("\w+"): (.*)`)

	for i := 0; i < len(dataSlice); i++ {
		data := rgex.FindStringSubmatch(dataSlice[i])
		if data != nil {
			key := strings.Replace(data[1], `"`, "", -1)
			value := strings.Replace(data[2], `"`, "", -1)
			bodyMap[key] = value
		}
	}

	if len(bodyMap) != 0 {
		returnedBytes, _ := json.Marshal(bodyMap)
		*c.JsonBody <- returnedBytes
		defer close(*c.JsonBody)
	}
}
