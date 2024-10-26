package models

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"os"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	FullName string `json:"fullName"`
	Email    string `json:"email"`
	Address  string `json:"address"`
}

type Account struct {
	User   User
	active bool
}

type AccountManager struct {
	Accounts map[string]Account
}

func (usr *User) Create(username, password, fullname, email, address string) {
	usr.Username = username
	usr.Password = password
	usr.FullName = fullname
	usr.Email = email
	usr.Address = address
}

func (manager *AccountManager) LogIn(username, password string) (Account, error) {
	account, ok := manager.Accounts[username]
	if ok {
		encodedPassword := base64.StdEncoding.EncodeToString([]byte(password))
		if account.User.Password == encodedPassword {
			account.active = true
			return account, nil
		}
	}
	return Account{}, errors.New("invalid username or password")
}

func (manager *AccountManager) SignUp(user User) error {
	_, ok := manager.Accounts[user.Username]
	if ok {
		return errors.New("username already exists")
	}
	user.Password = base64.StdEncoding.EncodeToString([]byte(user.Password))
	manager.Accounts[user.Username] = Account{User: user, active: false}
	return nil
}

func (manager AccountManager) Init() AccountManager {
	return AccountManager{Accounts: make(map[string]Account)}
}

func (manager *AccountManager) LogOut(username string) Account {
	account, ok := manager.Accounts[username]
	if ok {
		account.active = false
	}
	return Account{}
}

func (manager *AccountManager) Export() error {
	data, err := json.MarshalIndent(manager.Accounts, "", "  ")
	if err != nil {
		return err
	}
	err = os.WriteFile("accounts.json", data, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (manager *AccountManager) Import() error {
	data, err := os.ReadFile("accounts.json")
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, &manager.Accounts)
	if err != nil {
		return err
	}
	return nil
}
