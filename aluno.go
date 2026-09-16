package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Aluno struct {
	Matricula string `json:"matricula"`
	Nome      string `json:"nome"`
	Email     string `json:"email"`
}

var alunos = []Aluno{}

func procurarAluno(matricula string) int {
	for i, aluno := range alunos {
		if aluno.Matricula == matricula {
			return i
		}
	}
	return -1
}
// matrícula precisa ter exatamente 9 dígitos
func matriculaValida(matricula string) bool {
	if len(matricula) != 9 {
		return false
	}
	for i := 0; i < len(matricula); i++ {
		if matricula[i] < '0' || matricula[i] > '9' {
			return false
		}
	}
	return true
}

func listarAlunos(c *gin.Context) {
	c.JSON(http.StatusOK, alunos)
}

func buscarAluno(c *gin.Context) {
	i := procurarAluno(c.Param("matricula"))
	if i == -1 {
		c.JSON(http.StatusNotFound, gin.H{"erro": "Aluno não encontrado"})
		return
	}

	c.JSON(http.StatusOK, alunos[i])
}

func criarAluno(c *gin.Context) {
	var aluno Aluno
	err := c.ShouldBindJSON(&aluno)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "JSON inválido"})
		return
	}

	if !matriculaValida(aluno.Matricula) {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "Matrícula deve ter 9 dígitos"})
		return
	}

	if procurarAluno(aluno.Matricula) != -1 {
		c.JSON(http.StatusConflict, gin.H{"erro": "Já existe um aluno com essa matrícula"})
		return
	}

	alunos = append(alunos, aluno)

	c.JSON(http.StatusCreated, aluno)
}

func atualizarAluno(c *gin.Context) {
	i := procurarAluno(c.Param("matricula"))
	if i == -1 {
		c.JSON(http.StatusNotFound, gin.H{"erro": "Aluno não encontrado"})
		return
	}

	var alunoAtualizado Aluno
	err := c.ShouldBindJSON(&alunoAtualizado)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "JSON inválido"})
		return
	}

	alunos[i].Nome = alunoAtualizado.Nome
	alunos[i].Email = alunoAtualizado.Email

	c.JSON(http.StatusOK, alunos[i])
}

func excluirAluno(c *gin.Context) {
	matricula := c.Param("matricula")

	i := procurarAluno(matricula)
	if i == -1 {
		c.JSON(http.StatusNotFound, gin.H{"erro": "Aluno não encontrado"})
		return
	}

	for _, turma := range turmas {
		if alunoEstaNaTurma(turma, matricula) {
			c.JSON(http.StatusConflict, gin.H{"erro": "Aluno está matriculado na turma " + turma.Nome})
			return
		}
	}

	alunos = append(alunos[:i], alunos[i+1:]...)
	c.Status(http.StatusNoContent)
}
