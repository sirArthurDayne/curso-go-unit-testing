package controller

import (
	"catching-pokemons/models"
	"catching-pokemons/util"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

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

//test over the mux router
func TestGetPokemon(t *testing.T) {
    r, err := http.NewRequest(http.MethodGet, "/pokemon/{id}", nil)
    assert.NoError(t, err)

    w := httptest.NewRecorder()
    vars := map[string]string {
        "id": "pikachu",
    }

    r = mux.SetURLVars(r, vars)
    GetPokemon(w, r)

    var expectedPokemon models.Pokemon
    err = util.ReadFileAndUnmarshall(&expectedPokemon, "samples/api_response.json")
    assert.NoError(t ,err)
    var actualPokemon models.Pokemon
    err = json.Unmarshal(w.Body.Bytes(), &actualPokemon)
    assert.NoError(t , err)

    assert.Equal(t, http.StatusOK, w.Result().StatusCode)
    assert.Equal(t, expectedPokemon, actualPokemon)
}

func TestGetPokemonNotFound(t *testing.T) {
    r, err := http.NewRequest(http.MethodGet, "/pokemon/{id}", nil)
    assert.NoError(t, err)

    w := httptest.NewRecorder()
    vars := map[string]string {
        "id": "not-found-id-test",
    }

    r = mux.SetURLVars(r, vars)
    GetPokemon(w, r)

    assert.Equal(t, http.StatusNotFound, w.Result().StatusCode)
}
func TestGetPokemonFailure(t *testing.T) {
    r, err := http.NewRequest(http.MethodGet, "/pokemon/{id}", nil)
    assert.NoError(t, err)

    w := httptest.NewRecorder()
    vars := map[string]string {
        "id": "",
    }

    r = mux.SetURLVars(r, vars)
    GetPokemon(w, r)

    assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
}
