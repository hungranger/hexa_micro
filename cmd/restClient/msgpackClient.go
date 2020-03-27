package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"hexa_micro/pkg/shortenservice/container/logger"
	"hexa_micro/pkg/shortenservice/model"

	"github.com/vmihailenco/msgpack"
)

func main() {
	address := fmt.Sprintf("http://localhost%s", httpPort())
	redirect := model.Redirect{}
	redirect.URL = "https://github.com/tensor-programming?tab=repositories"

	body, err := msgpack.Marshal(&redirect)
	if err != nil {
		logger.Log.Fatalf("%+v", err)
	}

	resp, err := http.Post(address, "application/x-msgpack", bytes.NewBuffer(body))
	if err != nil {
		logger.Log.Fatalf("%+v", err)
	}
	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Log.Fatalf("%+v", err)
	}

	msgpack.Unmarshal(body, &redirect)

	logger.Log.Infof("%v\n", redirect)
}

func httpPort() string {
	port := "8080"
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}
	return fmt.Sprintf(":%s", port)
}
