package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func main() {
	// Declaration
	var budget float32
	var option int32
	var line string
	var tradeid int32

	//check the options
	commandargs := os.Args[1:]
	if len(commandargs) == 0 {
		fmt.Println("Choose :")
		fmt.Println("1. Buying stocks")
		fmt.Println("2. Checking Portfolio")
		fmt.Scan(&option)
	} else if len(commandargs) > 1 {
		option = 1
		line = commandargs[0]
		cache, _ := strconv.ParseFloat(commandargs[1], 32)
		budget = float32(cache)
	} else {
		option = 2
		number, _ := strconv.Atoi(commandargs[0])
		tradeid = int32(number)
	}
	if option == 1 {
		if len(commandargs) == 0 {
			fmt.Println("Enter the Request:")
			scanner := bufio.NewReader(os.Stdin)
			fmt.Println("eg : GOOG:50%,YHOO:35%..")
			line, _ = scanner.ReadString('\n')
			fmt.Print("Budget")
			fmt.Scan(&budget)
		}
		fmt.Println("Buying Stocks")
		// Marshal the data to be sent to the service
		var par [2]map[string]interface{}
		a := map[string]interface{}{"Stocksym": line, "Budget": budget}
		par[0] = a
		data, err := json.Marshal(map[string]interface{}{
			"method": "Api.Buyingstocks",
			"id":     "0",
			"params": par,
		})

		if err != nil {
			log.Fatalf("Marshal: %v", err)
		}
		//data sent via http Post
		resp, err := http.Post("http://127.0.0.1:8080/rpc",
			"application/json", strings.NewReader(string(data)))
		if err != nil {
			log.Fatalf("Post: %v", err)
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		log.Println("Body: ", string(body))
		if err != nil {
			log.Fatalf("ReadAll: %v", err)
		}
		result := make(map[string]interface{})
		// Unmarshal the data from the interface
		err = json.Unmarshal(body, &result)
		if err != nil {
			log.Fatalf("Unmarshal: %v", err)
		}
		trade := result["result"].(map[string]interface{})["Trade"]
		stocks := result["result"].(map[string]interface{})["Stocks"]
		unvested := result["result"].(map[string]interface{})["Unvested"]
		// Print the response
		fmt.Println("Response")
		fmt.Printf("tradeID %.0f \n", trade)
		fmt.Println("stocks", stocks)
		fmt.Printf("unvested Amount: %.0f \n ", unvested)
	}
	if option == 2 {
		if len(commandargs) == 0 {
			fmt.Println("Enter the Trade Id no:")
			fmt.Scan(&tradeid)
		}
		fmt.Println("Checking the Portfolio")
		fmt.Println("Trade Id", tradeid)
		// marshal the data for checking portfolio
		var par [2]map[string]interface{}
		a := map[string]interface{}{"Trade": tradeid}
		par[0] = a
		data, err := json.Marshal(map[string]interface{}{
			"method": "Api.Checkportfolio",
			"id":     "0",
			"params": par,
		})

		if err != nil {
			log.Fatalf("Marshal: %v", err)
		}
		resp, err := http.Post("http://127.0.0.1:8080/rpc",
			"application/json", strings.NewReader(string(data)))
		if err != nil {
			log.Fatalf("Post: %v", err)
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		log.Println("Body: ", string(body))
		if err != nil {
			log.Fatalf("ReadAll: %v", err)
		}
		result := make(map[string]interface{})
		err = json.Unmarshal(body, &result)
		if err != nil {
			log.Fatalf("Unmarshal: %v", err)
		}
		stocks := result["result"].(map[string]interface{})["Stocks"]
		unvested := result["result"].(map[string]interface{})["Unvested"]
		current := result["result"].(map[string]interface{})["Current"]
		//Print the response
		fmt.Println("Response")
		fmt.Println("Stocks", stocks)
		fmt.Printf("Current Market Value %.1f \n", current)
		fmt.Printf("unvested amount : %.0f \n", unvested)
	}
}
