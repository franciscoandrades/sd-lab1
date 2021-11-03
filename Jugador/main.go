package main

import (
	"context"
	pb "distribuidos/go-usermsg-grpc/usermsg"
	"fmt"
	"log"

	"google.golang.org/grpc"
)

const (
	address = "localhost: 50000"
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Did not connect: %v", err)
	}
	defer conn.Close()
	ServiceClient := pb.NewLiderServicesClient(conn)

	fmt.Println("-----------------------------------")
	fmt.Println("BIENVENIDO JUGADOR")
	fmt.Println("-----------------------------------")
	var resp string
	var name string
	var numero int32 = 0
	fmt.Printf("Desea ingresar al juego? [y/n]:  ")
	fmt.Scanf("%s \n", &resp)
	if resp == "y" {

		fmt.Printf("Ingrese su nombre: ")
		fmt.Scanf("%s \n", &name)
		r, err := ServiceClient.NewPlayer(context.Background(), &pb.Message{Name: name})
		if err != nil {
			log.Printf("Error to recive respond")
		}
		var ID1 int32 = r.GetID()
		fmt.Printf("Su numero de jugador es %d \n", ID1)
		fmt.Println("ETAPA 1: LUZ ROJA, LUZ VERDE")
		fmt.Println("-----------------------------------")
		fmt.Println("Reglas del juego: \n Escoja un número del 1 al 10, si es igual o mayor a que eliga el Lider será eliminado, tiene 4 turnos para formar que sus números sumen 21")

		ronda := 0
		for ronda < 4 {
			var suma int32 = 0
			fmt.Printf("Ingrese un número: ")
			fmt.Scanf("%d \n", &numero)
			suma = suma + numero
			r1, err1 := ServiceClient.Luz_Roja_Verde(context.Background(), &pb.Jugada_1{ID: ID1, NElegido: numero})
			ronda = int(r1.GetRonda())
			if err1 != nil {
				log.Printf("Error to recive respond")
			}
			if r1.GetBinario() == 1 {
				fmt.Println("USTED HA SIDO ELIMINADO")
				return
			}
			fmt.Printf("SOBREVIVIÓ A LA RONDA %d \n", r1.GetRonda())

		}
		fmt.Println("LUZ ROJA, LUZ VERDE TERMINADO")
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

		fmt.Println("ETAPA 2: TIRAR LA CUERDA")
		fmt.Println("-----------------------------------")

	}
}
