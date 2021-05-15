package main

import (
	"context"
	"encoding/json"
	"github.com/mahendrabagul/devsecops-meetup/employee"
	pb "github.com/mahendrabagul/devsecops-meetup/employee"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io/ioutil"
	"log"
	"net"
	"os"
	"path/filepath"
)

const (
	port         = ":50051"
	JsonFileName = "database/employees.json"
)

type server struct {
	employee.UnimplementedEmployeeServiceServer
}

// GetDetails implements employee.GetDetails
func (s *server) GetDetails(_ context.Context, in *pb.EmployeeRequest) (*pb.EmployeeResponse, error) {
	e := s.getEmployeeDetails(in.GetId())
	return &pb.EmployeeResponse{Message: e}, nil
}

func (s *server) getEmployeeDetails(id int32) *pb.EmployeeDetails {
	list := s.getAllEmployeeDetails()
	var e *pb.EmployeeDetails
	for i := 0; i < len(list); i++ {
		if id == list[i].Id {
			e = list[i]
			break
		}
	}
	return e
}

func (s *server) getAllEmployeeDetails() []*pb.EmployeeDetails {
	jsonFile, err := os.Open(JsonFileName)
	if err != nil {
		panic(err)
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var list []*pb.EmployeeDetails
	_ = json.Unmarshal(byteValue, &list)
	return list
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("failed to get wd: %v", err)
	}
	parent := filepath.Dir(wd)
	crt := filepath.Join(parent, "certificates", "serverCertificates", "grpc-server.crt")
	key := filepath.Join(parent, "certificates", "serverCertificates", "grpc-server.key")

	// Create the TLS credentials
	creds, err := credentials.NewServerTLSFromFile(crt, key)
	if err != nil {
		log.Fatalf("could not load TLS keys: %s", err)
	}

	s := grpc.NewServer(grpc.Creds(creds))

	employee.RegisterEmployeeServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
