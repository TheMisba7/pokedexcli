package model

type Pokemon struct {
	BaseExperience int    `json:"base_experience"`
	Weight         int    `json:"weight"`
	Height         int    `json:"height"`
	Name           string `json:"name"`
	Stats          []struct {
		BaseStat int `json:"base_stat"`
		Effort   int `json:"effort"`
		Stat     struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Slot int `json:"slot"`
		Type struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"type"`
	} `json:"types"`
}

type PokemonRS struct {
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}
type LocationAreaRS struct {
	Count    int
	Next     string
	Previous string
	Results  []struct {
		Name string
		URL  string
	}
}

type Command struct {
	Name        string
	Description string
	Callback    func([]string) error
	Config      *Config
}

type Config struct {
	Next     string
	Previous string
}

type PokeApi interface {
	GetLocationArea(url string) (LocationAreaRS, error)
	GetListOfPokemon(area string) (PokemonRS, error)
	GetPokemon(url string) (Pokemon, error)
}
