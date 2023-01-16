package main

import "github.com/byzk-org/bypt/cmd"

func main() {
	defer func() { recover() }()
	cmd.Execute()
	//socket.GetClientConn()
}
