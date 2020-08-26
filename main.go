package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

const POOLID = "paulainsley1betanet.stakehouse.betanet"
const NETWORK = "betanet"
const ACCOUNTID = "paulainsley1.betanet"
const UPPER_BID_THRESHOLD = 1.8
const LOWER_BID_THRESHOLD = 1.1
const SEAT_PRICE_PERCENTAGE = 1.3


type BetanetNear struct {
	Jsonrpc string `json:"jsonrpc"`
	Result ResultJson `json:"result"`
}

type ResultJson struct {
	Version string `json:"version"`
	Validators []ValidatorsJson `json:"validators"`
}

type ValidatorsJson struct {
	ACCOUNTID string `json:"account_id"`
	IsSlashed bool `json:"is_slashed"`
}

func proposalPingTest() {
	cmd := exec.Command("near", POOLID, "ping", "--ACCOUNTID", ACCOUNTID)
	if _, err := cmd.Output(); err != nil {
		log.Fatalln("Unable to ping the Near contracts")
	} else {
		log.Println("Near proposals Ping successful")
	}
}

func verifyValidator(url string) {
	values := map[string]string{"jsonrpc": "2.0", "method": "status", "id": POOLID}
	jsonValue, _ := json.Marshal(values)
	response, err := http.Post(url, "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		panic(err.Error())
	}
	betanet_near := new(BetanetNear)
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err.Error())
	}
	_ = json.Unmarshal([]byte(body), &betanet_near)
	for _, validator := range betanet_near.Result.Validators {
		if validator.ACCOUNTID == POOLID {
			log.Printf("Stake Pool ID exists %s \n", POOLID)
		} else {
			log.Fatalf("No such Stake POOLID exists \n")
		}
	}
}

func getNextSlotPrice() float64 {
	cmd := exec.Command("near", "proposals")
	if output, err := cmd.Output(); err != nil {
		log.Fatalln("Failed to get Next SLOT Error: ", err.Error())
		return float64(0)
	} else {
		next_slot_price_string := strings.Split(strings.Split(string(output), "seat price = ")[1], ")")[0]
		integer_value, _ := strconv.Atoi(strings.ReplaceAll(next_slot_price_string, ",", ""))
		integer_value = integer_value * 10
		next_slot_price := math.Pow(float64(integer_value), 24)
		log.Printf("The next slotprice will most likely be %f yocto nears\n", next_slot_price)
		return next_slot_price
	}
}

func getCurrentBid() float64 {
	cmd := exec.Command("near", "proposals")
	if output, err := cmd.Output(); err != nil {
		log.Fatalln("Failed to get Next SLOT Error: ", err.Error())
		return float64(0)
	} else {
		base_string := strings.Split(strings.Split(string(output), POOLID)[1], "|")[1]
		sec_base_string := strings.Split(base_string, "=>")[0]
		replace_string := strings.ReplaceAll(strings.ReplaceAll(sec_base_string, ",", ""), " ", "")
		integer_value, _ := strconv.Atoi(replace_string)
		integer_value = integer_value * 10
		locked_amount := math.Pow(float64(integer_value), 24)
		log.Printf("%s has currently %f bid in auction \n", POOLID, locked_amount)
		return locked_amount
	}
}

func reduceStake(stake float64, next_slot_price float64){
	amount_to_unstake := int64(stake - (next_slot_price * SEAT_PRICE_PERCENTAGE))
	cmd := exec.Command("near", "call", POOLID, "unstake", fmt.Sprint("'{{\"amount\": \"%f\"}}'", amount_to_unstake), "--amountId", ACCOUNTID)
	if err := cmd.Run(); err != nil {
		log.Fatalln("Unstaking the currently staked Funds failed! Error: ", err.Error())
	}
}
func increaseStake(stake float64, next_slot_price float64){
	amount_to_unstake := int64(stake - (next_slot_price * SEAT_PRICE_PERCENTAGE))
	cmd := exec.Command("near", "call", POOLID, "stake", fmt.Sprint("'{{\"amount\": \"%f\"}}'", amount_to_unstake), "--amountId", ACCOUNTID)
	if err := cmd.Run(); err != nil {
		log.Fatalln("Staking the Funds failed! Error: ", err.Error())
	}
}


func adaptStake(bid_price float64, next_slot_price float64) {
	if bid_price > (next_slot_price * UPPER_BID_THRESHOLD) {
		reduceStake(bid_price, next_slot_price)
	} else if bid_price < (next_slot_price * LOWER_BID_THRESHOLD) {
		increaseStake(bid_price, next_slot_price)
	} else {
		log.Println("Current bid is fine")
	}
}

func run() {
	proposalPingTest()
	verifyValidator(fmt.Sprintf("https://rpc.%s.near.org", NETWORK))
	next_slot_price := getNextSlotPrice()
	my_bid := getCurrentBid()
	adaptStake(my_bid, next_slot_price)
}

func waitUntilNextEpoch()  {
	log.Println("Waiting for 10 minutes")
	time.Sleep(10 * 60)
}

func main()  {
	for true {
		run()
		waitUntilNextEpoch()
	}
}