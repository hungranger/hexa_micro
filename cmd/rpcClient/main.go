package main

import (
	"hexa_micro/pkg/shortenservice/container/logger"
	"hexa_micro/pkg/shortenservice/model"
	"net/rpc"
)

func main() {
	var result []model.Redirect
	var reply model.Redirect

	client, err := rpc.DialHTTP("tcp", "localhost:4040")
	if err != nil {
		logger.Log.Fatalf("%+v", err)
	}

	git := model.Redirect{URL: "https://github.com"}
	vnexpress := model.Redirect{URL: "https://vnexpress.net"}
	goolge := model.Redirect{URL: "https://google.com"}

	err = client.Call("RPCHandler.Store", &git, &reply)
	result = append(result, reply)
	err = client.Call("RPCHandler.Store", &vnexpress, &reply)
	result = append(result, reply)
	err = client.Call("RPCHandler.Store", &goolge, &reply)
	result = append(result, reply)

	for _, item := range result {
		client.Call("RPCHandler.Find", item.Code, &reply)
		logger.Log.Infof("%s => %s\n", reply.URL, reply.Code)
	}
}
