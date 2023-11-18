package utils

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"sync"

	"github.com/ibraheemacamara/merkletree"
)

var gzipReaderPool sync.Pool

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

func UngzipData(data []byte) ([]byte, error) {
	b := bytes.NewReader(data)
	gr, err := getGzipReader(b)
	if err != nil {
		return nil, err
	}
	defer putGzipReader(gr)
	var buf bytes.Buffer
	_, err = io.Copy(&buf, gr)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func putGzipReader(gr *gzip.Reader) {
	gzipReaderPool.Put(gr)
}

func getGzipReader(r io.Reader) (*gzip.Reader, error) {
	if v := gzipReaderPool.Get(); v != nil {
		gr := v.(*gzip.Reader)
		gr.Reset(r)
		return gr, nil
	}
	return gzip.NewReader(r)
}
