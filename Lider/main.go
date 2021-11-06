package main

import (
	"context"
	pb "distribuidos/go-usermsg-grpc/usermsg"
	"log"
	"math/rand"
	"net"
	"time"

	"google.golang.org/grpc"
)

var n_etapa1 int32 = 0
var user_id int32 = 0

const (
	port              = "10.6.40.181: 50000"
	address_pozo      = "10.6.40.181: 50011"
	address_name_node = "10.6.40.181: 50020"
)

func choose_number() {
	rand.Seed(time.Now().UTC().UnixNano())
	n_etapa1 = int32(rand.Intn(4) + 6)
	//return int32(rand.Intn(4) + 6)
}

type UserManagementServer struct {
	pb.UnimplementedLiderServicesServer //UnimplementedLiderServices está en el usermsg_grpc.pb, aquí se debe implementar
}

func (s *UserManagementServer) Play(ctx context.Context, in *pb.Message) (*pb.User, error) { //implementar del método Play
	user_id = user_id + 1
	log.Printf("Jugador %d acepta jugar", user_id)
	return &pb.User{ID: user_id}, nil
}

func (s *UserManagementServer) Etapa1(ctx context.Context, in *pb.Jugada1) (*pb.Resp, error) { //implementacion del método Etapa1
	var bin int32 = 1
	var jugada int32 = in.GetJugada()
	choose_number()
	partida := 1
	log.Printf("El Lider eligió %d y la persona eligio %d", n_etapa1, jugada)
	if jugada >= n_etapa1 {
		bin = 0
	}
	ronda := in.GetRonda() + 1

	conn, err := grpc.Dial(address_name_node, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Did not connect: %v", err)
	}
	defer conn.Close()
	ServiceClient := pb.NewNameNodeClient(conn)
	_, err = ServiceClient.JugadaPlayer(context.Background(), &pb.Jugada{ID: in.GetID(), Jugada: jugada, Ronda: ronda, Etapa: in.GetEtapa()})

	if ronda == 4 {
		if int(in.GetSuma()) < 21 {
			bin = 0
		}
		partida = 0
	}
	return &pb.Resp{Survive: bin, Partida: int32(partida), Juego: 1, Ronda: ronda, Etapa: in.GetEtapa(), Suma: in.GetSuma()}, nil
}

func (s *UserManagementServer) Pozo(ctx context.Context, in *pb.Req) (*pb.Monto, error) {
	conn, err := grpc.Dial(address_pozo, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Did not connect: %v", err)
	}
	defer conn.Close()
	ServiceClient := pb.NewPozoServicesClient(conn)
	quest := in.GetReq()
	r, err := ServiceClient.MontoPozo(context.Background(), &pb.Req{Req: quest})
	conn.Close()
	return &pb.Monto{Monto: r.GetMonto()}, nil
}

func main() {
	listner, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterLiderServicesServer(grpcServer, &UserManagementServer{})
	log.Printf("server listening at %v", listner.Addr())

	if err = grpcServer.Serve(listner); err != nil {
		log.Fatalf("Failed to listen on port 50000: %v", err)
	}
}
