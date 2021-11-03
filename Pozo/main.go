package main

import (
	"context"
	pb "distribuidos/go-usermsg-grpc/usermsg"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"strconv"
	"strings"

	"google.golang.org/grpc"
)

const (
	port = ": 50011"
)

var pozo = 30000

func leer_pozo() int {
	Bytes, err := ioutil.ReadFile("Pozo/pozo.txt")
	if err != nil {
		log.Fatal(err)
	}
	datos := string(Bytes)
	array := strings.Fields(datos)
	pozo_string := array[len(array)-1]
	fmt.Printf("%q", array)
	pozo, err = strconv.Atoi(pozo_string)
	return pozo

}

type UserManagementServer struct {
	pb.UnimplementedPozoServicesServer
}

func (s *UserManagementServer) MontoPozo(ctx context.Context, in *pb.Req) (*pb.Monto, error) {

	return &pb.Monto{Monto: int32(leer_pozo())}, nil
}

func main() {
	//b:= []byte("Jugador_1 Ronda_1 40000")
	//err := ioutil.WriteFile("pozo.txt", b, 0644)

	listner, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterPozoServicesServer(grpcServer, &UserManagementServer{})
	log.Printf("server listening at %v", listner.Addr())

	if err = grpcServer.Serve(listner); err != nil {
		log.Fatalf("Failed to listen on port 50011: %v", err)
	}
}
