package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/kenta-kenta/gaishi-serverside/models"
	"github.com/kenta-kenta/gaishi-serverside/services"
)

func Hello(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Hello, World")
}

func CreateMemo(w http.ResponseWriter, r *http.Request) {
	var req models.CreateMemoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("decode err: %v", err), http.StatusBadRequest)
		return
	}

	memo, err := services.CreateMemo(req)
	if err != nil {
		http.Error(w, fmt.Sprintf("service err: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(memo); err != nil {
		http.Error(w, fmt.Sprintf("encode err: %v", err), http.StatusInternalServerError)
		return
	}
}

func GetMemoByID(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	parts := strings.Split(path, "/")
	if len(parts) < 3 {
		http.Error(w, "invalid URL path", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(parts[3])
	if err != nil {
		http.Error(w, "invalid memo ID", http.StatusBadRequest)
		return
	}

	memo, err := services.GetMemoByID(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("service err: %v", err), http.StatusInternalServerError)
		return
	}
	if memo == nil {
		http.Error(w, "memo not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := json.NewEncoder(w).Encode(memo); err != nil {
		http.Error(w, fmt.Sprintf("encode err: %v", err), http.StatusInternalServerError)
		return
	}
}

func UpdateMemo(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	parts := strings.Split(path, "/")
	if len(parts) < 3 {
		http.Error(w, "invalid URL path", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(parts[3])
	if err != nil {
		http.Error(w, "invalid memo ID", http.StatusBadRequest)
		return
	}

	var req models.UpdateMemoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("decode err: %v", err), http.StatusBadRequest)
		return
	}

	memo, err := services.UpdateMemo(id, req)
	if err != nil {
		http.Error(w, fmt.Sprintf("service err: %v", err), http.StatusInternalServerError)
		return
	}
	if memo == nil {
		http.Error(w, "memo not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := json.NewEncoder(w).Encode(memo); err != nil {
		http.Error(w, fmt.Sprintf("encode err: %v", err), http.StatusInternalServerError)
		return
	}
}

func DeleteMemo(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	parts := strings.Split(path, "/")
	if len(parts) < 3 {
		http.Error(w, "invalid URL path", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(parts[3])
	if err != nil {
		http.Error(w, "invalid memo ID", http.StatusBadRequest)
		return
	}

	if err := services.DeleteMemo(id); err != nil {
		http.Error(w, fmt.Sprintf("service err: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
