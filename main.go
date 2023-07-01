package main

import (
	"encoding/json"
	"fmt"
	"github.com/NethermindEth/neth-monitor-back/models"
	"log"
	"net/http"
	"os"
	"path"
	"runtime"
)

func main() {
	username := os.Getenv("MONGO_USER")
	password := os.Getenv("MONGO_PASSWORD")

	models.InitMongo(fmt.Sprintf("mongodb://%s:%s@localhost:27017", username, password))
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "static")

	fs := http.FileServer(http.Dir(dir))
	http.Handle("/", fs)

	http.HandleFunc("/api/subscribeNode", subscribeNodeHandler)
	http.HandleFunc("/api/updateNode", saveDataHandler)

	http.HandleFunc("/api/getData", getDataHandler)

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func subscribeNodeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("calling /api/subscribeNode")
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var data models.Node
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Invalid data", http.StatusBadRequest)
		return
	}
	if data.Data == nil {
		data.Data = make([]models.EthereumNodeData, 0)
	}
	err = models.CreateNode(&data)
	if err != nil {
		http.Error(w, "Unable to store node data", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

type request struct {
	Enode string                  `json:"enode"`
	Data  models.EthereumNodeData `json:"data"`
}

func saveDataHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("calling /api/updateNode")
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var data request
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Invalid data", http.StatusBadRequest)
		return
	}

	err = models.AddEntry(data.Enode, data.Data)
	if err != nil {
		http.Error(w, "Unable to store node data", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func getDataHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("calling /api/getData")
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Retrieve data
	data := models.GetAllData()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
