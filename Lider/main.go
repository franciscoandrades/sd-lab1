package main

import (
	"context"
	pb "distribuidos/go-usermsg-grpc/usermsg"
	"log"
	"math/rand"
	"net"
	"strconv"
	"time"

	"google.golang.org/grpc"
)

var rondas_luz_verde int32 = 0
var n_etapa1 int32 = 0
var user_id int32 = 0

const (
	port              = ":50000"
	address_pozo      = "localhost: 50011"
	address_name_node = "localhost: 50020"
)

func choose_number() {
	rand.Seed(time.Now().UTC().UnixNano())
	n_etapa1 = int32(rand.Intn(4) + 6)
}

type UserManagementServer struct {
	pb.UnimplementedLiderServicesServer //UnimplementedLiderServices está en el usermsg_grpc.pb, aquí se debe implementar
}

func (s *UserManagementServer) NewPlayer(ctx context.Context, in *pb.Message) (*pb.User, error) { //implementar el método NewPlayer
	log.Printf("Nombre del Usuario: %v", in.GetName())
	user_id = user_id + 1
	return &pb.User{Name: in.GetName(), ID: user_id}, nil
}

func (s *UserManagementServer) Luz_Roja_Verde(ctx context.Context, in *pb.Jugada_1) (*pb.Resp_1, error) {
	var bin int32 = 1
	var n_persona int32 = in.GetNElegido()
	choose_number()
	rondas_luz_verde = rondas_luz_verde + 1
	log.Printf("El Lider eligió %d y la persona eligio %d", n_etapa1, n_persona)
	if n_persona >= n_etapa1 {
		bin = 0
	}
	conn, err := grpc.Dial(address_name_node, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Did not connect: %v", err)
	}
	defer conn.Close()
	ServiceClient := pb.NewNameNodeClient(conn)
	_, err = ServiceClient.JugadaPlayer(context.Background(), &pb.Jugada{ID: in.ID, Juego: "1", Ronda: strconv.Itoa(int(rondas_luz_verde)), Jugada: in.GetNElegido()})

	conn.Close()
	return &pb.Resp_1{Binario: bin, Ronda: rondas_luz_verde, EstJuego: "Se juega"}, nil
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
