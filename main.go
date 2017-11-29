package main

import (
	"fmt"
	"io"
	"net"
	"os"
)

var addr net.IP
var port string

//var listener = net.Listener

func echoServer(l net.Listener) {
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println(os.Stderr, err)
			os.Exit(1)
		}
		go func(c net.Conn) {
			io.Copy(c, c)
			c.Close()
		}(conn)

	}
}



func init() {

	if ( len(os.Args) > 2  ) {
		fmt.Fprintf(os.Stderr, "Usage: kgecho /path/to/config.json\n")
		fmt.Fprintf(os.Stderr, "Or, leave the path blank for default\n")
		fmt.Println(len(os.Args))
		os.Exit(1)
	}

	configFile := ""
	if len(os.Args) == 2 {
		configFile = os.Args[1]
	}
	loadConfigFromFile(configFile)

}

func main() {
	composeListenAddress := addr.String() + ":" + port
	fmt.Println("Composed address is: ", composeListenAddress, "\n")
	listener, err := net.Listen("tcp", composeListenAddress)
	if err != nil {
		fmt.Println(os.Stderr, err)
		os.Exit(1)
	}
	defer listener.Close()
	echoServer(listener)

	fmt.Println("The address is ", addr.String(), "\n")
	fmt.Println("exiting\n")
	os.Exit(0)

}
