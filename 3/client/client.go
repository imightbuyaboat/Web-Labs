package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

const (
	CONFIG_FILE = "config.json"
	LOG_FILE    = "log.txt"
	BUF_SIZE    = 1024
)

func readConfig() (string, error) {
	var config struct {
		ServerAddress string `json:"server_address"`
	}

	file, err := os.OpenFile(CONFIG_FILE, os.O_RDONLY, 0777)
	if err != nil {
		return "", err
	}
	defer file.Close()

	if err = json.NewDecoder(file).Decode(&config); err != nil {
		return "", err
	}

	if config.ServerAddress == "" {
		return "", fmt.Errorf("server_address is empty")
	}

	return config.ServerAddress, nil
}

func main() {
	serverAddress, err := readConfig()
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.OpenFile("log.txt", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	logFileLogger := log.Logger{}
	logFileLogger.SetOutput(file)

	conn, err := net.Dial("tcp", serverAddress)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	logFileLogger.Printf("%v: Connect to server %s", time.Now(), serverAddress)
	time.Sleep(2 * time.Second)

	msg := "Вовк Илья Богданович"
	_, err = conn.Write([]byte(msg))
	if err != nil {
		log.Fatal(err)
	}
	logFileLogger.Printf("%v: Send message: %s", time.Now(), msg)

	buf := make([]byte, BUF_SIZE)
	n, err := conn.Read(buf)
	if err != nil {
		log.Fatal(err)
	}
	logFileLogger.Printf("%v: Recieve message: %s", time.Now(), string(buf[:n]))

	conn.Read(buf)
}
