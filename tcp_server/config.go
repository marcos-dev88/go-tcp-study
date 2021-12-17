package tcp_server

import (
	"net"
	"regexp"
	"strings"
	"sync"
)

type HTTPConfig interface{
	FormatBody(inputData <-chan []byte, outputData chan string, newWg *sync.WaitGroup)
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