package types

import "github.com/ibraheemacamara/merkletree"

type ServerResponse struct {
	File  []byte
	Proof merkletree.Proof
}

var (
	BUFFER_SIZE     int64 = 1024
	CMD_CLIENT_SEND       = "send"
	CMD_CLIENT_GET        = "get"
)
