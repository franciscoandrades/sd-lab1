package main

import (
	"context"
	pb "distribuidos/go-usermsg-grpc/usermsg"
	"fmt"
	"log"
	"math/rand"
	"net"
	"strconv"
	"time"

	"github.com/streadway/amqp"
	"google.golang.org/grpc"
)

var n_etapa1 int32 = 0
var user_id int32 = 0
var cont = 0

const (
	port         = ":50000"
	address_pozo = "10.6.40.184:50011"
	//address_pozo = "localhost:50011"
	address_pozo2 = "10.6.40.184:50012"
	//address_pozo2 = "localhost:50012"
	address_name_node = "10.6.40.183:50020"
	//address_name_node = "localhost:50020"
)

type InfoJugadores struct {
	ID         int32
	equipo     int
	ID_rival   int32
	alive      int
	ult_jugada int32
	ronda_et1  int
	suma       int32
}

var Jugadores [16]InfoJugadores

func choose_number() {
	rand.Seed(time.Now().UTC().UnixNano())
	n_etapa1 = int32(rand.Intn(4) + 6)
	//return int32(rand.Intn(4) + 6)
}

type UserManagementServer struct {
	pb.UnimplementedLiderServicesServer //UnimplementedLiderServices está en el usermsg_grpc.pb, aquí se debe implementar
}

func PlayersAlive() (alive int) {
	for i := 0; i < 16; i++ {
		if Jugadores[i].alive == 1 {
			alive++
		}
	}
	return alive
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func mandar_pozo(id_jugador int, etapa int) {
	conn, err := amqp.Dial("amqp://guest:guest@10.6.40.184:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")

	body := "Jugador_" + strconv.Itoa(id_jugador) + " Ronda_" + strconv.Itoa(etapa)

	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	failOnError(err, "Failed to publish a message")
}

func (s *UserManagementServer) Play(ctx context.Context, in *pb.Message) (*pb.User, error) { //implementar del método Play
	user_id = user_id + 1
	log.Printf("Jugador %d acepta jugar", user_id)
	Jugador := InfoJugadores{
		ID:         user_id,
		equipo:     0,
		ID_rival:   0,
		alive:      1,
		ult_jugada: 0,
		ronda_et1:  0,
		suma:       0,
	}
	Jugadores[user_id-1] = Jugador
	return &pb.User{ID: user_id}, nil
}

func (s *UserManagementServer) Etapa1(ctx context.Context, in *pb.Jugada1) (*pb.Resp, error) { //implementacion del método Etapa1
	var bin int32 = 1
	if PlayersAlive() == 1 {
		fmt.Printf("El Jugador %d ha ganado el Juego del CALAMAR \n", in.GetID())
		return &pb.Resp{Survive: bin, Partida: int32(0), Juego: 0, Etapa: in.GetEtapa()}, nil
	}
	var jugada int32 = in.GetJugada()
	choose_number()
	partida := 1
	fmt.Printf("El Lider eligió %d y el jugador %d eligió %d \n", n_etapa1, in.GetID(), jugada)
	if jugada >= n_etapa1 {
		bin = 0
		Jugadores[in.GetID()-1].alive = 0
		fmt.Printf("El Jugador %d ha muerto \n", in.GetID())
		// mandar jugador al pozo
		//mandar_pozo(int(in.GetID()), int(in.GetEtapa()))

	}
	fmt.Printf("----------------------------------------------------------------------------- \n")
	ronda := Jugadores[in.GetID()-1].ronda_et1 + 1
	Jugadores[in.GetID()-1].ronda_et1 = ronda
	Jugadores[in.GetID()-1].suma += jugada
	conn, err := grpc.Dial(address_name_node, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Did not connect: %v", err)
	}
	defer conn.Close()
	ServiceClient := pb.NewNameNodeClient(conn)
	_, err = ServiceClient.JugadaPlayer(context.Background(), &pb.Jugada{ID: in.GetID(), Jugada: jugada, Ronda: int32(ronda), Etapa: in.GetEtapa()})

	if ronda == 4 {
		if Jugadores[in.GetID()-1].suma < 21 {
			bin = 0
			Jugadores[in.GetID()-1].alive = 0
			fmt.Printf("El Jugador %d ha muerto por no sumar 21 en 4 rondas \n", in.GetID())
			fmt.Printf("----------------------------------------------------------------------------- \n")
			//mandar jugador al pozo
		}
		partida = 0
	}
	return &pb.Resp{Survive: bin, Partida: int32(partida), Juego: 1, Ronda: int32(ronda), Etapa: in.GetEtapa()}, nil
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
