package client

import (
	"io"
	"fmt"
	"cli-chat/pkg/utils"
	"cli-chat/pkg/config"
	"cli-chat/pkg/chatmp"
)

func manageIn(client *chatmp.Client) {
	for {
		f, e := client.RecvFrame()
		if e == io.EOF {
			break
		} else if e != nil {
			utils.FatalError("chatmp.Client.RecvFrame", e)
		}
		fmt.Println(string(f.Body))
	}
}

func manageOut(client *chatmp.Client, outChan chan string) {
	for {
		userMsg := <-outChan
		fmt.Println("Sending:", userMsg)
		e := client.SendText(userMsg)
		utils.CheckError(e)
	}
}

func manageStdin(outChan chan string) {
	for utils.Stdin.Scan() {
		outChan <-utils.Stdin.Text()
	}
	if e := utils.Stdin.Err(); e != io.EOF {
		utils.CheckError(e)
	}
}

func Run() {
	config, e := config.LoadOrPromptConfig()
	if e == io.EOF {
		return
	} else if e != nil {
		utils.FatalError("config.LoadOrPromptConfig", e)
	}
	fmt.Println(config)

	outChan := make(chan string)
	client, e := chatmp.NewClient(config.Username, "127.0.0.1:8000")
	utils.CheckError(e)
	defer client.Close()

	if e = client.ClaimUsername(); e != nil {
		utils.FatalError("chatmp.Client.ClaimUsername", e)
	} else {
		fmt.Println("Using username:", client.Username)
	}

	go manageIn(client)
	go manageOut(client, outChan)
	manageStdin(outChan)
}
