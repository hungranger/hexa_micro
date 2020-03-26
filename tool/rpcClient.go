package main

import (
	"fmt"
	"hexa_micro/model"
	"log"
	"net/rpc"
)

func main() {
	var result []model.Redirect
	var reply model.Redirect

	client, err := rpc.DialHTTP("tcp", "localhost:4040")
	if err != nil {
		log.Fatal("Connection Error: ", err)
	}

	git := model.Redirect{URL: "https://github.com"}
	vnexpress := model.Redirect{URL: "https://vnexpress.net"}
	goolge := model.Redirect{URL: "https://google.com"}

	err = client.Call("SimpleRPC.Store", &git, &reply)
	result = append(result, reply)
	err = client.Call("SimpleRPC.Store", &vnexpress, &reply)
	result = append(result, reply)
	err = client.Call("SimpleRPC.Store", &goolge, &reply)
	result = append(result, reply)

	for _, item := range result {
		client.Call("SimpleRPC.Find", item.Code, &reply)
		fmt.Printf("%s => %s\n", reply.Code, reply.URL)
	}
}
