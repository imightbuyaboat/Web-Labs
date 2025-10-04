package main

import (
	"log"
	"net"
	"os"
	"strings"
	"time"
	"unicode"
)

const (
	BUF_SIZE = 1024
)

var (
	logFileLogger log.Logger
)

func handleConn(conn net.Conn) {
	defer func() {
		conn.Close()
		logFileLogger.Printf("%v: Connection from %s closed", time.Now(), conn.RemoteAddr().String())
	}()

	conn.SetDeadline(time.Now().Add(10 * time.Second))

	buf := make([]byte, BUF_SIZE)

	for {
		n, err := conn.Read(buf)
		if err != nil {
			return
		}
		logFileLogger.Printf("%v: Recieve new message from %s: %s", time.Now(), conn.RemoteAddr().String(), string(buf[:n]))

		inverse := inverseString(string(buf[:n]))
		time.Sleep(1 * time.Second)

		resp := append([]byte(inverse), []byte(". Сервер написан Вовком И.Б. М3О-425Бк-22")...)
		_, err = conn.Write(resp)
		if err != nil {
			return
		}
		logFileLogger.Printf("%v: Sent message to %s: %s", time.Now(), conn.RemoteAddr().String(), string(resp))
	}
}

func inverseString(s string) string {
	words := strings.Split(s, " ")
	for i, w := range words {
		runes := []rune(w)

		for l, r := 0, len(runes)-1; l < r; l, r = l+1, r-1 {
			runes[l], runes[r] = runes[r], runes[l]
		}

		for j, r := range runes {
			if j == 0 {
				runes[j] = unicode.ToUpper(r)
			} else {
				runes[j] = unicode.ToLower(r)
			}
		}

		words[i] = string(runes)
	}
	return strings.Join(words, " ")
}

func main() {
	file, err := os.OpenFile("log.txt", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	logFileLogger.SetOutput(file)

	listener, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		log.Fatal(err)
	}
	logFileLogger.Printf("%s: Server starting on :8080", time.Now())

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		logFileLogger.Printf("%v: Accepted new client: %s", time.Now(), conn.RemoteAddr().String())

		go handleConn(conn)
	}
}
