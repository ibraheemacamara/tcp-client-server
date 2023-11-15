package main

import (
	"archive/zip"
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"strings"

	"github.com/ibraheemacamara/merkletree"
	"github.com/ibraheemacamara/tcp-client-server/types"
	"github.com/ibraheemacamara/tcp-client-server/utils"
)

var defaultDataDir = filepath.Join("$HOME", ".client")

func main() {

	//start client
	// make client --server-host="localhost" --server-port=8080
	//default localhost:8080

	log.Printf("starting client")

	host := flag.String("server-host", "localhost", "server hostname, if empty localhost will be used")
	port := flag.Int("server-port", 8080, "server hostname, if empty 8080 will be used")
	flag.Parse()

	address := fmt.Sprintf("%v:%v", *host, *port)
	fmt.Println(address)
	connection, err := net.Dial("tcp", address)
	if err != nil {
		log.Fatalf("failed to connect cllient to server: %v", err)
	}

	//read command from user
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Please enter 'get <filename>' or 'send <dirname>' to transfer files to the server\n\n")
	input, _ := reader.ReadString('\n')
	cmds := strings.Split(input, " ")

	if cmds[0] == types.CMD_CLIENT_GET {
		fileBytes, proof, err := getFileAndProofFromServer(cmds[1], connection)
		fmt.Println(err)
		if err != nil {
			log.Printf("failed to get file: %v", err.Error())
			return
		}

		if isValid(proof, fileBytes) {
			log.Println("File is correct !!!")
		} else {
			log.Println("File is not correct !!!")
		}

	} else if cmds[0] == types.CMD_CLIENT_SEND {
		sendFileToServer(cmds[1], connection)
	} else {
		fmt.Println("Bad Command")
	}
}

func getFileAndProofFromServer(fileName string, con net.Conn) ([]byte, merkletree.Proof, error) {
	defer con.Close()

	return nil, merkletree.Proof{}, nil
}

func sendFileToServer(dirname string, con net.Conn) {
	defer con.Close()
	filesDate, merkleRootHash, err := readSetOfFilesFromDir(dirname)
	if err != nil {
		log.Printf("failed to send files to server: %v \n", err.Error())
		return
	}
	//save merkle root on disk
	err = saveMerkleRootOnDisk(merkleRootHash)
	if err != nil {
		log.Printf("failed to send files to server: %v \n", err.Error())
		return
	}

	//send files data to server
	con.Write([]byte("send "))
	con.Write(filesDate)
}

// This function read files from directory
// compressa data
// generate merkle root hash
// return zip data, merkle root hash, error
func readSetOfFilesFromDir(dirpath string) ([]byte, []byte, error) {
	//data := map[string][]byte{}

	var files [][]byte

	dir, err := os.ReadDir(dirpath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read files from directory: %v", err.Error())
	}

	var buf bytes.Buffer
	zipWriter := zip.NewWriter(&buf)
	defer zipWriter.Close()

	for _, file := range dir {
		if file.IsDir() { //TODO allow read sub dir
			return nil, nil, fmt.Errorf("sub dictory not allowed")
		}
		fileData, err := os.ReadFile(file.Name())
		if err != nil {
			return nil, nil, fmt.Errorf("failed to read data from file: %v", err.Error())
		}
		files = append(files, fileData)

		w, err := zipWriter.Create(file.Name())
		if err != nil {
			return nil, nil, fmt.Errorf("failed to create zip file: %v", err.Error())
		}
		fileReader, err := os.Open(file.Name())
		if err != nil {
			return nil, nil, fmt.Errorf("failed to open file %v: %v", file.Name(), err.Error())
		}
		_, err = io.Copy(w, fileReader)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to copy file content to archive: %v", err.Error())
		}
	}

	tree, err := utils.ComputeMerkleTree(files)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create merkle tree: %v", err.Error())
	}

	return buf.Bytes(), tree.RootNode.Hash, nil
}

func isValid(proof merkletree.Proof, data []byte) bool {
	rootHash, err := readMerkleRootHashFromDisk()
	if err != nil {
		return false
	}
	return merkletree.VerifyProof(rootHash, data, proof.Path, proof.Idxs)
}

func saveMerkleRootOnDisk(data []byte) error {
	path := filepath.Join(defaultDataDir, "data.txt")
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to save data on disk: %v", err.Error())
	}

	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		return fmt.Errorf("failed to save data on disk: %v", err.Error())
	}

	return nil
}

func readMerkleRootHashFromDisk() ([]byte, error) {
	path := filepath.Join(defaultDataDir, "data.txt")
	data, err := os.ReadFile(path)
	if err != nil {
		return []byte{}, fmt.Errorf("failed to read data from disk: %v", err.Error())
	}

	return data, nil
}
