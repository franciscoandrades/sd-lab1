package main

import (
	"context"
	pb "distribuidos/go-usermsg-grpc/usermsg"
	"io/ioutil"
	"log"
	"net"
	"strconv"

	"google.golang.org/grpc"
)

const (
	port = ":50023"
)

type UserManagementServer struct {
	pb.UnimplementedDataNodeServer
}

func (s *UserManagementServer) RegistrarInfo(ctx context.Context, in *pb.Jugada) (*pb.Check, error) {
	b := []byte(strconv.Itoa(int(in.GetJugada())))
	err := ioutil.WriteFile("DataNode/Jugador_"+strconv.Itoa(int(in.GetID()))+" Ronda_"+in.GetRonda()+".txt", b, 0644)
	if err != nil {
		log.Fatalf("Failed to write in Registro.txt")
	}
	return &pb.Check{Check: "OK"}, nil
}
func main() {

	listner, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterDataNodeServer(grpcServer, &UserManagementServer{})
	log.Printf("server listening at %v", listner.Addr())

	if err = grpcServer.Serve(listner); err != nil {
		log.Fatalf("Failed to listen on port 50011: %v", err)
	}
}
