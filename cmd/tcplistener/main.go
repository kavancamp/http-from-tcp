package main

import (
	"fmt"
	"io"
	"net"
	"os"
)

func getLinesChannel(r io.ReadCloser) <-chan string {
	out := make(chan string)

	go func() {
		defer r.Close()
		defer close(out)

		const chunkSize = 8
		buf := make([]byte, chunkSize)
		acc := make([]byte, 0, 256)

		emit := func() {
			if len(acc) == 0 {
				return
			}
			// Trim trailing '\r' for CRLF
			if acc[len(acc)-1] == '\r' {
				out <- string(acc[:len(acc)-1])
			} else {
				out <- string(acc)
			}
			acc = acc[:0]
		}

		for {
			n, err := r.Read(buf)
			if n > 0 {
				chunk := buf[:n]
				start := 0
				for i, b := range chunk {
					if b == '\n' {
						acc = append(acc, chunk[start:i]...)
						emit()
						start = i + 1
					}
				}
				if start < len(chunk) {
					acc = append(acc, chunk[start:]...)
				}
			}
			if err != nil {
				if err == io.EOF {
					if len(acc) > 0 {
						emit()
					}
				}
				return
			}
		}
	}()

	return out
}

func main() {
	ln, err := net.Listen("tcp", ":42069")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error getting listener:", err)
		return
	}
	defer ln.Close()

	fmt.Fprintln(os.Stderr, "Listening on :42069")

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Accept error:", err)
			continue
		}
		addr := conn.RemoteAddr().String()
		fmt.Fprintln(os.Stderr, "New connection accepted:", addr)

		go func(c net.Conn, who string) {
			lines := getLinesChannel(c)
			for line := range lines {
				// Print to stdout, exactly one line per message.
				fmt.Println(line)
			}
			fmt.Fprintln(os.Stderr, "Connection closed:", who)
		}(conn, addr)
	}
}
