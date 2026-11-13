package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
)

type Barang struct {
	ID     int    `json:"id"`
	Nama   string `json:"nama"`
	Lokasi string `json:"lokasi"`
	Status string `json:"status"`
}

type ResponseList struct {
	Jumlah int      `json:"jumlah"`
	Data   []Barang `json:"data"`
}

type DataStore struct {
	mu     sync.Mutex
	data   map[int]Barang
	nextID int
}

var store = DataStore{
	data:   make(map[int]Barang),
	nextID: 1,
}

func sendJSON(w http.ResponseWriter, status int, val any) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(val)
}

func sendError(w http.ResponseWriter, status int, message string) {
	sendJSON(w, status, map[string]string{"error": message})
}

func handleHealthz(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		sendError(w, http.StatusMethodNotAllowed, "Method not Allowed")
		return
	}
	sendJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func handleBarang(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var input struct {
			Nama   string `json:"nama"`
			Lokasi string `json:"lokasi"`
		}
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			sendError(w, http.StatusBadRequest, "Invalid Payload or Damaged JSON")
			return
		}

		nama := strings.TrimSpace(input.Nama)
		lokasi := strings.TrimSpace(input.Lokasi)

		if nama == "" || lokasi == "" {
			sendError(w, http.StatusBadRequest, "Name and Location cannot be empty")
			return
		}

		store.mu.Lock()
		defer store.mu.Unlock()

		for _, b := range store.data {
			if b.Status == "tersimpan" && strings.EqualFold(b.Nama, nama) && strings.EqualFold(b.Lokasi, lokasi) {
				sendError(w, http.StatusConflict, "Same data still save")
				return
			}
		}

		b := Barang{
			ID:     store.nextID,
			Nama:   nama,
			Lokasi: lokasi,
			Status: "tersimpan",
		}
		store.data[store.nextID] = b
		store.nextID++

		sendJSON(w, http.StatusCreated, b)

	case http.MethodGet:
		statusQuery := r.URL.Query().Get("status")
		if statusQuery != "" && statusQuery != "tersimpan" && statusQuery != "diklaim" {
			sendError(w, http.StatusBadRequest, "Status unknown")
			return
		}

		store.mu.Lock()
		defer store.mu.Unlock()

		var list []Barang
		for i := 1; i < store.nextID; i++ {
			if b, ok := store.data[i]; ok {
				if statusQuery == "" || b.Status == statusQuery {
					list = append(list, b)
				}
			}
		}

		if list == nil {
			list = []Barang{}
		}

		sendJSON(w, http.StatusOK, ResponseList{
			Jumlah: len(list),
			Data:   list,
		})

	default:
		sendError(w, http.StatusMethodNotAllowed, "Not allowed Method")
	}
}

func handleBarangByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		sendError(w, http.StatusBadRequest, "Id must contain number")
		return
	}

	store.mu.Lock()
	defer store.mu.Unlock()

	b, exists := store.data[id]
	if !exists {
		sendError(w, http.StatusNotFound, "Data Not Found")
		return
	}

	switch r.Method {
	case http.MethodGet:
		sendJSON(w, http.StatusOK, b)

	case http.MethodDelete:
		if b.Status == "claimed" {
			sendError(w, http.StatusConflict, "Completed data must not be deleted")
			return
		}
		delete(store.data, id)
		sendJSON(w, http.StatusOK, map[string]string{"status": "deleted"})

	default:
		sendError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

func handleBarangKlaim(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		sendError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		sendError(w, http.StatusBadRequest, "Id must contain number")
		return
	}

	store.mu.Lock()
	defer store.mu.Unlock()

	b, exists := store.data[id]
	if !exists {
		sendError(w, http.StatusNotFound, "Data Not Found")
		return
	}

	if b.Status == "claimed" {
		sendError(w, http.StatusConflict, "Data already done")
		return
	}

	b.Status = "claimed"
	store.data[id] = b

	sendJSON(w, http.StatusOK, b)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", handleHealthz)
	mux.HandleFunc("/barang", handleBarang)
	mux.HandleFunc("/barang/{id}", handleBarangByID)
	mux.HandleFunc("/barang/{id}/klaim", handleBarangKlaim)

	addr := ":" + port
	log.Printf("Server is running on http://localhost%s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
