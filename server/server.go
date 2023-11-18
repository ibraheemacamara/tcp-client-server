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
	"github.com/ibraheemacamara/tcp-client-server/utils"
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

	//log.Printf("received command from client: %v", string(buffer))

	cleanedBuffer := bytes.Trim(buffer, "\x00")
	cmdStrings := strings.TrimSpace(string(cleanedBuffer))
	cmdList := strings.Split(cmdStrings, " ")

	if cmdList[0] == types.CMD_CLIENT_SEND {
		log.Println("received save file command")
		saveClientFiles(cmdList[1], buffer, con)
	} else if cmdList[0] == types.CMD_CLIENT_GET {
		log.Println("received get file command")
		sendFileToClient(cmdList[1], con)
	} else {
		log.Println("command not valid")
		return
	}
}

func saveClientFiles(dirname string, data []byte, con net.Conn) {
	defer con.Close()

	clientId := con.RemoteAddr()
	key := []byte(fmt.Sprintf("%v%v", clientId, dirname))
	err := dbClient.Put(key, data)
	if err != nil {
		log.Printf("failed to save set of files %v of client %v into db: %v", dirname, clientId, err.Error())
		return
	}

	log.Printf("files successfuly saved for client: %v", clientId)
}

func sendFileToClient(fileId string, con net.Conn) {
	log.Printf("get file %v from db", fileId)
	clientId := con.RemoteAddr()
	key := []byte(fmt.Sprintf("%v%v", clientId, fileId))
	dbClient.Put(key, []byte{})
	data, err := dbClient.Get(key)
	if err != nil {
		log.Println(err)
		return
	}

	unzipData, err := utils.UngzipData(data)

	fmt.Printf("UnZip Dataaaa: %v", string(unzipData))

}
