package main

import (
	"bufio"
	"fmt"
	"golang-book/assign1/api"
	"log"
	"net"
	"net/rpc/jsonrpc"
	"os"
	"strconv"
)

func main() {
	// Declaration
	var replies api.Reply
	var args *api.Args
	var budget float32
	var option int32
	var line string
	var tradeid int32

	conn, err := net.Dial("tcp", "127.0.0.1:1234")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	// Establish connection thru JSON rpc
	c := jsonrpc.NewClient(conn)

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
		// Send the data to api
		args = &api.Args{line, budget, 0}
		//Call the servcie
		err = c.Call("Api.Buyingstocks", args, &replies)
		if err != nil {
			log.Fatal("error ", err)
		}
		// Print the response
		fmt.Println("Response")
		fmt.Println("tradeID", replies.Trade)
		fmt.Println("stocks", replies.Stocks)
		fmt.Printf("unvested Amount: %.0f \n ", replies.Unvested)
	}
	if option == 2 {
		if len(commandargs) == 0 {
			fmt.Println("Enter the Trade Id no:")
			fmt.Scan(&tradeid)
		}
		fmt.Println("Checking the Portfolio")
		fmt.Println("Trade Id", tradeid)
		args = &api.Args{" ", 0, tradeid}
		//call the API
		err = c.Call("Api.Checkportfolio", args, &replies)
		if err != nil {
			log.Fatal("error ", err)
		}
		//Print the response
		fmt.Println("Response")
		fmt.Println("Stocks", replies.Stocks)
		fmt.Println("Current Market Value", replies.Current)
		fmt.Printf("unvested amount : %.0f \n", replies.Unvested)
	}

}
