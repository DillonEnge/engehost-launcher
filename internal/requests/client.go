package requests

import (
	"encoding/json"
	"net/http"
)

type Client struct {
	URL string
}

type Game struct {
	ID                 int    `json:"id"`
	Name               string `json:"name"`
	RepoName           string `json:"repo_name"`
	RepoOwner          string `json:"repo_owner"`
	IconURL            string `json:"icon_url"`
	BackgroundImageURL string `json:"background_image_url"`
}

func (c *Client) GetGames() ([]Game, error) {
	var games []Game

	res, err := http.Get(c.URL + "/games")
	if err != nil {
		return nil, err
	}

	json.NewDecoder(res.Body).Decode(&games)

	return games, nil
}
