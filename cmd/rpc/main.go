package main

import (
	"hexa_micro/api"
	"hexa_micro/repository"
	"hexa_micro/repository/RedirectRepository/inmemory"
	mongo "hexa_micro/repository/RedirectRepository/mongodb"
	"hexa_micro/repository/RedirectRepository/redis"
	"hexa_micro/usecase/shortenURL"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"strconv"
)

func main() {
	repo := chooseRepo()
	shortenURLUseCase := shortenURL.NewShortenURLUseCase(repo)
	rpcHandler := api.NewRPCHandler(shortenURLUseCase)
	err := rpc.Register(rpcHandler)
	if err != nil {
		log.Fatal("error registering API ", err)
	}

	rpc.HandleHTTP()
	listener, err := net.Listen("tcp", ":4040")
	if err != nil {
		log.Fatal("Listener error ", err)
	}

	log.Printf("serving rpc on port %d", 4040)
	http.Serve(listener, nil)
}

func chooseRepo() repository.IRedirectRepository {
	switch os.Getenv("URL_DB") {
	case "redis":
		redisURL := os.Getenv("REDIS_URL")
		redisPassword := os.Getenv("REDIS_PASSWORD")
		repo, err := redis.NewRedisRepository(redisURL, redisPassword)
		if err != nil {
			log.Fatal(err)
		}
		return repo
	case "mongo":
		mongoURL := os.Getenv("MONGO_URL")
		mongodb := os.Getenv("MONGO_DB")
		mongoTimeout, _ := strconv.Atoi(os.Getenv("MONGO_TIMEOUT"))
		repo, err := mongo.NewMongoRepository(mongoURL, mongodb, mongoTimeout)
		if err != nil {
			log.Fatal(err)
		}
		return repo
	case "inmemory":
		repo := inmemory.NewInmemoryRepository()
		return repo
	default:
		log.Fatal("Cannot choose repository")
	}
	return nil
}
