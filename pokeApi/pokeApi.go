package pokeApi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"sync"
)

type ResponseAPI struct {
	Count    int
	Next     string
	Previous string
	Results  []struct {
		Name string
		URL  string
	}
}

type Pokemon struct {
	Name    string
	ID      int
	Sprites struct {
		FrontDefault string
	}
	Types []struct {
		Type struct {
			Name string
		}
	}
}

type PokemonResponse struct {
	Name  string `json:"name"`
	ID    int    `json:"id"`
	Img   string `json:"img"`
	Types string `json:"types"`
}

func closeBody(body io.ReadCloser) {
	err := body.Close()
	if err != nil {
		fmt.Println("Erro ao fechar o corpo da resposta:", err)
	}
}

func getListPokemons() (*ResponseAPI, error) {
	apiURL := "https://pokeapi.co/api/v2/pokemon"
	response, err := http.Get(apiURL)
	if err != nil {
		return nil, err
	}
	defer closeBody(response.Body)

	if response.StatusCode != http.StatusOK {
		return nil, errors.New("erro na solicitação HTTP")
	}

	bodyResponse, erro := io.ReadAll(response.Body)
	if erro != nil {
		return nil, erro
	}

	jsonDataAll := ResponseAPI{}
	err = json.Unmarshal(bodyResponse, &jsonDataAll)
	if err != nil {
		return nil, err
	}
	return &jsonDataAll, nil
}

func HandleListPokemons() (*[]PokemonResponse, error) {
	var wg sync.WaitGroup
	var mu sync.Mutex
	baseImg := "https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/other/official-artwork/"

	listPokemons, err := getListPokemons()
	if err != nil {
		return nil, err
	}

	var pokemonsUrl []string
	for _, pokemon := range listPokemons.Results {
		pokemonsUrl = append(pokemonsUrl, pokemon.URL)
	}

	var pokemonResponse = make([]PokemonResponse, len(pokemonsUrl))

	for index, url := range pokemonsUrl {
		wg.Add(1)
		go func(index int, url string) {
			defer wg.Done()

			response, err := http.Get(url)
			if err != nil {
				return
			}
			defer closeBody(response.Body)

			bodyResponse, erro := io.ReadAll(response.Body)
			if erro != nil {
				return
			}

			jsonDataPokemon := Pokemon{}
			err = json.Unmarshal(bodyResponse, &jsonDataPokemon)
			if err != nil {
				return
			}

			mu.Lock()
			defer mu.Unlock()
			pkmTemp := PokemonResponse{
				Name:  jsonDataPokemon.Name,
				ID:    jsonDataPokemon.ID,
				Img:   baseImg + strconv.Itoa(jsonDataPokemon.ID) + ".png",
				Types: jsonDataPokemon.Types[0].Type.Name,
			}
			pokemonResponse[index] = pkmTemp
		}(index, url)
	}
	wg.Wait()

	return &pokemonResponse, nil
}
