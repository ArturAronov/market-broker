package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func compress(input string) ([]byte, error) {
	u, err := uuid.Parse(input)
	if err != nil {
		return nil, err
	}
	return u[:], nil
}

func main() {
	// body := []byte(`{"m": "h"}`)
	// req, err := http.NewRequest("GET", "http://localhost:8080", bytes.NewReader(body))
	tranactionType := 1                                     //1 byte
	transactionMethod := 1                                  //1 byte
	orderType := 1                                          //1 byte
	ticker := "XYZQ"                                        //4 bytes
	quantity := 4294967295                                  //?? bytes
	price := 770                                            // ?? bytes
	var orderDate uint64 = uint64(time.Now().Unix())        //?? bytes
	var goodUntil uint64 = 18446744073709551615             //?? bytes
	traderId := "46090691-dc41-4280-960f-3428806c3c5e"      //36 bytes
	clientOrderId := "2527ee99-4f5c-4adf-8a10-dc5c7b7ee287" //36 bytes

	tickerByte := []byte(ticker)
	quantityByte := byte(quantity)                  //?? bytes
	priceByte := byte(price)                        //?? bytes
	orderDateByte := byte(orderDate)                //?? bytes
	goodUntilByte := byte(goodUntil)                //?? bytes
	traderIdByte, _ := compress(traderId)           //?? bytes
	clientOrderIdByte, _ := compress(clientOrderId) //?? bytes

	fmt.Printf("tranactionType: %v\n", tranactionType)
	fmt.Printf("transactionMethod: %v\n", transactionMethod)
	fmt.Printf("orderType: %v\n", orderType)
	fmt.Printf("ticker: %v\n", tickerByte)
	fmt.Printf("quantityByte: %v\n", quantityByte)
	fmt.Printf("priceByte: %v\n", priceByte)
	fmt.Printf("orderDateByte: %v %v\n", orderDate, orderDateByte)
	fmt.Printf("goodUntilByte: %v\n", goodUntilByte)
	fmt.Printf("traderIdByte: %v\n", traderIdByte)
	fmt.Printf("clientOrderIdByte: %v %d\n", clientOrderIdByte, len(clientOrderIdByte))

	req, err := http.NewRequest("GET", "http://localhost:8080", nil)
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}

	// req.Header.Set("User-Agent", "_")
	// req.Header.Set("Accept", "*/*")
	// req.Header.Set("Content-Type", "application/json")
	// Send the request

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("Failed to send request: %v", err)
	}

	defer resp.Body.Close()

	log.Printf("Response: %s, protocol: %s", resp.Status, resp.Proto)

}
