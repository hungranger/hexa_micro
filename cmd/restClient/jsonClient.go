package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"hexa_micro/pkg/shortenservice/model"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func main() {
	address := fmt.Sprintf("http://localhost%s", httpPort())
	redirect := model.Redirect{}
	redirect.URL = "https://github.com/tensor-programming?tab=repositories"

	body, err := json.Marshal(&redirect)
	if err != nil {
		log.Fatalln(err)
	}

	resp, err := http.Post(address, "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	json.Unmarshal(body, &redirect)

	log.Printf("%v\n", redirect)
}

func httpPort() string {
	port := "8080"
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}
	return fmt.Sprintf(":%s", port)
}
