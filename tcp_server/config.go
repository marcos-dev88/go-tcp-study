package tcp_server

import (
	"regexp"
	"strings"
	"sync"
)

type HTTPConfig interface{
	FormatBody(inputData <-chan []byte, outputData chan string, newWg *sync.WaitGroup)
}

type Config struct {
	connHTTP ConnHTTP
}

func NewConfig(connHTTP ConnHTTP) Config {
	return Config{connHTTP: connHTTP}
}


// FormatBody - Receiving all HTTP POST Request data,
// it "filter" the body from request and send it to an output string channel
func (c Config) FormatBody(newWg *sync.WaitGroup) {
	defer newWg.Done()

	var rgex = regexp.MustCompile(`Length: \d+(.*)`)

	stringLine := string(<-*c.connHTTP.Body)
	dataBody := rgex.FindStringSubmatch(stringLine)
	if dataBody != nil {
		returnData := strings.Replace(dataBody[1], "{", "", -1)
		returnData = strings.Replace(returnData, "}", "", -1)
		*c.connHTTP.StringBody <- returnData
	}
	close(*c.connHTTP.StringBody)
}