package models

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"os"
)

type Player struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	FullName  string `json:"fullname"`
	Email     string `json:"email"`
	Address   string `json:"address"`
	BestScore int    `json:"best"`
}

type Players struct {
	Records map[string]Player `json:"players"`
}

func (p *Player) Create(username, password, fullname string, emails, addresses string) {
	p.Username = username
	p.Password = base64.StdEncoding.EncodeToString([]byte(password))
	p.FullName = fullname
	p.Email = emails
	p.Address = addresses
	p.BestScore = math.MaxInt32
}

func (p *Players) Create() {
	p.Records = make(map[string]Player)
}

func (list *Players) UpdateBestAttempt(player Player, score int) {
	player.BestScore = min(player.BestScore, score)
	list.Records[player.Username] = player
	list.Export()
}

func (p *Players) Add(username, password, fullname, email, address string) {
	player := Player{}
	player.Create(username, password, fullname, email, address)
	if p.Records == nil {
		p.Records = make(map[string]Player)
	}
	p.Records[username] = player
}

func (p *Players) SignIn(username, password string) (Player, error) {
	if player, ok := p.Records[username]; ok {
		if base64.StdEncoding.EncodeToString([]byte(password)) == player.Password {
			return player, nil
		}
	}
	return Player{}, errors.New("invalid username or password")
}

func (p *Players) Export() error {
	data, err := json.MarshalIndent(p.Records, "", "  ")
	if err != nil {
		return err
	}
	file, err := os.Create("score.json")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return err
	}
	defer file.Close()
	_, err = file.Write(data)
	if err != nil {
		fmt.Println("Error writing JSON to file:", err)
		return err
	}
	fmt.Println("JSON data has been updated to score.json")
	return nil
}

func (p *Players) Import() error {
	file, err := os.Open("score.json")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return err
	}
	defer file.Close()

	data := make([]byte, 1024)
	count, err := file.Read(data)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return err
	}

	err = json.Unmarshal(data[:count], &p.Records)
	if err != nil {
		fmt.Println("Error unmarshalling JSON data:", err)
		return err
	}

	fmt.Println("JSON data has been read from score.json")
	return nil
}
