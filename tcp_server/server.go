package tcp_server

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"net"
	"regexp"
	"strings"
)

type Message struct {
	Value     string `json:"value"`
	IsWorking bool   `json:"is_working"`
}

func Handle(conn net.Conn) {

	defer conn.Close()

	var dataBytes = make([]byte, 3072)
	sc := bufio.NewScanner(conn)

	for sc.Scan() {
		dataBytes = append(dataBytes, sc.Bytes()...)
		if bytes.Contains(dataBytes, []byte("}")) {
			break
		}
	}

	formatedLine := formatLine(dataBytes)
	body := getDataFromLine(formatedLine)

	if len(body) != 0 {
		fmt.Printf("%v | %T", body["value"], body["value"])
	}

	if err := sc.Err(); err != nil {
		log.Fatal(err)
	}

}

func formatLine(line []byte) string {
	var rgex = regexp.MustCompile(`Length: \d+(.*)`)

	stringLine := string(line)
	data := rgex.FindStringSubmatch(stringLine)
	if data != nil {
		returnData := strings.Replace(data[1], "{", "", -1)
		returnData = strings.Replace(returnData, "}", "", -1)
		return returnData
	}
	return ""
}

func getDataFromLine(line string) map[string]interface{} {
	s := strings.Split(line, ",")

	var bodyMap = make(map[string]interface{})

	var rgex = regexp.MustCompile(`("\w+"): (.*)`)

	for i := 0; i < len(s); i++ {
		data := rgex.FindStringSubmatch(s[i])
		if data != nil {
			key := strings.Replace(data[1], `"`, "", -1)
			bodyMap[key] = data[2]
		}
	}
	return bodyMap
}
