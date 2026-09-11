package main

import (
	"encoding/json"
	"net/http"
)

type Aluno struct {
	ID        int    `json:"id"`
	Nome      string `json:"nome"`
	Matricula string `json:"matricula"`
	TurmaID   int    `json:"turmaId"`
}

var alunos = []Aluno{}
var proximoIDAluno = 1

func handleAlunos(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		listarAlunos(w, r)
	case http.MethodPost:
		criarAluno(w, r)
	case http.MethodPut:
		atualizarAluno(w, r)
	case http.MethodDelete:
		excluirAluno(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Método não permitido"))
	}
}

func listarAlunos(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(alunos)
}

func criarAluno(w http.ResponseWriter, r *http.Request) {
	var aluno Aluno
	err := json.NewDecoder(r.Body).Decode(&aluno)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("JSON inválido"))
		return
	}

	aluno.ID = proximoIDAluno
	proximoIDAluno++
	alunos = append(alunos, aluno)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(aluno)
}

func atualizarAluno(w http.ResponseWriter, r *http.Request) {
	var alunoAtualizado Aluno
	err := json.NewDecoder(r.Body).Decode(&alunoAtualizado)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("JSON inválido"))
		return
	}

	for i, aluno := range alunos {
		if aluno.ID == alunoAtualizado.ID {
			alunos[i] = alunoAtualizado
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(alunoAtualizado)
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("Aluno não encontrado"))
}

func excluirAluno(w http.ResponseWriter, r *http.Request) {
	var alunoExcluido Aluno
	err := json.NewDecoder(r.Body).Decode(&alunoExcluido)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("JSON inválido"))
		return
	}

	for i, aluno := range alunos {
		if aluno.ID == alunoExcluido.ID {
			alunos = append(alunos[:i], alunos[i+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("Aluno não encontrado"))
}
