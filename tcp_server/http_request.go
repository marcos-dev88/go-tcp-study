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
		GetBody(inputData chan []byte, newWg *sync.WaitGroup)
		GetJsonBody(body <-chan string, bodyJSON chan []byte, newWg *sync.WaitGroup)
	}

	HTTPData interface {
		GetURL() string
		GetHeaders() map[string]string
	}
)

type ConnHTTP struct {
	Conn net.Conn
}

func NewConnHTTP(conn net.Conn) ConnHTTP {
	return ConnHTTP{Conn: conn}
}


// GetBody - It gets all scanned data sent by HTTP POST request and sends to channel what accepts []byte
func (c ConnHTTP) GetBody(inputData chan []byte, newWg *sync.WaitGroup) {
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

	inputData <- dataBytes
	close(inputData)
}

// GetJsonBody - It gets data from body already formatted and sends data to channel what accepts []byte
func (c ConnHTTP) GetJsonBody(body <-chan string, bodyJSON chan []byte, newWg *sync.WaitGroup) {
	defer newWg.Done()

	dataSlice := strings.Split(<-body, ",")

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
		bodyJSON <- returnedBytes
		defer close(bodyJSON)
	}
}
