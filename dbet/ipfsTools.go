package main

import (
	"os"
	"encoding/json"
	"io/ioutil"
	"os/exec"
	"strings"
	"fmt"
	"errors"
)

var ipfsHashes map[string]ipfsHash

type ipfsHash struct {
	Data     string `json:"d"`
	KeyMov   string `json:"k"`
	Combined string `json:"c"`
}

func init() {
	file, err := os.Open("./ipfsHashes.json")
	if os.IsNotExist(err) {
		ipfsHashes = make(map[string]ipfsHash, 100)
		return
	}

	if err != nil {
		panic(err)
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&ipfsHashes)
	if err != nil {
		panic(err)
	}
}

func saveIpfsHashes() error {
	b, _ := json.MarshalIndent(ipfsHashes, "", " ")
	return ioutil.WriteFile("ipfsHashes.json", b, 0644)
}

func ipfsPinPath(path string, name string) (string, error) {
	fmt.Println("Pinning " + name)
	bin := "ipfs"
	args := []string{"add", "-r", "-p=false", "--nocopy", path}

	ial := exec.Command(bin, args...)
	ial.Env = append(ial.Env, "IPFS_PATH=/services/tomography/.ipfs")
	out, err := ial.CombinedOutput()
	if err != nil {
		return string(out), err
	}

	lines := strings.Split(string(out), "\n")
	last := lines[len(lines)-2]
	words := strings.Split(last, " ")

	if words[0] == "added" && words[2] == name {
		fmt.Println("Pinned. " + words[1])
		return words[1], nil
	} else {
		fmt.Println(words)
		return string(out), errors.New("ipfs hash not found")
	}
}

func ipfsAddLink(dirHash string, name string, link string) (string, error) {
	bin := "ipfs"
	args := []string{"object", "patch", "add-link", dirHash, name, link}

	ial := exec.Command(bin, args...)
	ial.Env = append(ial.Env, "IPFS_PATH=/services/tomography/.ipfs")
	out, err := ial.CombinedOutput()
	if err != nil {
		return string(out), err
	}

	return strings.TrimSpace(string(out)), nil
}