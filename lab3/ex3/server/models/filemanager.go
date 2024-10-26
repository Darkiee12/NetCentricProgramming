package models

import (
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type File struct {
	id   int64
	name string
	path string
	size int64 // in bytes
}

func (f *File) String() string {
	mib := f.size / 1024 / 1024
	return fmt.Sprintf("File: %s\nPath: %s\nSize: %d MiB\n", f.name, f.path, mib)
}

type FileManager struct {
	Am    AccountManager
	Files map[int64]File
}

func (fm *FileManager) Init() {
	fm.Am = AccountManager{}.Init()
	fm.Files = make(map[int64]File)
	if err := fm.Update(); err != nil {
		fmt.Println("Error updating files:", err)
		return
	}
	fmt.Println("File manager is initialized successfully")
}

// Update scans the asset directory and updates the file list
func (fm *FileManager) Update() error {
	assetDir := "./assets"
	err := filepath.Walk(assetDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			file := File{
				id:   time.Now().UnixNano(),
				name: info.Name(),
				path: path,
				size: info.Size(),
			}
			fm.Files[file.id] = file
			fmt.Println("Inserting file:", info.Name())
			time.Sleep(1 * time.Microsecond)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("error scanning directory: %w", err)
	}
	fmt.Println("Successfully updated files")
	return nil
}

func (fm *FileManager) ListFiles() string {
	fileList := "Available files:\n"
	for _, file := range fm.Files {
		fileList += fmt.Sprintf("%d. %s\n", file.id, file.name)
	}
	return fileList
}

func (fm *FileManager) Get(id string) (File, error) {
	fileId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return File{}, fmt.Errorf("error parsing id: %w", err)
	}
	file, exists := fm.Files[fileId]
	if !exists {
		return File{}, fmt.Errorf("file not found with id: %d", fileId)
	}

	return file, nil
}

// Method to handle file transfer to the client
func (fm *FileManager) SendFile(conn net.Conn, file File) {
	defer conn.Close()
	fmt.Println("Client connected:", conn.RemoteAddr().String())
	err := fm.sendFileHelper(conn, file)
	if err != nil {
		fmt.Println("Error sending file:", err)
		return
	}

	fmt.Println("[Done] File sent successfully:", file.name)
}

// Method to send the specified file to the client
func (fm *FileManager) sendFileHelper(conn net.Conn, file File) error {
	conn.Write([]byte(fmt.Sprintf("[Download] Sending file: %s\n", file.name)))
	f, err := os.Open(file.path)
	if err != nil {
		return fmt.Errorf("error opening file: %w", err)
	}
	defer f.Close()

	if _, err := conn.Write([]byte(fmt.Sprintf("%d\n", file.size))); err != nil {
		return fmt.Errorf("error sending file size: %w", err)
	}

	if _, err := io.Copy(conn, f); err != nil {
		return fmt.Errorf("error sending file content: %w", err)
	}
	return nil
}

func (fm *FileManager) Menu(conn net.Conn, account *Account) error {
	for {
		if account.active {
			conn.Write([]byte(fmt.Sprintf("Welcome, %s!\n", account.User.FullName)))
			conn.Write([]byte("1. Download a file\n2. Log out\n3. Export accounts\n4. Exit\n"))
			conn.Write([]byte("[Response] Choose option:"))
			option, err := handleInput(conn)
			if err != nil {
				return fmt.Errorf("error handling input: %w", err)
			}
			switch option {
			case "1":
				conn.Write([]byte(fm.ListFiles()))
				conn.Write([]byte("[Response] File ID:"))
				fileID, err := handleInput(conn)
				if err != nil {
					return fmt.Errorf("error handling input: %w", err)
				}
				file, err := fm.Get(fileID)
				if err != nil {
					conn.Write([]byte(fmt.Sprintf("Error: %s\n", err.Error())))
					continue
				}
				conn.Write([]byte(file.String()))
				conn.Write([]byte("[Response] Confirm download? (y/n):"))
				confirm, err := handleInput(conn)
				if err != nil {
					return fmt.Errorf("error handling input: %w", err)
				}
				if confirm == "y" {
					fm.SendFile(conn, file)
					continue
				}
			case "2":
				fm.Am.LogOut(account.User.Username)
				return nil // Exit the menu loop after logout
			case "3":
				err := fm.Am.Export()
				if err != nil {
					conn.Write([]byte(fmt.Sprintf("Error: %s\n", err.Error())))
				}
			case "4":
				return nil
			default:
				conn.Write([]byte("Invalid option. Please try again.\n"))
			}
		} else {
			conn.Write([]byte("Welcome to the file manager. Please sign in to download a file\n"))
			conn.Write([]byte("1. Sign in\n2. Sign up\n3. Import accounts\n4. Exit\n"))
			conn.Write([]byte("[Response] Choose option:"))
			option, err := handleInput(conn)
			if err != nil {
				return fmt.Errorf("error handling input: %w", err)
			}
			switch option {
			case "1":
				account, err := fm.signIn(conn)
				if err != nil {
					conn.Write([]byte(fmt.Sprintf("Error: %s\n", err.Error())))
					continue
				}
				return fm.Menu(conn, &account)

			case "2":
				err := fm.signUp(conn)
				if err != nil {
					conn.Write([]byte(fmt.Sprintf("Error: %s\n", err.Error())))
					continue
				}
				conn.Write([]byte("Account created successfully!\n"))
			case "3":
				err := fm.Am.Import()
				if err != nil {
					conn.Write([]byte(fmt.Sprintf("Error: %s\n", err.Error())))
					continue
				}
			case "4":
				return nil
			default:
				conn.Write([]byte("Invalid option. Please try again.\n"))
			}
		}
	}
}

func (fm *FileManager) signIn(conn net.Conn) (Account, error) {
	conn.Write([]byte("[Response] Username:"))
	username, err := handleInput(conn)
	if err != nil {
		return Account{}, fmt.Errorf("error handling input: %w", err)
	}
	conn.Write([]byte("[Response] Password:"))
	password, err := handleInput(conn)
	if err != nil {
		return Account{}, fmt.Errorf("error handling input: %w", err)
	}
	account, err := fm.Am.LogIn(username, password)
	if err != nil {
		return Account{}, err
	}
	return account, nil
}

func (fm *FileManager) signUp(conn net.Conn) error {
	conn.Write([]byte("[Response] Full Name:"))
	fullName, err := handleInput(conn)
	if err != nil {
		return fmt.Errorf("error handling input: %w", err)
	}
	conn.Write([]byte("[Response] Email:"))
	email, err := handleInput(conn)
	if err != nil {
		return fmt.Errorf("error handling input: %w", err)
	}
	conn.Write([]byte("[Response] Address:"))
	address, err := handleInput(conn)
	if err != nil {
		return fmt.Errorf("error handling input: %w", err)
	}
	conn.Write([]byte("[Response] Username:"))
	username, err := handleInput(conn)
	if err != nil {
		return fmt.Errorf("error handling input: %w", err)
	}
	conn.Write([]byte("[Response] Password:"))
	password, err := handleInput(conn)
	if err != nil {
		return fmt.Errorf("error handling input: %w", err)
	}
	user := User{
		Username: username,
		Password: password,
		FullName: fullName,
		Email:    email,
		Address:  address,
	}
	return fm.Am.SignUp(user)
}

func (fm FileManager) Loop() error {
	fm.Init()
	return nil
}

func handleInput(conn net.Conn) (string, error) {
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		return "", fmt.Errorf("error reading from connection: %w", err)
	}
	return string(buffer[:n]), nil
}
