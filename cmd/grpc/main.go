package main

import (
	"hexa_micro/api"
	"hexa_micro/repository/inmemory"
	mongo "hexa_micro/repository/mongodb"
	"hexa_micro/repository/redis"
	"hexa_micro/serializer/protobuf"
	"log"
	"net"
	"os"
	"strconv"

	shortener "hexa_micro/shotener"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	repo := chooseRepo()
	service := shortener.NewRedirectService(repo)
	grpcHandler := api.NewGRPCHandler(service)

	listener, err := net.Listen("tcp", ":4040")
	if err != nil {
		log.Panicln(err)
	}

	srv := grpc.NewServer()
	protobuf.RegisterShortenerServiceServer(srv, grpcHandler)
	reflection.Register(srv)

	log.Printf("serving grpc on port %d\n", 4040)
	if err := srv.Serve(listener); err != nil {
		log.Panic(err)
	}
}

func chooseRepo() shortener.RedirectRepository {
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
