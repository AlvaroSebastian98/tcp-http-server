package main

import (
	"fmt"
	"main/app"
	r "main/router"
	"net"
	"os"
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

		router := r.BuildRouter(conn)

		go func ()  {
			defer conn.Close()
			app.HandleRequest(router)
		}()
	}

}
