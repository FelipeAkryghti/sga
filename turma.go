package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type Turma struct {
	ID         int      `json:"id"`
	Nome       string   `json:"nome"`
	Disciplina string   `json:"disciplina"`
	Professor  string   `json:"professor"`
	Matriculas []string `json:"matriculas"`
	QtdeAlunos int      `json:"qtdeAlunos"`
	Alocada    bool     `json:"alocada"`
	SalaID     int      `json:"salaId"`
	DiaSemana  string   `json:"diaSemana"`
	HoraInicio string   `json:"horaInicio"`
	HoraFim    string   `json:"horaFim"`
}

// corpo do POST /turmas/:id/alunos
type PedidoMatricula struct {
	Matricula string `json:"matricula"`
}

// corpo do POST /turmas/:id/alocar
type PedidoAlocacao struct {
	SalaID     int    `json:"salaId"`
	DiaSemana  string `json:"diaSemana"`
	HoraInicio string `json:"horaInicio"`
	HoraFim    string `json:"horaFim"`
}

var turmas = []Turma{}
var proximoIDTurma = 1

var diasDaSemana = []string{"segunda", "terca", "quarta", "quinta", "sexta", "sabado", "domingo"}

// devolve a posição da turma na lista, ou -1 se não existir
func procurarTurma(id int) int {
	for i, turma := range turmas {
		if turma.ID == id {
			return i
		}
	}
	return -1
}

func alunoEstaNaTurma(turma Turma, matricula string) bool {
	for _, m := range turma.Matriculas {
		if m == matricula {
			return true
		}
	}
	return false
}

func diaValido(dia string) bool {
	for _, d := range diasDaSemana {
		if d == dia {
			return true
		}
	}
	return false
}

// horário no formato HH:MM, ex: 19:00
func horaValida(hora string) bool {
	_, err := time.Parse("15:04", hora)
	return err == nil
}

// dois horários se sobrepõem quando um começa antes do outro terminar
// (como o formato é sempre HH:MM, dá pra comparar direto como texto)
func horariosConflitam(inicio1, fim1, inicio2, fim2 string) bool {
	return inicio1 < fim2 && fim1 > inicio2
}

func listarTurmas(c *gin.Context) {
	c.JSON(http.StatusOK, turmas)
}

func criarTurma(c *gin.Context) {
	var dados Turma
	err := c.ShouldBindJSON(&dados)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "JSON inválido"})
		return
	}

	// só nome, disciplina e professor vêm do cliente; o resto começa vazio
	turma := Turma{
		ID:         proximoIDTurma,
		Nome:       dados.Nome,
		Disciplina: dados.Disciplina,
		Professor:  dados.Professor,
		Matriculas: []string{},
	}
	proximoIDTurma++
	turmas = append(turmas, turma)

	c.JSON(http.StatusCreated, turma)
}

func atualizarTurma(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "ID inválido"})
		return
	}

	i := procurarTurma(id)
	if i == -1 {
		c.JSON(http.StatusNotFound, gin.H{"erro": "Turma não encontrada"})
		return
	}

	var turmaAtualizada Turma
	err = c.ShouldBindJSON(&turmaAtualizada)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "JSON inválido"})
		return
	}

	// matrículas e alocação não mudam por aqui, só os dados básicos
	turmas[i].Nome = turmaAtualizada.Nome
	turmas[i].Disciplina = turmaAtualizada.Disciplina
	turmas[i].Professor = turmaAtualizada.Professor

	c.JSON(http.StatusOK, turmas[i])
}

func excluirTurma(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "ID inválido"})
		return
	}

	i := procurarTurma(id)
	if i == -1 {
		c.JSON(http.StatusNotFound, gin.H{"erro": "Turma não encontrada"})
		return
	}

	turmas = append(turmas[:i], turmas[i+1:]...)
	c.Status(http.StatusNoContent)
}

func matricularAluno(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "ID inválido"})
		return
	}

	t := procurarTurma(id)
	if t == -1 {
		c.JSON(http.StatusNotFound, gin.H{"erro": "Turma não encontrada"})
		return
	}

	var pedido PedidoMatricula
	err = c.ShouldBindJSON(&pedido)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "JSON inválido"})
		return
	}

	if procurarAluno(pedido.Matricula) == -1 {
		c.JSON(http.StatusNotFound, gin.H{"erro": "Aluno não encontrado"})
		return
	}

	// regra 1: o aluno não pode estar duas vezes na mesma turma
	if alunoEstaNaTurma(turmas[t], pedido.Matricula) {
		c.JSON(http.StatusConflict, gin.H{"erro": "Aluno já está matriculado nesta turma"})
		return
	}

	// as regras 2 e 3 só fazem sentido se a turma já tem sala e horário
	if turmas[t].Alocada {
		// regra 2: não pode passar da capacidade da sala
		s := procurarSala(turmas[t].SalaID)
		if turmas[t].QtdeAlunos+1 > salas[s].Capacidade {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"erro": "Sala lotada, capacidade máxima é " + strconv.Itoa(salas[s].Capacidade)})
			return
		}

		// regra 3: o aluno não pode ter outra turma no mesmo dia e horário
		for _, outra := range turmas {
			if outra.ID != id && outra.Alocada && outra.DiaSemana == turmas[t].DiaSemana && alunoEstaNaTurma(outra, pedido.Matricula) {
				if horariosConflitam(turmas[t].HoraInicio, turmas[t].HoraFim, outra.HoraInicio, outra.HoraFim) {
					c.JSON(http.StatusConflict, gin.H{"erro": "Aluno já tem aula nesse horário na turma " + outra.Nome})
					return
				}
			}
		}
	}

	turmas[t].Matriculas = append(turmas[t].Matriculas, pedido.Matricula)
	turmas[t].QtdeAlunos = len(turmas[t].Matriculas)

	c.JSON(http.StatusCreated, turmas[t])
}

func listarAlunosDaTurma(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "ID inválido"})
		return
	}

	t := procurarTurma(id)
	if t == -1 {
		c.JSON(http.StatusNotFound, gin.H{"erro": "Turma não encontrada"})
		return
	}

	lista := []Aluno{}
	for _, matricula := range turmas[t].Matriculas {
		a := procurarAluno(matricula)
		if a != -1 {
			lista = append(lista, alunos[a])
		}
	}

	c.JSON(http.StatusOK, lista)
}

func alocarSala(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "ID inválido"})
		return
	}

	t := procurarTurma(id)
	if t == -1 {
		c.JSON(http.StatusNotFound, gin.H{"erro": "Turma não encontrada"})
		return
	}

	var pedido PedidoAlocacao
	err = c.ShouldBindJSON(&pedido)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "JSON inválido"})
		return
	}

	s := procurarSala(pedido.SalaID)
	if s == -1 {
		c.JSON(http.StatusNotFound, gin.H{"erro": "Sala não encontrada"})
		return
	}

	if !diaValido(pedido.DiaSemana) {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "Dia da semana inválido, use: segunda, terca, quarta, quinta, sexta, sabado ou domingo"})
		return
	}

	if !horaValida(pedido.HoraInicio) || !horaValida(pedido.HoraFim) {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "Horário inválido, use o formato HH:MM"})
		return
	}

	if pedido.HoraInicio >= pedido.HoraFim {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "Horário de início deve ser antes do horário de término"})
		return
	}

	// regra 1: a sala precisa caber todos os alunos já matriculados
	if salas[s].Capacidade < turmas[t].QtdeAlunos {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"erro": "Sala não comporta a turma, capacidade é " + strconv.Itoa(salas[s].Capacidade) + " e a turma tem " + strconv.Itoa(turmas[t].QtdeAlunos) + " alunos"})
		return
	}

	// regra 2: a sala não pode ter outra turma no mesmo dia com horário sobreposto
	for _, outra := range turmas {
		if outra.ID != id && outra.Alocada && outra.SalaID == pedido.SalaID && outra.DiaSemana == pedido.DiaSemana {
			if horariosConflitam(pedido.HoraInicio, pedido.HoraFim, outra.HoraInicio, outra.HoraFim) {
				c.JSON(http.StatusConflict, gin.H{"erro": "Sala já está ocupada nesse horário pela turma " + outra.Nome})
				return
			}
		}
	}

	turmas[t].Alocada = true
	turmas[t].SalaID = pedido.SalaID
	turmas[t].DiaSemana = pedido.DiaSemana
	turmas[t].HoraInicio = pedido.HoraInicio
	turmas[t].HoraFim = pedido.HoraFim

	c.JSON(http.StatusOK, turmas[t])
}
