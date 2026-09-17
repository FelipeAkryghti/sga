package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.New()

	// Uso dos Middlewares globais nativos e personalizados
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// 4. Mapeamento de Rotas sob Grupo Versionado
	v1 := r.Group("/api/v1")
	{
		// Monitoramento da API
		v1.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"status":    "healthy",
				"timestamp": time.Now(),
				"version":   "1.0.0",
			})
		})

		// Domínio de Salas
		v1.POST("/salas", criarSala)
		v1.GET("/salas", listarSalas)
		v1.PUT("/salas/:id", atualizarSala)
		v1.DELETE("/salas/:id", excluirSala)

		// Domínio de Alunos
		v1.POST("/alunos", criarAluno)
		v1.GET("/alunos", listarAlunos)
		v1.GET("/alunos/:matricula", buscarAluno)
		v1.PUT("/alunos/:matricula", atualizarAluno)
		v1.DELETE("/alunos/:matricula", excluirAluno)

		// Domínio de Turmas (Classes)
		v1.POST("/turmas", criarTurma)
		v1.GET("/turmas", listarTurmas)
		v1.PUT("/turmas/:id", atualizarTurma)
		v1.DELETE("/turmas/:id", excluirTurma)
		v1.POST("/turmas/:id/alunos", matricularAluno)
		v1.GET("/turmas/:id/alunos", listarAlunosDaTurma)
		v1.POST("/turmas/:id/alocar", alocarSala)
	}

	r.Run(":8080")
}
