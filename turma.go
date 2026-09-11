package main

import (
	"encoding/json"
	"net/http"
)

type Turma struct {
	ID         int    `json:"id"`
	Nome       string `json:"nome"`
	QtdeAlunos int    `json:"qtdeAlunos"`
	SalaID     int    `json:"salaId"`
}

var turmas = []Turma{}
var proximoIDTurma = 1

func handleTurmas(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		listarTurmas(w, r)
	case http.MethodPost:
		criarTurma(w, r)
	case http.MethodPut:
		atualizarTurma(w, r)
	case http.MethodDelete:
		excluirTurma(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Método não permitido"))
	}
}

func listarTurmas(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(turmas)
}

func criarTurma(w http.ResponseWriter, r *http.Request) {
	var turma Turma
	err := json.NewDecoder(r.Body).Decode(&turma)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("JSON inválido"))
		return
	}

	turma.ID = proximoIDTurma
	proximoIDTurma++
	turmas = append(turmas, turma)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(turma)
}

func atualizarTurma(w http.ResponseWriter, r *http.Request) {
	var turmaAtualizada Turma
	err := json.NewDecoder(r.Body).Decode(&turmaAtualizada)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("JSON inválido"))
		return
	}

	for i, turma := range turmas {
		if turma.ID == turmaAtualizada.ID {
			turmas[i] = turmaAtualizada
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(turmaAtualizada)
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("Turma não encontrada"))
}

func excluirTurma(w http.ResponseWriter, r *http.Request) {
	var turmaExcluida Turma
	err := json.NewDecoder(r.Body).Decode(&turmaExcluida)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("JSON inválido"))
		return
	}

	for i, turma := range turmas {
		if turma.ID == turmaExcluida.ID {
			turmas = append(turmas[:i], turmas[i+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("Turma não encontrada"))
}
