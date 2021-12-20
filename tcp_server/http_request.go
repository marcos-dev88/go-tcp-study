package tcp_server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"regexp"
	"strings"
)

type (
	HTTPJson interface {
		GetJsonBody() []byte
	}

	DataConnHTTP interface {
		GetHTTPData() HTTPData
	}

	HTTPParams interface {
		GetURL() string
		GetBody() []byte
		GetHeaders() map[string]string
	}
)

type (
	ConnHTTP struct {
		Conn net.Conn
	}

	HTTPData struct {
		Data []byte
	}
)

func NewConnHTTP(conn net.Conn) ConnHTTP {
	return ConnHTTP{Conn: conn}
}

func NewHTTPRequest(data []byte) HTTPData {
	return HTTPData{Data: data}
}

// GetHTTPData - It gets all scanned data sent by HTTP POST request and sends to channel what accepts []byte
func (c ConnHTTP) GetHTTPData() HTTPData {
	buffer := make([]byte, 3072)

	for {
		n, err := c.Conn.Read(buffer)

		if err != nil {
			n = 0
			if errors.Is(err, io.EOF) {
				break
			}
			log.Printf("Read error: %s", err)
		}

		if len(buffer) > n {
			break
		}

	}
	return NewHTTPRequest(buffer)
}

func (h HTTPData) GetURL() string {

	httpData := h.GetBody()

	fmt.Println(httpData)

	return ""
}

func (h HTTPData) GetBody() []byte {
	var rgex = regexp.MustCompile(`{([\s\S]*)$`)

	stringLine := string(h.Data)

	dataBody := rgex.FindStringSubmatch(stringLine)
	if dataBody != nil {
		returnData := strings.Replace(dataBody[1], "{", "", -1)
		returnData = strings.Replace(returnData, "}", "", -1)
		return []byte(returnData)
	}
	return nil
}

// GetJsonBody - It gets data from body already formatted and sends data to channel what accepts []byte
func (h HTTPData) GetJsonBody() []byte {

	stringedBody := string(h.GetBody())

	dataSlice := strings.Split(stringedBody, ",")
	var bodyMap = make(map[string]interface{})
	var bodyMapChan = make(chan map[string]interface{})

	var rgex = regexp.MustCompile(`("\w+"): (.*)`)

	go func() {
		for i := 0; i < len(dataSlice); i++ {
			data := rgex.FindStringSubmatch(dataSlice[i])
			if data != nil {
				key := strings.Replace(data[1], `"`, "", -1)
				value := strings.Replace(data[2], `"`, "", -1)
				bodyMap[key] = value
			}
		}
		bodyMapChan <- bodyMap
		close(bodyMapChan)
	}()

	returnedBytes, _ := json.Marshal(<-bodyMapChan)

	return returnedBytes
}

func (h HTTPData) GetHeaders() map[string]string {
	some := map[string]string{}
	return some
}
