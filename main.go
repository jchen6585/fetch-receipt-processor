package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/google/uuid"

	"github.com/jchen6585/fetch-receipt-processor/utils"
)

var (
	db    = make(map[string]int)
	mutex sync.Mutex
)

func main() {
	http.HandleFunc("/receipts/process", postHandler)
	http.HandleFunc("/receipts/{id}/points", getHandler)

	log.Print("Server is starting up at localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Server did not start up succesfully")
	}
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		log.Print("Only POST requests are processed at the /receipts/process endpoints")
		return
	}

	mutex.Lock()
	defer mutex.Unlock()

	var receipt utils.Receipt

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&receipt)
	if err != nil {
		log.Print(err)
		return
	}

	if !utils.ValidateReceipt(receipt) {
		log.Print("Receipt is invalid")
		return
	}

	points := utils.CalculatePoints(receipt)

	id := uuid.New()
	db[id.String()] = points

	json.NewEncoder(w).Encode(id.String())

}

func getHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		log.Print("Only GET requests are processed at the /receipts/{id}/points endpoints")
		return
	}

	mutex.Lock()
	defer mutex.Unlock()

	id := r.PathValue("id")

	val, ok := db[id]
	if !ok {
		log.Printf("Receipt with the id %s does not exist", id)
		return
	}

	json.NewEncoder(w).Encode(val)
}
