package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	servAddress, err := net.ResolveUDPAddr("udp", ":42069")
	if err != nil {
		fmt.Println("Error resolving UDP address:", err)
		return
	}
	conn, err := net.DialUDP("udp", nil, servAddress)
	if err != nil {
		fmt.Println("Error dialing UDP address:", err)
		return
	}
	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)

	//fmt.Print("Enter message: ")

	for {
		fmt.Println(">")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error receiving response:", err)
			return
		}

		_, err = conn.Write([]byte(input))
		if err != nil {
			fmt.Println("Error sending UDP message:", err)
			return
		}
		fmt.Println("Message sent:", input)
	}
}
