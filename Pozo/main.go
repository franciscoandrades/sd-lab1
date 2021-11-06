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
	port = ":50011"
)

var pozo = 30000

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

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
	//fmt.Printf("%q", pozo)
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

	/*go func() {
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
	}()*/

	/*conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
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

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	go_to := make(chan bool)
	go func() {
		for d := range msgs {
			var pozo = int(leer_pozo())
			var body = string(d.Body) + " " + strconv.Itoa(pozo)
			log.Printf("Se recibe... %s", body)

		}

	}()
	log.Printf(" [*] Pozo recibiendo mensajes en puerto 5672")
	<-go_to*/
}
