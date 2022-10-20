package main

import (
	"fmt"
	"net"
	"strconv"
	"log"
	grpcserver "xqledger/grpcserver/grpc"
	pb "xqledger/grpcserver/protobuf"
	"google.golang.org/grpc"
)

const componentMessage = "Main process"
const port = 50053

func main() {
	go startRecordHistoryService()
	startQueryService()
}

func startRecordHistoryService() {
	port := 50053
	listener, listenerErr := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if listenerErr != nil {
		log.Fatalln(listenerErr, componentMessage, "Starting", "Error")
	}
	log.Println(componentMessage, "Starting", "Starting RecordHistoryService gRPC  on port "+strconv.Itoa(port))
	service := pb.RecordHistoryServiceServer(&grpcserver.RecordHistoryService{})
	server := grpc.NewServer()
	pb.RegisterRecordHistoryServiceServer(server, service)

	if err := server.Serve(listener); err != nil {
		log.Fatalln(listenerErr, componentMessage, "Grpc RecordHistoryService start", "Error")
	}
}


func startQueryService() {
	port := 50054
	listener, listenerErr := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if listenerErr != nil {
		log.Fatalln(listenerErr, componentMessage, "Starting", "Error")
	}
	log.Println(componentMessage, "Starting", "Starting RecordQueryService gRPC services on port "+strconv.Itoa(port))
	service := pb.RDBQueryServiceServer(&grpcserver.RecordQueryService{})
	server := grpc.NewServer()
	pb.RegisterRDBQueryServiceServer(server, service)

	if err := server.Serve(listener); err != nil {
		log.Fatalln(listenerErr, componentMessage, "Grpc RecordQueryService start", "Error")
	}
}