package controller

import (
	"catching-pokemons/models"
	"catching-pokemons/util"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

// func TestGetPokemonFromPokeApi(t *testing.T) {
// 	pokemon, err := GetPokemonFromPokeApi("pikachu")
// 	assert.NoError(t, err)
//
// 	var expected models.PokeApiPokemonResponse
// 	err = util.ReadFileAndUnmarshall(&expected, "samples/poke_api_readed.json")
// 	assert.NoError(t, err)
//
// 	assert.Equal(t, expected, pokemon)
// }

func TestGetPokemonFromPokeApiSuccess(t *testing.T) {
	httpmock.Activate() // enable before running TestGetPokemonFromPokeApiMocks
	defer httpmock.DeactivateAndReset()

	body, err := ioutil.ReadFile("samples/pokeapi_response.json")
	httpmock.RegisterResponder(http.MethodGet, "https://pokeapi.co/api/v2/pokemon/pikachu", httpmock.NewBytesResponder(http.StatusOK, body))

	pokemon, err := GetPokemonFromPokeApi("pikachu")
	assert.NoError(t, err)

	var expected models.PokeApiPokemonResponse
	err = util.ReadFileAndUnmarshall(&expected, "samples/pokeapi_response.json")
	assert.NoError(t, err)

	assert.Equal(t, expected, pokemon)
}

func TestGetPokemonFromPokeApiNotFound(t *testing.T) {
	httpmock.Activate() // enable before running tests
	defer httpmock.DeactivateAndReset()

	body, err := ioutil.ReadFile("samples/pokeapi_response.json")
	httpmock.RegisterResponder(http.MethodGet, "https://pokeapi.co/api/v2/pokemon/not-found-id-test", httpmock.NewBytesResponder(http.StatusNotFound, body))

	_, err = GetPokemonFromPokeApi("not-found-id-test")
	assert.Error(t, err)

	assert.EqualError(t, err, ErrPokemonNotFound.Error())
}

func TestGetPokemonFromPokeApiFailure(t *testing.T) {
	httpmock.Activate() // enable before running tests
	defer httpmock.DeactivateAndReset()

	body, err := ioutil.ReadFile("samples/pokeapi_response.json")
	httpmock.RegisterResponder(http.MethodGet, "https://pokeapi.co/api/v2/pokemon/", httpmock.NewBytesResponder(http.StatusInternalServerError, body))

	_, err = GetPokemonFromPokeApi("")
	assert.Error(t, err)

	assert.EqualError(t, err, ErrPokemonApiFailure.Error())
}
