package main

import (
	"encoding/base64"
	"encoding/binary"
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

func uint64ToBytes(input uint64) []byte {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, input)

	return buf
}

func uint32ToBytes(input uint32) []byte {
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, input)

	return buf
}

func main() {
	timeNow := time.Now()

	tranactionType := 1                                            // 1 byte
	transactionMethod := 1                                         // 1 byte
	orderType := 1                                                 // 1 byte
	ticker := "XYZQ"                                               // 4 bytes
	quantity := 4294967295                                         // ?? bytes
	price := 770                                                   // ?? bytes
	var orderDate uint64 = uint64(timeNow.AddDate(1, 0, 0).Unix()) // ?? bytes
	var goodUntil uint64 = 18446744073709551615                    // ?? bytes
	traderId := "46090691-dc41-4280-960f-3428806c3c5e"             // 36 bytes
	clientOrderId := "2527ee99-4f5c-4adf-8a10-dc5c7b7ee287"        // 36 bytes

	tickerByte := []byte(ticker)                    // 4 bytes
	quantityByte := make([]byte, 4)                 // 4 bytes
	priceByte := make([]byte, 4)                    // 4 bytes
	orderDateByte := uint64ToBytes(orderDate)       // 8 bytes
	goodUntilByte := uint64ToBytes(goodUntil)       // 8 bytes
	traderIdByte, _ := compress(traderId)           // 16 bytes
	clientOrderIdByte, _ := compress(clientOrderId) // 16 bytes

	binary.BigEndian.PutUint32(quantityByte, uint32(quantity))
	binary.BigEndian.PutUint32(priceByte, uint32(price))

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

	orderPayload := []byte{}

	orderPayload = append(orderPayload, byte(tranactionType))
	orderPayload = append(orderPayload, byte(transactionMethod))
	orderPayload = append(orderPayload, byte(orderType))
	orderPayload = append(orderPayload, tickerByte...)
	orderPayload = append(orderPayload, quantityByte...)
	orderPayload = append(orderPayload, priceByte...)
	orderPayload = append(orderPayload, orderDateByte...)
	orderPayload = append(orderPayload, goodUntilByte...)
	orderPayload = append(orderPayload, traderIdByte...)
	orderPayload = append(orderPayload, clientOrderIdByte...)
	orderPayloadUrlEncoded := base64.URLEncoding.EncodeToString(orderPayload)
	// decodedBytes, _ := base64.URLEncoding.DecodeString(orderPayloadUrlEncoded)

	fmt.Println(orderPayload, len(orderPayload))
	// fmt.Println(decodedBytes)
	fmt.Println(orderPayloadUrlEncoded)

	req, err := http.NewRequest("GET", "http://localhost:8080/"+orderPayloadUrlEncoded, nil)
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}

	req.Header.Set("User-Agent", "market-broker")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("Failed to send request: %v", err)
	}

	defer resp.Body.Close()

	log.Printf("Response: %s, protocol: %s", resp.Status, resp.Proto)

}
