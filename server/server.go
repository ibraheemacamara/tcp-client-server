package main

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"github.com/ibraheemacamara/tcp-client-server/db"
	"github.com/ibraheemacamara/tcp-client-server/types"
)

var dbClient = &db.DbClient{}

func main() {
	log.Println("Starting server...")

	//Init data base
	var err error
	dbClient, err = db.InitDblient()
	if err != nil {
		log.Fatalf("failed to start db client: %v", err)
	}

	host := os.Getenv("SERVER_HOST")
	port := os.Getenv("SERVER_PORT")

	address := fmt.Sprintf("%v:%v", host, port)
	server, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("failed to start server: %v", err)
	}

	defer server.Close()

	log.Printf("server is serving at: %v", server.Addr())

	for {
		con, err := server.Accept()
		if err != nil {
			log.Printf("an error occured with the connection: %v \n", err.Error())
			return
		}

		log.Printf("connection established with %v\n", con.RemoteAddr())
		//handle client connection in separate thread
		go handleClientConnection(con)
	}
}

func handleClientConnection(con net.Conn) {
	defer con.Close()

	buffer := make([]byte, types.BUFFER_SIZE)

	_, err := con.Read(buffer)
	if err != nil {
		log.Printf("failed to read connection data: %v", err.Error())
		return
	}

	log.Printf("received command from client: %v", string(buffer))

	cleanedBuffer := bytes.Trim(buffer, "\x00")
	cmdStrings := strings.TrimSpace(string(cleanedBuffer))
	cmdList := strings.Split(cmdStrings, " ")

	if cmdList[0] == types.CMD_CLIENT_SEND {
		log.Println("received save file command")
		saveFile(cmdList[1], []byte(cmdList[2]), con)
	} else if cmdList[0] == types.CMD_CLIENT_GET {
		log.Println("received get file command")
		getFile(cmdList[1], con)
	} else {
		log.Println("command not valid")
		return
	}
}

func saveFile(dirname string, data []byte, con net.Conn) {
	defer con.Close()

	clientId := con.RemoteAddr()
	key := []byte(fmt.Sprintf("%v%v", clientId, dirname))
	err := dbClient.Put(key, data)
	if err != nil {
		log.Printf("failed to save set of files %v of client %v into db: %v", dirname, clientId, err.Error())
		return
	}
}

func getFile(fileId string, con net.Conn) (data []byte, proof [][]byte, idxs int) {
	clientId := con.RemoteAddr()
	key := []byte(fmt.Sprintf("%v%v", clientId, fileId))
	dbClient.Put(key, []byte{}) //TODO save data
	return nil, nil, 0
}
