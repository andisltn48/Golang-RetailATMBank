package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	// "fmt"
	"net/http"

	"github.com/urfave/cli"
)

type CurrAccount struct {
	Name string
}

func main() {

	app := cli.NewApp()
	app.Name = "ATM - Retail Bank"

	app.Commands = []cli.Command{
		{
			Name:     "login",
			HelpName: "login",
			Action:   loginUser,
		},
		{
			Name:     "deposit",
			HelpName: "deposit",
			Action:   deposit,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}

func loginUser(c *cli.Context) {
	name := c.Args()[0]

	var currentAccount CurrAccount
	currentAccount.Name = name
	postBody, _ := json.Marshal(map[string]string{
		"name": name,
	})
	responseBody := bytes.NewBuffer(postBody)

	resp, err := http.Post("http://127.0.0.1:8081/login", "application/json", responseBody)
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	type Response struct {
		Name    string `json:"name"`
		Balance int    `json:"balance"`
	}
	var responseObject Response
	json.Unmarshal(body, &responseObject)

	if responseObject.Name == "" {
		registerUser(name)
		responseObject.Name = name
		responseObject.Balance = 0
	}

	log.Println("Hello", responseObject.Name, "!")
	log.Println("CurrName", currentAccount.Name, "!")
	log.Println("Your balance is $", responseObject.Balance)
}

func registerUser(name string) {
	postBody, _ := json.Marshal(map[string]string{
		"name":    name,
		"balance": "0",
	})
	responseBody := bytes.NewBuffer(postBody)

	_, err := http.Post("http://127.0.0.1:8081/register", "application/json", responseBody)
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}
}

func deposit(c *cli.Context) {
	name := c.Args()[0]
	var currentAccount CurrAccount
	log.Println("Your balance is $", currentAccount.Name)
	balance := c.Args()[1]
	postBody, _ := json.Marshal(map[string]string{
		"name":    name,
		"balance": balance,
	})
	responseBody := bytes.NewBuffer(postBody)

	resp, err := http.Post("http://127.0.0.1:8081/add-balance", "application/json", responseBody)
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	type Response struct {
		Balance int `json:"balance"`
	}

	var responseObject Response
	json.Unmarshal(body, &responseObject)

	log.Println("Your balance is $", responseObject.Balance)
}
