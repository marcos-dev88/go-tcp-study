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

type BodyHTTPConfig interface{
	FormatBody(inputData <-chan []byte, outputData chan string, newWg *sync.WaitGroup)
	GetHTTPPostBody(inputData chan []byte, newWg *sync.WaitGroup)
	GetJsonData(body <-chan string, bodyJSON chan []byte, newWg *sync.WaitGroup)
}

type Config struct {
	Conn net.Conn
}

func NewConfig(conn net.Conn) Config {
	return Config{Conn: conn}
}


// FormatBody - Receiving all HTTP POST Request data,
// it "filter" the body from request and send it to an output string channel
func (c Config) FormatBody(inputData <-chan []byte, outputData chan string, newWg *sync.WaitGroup) {
	defer newWg.Done()

	var rgex = regexp.MustCompile(`Length: \d+(.*)`)

	stringLine := string(<-inputData)
	dataBody := rgex.FindStringSubmatch(stringLine)
	if dataBody != nil {
		returnData := strings.Replace(dataBody[1], "{", "", -1)
		returnData = strings.Replace(returnData, "}", "", -1)
		outputData <- returnData
	}
	close(outputData)
}

// GetHTTPPostBody - It gets all scanned data sent by HTTP POST request and sends to channel what accepts []byte
func (c Config) GetHTTPPostBody(inputData chan []byte, newWg *sync.WaitGroup) {
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

// GetJsonData - It gets data from body already formatted and sends data to channel what accepts []byte
func (c Config) GetJsonData(body <-chan string, bodyJSON chan []byte, newWg *sync.WaitGroup) {
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