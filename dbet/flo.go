package main

import (
	"github.com/bitspill/flojson"
	"errors"
	"fmt"
	"time"
)

var (
	id     int64
	user   string
	pass   string
	server string
)

func init() {
	id = 0 // id is static at 0, for "proper" json-rpc increment with each call
	user = config.FloConfiguration.RpcUser
	pass = config.FloConfiguration.RpcPass
	server = config.FloConfiguration.RpcAddress
}

func signMessage(address, message string) (string, error) {
	cmd, err := flojson.NewSignMessageCmd(id, address, message)
	if err != nil {
		return "", err
	}

	reply, err := sendRPC(cmd)
	if err != nil {
		return "", err
	}

	if signature, ok := reply.Result.(string); ok {
		return signature, nil
	}

	return "", errors.New("unexpected rpc error")
}

func sendToAddress(address string, amount float64, floData string) (string, error) {
	satoshi := int64(1e8 * amount)
	cmd, err := flojson.NewSendToAddressCmd(id, address, satoshi, "", "", floData)
	if err != nil {
		return "", err
	}

	reply, err := sendRPC(cmd)
	if err != nil {
		return "", err
	}
	if reply.Error != nil {
		return "", reply.Error
	}
	return reply.Result.(string), nil
}

func sendRPC(cmd flojson.Cmd) (flojson.Reply, error) {
	t := 0
	for true {
		reply, err := flojson.RpcSend(user, pass, server, cmd)
		if err != nil {
			fmt.Println(reply, err)
			return reply, err
		}
		if reply.Error != nil {
			if reply.Error.Code == -6 && reply.Error.Message == "Insufficient funds" {
				if t > 20 {
					fmt.Println("It's been 10 minutes, perhaps you're really out of funds")
					return reply, reply.Error
				}
				t++
				fmt.Println("Sleeping 30s for a block to re-confirm balance")
				time.Sleep(30 * time.Second)
				continue
			}
			return reply, reply.Error
		}
		return reply, nil
	}
	panic("the above loop didn't return, something terrible has gone wrong")
}
