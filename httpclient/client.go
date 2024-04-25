package httpclient

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type HttpClient struct {
	Client    *http.Client
	AuthToken string
}

type GameStatus struct {
	Nick           string   `json:"nick"`
	GameStatus     string   `json:"game_status"`
	LastGameStatus string   `json:"last_game_status"`
	Opponent       string   `json:"opponenet"`
	OpponentShots  []string `json:"opp_shots"`
}

func (c *HttpClient) GetGameStatus() (GameStatus, int) {

	req, err := http.NewRequest("GET", "https://go-pjatk-server.fly.dev/api/game", nil)
	if err != nil {
		log.Println("Error creating a request to check game status: ", err)
	}
	req.Header.Add("X-Auth-Token", c.AuthToken)

	resp, err := c.Client.Do(req)
	if err != nil {
		log.Println("Error sending request while checking game status")
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading body while checking game status: ", err)
	}
	var status GameStatus
	err = json.Unmarshal(body, &status)
	if err != nil {
		log.Println("Error unmarshaling game status: ", err)
	}
	return status, resp.StatusCode
}

type GameConfig struct {
	Wpbot  bool   `json:"wpbot"`
	Desc   string `json:"desc"`
	Nick   string `json:"nick"`
	Coords []byte `json:"coords"`
	Target string `json:"target_nick"`
}

func (c *HttpClient) GetAuthToken(cfg *GameConfig) (string, int) {

	bm, err := json.Marshal(cfg)
	if err != nil {
		log.Fatal("Error marshaling", err)
	}

	req, err := http.NewRequest("POST", "https://go-pjatk-server.fly.dev/api/game", bytes.NewReader(bm))
	if err != nil {
		log.Fatal("Error creating a request", err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.Client.Do(req)
	if err != nil {
		log.Fatal("Error requesting /game", err)
	}
	return resp.Header.Get("X-Auth-Token"), resp.StatusCode

}

func (c *HttpClient) GetGameBoard() ([]string, int) {
	req, err := http.NewRequest("GET", "https://go-pjatk-server.fly.dev/api/game/board", nil)
	if err != nil {
		log.Println("Error creating request while getting game board", err)
	}
	req.Header.Add("X-Auth-Token", c.AuthToken)

	resp, err := c.Client.Do(req)
	if err != nil {
		log.Println("Error requesting game board ", err)
	}

	type board struct {
		Brd []string `json:"board"`
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading response body while getting game board ", err)
	}
	var brd board
	err = json.Unmarshal(body, &brd)
	if err != nil {
		log.Println("Error unmarshaling response body while getting game board ", err)
	}

	return brd.Brd, resp.StatusCode
}

func (c *HttpClient) Fire(toFire string) (string, int) {
	type coord struct {
		Coord string `json:"coord"`
	}

	crd := &coord{
		Coord: toFire,
	}

	crdm, err := json.Marshal(crd)
	if err != nil {
		log.Println("Error marshaling fire coords")
	}

	req, err := http.NewRequest("POST", "https://go-pjatk-server.fly.dev/api/game/fire", bytes.NewReader(crdm))
	if err != nil {
		log.Println("Error creating request while firing", err)
	}
	req.Header.Add("X-Auth-Token", c.AuthToken)

	resp, err := c.Client.Do(req)
	if err != nil {
		log.Println("Error firing", err)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading response body while firing", err)
	}
	type result struct {
		Result string `json:"result"`
	}
	var res result
	err = json.Unmarshal(body, &res)
	if err != nil {
		log.Println("Error unmarshaling response while firing", err)
	}

	return res.Result, resp.StatusCode

}
