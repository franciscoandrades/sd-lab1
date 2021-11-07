# Laboratorio 2 - Sistemas Distribuidos
# Squid Game

----- Integrantes ----
* Nicolás Puente 201873618-K
* Francisco Andrades 201673584-4
* Lucas Díaz Aravena 201673524-0
----------------------

# Información
1) Los procesos están distribuidos según:
	
a) Máquina Virtual 10.6.40.181
- Proceso: Líder y DataNode2
- Puertos: Líder escucha en puerto 50000 y DataNode en puerto 50023

b) Máquina Virtual 10.6.40.182
- Proceso: Jugadores y DataNode1
- Puertos: El DataNode1 escucha en el puerto 50023

c) Máquina Virtual 10.6.40.183
- Proceso: NameNode
- Puertos: 50020

d) Máquina Virtual 10.6.40.184
- Proceso: Pozo y DataNode3
- Puertos: El pozo escucha en el puerto 50011 y DataNode3 en puerto 50023

# Instrucciones ejecución

* dist41
Abrir 2 consolas .

(i)
cd sd-lab1
make liderx
(ii)
cd sd-lab1
make datanodex

* dist43
Abrir 1 consola.

(i)
cd sd-lab1
make namenodex

* dist44
Abrir 2 consolas.

(i)
cd sd-lab1
make datanodex

(ii)
cd sd-lab1
make pozox

* dist42
Abrir 2 consolas.
Importante: Notar que jugador debe ejecutarse al final de todo.

(i)
cd sd-lab1
make datanodex

(ii)
cd sd-lab1
make jugadorx



* En la dist44, el server de RabbitMQ debiese estar corriendo. Si no es así, por favor ejecutar el siguiente comando en otra consola de dicha máquina virtual:

rabbitmq-server