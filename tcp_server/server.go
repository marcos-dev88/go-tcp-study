package tcp_server

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"net"
	"regexp"
	"strings"
	"sync"
)

type Message struct {
	Value     string `json:"value"`
	IsWorking bool   `json:"is_working"`
}

func Handle(conn net.Conn) {

	byteDataCH, formattedBodyCH := make(chan []byte, 3072), make(chan string, 2)

	wg := sync.WaitGroup{}
	wg.Add(2)

	defer wg.Wait()
	defer conn.Close()

	go getByteDataFromConn(conn, byteDataCH, &wg)
	go formatData(byteDataCH, formattedBodyCH, &wg)

	body := getDataFromBody(<-formattedBodyCH)

	if len(body) != 0 {
		fmt.Printf("%v | %T", body["value"], body["value"])
	}

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

func getDataFromBody(body string) map[string]interface{} {
	dataSlice := strings.Split(body, ",")

	var bodyMap = make(map[string]interface{})
	var rgex = regexp.MustCompile(`("\w+"): (.*)`)

	for i := 0; i < len(dataSlice); i++ {
		data := rgex.FindStringSubmatch(dataSlice[i])
		if data != nil {
			key := strings.Replace(data[1], `"`, "", -1)
			bodyMap[key] = data[2]
		}
	}
	return bodyMap
}
