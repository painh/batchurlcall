package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type RequestReturn struct {
	Httpcode int    `json:"httpcode"`
	Body     string `json:"body"`
}

func request(url string) (RequestReturn, error) {
	reader := strings.NewReader(``)
	request, err := http.NewRequest("GET", url, reader)
	if err != nil {
		return RequestReturn{}, err
	}

	client := &http.Client{}
	resp, err := client.Do(request)

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	return RequestReturn{Httpcode: resp.StatusCode, Body: string(data)}, nil
}

func main() {
	b, err := ioutil.ReadFile("./config.json")
	if err != nil {
		fmt.Println(err)
		return
	}

	var result map[string]interface{}
	json.Unmarshal([]byte(b), &result)

	port, _ := strconv.Atoi(result["port"].(string))

	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {

		retMap := make(map[string]RequestReturn)

		if rec, ok := result["urls"].(map[string]interface{}); ok {
			for key, val := range rec {
				ret, _ := request(val.(string))
				retMap[key] = ret
			}
		}

		jsonString, _ := json.Marshal(retMap)

		w.Write([]byte(jsonString))
	})

	fmt.Printf("server listening on port %v\n", port)

	http.ListenAndServe(":"+strconv.Itoa(port), nil)

}
