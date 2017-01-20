package api

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
)

type Api string

// Structure for recieving requests and sending response
type Args struct {
	Stocksym string
	Budget   float32
	Trade    int32
}
type Reply struct {
	Trade    int32
	Stocks   string
	Unvested float32
	Current  float32
}
type Values struct {
	Stocks   string
	Unvested float32
}

// Define in memory storage
var hash map[int32]Values

func (t *Api) Buyingstocks(r *http.Request, args *Args, replies *Reply) error {

	var temperc int

	if args.Trade == 0 {
		//check the budget is entered
		if args.Budget == 0 {
			return errors.New("Budget not entered. Please enter again")
		}
		//error handling for percentage
		errinp := strings.Split(args.Stocksym, " ")
		errper := strings.Split(errinp[0], ",")
		for i := 0; i < len(errper); i++ {
			errpercen := strings.Split(errper[i], ":")
			errpercent := strings.Split(errpercen[1], "%")
			errint, _ := strconv.Atoi(errpercent[0])

			temperc = temperc + errint
			if temperc > 100 {
				return errors.New("Percentage entered is greater than 100. Please reenter")
			}
		}

		// Define in memory storage
		if hash == nil {
			hash = make(map[int32]Values)
		}

		// clean the incoming data
		input := strings.Split(args.Stocksym, " ")
		incom := strings.Split(input[0], ",")

		total := args.Budget
		for i := 0; i < len(incom); i++ {
			stocksym := strings.Split(incom[i], ":")
			con := "https://query.yahooapis.com/v1/public/yql?q=select%20*%20from%20yahoo.finance.quotes%20where%20symbol%20in%20(%22"
			con += stocksym[0]
			con += "%22)%0A%09%09&format=json&diagnostics=true&env=http%3A%2F%2Fdatatables.org%2Falltables.env&callback="
			// HTTP get to get the data from yahoo finance api
			resp, err := http.Get(con)
			if err == nil {
				body, err := ioutil.ReadAll(resp.Body)
				if err == nil {
					var data map[string]map[string]interface{}
					var buff interface{}
					json.Unmarshal(body, &data)
					buff = data["query"]["results"].(map[string]interface{})["quote"].(map[string]interface{})["Ask"]
					// calculations for buying stocks
					pers := strings.Split(stocksym[1], "%")
					per, _ := strconv.ParseFloat(pers[0], 32)
					percent := per / 100
					//convert ask to float
					ask := (buff).(string)
					amt, _ := strconv.ParseFloat(ask, 32)
					budgetalloc := args.Budget * float32(percent)
					var stocks string
					var st float64
					amtprec := float32(amt)
					if budgetalloc > amtprec {
						stockbought := budgetalloc / amtprec
						st = float64(stockbought)
						stocks = strconv.FormatFloat(st, 'f', 6, 64)
					}
					notused := total - float32((st * amt))
					total = notused
					replies.Unvested = notused
					if len(replies.Stocks) > 0 {
						replies.Stocks += ","
					}
					replies.Stocks += stocksym[0]
					replies.Stocks += ":"
					replies.Stocks += stocks
					replies.Stocks += ":$"
					replies.Stocks += ask
				}
			}
		}
		//generate random trade id
		tradeid := rand.Intn(9223372036854775807)
		replies.Trade = int32(tradeid)
		//in memory storage of the data in hash maps
		hash[replies.Trade] = Values{replies.Stocks, replies.Unvested}
	}
	return nil
}
func (t *Api) Checkportfolio(r *http.Request, args *Args, replies *Reply) error {
	// Initialize market value
	var marketvalue float32
	marketvalue = 0
	var stocksinmap string
	if args.Trade != 0 {
		//get the stocks purchased from the map
		for key, values := range hash {
			if key == args.Trade {
				stocksinmap = values.Stocks
				// get the stocks name
				name := strings.Split(stocksinmap, ",")
				for i := 0; i < len(name); i++ {
					// get the stocks symbol
					sym := strings.Split(name[i], ":")
					con := "https://query.yahooapis.com/v1/public/yql?q=select%20*%20from%20yahoo.finance.quotes%20where%20symbol%20in%20(%22"
					con += sym[0]
					con += "%22)%0A%09%09&format=json&diagnostics=true&env=http%3A%2F%2Fdatatables.org%2Falltables.env&callback="
					resp, err := http.Get(con)
					if err == nil {
						body, err := ioutil.ReadAll(resp.Body)
						if err == nil {
							var data map[string]map[string]interface{}
							var buff interface{}
							json.Unmarshal(body, &data)
							buff = data["query"]["results"].(map[string]interface{})["quote"].(map[string]interface{})["Ask"]

							cachestring := (buff).(string)
							current, _ := strconv.ParseFloat(cachestring, 32)
							//calculate raise or loss
							last := strings.Split(sym[2], "$")
							//convert last to float
							prev, _ := strconv.ParseFloat(last[1], 32)
							//convert no of stock to float
							number, _ := strconv.ParseFloat(sym[1], 32)
							marketvalue = marketvalue + float32(number*current)
							replies.Unvested = values.Unvested
							replies.Current = marketvalue
							if len(replies.Stocks) > 0 {
								replies.Stocks += ","
							}
							//return response
							replies.Stocks += sym[0]
							replies.Stocks += ":"
							replies.Stocks += sym[1]
							replies.Stocks += ":"
							if prev > current {
								replies.Stocks += "-"
							} else if prev < current {
								replies.Stocks += "+"
							}
							replies.Stocks += "$"
							replies.Stocks += cachestring
						}
					}
				}
			}
		}
	}
	return nil
}
