package main

import (
	"context"
	pb "distribuidos/go-usermsg-grpc/usermsg"
	"io/ioutil"
	"log"
	"net"
	"strconv"
	"strings"

	"google.golang.org/grpc"
)

const (
	port = ":50023"
)

type UserManagementServer struct {
	pb.UnimplementedDataNodeServer
}

func suma(id int, etapa int) (suma int) {
	suma = 0
	Bytes, err := ioutil.ReadFile("DataNode/Jugador_" + strconv.Itoa(id) + " Etapa_" + strconv.Itoa(etapa) + ".txt")
	if err != nil {
		log.Fatal(err)
	}
	datos := string(Bytes)
	arr := strings.Fields(datos)
	for i := 0; i < len(arr); i++ {
		num, _ := strconv.Atoi(arr[i])
		suma = suma + num
	}
	return suma
}

func (s *UserManagementServer) RegistrarInfo(ctx context.Context, in *pb.Jugada) (*pb.Check, error) {
	b := []byte(strconv.Itoa(int(in.GetJugada())))
	err := ioutil.WriteFile("DataNode/Jugador_"+strconv.Itoa(int(in.GetID()))+"__Etapa_"+strconv.Itoa(int(in.GetEtapa()))+".txt", b, 0644)
	if err != nil {
		log.Fatalf("Failed to write in Registro.txt")
	}
	return &pb.Check{Check: "OK"}, nil
}

func (s *UserManagementServer) PlayerInfo(ctx context.Context, in *pb.User) (*pb.Data, error) {
	archivos, err := ioutil.ReadDir("./DataNode")
	if err != nil {
		log.Fatal(err)
	}
	var arr []string
	var etapa string
	for _, archivo := range archivos {
		if strings.Contains(archivo.Name(), "Jugador_"+strconv.Itoa(int(in.GetID()))) {
			array := strings.Split(archivo.Name(), "__")
			etapa = strings.ReplaceAll(array[1], "Etapa_", "")
			Bytes, _ := ioutil.ReadFile(archivo.Name())
			datos := string(Bytes)
			arr = strings.Fields(datos)
		}
	}
	etapa_int, _ := strconv.Atoi(etapa)
	return &pb.Data{Etapa: int32(etapa_int), Ronda: int32(len(arr)), Jugada: 0}, nil
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
		log.Fatalf("Failed to listen on port 50023: %v", err)
	}
}
