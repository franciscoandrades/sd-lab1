package main

import (
	"context"
	pb "distribuidos/go-usermsg-grpc/usermsg"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"time"

	"google.golang.org/grpc"
)

var address_1 = []string{"10.6.40.182:50023", "10.6.40.181:50023", "10.6.40.184:50023"}

const (
	port = "10.6.40.181: 50020"
)

func choose_number() int32 {
	rand.Seed(time.Now().UTC().UnixNano())
	elec := int32(rand.Intn(3))
	return elec
}

type UserManagementServer struct {
	pb.UnimplementedNameNodeServer
}

func (s *UserManagementServer) JugadaPlayer(ctx context.Context, in *pb.Jugada) (*pb.Registro, error) {
	add := address_1[choose_number()]
	ip := strings.Trim(add, ":50023")
	b := []byte("Jugador_" + strconv.Itoa(int(in.GetID())) + " Ronda_" + strconv.Itoa(int(in.GetRonda())) + " " + ip)
	err := ioutil.WriteFile("NameNode/Registro.txt", b, 0644)
	if err != nil {
		log.Fatalf("Failed to write in Registro.txt")
	}
	conn, err := grpc.Dial("10.6.40.181:50023", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Did not connect: %v", err)
	}
	defer conn.Close()
	ServiceClient := pb.NewDataNodeClient(conn)
	_, err = ServiceClient.RegistrarInfo(context.Background(), &pb.Jugada{ID: in.GetID(), Etapa: in.GetEtapa(), Jugada: in.GetJugada()})

	return &pb.Registro{Response: ""}, nil
}

func main() {
	listner, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterNameNodeServer(grpcServer, &UserManagementServer{})
	log.Printf("server listening at %v", listner.Addr())

	if err = grpcServer.Serve(listner); err != nil {
		log.Fatalf("Failed to listen on port 50011: %v", err)
	}
}
