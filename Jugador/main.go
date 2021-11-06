package main

import (
	"context"
	pb "distribuidos/go-usermsg-grpc/usermsg"
	"fmt"
	"log"
	"math/rand"
	"time"

	"google.golang.org/grpc"
)

var etapas = []string{"ETAPA 1: LUZ ROJA, LUZ VERDE", "ETAPA 2: TIRAR LA CUERDA", "ETAPA 3: TODO O NADA"}
var reglas = []string{"Reglas del juego: \n Escoja un número entre el 1 y el 10, si es igual o mayor a que \neliga el Lider será eliminado, tiene 4 turnos para formar que sus números sumen 21.",
	"Reglas del juego: \n Escoja un número entre el 1 y el 4, su elección será sumada \ncon las elecciones de sus compañeros de equipo, si dicha suma tiene la misma paridad que la elección tiene la misma paridad que el lider pasan a la siguiente ronda.",
	"Reglas del juego: \n Escoja un número entre el 1 y el 10, si su valor es el más \ncercano al número elegido por el lider (se considerará el valor absoluto) gana, en caso contrario gana su rival."}

const (
	address = "10.6.40.181:50000"
	//address = "localhost:50000"
)

type Bot struct {
	ID      int32
	jugada  int32
	survive int
}

const cant_bots = 16

var Jugadores [cant_bots]Bot

func interfaz() {
	fmt.Println("-----------------------------------")
	fmt.Println("BIENVENIDO JUGADOR")
	fmt.Println("-----------------------------------")

	fmt.Printf("Desea ingresar al juego? [y/n]:  ")
}

func luzverdeluzroja() {
	fmt.Println("ETAPA 1: LUZ ROJA, LUZ VERDE")
	fmt.Println("-----------------------------------")
	fmt.Println("Reglas del juego: \n Escoja un número del 1 al 10, si es igual o mayor al que elija el Lider será eliminado, tiene 4 turnos para formar que sus números sumen 21")
}

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
	var jugada int32
	var resp string
	fmt.Printf("Desea ingresar al juego? [y/n]:  ")
	fmt.Scanf("%s \n", &resp)
	if resp == "y" {

		for i := 0; i < cant_bots; i++ {
			r, err := ServiceClient.Play(context.Background(), &pb.Message{Play: "y"})
			if err != nil {
				log.Printf("Error to recive respond")
			}
			Bot := Bot{
				ID:      r.GetID(),
				jugada:  int32(0),
				survive: 1,
			}

			Jugadores[i] = Bot

		}
		for n < 3 {
			fmt.Println(etapas[n])
			fmt.Println(reglas[n])
			var partida = 1
			var juego = 1
			for partida == 1 {
				for i := 0; i < 16; i++ {
					if Jugadores[i].survive == 1 {
						var ID = Jugadores[i].ID
						if ID == 16 {
							fmt.Printf("Ingrese elección: ")
							fmt.Scanf("%d \n", &jugada)
						} else {
							rand.Seed(time.Now().UTC().UnixNano())
							jugada = int32(rand.Intn(7) + 1)
						}
						if n == 0 {
							//fmt.Printf("El jugador %d jugó %d \n", ID, jugada)
							r, err := ServiceClient.Etapa1(context.Background(), &pb.Jugada1{ID: ID, Jugada: jugada, Etapa: int32(n + 1)})
							if err != nil {
								log.Fatal("Error to recive response")
							}
							Jugadores[i].survive = int(r.GetSurvive())
							partida = int(r.GetPartida())
							juego = int(r.GetJuego())
							if juego == 0 {
								return
							}
						}
						if n == 1 {
							//fmt.Printf("El jugador %d jugó %d \n", ID, jugada)
							r, err := ServiceClient.Etapa2(context.Background(), &pb.Jugada{ID: ID, Jugada: jugada, Etapa: int32(n + 1)})
							if err != nil {
								log.Fatal("Error to recive response")
							}
							Jugadores[i].survive = int(r.GetSurvive())
							partida = int(r.GetPartida())
							juego = int(r.GetJuego())
							if juego == 0 {
								return
							}
						}
						if n == 2 {
							//fmt.Printf("El jugador %d jugó %d \n", ID, jugada)
							r, err := ServiceClient.Etapa3(context.Background(), &pb.Jugada{ID: ID, Jugada: jugada, Etapa: int32(n + 1)})
							if err != nil {
								log.Fatal("Error to recive response")
							}
							Jugadores[i].survive = int(r.GetSurvive())
							partida = int(r.GetPartida())
							juego = int(r.GetJuego())
							if juego == 0 {
								return
							}
						}
					}

				}
				var option = 1
				for option == 1 {
					fmt.Println("1.Ver monto del pozo")
					fmt.Println("2.Avanzar a la siguiente ronda")
					fmt.Scanf("%d", &option)
					if option == 1 {
						r, err := ServiceClient.Pozo(context.Background(), &pb.Req{Req: "Hi"})
						if err != nil {
							log.Fatal("Error to recive response Pozo")
						}
						fmt.Printf("El monto del pozo es %d \n", r.GetMonto())
					}
				}
			}
			n++
		}
	}
}
