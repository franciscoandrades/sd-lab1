package main

import (
	"context"
	pb "distribuidos/go-usermsg-grpc/usermsg"
	"fmt"
	"log"

	"google.golang.org/grpc"
)

var etapas = []string{"ETAPA 1: LUZ ROJA, LUZ VERDE", "ETAPA 2: TIRAR LA CUERDA", "ETAPA 3: TODO O NADA"}
var reglas = []string{"Reglas del juego: \n Escoja un número entre el 1 y el 10, si es igual o mayor a que \neliga el Lider será eliminado, tiene 4 turnos para formar que sus números sumen 21.",
	"Reglas del juego: \n Escoja un número entre el 1 y el 4, su elección será sumada \ncon las elecciones de sus compañeros de equipo, si dicha suma tiene la misma paridad que la elección tiene la misma paridad que el lider pasan a la siguiente ronda.",
	"Reglas del juego: \n Escoja un número entre el 1 y el 10, si su valor es el más \ncercano al número elegido por el lider (se considerará el valor absoluto) gana, en caso contrario gana su rival."}

const (
	address = "localhost: 50000"
)

func main() {
	fmt.Println("-----------------------------------")
	fmt.Println("BIENVENIDO JUGADOR")
	fmt.Println("-----------------------------------")
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	n := 0 //corresponde a la etapa que jugará
	if err != nil {
		log.Fatalf("Did not connect: %v", err)
	}
	defer conn.Close()
	ServiceClient := pb.NewLiderServicesClient(conn)

	var resp string
	var numero int32 = 0
	var suma int32 = 0
	fmt.Printf("Desea ingresar al juego? [y/n]:  ")
	fmt.Scanf("%s \n", &resp)
	if resp == "y" {
		r, err := ServiceClient.Play(context.Background(), &pb.Message{Play: resp})
		survive := 1
		est_juego := 0
		for survive == 1 || est_juego != 0 {
			if err != nil {
				log.Printf("Error to recive respond")
			}
			var ID1 int32 = r.GetID()
			fmt.Printf("Su numero de jugador es %d \n", ID1)
			fmt.Println(etapas[n])
			fmt.Println("-----------------------------------")
			fmt.Println(reglas[n])
			ronda := 0
			partida := 1
			for partida == 1 && survive == 1 {
				fmt.Printf("Ingrese un número: ")
				fmt.Scanf("%d \n", &numero)
				if n == 0 {
					r1, err1 := ServiceClient.Etapa1(context.Background(), &pb.Jugada1{ID: ID1, Jugada: numero, Ronda: int32(ronda), Etapa: int32(n + 1), Suma: suma})
					survive = int(r1.GetSurvive())
					ronda = int(r1.GetRonda())
					suma = r1.GetSuma()
					if err1 != nil {
						log.Printf("Error to recive respond")
					}
				}
				if n == 1 {
					r1, err1 := ServiceClient.Etapa2(context.Background(), &pb.Jugada{ID: ID1, Jugada: numero})
					survive = int(r1.GetSurvive())
					if err1 != nil {
						log.Printf("Error to recive respond")
					}
				}
				if n == 3 {
					r1, err1 := ServiceClient.Etapa3(context.Background(), &pb.Jugada{ID: ID1, Jugada: numero})
					survive = int(r1.GetSurvive())
					if err1 != nil {
						log.Printf("Error to recive respond")
					}
				}

			}
			n++
			fmt.Println("ETAPA TERMINADA")
			fmt.Println("-----------------------------------")
			var elec int32 = 0
			for int(elec) != 2 {
				fmt.Println("1. Ver el monto acumulado en el pozo \n2. Avanzar a la siguiente etapa")
				fmt.Scanf("%d", &elec)
				fmt.Printf("Elec es %d \n", elec)
				if int(elec) == 1 {
					r2, err2 := ServiceClient.Pozo(context.Background(), &pb.Req{Req: "POZO"})
					if err2 != nil {
						log.Println("Error to recive respond")

					}
					fmt.Printf("El pozo actual es de: %d \n", r2.GetMonto())
				}
			}
			_, err2 := ServiceClient.Continue(context.Background(), &pb.Message{Play: "READY"})
			if err2 != nil {
				log.Println("Error to recive respond")

			}

		}
	}
}
