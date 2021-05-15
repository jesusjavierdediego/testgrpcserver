package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"github.com/mahendrabagul/devsecops-meetup/employee"
	pb "github.com/mahendrabagul/devsecops-meetup/employee"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
	"io/ioutil"
	"log"
	"net"
	"os"
	"path/filepath"
)

const (
	port         = ":50051"
	JsonFileName = "database/employees.json"
	CertChain    = "grpc-root-ca-and-grpc-server-ca-and-grpc-client-ca-chain.crt"
	ServerCert   = "grpc-server.crt"
	ServerKey    = "grpc-server.key"
)

type server struct {
	employee.UnimplementedEmployeeServiceServer
}

// GetDetails implements employee.GetDetails
func (s *server) GetDetails(ctx context.Context, in *pb.EmployeeRequest) (*pb.EmployeeResponse, error) {
	err, invalid := isInvalidCertificate(ctx)
	if invalid {
		return nil, err
	}
	e := s.getEmployeeDetails(in.GetId())
	return &pb.EmployeeResponse{Message: e}, nil
}

func isInvalidCertificate(ctx context.Context) (error, bool) {
	p, ok := peer.FromContext(ctx)
	if !ok {
		err := status.Error(codes.Unauthenticated, "no peer found")
		return err, true
	}
	tlsAuth, ok := p.AuthInfo.(credentials.TLSInfo)
	if !ok {
		err := status.Error(codes.Unauthenticated, "unexpected peer transport credentials")
		return err, true
	}
	if len(tlsAuth.State.VerifiedChains) == 0 || len(tlsAuth.State.VerifiedChains[0]) == 0 {
		err := status.Error(codes.Unauthenticated, "could not verify peer certificate")
		return err, true
	}
	// Check subject common name against configured username
	if !contains(tlsAuth.State.VerifiedChains[0][0].Subject.CommonName) {
		err := status.Error(codes.Unauthenticated, "invalid subject common name : Unauthenticated Client")
		return err, true
	}
	return nil, false
}

func contains(e string) bool {
	var validClients = []string{"node-grpc-client1", "node-grpc-client2", "node-grpc-client3"}
	for _, a := range validClients {
		if a == e {
			return true
		}
	}
	return false
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
	tc := getTlsConfig()
	serverOption := grpc.Creds(credentials.NewTLS(tc))
	s := grpc.NewServer(serverOption)
	employee.RegisterEmployeeServiceServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func getTlsConfig() *tls.Config {
	parent := getParent()
	certificate := getCertificate(parent)
	certPool := getCertPool(parent)
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{certificate},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    certPool,
	}
	return tlsConfig
}

func getCertPool(parent string) *x509.CertPool {
	certPool := x509.NewCertPool()
	bs, err := ioutil.ReadFile(filepath.Join(parent, "certificates", "certificatesChain", CertChain))
	if err != nil {
		log.Fatalf("failed to read certificates chain: %s", err)
	}
	ok := certPool.AppendCertsFromPEM(bs)
	if !ok {
		log.Fatalf("failed to append certs")
	}
	return certPool
}

func getCertificate(parent string) tls.Certificate {
	crt := filepath.Join(parent, "certificates", "serverCertificates", ServerCert)
	key := filepath.Join(parent, "certificates", "serverCertificates", ServerKey)
	certificate, _ := tls.LoadX509KeyPair(crt, key)
	return certificate
}

func getParent() string {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("failed to get wd: %v", err)
	}
	parent := filepath.Dir(wd)
	return parent
}
