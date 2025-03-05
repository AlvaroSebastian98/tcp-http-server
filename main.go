package main

import (
	"fmt"
	"main/app"
	"main/utils"
	"net"
	"os"
	"strings"
	// "time"
)

func main() {

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		buff := make([]byte, 1024)
		n, _ := conn.Read(buff)
		req := utils.ParseRequest(strings.Split(string(buff[:n]), "\r\n"))
		res := utils.BuildResponse(conn)
		router := utils.BuildRouter(req)

		go func ()  {
			defer conn.Close()
			app.HandleRequest(router, req, res)
		}()
	}

}
