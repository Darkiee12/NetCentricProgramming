package models

import (
	"net"
	"time"
)

type Session struct {
	Conn      net.Conn
	LoggedIn  bool
	Player    Player
	SessionId int64
}

func (c *Session) LogIn(player Player) {
	c.LoggedIn = true
	c.Player = player
}

func (c *Session) New(conn net.Conn) {
	c.Conn = conn
	c.LoggedIn = false
	c.Player = Player{}
	c.SessionId = time.Now().UnixMilli()
}

func (c *Session) LogOut() {
	c.LoggedIn = false
	c.Player = Player{}
}

func (c *Session) IsLoggedIn() bool {
	return c.LoggedIn
}
