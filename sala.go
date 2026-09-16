package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Sala struct {
	ID         int      `json:"id"`
	Nome       string   `json:"nome"`
	Capacidade int      `json:"capacidade"`
	Recursos   []string `json:"recursos"`
}

var salas = []Sala{}
var proximoIDSala = 1

// devolve a posição da sala na lista, ou -1 se não existir
func procurarSala(id int) int {
	for i, sala := range salas {
		if sala.ID == id {
			return i
		}
	}
	return -1
}

func listarSalas(c *gin.Context) {
	c.JSON(http.StatusOK, salas)
}

func criarSala(c *gin.Context) {
	var sala Sala
	err := c.ShouldBindJSON(&sala)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "JSON inválido"})
		return
	}

	if sala.Capacidade <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "Capacidade deve ser maior que zero"})
		return
	}

	sala.ID = proximoIDSala
	proximoIDSala++
	salas = append(salas, sala)

	c.JSON(http.StatusCreated, sala)
}

func atualizarSala(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "ID inválido"})
		return
	}

	i := procurarSala(id)
	if i == -1 {
		c.JSON(http.StatusNotFound, gin.H{"erro": "Sala não encontrada"})
		return
	}

	var salaAtualizada Sala
	err = c.ShouldBindJSON(&salaAtualizada)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "JSON inválido"})
		return
	}

	if salaAtualizada.Capacidade <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "Capacidade deve ser maior que zero"})
		return
	}

	// o ID não muda, só os dados da sala
	salas[i].Nome = salaAtualizada.Nome
	salas[i].Capacidade = salaAtualizada.Capacidade
	salas[i].Recursos = salaAtualizada.Recursos

	c.JSON(http.StatusOK, salas[i])
}

func excluirSala(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "ID inválido"})
		return
	}

	i := procurarSala(id)
	if i == -1 {
		c.JSON(http.StatusNotFound, gin.H{"erro": "Sala não encontrada"})
		return
	}

	// não deixa excluir uma sala que está em uso por alguma turma
	for _, turma := range turmas {
		if turma.Alocada && turma.SalaID == id {
			c.JSON(http.StatusConflict, gin.H{"erro": "Sala está alocada na turma " + turma.Nome})
			return
		}
	}

	salas = append(salas[:i], salas[i+1:]...)
	c.Status(http.StatusNoContent)
}
