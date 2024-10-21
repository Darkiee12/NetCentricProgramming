package models

import (
	"errors"
	"fmt"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"time"
)

type Game struct {
	Leaderboard Players
	Sessions    map[int64]Session
}

func (g *Game) Gameloop(conn net.Conn) {
	g.Leaderboard.Create()
	g.Sessions = make(map[int64]Session)
	session := &Session{}
	session.New(conn)
	g.Sessions[session.SessionId] = *session
	g.MenuBoard(session)
}

func (g *Game) SignUp(session *Session) error {
	session.Conn.Write([]byte("[Response] Enter a username: "))
	username := handleInput(session.Conn)

	if _, exists := g.Leaderboard.Records[username]; exists {
		return errors.New("username already taken")
	}

	session.Conn.Write([]byte("[Response] Enter a password: "))
	password := handleInput(session.Conn)

	session.Conn.Write([]byte("[Response] Enter your full name: "))
	fullname := handleInput(session.Conn)

	session.Conn.Write([]byte("[Response] Enter your email address: "))
	email := handleInput(session.Conn)

	session.Conn.Write([]byte("[Response] Enter your address: "))
	address := handleInput(session.Conn)

	g.Leaderboard.Add(username, password, fullname, email, address)
	session.Conn.Write([]byte("Account created successfully. You can now sign in.\n"))

	return nil
}

func (g *Game) MainGame(session *Session) {
	rng := rand.Intn(100) + 1
	session.Conn.Write([]byte("[Response] The number has been generated (1-100)! Guess the number: "))
	score := 0
	for {
		input := handleInput(session.Conn)
		guess, err := strconv.Atoi(input)
		if err != nil {
			session.Conn.Write([]byte("Please enter a valid number."))
			continue
		}
		if guess < rng {
			session.Conn.Write([]byte(fmt.Sprintf("[Response] %d is too low. Try again!", guess)))
		} else if guess > rng {
			session.Conn.Write([]byte(fmt.Sprintf("[Response] %d is too high. Try again!", guess)))
		} else {
			score++
			session.Conn.Write([]byte(fmt.Sprintf("Attempts: %d\n", score)))
			session.Conn.Write([]byte("Congratulations! You guessed the correct number.\n"))
			g.Leaderboard.UpdateBestAttempt(session.Player, score)
			session.Conn.Write([]byte("[Response] Do you want to play again? (yes/no): "))
			playAgain := handleInput(session.Conn)
			if strings.ToLower(playAgain) == "yes" {
				g.MainGame(session)
			} else {
				session.Conn.Write([]byte("Thanks for playing!"))
				g.MenuBoard(session)
			}
			break
		}
		score++
	}
}

func handleInput(conn net.Conn) string {
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Error reading from client:", err)
		conn.Close()
		return ""
	}

	return string(buffer[:n])
}

func (g *Game) SignIn(session *Session) (Player, error) {
	session.Conn.Write([]byte("[Response] Enter your username: "))
	username := handleInput(session.Conn)
	session.Conn.Write([]byte("[Response] Enter your password: "))
	password := handleInput(session.Conn)

	player, err := g.Leaderboard.SignIn(username, password)
	if err != nil {
		session.Conn.Write([]byte("Invalid username or password.\n"))
		return Player{}, err
	}
	return player, nil
}

func (g *Game) MenuBoard(session *Session) {
	if session.IsLoggedIn() {
		session.Conn.Write([]byte(fmt.Sprintf("Welcome, %s!\n1. Play\n2. Sign out\n3. View leaderboard\n4. Export leaderboard\n5. Import leaderboard\n6. Exit\n", session.Player.FullName)))
		session.Conn.Write([]byte("[Response] Enter your choice: "))
		option := handleInput(session.Conn)
		switch option {
		case "1":
			debug(fmt.Sprintf("%s starts the game", session.Player.Username))
			g.MainGame(session)
		case "2":
			debug(fmt.Sprintf("%s logs out", session.Player.Username))
			session.LogOut()
			session.Conn.Write([]byte("Successfully logged out."))
			session.Conn.Write([]byte("[Response] Press any key to return to Menu..."))
			handleInput(session.Conn)
			g.MenuBoard(session)
		case "3":
			debug(fmt.Sprintf("%s views the leaderboard", session.Player.Username))
			g.ViewLeaderboard(session)
		case "4":
			debug(fmt.Sprintf("%s exports the leaderboard", session.Player.Username))
			g.Leaderboard.Export()
			session.Conn.Write([]byte("Leaderboard exported successfully.\n"))
			session.Conn.Write([]byte("[Response] Press any key to return to Menu..."))
			handleInput(session.Conn)
			g.MenuBoard(session)
		case "5":
			debug(fmt.Sprintf("%s imports the leaderboard", session.Player.Username))
			g.Leaderboard.Import()
			session.Conn.Write([]byte("Leaderboard imported successfully."))
			session.Conn.Write([]byte("[Response] Press any key to return to Menu...\n"))
			handleInput(session.Conn)
			g.MenuBoard(session)
		case "6":
			debug(fmt.Sprintf("%s exits the game", session.Player.Username))
			session.Conn.Write([]byte("Exiting the game. Goodbye!"))
			g.Leaderboard.Export() // Save leaderboard before exiting
			session.Conn.Close()
		default:
			session.Conn.Write([]byte("Invalid option. Please try again."))
			g.MenuBoard(session)
		}

	} else {
		session.Conn.Write([]byte("Welcome to the game!\n1. Sign in\n2. Sign up\n3. View leaderboard\n4. Import leaderboard\n5. Exit\n"))
		session.Conn.Write([]byte("[Response] Enter your choice: "))
		option := handleInput(session.Conn)
		switch option {
		case "1":
			player, err := g.SignIn(session)
			if err != nil {
				session.Conn.Write([]byte("Invalid username or password. Please try again.\n"))
				g.MenuBoard(session)
			}
			session.LogIn(player)
			debug(fmt.Sprintf("%s logs in", session.Player.Username))
			g.MenuBoard(session)
		case "2":
			err := g.SignUp(session)
			if err != nil {
				session.Conn.Write([]byte("Username already taken. Please try again.\n"))
			}
			g.MenuBoard(session)
		case "3":
			g.ViewLeaderboard(session)
		case "4":
			g.Leaderboard.Import()
			session.Conn.Write([]byte("Leaderboard imported successfully.\n"))
			session.Conn.Write([]byte("[Response] Press any key to return to Menu...\n"))
			handleInput(session.Conn)
			g.MenuBoard(session)
		case "5":
			session.Conn.Write([]byte("Exiting the game. Goodbye!\n"))
			session.Conn.Close()
		default:
			session.Conn.Write([]byte("Invalid option. Please try again.\n"))
			g.MenuBoard(session)
		}
	}
}

func (g *Game) ViewLeaderboard(session *Session) {
	var leaderboard string
	for _, player := range g.Leaderboard.Records {
		leaderboard += fmt.Sprintf("Username: %s, Best Score: %d\n", player.Username, player.BestScore)
	}
	if leaderboard == "" {
		leaderboard = "No players in the leaderboard.\n"
	}
	session.Conn.Write([]byte(leaderboard))
	session.Conn.Write([]byte("[Response] Press any key to return to Menu...\n"))
	handleInput(session.Conn)
	g.MenuBoard(session)
}

func debug(message string) {
	now := time.Now()
	fmt.Printf("[DEBUG] %s: %s\n", now.Format(time.RFC3339), message)
}
