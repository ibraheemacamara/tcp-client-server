package utils

import (
	"fmt"

	"github.com/ibraheemacamara/merkletree"
)

func ComputeMerkleTree(data [][]byte) (*merkletree.MerkleTreee, error) {
	tree, err := merkletree.NewMerkleTree(data)
	if err != nil {
		return nil, fmt.Errorf("failed to compute merkle tree: %v", err)
	}

	return tree, nil
}

func GetMerkleProof(tree *merkletree.MerkleTreee, data []byte) (merkletree.Proof, error) {
	proof, err := tree.GetProof(data)
	if err != nil {
		return merkletree.Proof{}, fmt.Errorf("failed to get proof: %v", err.Error())
	}

	return proof, nil
}
