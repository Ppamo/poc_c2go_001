package utils

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"io/ioutil"
	"math/rand"
	"os"
	"time"
)

const (
	TIME_FORMAT = "2006-01-02 15:04:05.000000"
	configPath  = "../config/server.json"
)

var (
	timePrint      func(format string, a ...interface{}) = color.New(color.FgCyan).PrintfFunc()
	separatorPrint func(format string, a ...interface{}) = color.New(color.FgYellow).PrintfFunc()
	config         *ConfigStruct
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func GetConfig(debug bool) (*ConfigStruct, error) {
	var err error
	if config == nil {
		PrintDebug(debug, "Loading configuration")
		file, err := os.Open(configPath)
		if err != nil {
			PrintDebug(true, "Error opening config file from '%s'\n%v", configPath, err)
			return nil, err
		}
		defer file.Close()
		bytes, err := ioutil.ReadAll(file)
		if err != nil {
			PrintDebug(true, "Error reading bytes from config file\n%v", err)
			return nil, err
		}
		PrintDebug(debug, "Using config:\n%s", string(bytes))
		err = json.Unmarshal(bytes, &config)
	}
	return config, err
}

func Print(message string, args ...interface{}) {
	if args == nil {
		PrintDebug(true, message)
	} else {
		PrintDebug(true, message, args)
	}
}

func PrintlnDebug(debug bool, message string, args ...interface{}) {
	fmt.Println()
	if args == nil {
		PrintDebug(debug, message)
	} else {
		PrintDebug(debug, message, args)
	}
}

func PrintDebug(debug bool, message string, args ...interface{}) {
	if debug {
		now := time.Now()
		timePrint("%s", now.Format(TIME_FORMAT))
		separatorPrint(" > ")
		fmt.Printf(message, args...)
		fmt.Println()
	}
}

func RandStringRunes(n int) string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func EncodeString(debug bool, text string) (string, error) {
	encoded := base64.StdEncoding.EncodeToString([]byte(text))
	randomText := RandStringRunes(300)
	randomText = base64.StdEncoding.EncodeToString([]byte(randomText))
	pl := int(randomText[0:1][0])
	sl := int(randomText[len(randomText)-1:][0])
	encoded = randomText[:pl] + encoded
	encoded = encoded + randomText[len(randomText)-sl:]
	encoded = base64.StdEncoding.EncodeToString([]byte(encoded))
	return encoded, nil
}

func DecodeString(debug bool, encoded string) (string, error) {
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		PrintDebug(debug, "Error decoding from base64 first stage\n%v", err)
		return "", err
	}
	pl := int(decoded[0:1][0])
	sl := int(decoded[len(decoded)-1:][0])
	decoded = decoded[pl:]
	decoded = decoded[0 : len(decoded)-sl]
	decoded, err = base64.StdEncoding.DecodeString(string(decoded))
	if err != nil {
		PrintDebug(debug, "Error decoding from base64 second stage\n%v", err)
		return "", err
	}
	return string(decoded), nil
}

func DecodeConfig(debug bool, code string) error {
	PrintDebug(debug, "Deobfuscating configuration")
	data, err := DecodeString(debug, code)
	if err != nil {
		PrintDebug(debug, "Error decoding string\n%v", err)
		return err
	}
	PrintDebug(debug, "Decoded configuration:\n%s", data)
	err = json.Unmarshal([]byte(data), &config)
	if err != nil {
		PrintDebug(debug, "Error unmarshalling decoded data\n%v", err)
		return err
	}
	return nil
}
