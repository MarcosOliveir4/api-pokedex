package main

import (
	"api-pokedex/pokeApi"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:  []string{"*"},
		AllowMethods:  []string{"GET"},
		AllowHeaders:  []string{"Origin"},
		ExposeHeaders: []string{"Content-Length"},
	}))

	r.GET("/", func(c *gin.Context) {
		pokemonList, err := pokeApi.HandleListPokemons()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  http.StatusBadRequest,
				"message": err.Error(),
			})
			c.Abort()
			return
		}
		params := c.Request.URL.Query()
		fmt.Println(params)
		c.JSON(http.StatusOK, pokemonList)
	})

	errGin := r.Run()
	if errGin != nil {
		fmt.Println("Erro ao iniciar o servidor:", errGin)
		return
	}
}
