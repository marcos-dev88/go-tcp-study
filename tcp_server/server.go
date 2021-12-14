package tcp_server

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"regexp"
	"strings"
	"sync"
)

type Message struct {
	Value     string `json:"value"`
	IsWorking string   `json:"is_working"`
}

func Handle(conn net.Conn) {

	byteDataCH, formattedBodyCH, bodyJsonData := make(chan []byte, 3072), make(chan string, 2), make(chan []byte, 3072)

	wg := sync.WaitGroup{}
	wg.Add(3)
	defer wg.Wait()

	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(conn)

	go getByteDataFromConn(conn, byteDataCH, &wg)
	go formatData(byteDataCH, formattedBodyCH, &wg)
	go getDataFromBody(formattedBodyCH, bodyJsonData, &wg)

	var m Message
	err := json.Unmarshal(<-bodyJsonData, &m)

	if err != nil {

	}

	fmt.Printf("%v\n", m)

}

func formatData(inputData <-chan []byte, outputData chan string, newWg *sync.WaitGroup) {
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

func getByteDataFromConn(conn net.Conn, inputData chan []byte, newWg *sync.WaitGroup) {
	defer newWg.Done()

	sc := bufio.NewScanner(conn)

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

func getDataFromBody(body <-chan string, bodyData chan []byte, newWg *sync.WaitGroup) {
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
		bodyData <- returnedBytes
		defer close(bodyData)
	}
}
