package main

import (
	"encoding/json"
	"net/http"
)

type Sala struct {
	ID         int      `json:"id"`
	Nome       string   `json:"nome"`
	Capacidade int      `json:"capacidade"`
	Recursos   []string `json:"recursos"`
}

var salas = []Sala{}
var proximoIDSala = 1

func handleSalas(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		listarSalas(w, r)
	case http.MethodPost:
		criarSala(w, r)
	case http.MethodPut:
		atualizarSala(w, r)
	case http.MethodDelete:
		excluirSala(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Método não permitido"))
	}
}

func listarSalas(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(salas)
}

func criarSala(w http.ResponseWriter, r *http.Request) {
	var sala Sala
	err := json.NewDecoder(r.Body).Decode(&sala)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("JSON inválido"))
		return
	}

	sala.ID = proximoIDSala
	proximoIDSala++
	salas = append(salas, sala)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(sala)
}

func atualizarSala(w http.ResponseWriter, r *http.Request) {
	var salaAtualizada Sala
	err := json.NewDecoder(r.Body).Decode(&salaAtualizada)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("JSON inválido"))
		return
	}

	for i, sala := range salas {
		if sala.ID == salaAtualizada.ID {
			salas[i] = salaAtualizada
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(salaAtualizada)
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("Sala não encontrada"))
}

func excluirSala(w http.ResponseWriter, r *http.Request) {
	var salaExcluida Sala
	err := json.NewDecoder(r.Body).Decode(&salaExcluida)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("JSON inválido"))
		return
	}

	for i, sala := range salas {
		if sala.ID == salaExcluida.ID {
			salas = append(salas[:i], salas[i+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("Sala não encontrada"))
}
