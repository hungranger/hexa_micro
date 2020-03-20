package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	shortener "hexa_micro/shotener"

	"github.com/vmihailenco/msgpack"
)

func main() {
	address := fmt.Sprintf("http://localhost%s", httpPort())
	redirect := shortener.Redirect{}
	redirect.URL = "https://github.com/tensor-programming?tab=repositories"

	body, err := msgpack.Marshal(&redirect)
	if err != nil {
		log.Fatalln(err)
	}

	resp, err := http.Post(address, "application/x-msgpack", bytes.NewBuffer(body))
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	msgpack.Unmarshal(body, &redirect)

	log.Printf("%v\n", redirect)
}

func httpPort() string {
	port := "8080"
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}
	return fmt.Sprintf(":%s", port)
}
