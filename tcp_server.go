package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
	"sync"

	"github.com/Computarfant/funtemps/conv"
	"github.com/Computarfant/is105sem03/mycrypt"
)

func main() {

	var wg sync.WaitGroup

	server, err := net.Listen("tcp", "172.17.0.2:8000")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("bundet til %s", server.Addr().String())
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			log.Println("før server.Accept() kallet")
			conn, err := server.Accept()
			if err != nil {
				return
			}
			go func(c net.Conn) {
				defer c.Close()
				for {
					buf := make([]byte, 1024)
					n, err := c.Read(buf)
					if err != nil {
						if err != io.EOF {
							log.Println(err)
						}
						return // from for loop
					}
					dekryptertMelding := mycrypt.Krypter([]rune(string(buf[:n])), mycrypt.ALF_SEM03, -4)
					log.Println("Dekrypter melding: ", string(dekryptertMelding))
					switch msg := string(dekryptertMelding); msg {
					case "ping":
						svar := mycrypt.Krypter([]rune("pong"), mycrypt.ALF_SEM03, 4)
						_, err = c.Write([]byte(svar))
					case strings.HasPrefix(msg, "Kjevik"):
						fields := strings.Split(msg, ";")
						if len(fields) != 4 || fields[3] == "" {
							continue
						}
						celsius, err := strconv.ParseFloat(fields[3], 64)
						if err != nil {
							return
						}
						fahrenheit := conv.CelsiusToFahrenheit(celsius)
						svar := mycrypt.Krypter([]rune(fmt.Sprintf("%s;%s;%s;%.1f\n", fields[0], fields[1], fields[2], fahrenheit)), mycrypt.ALF_SEM03, 4)
						_, err = c.Write([]byte(svar))
					default:
						svar := mycrypt.Krypter([]rune(msg), mycrypt.ALF_SEM03, 4)
						_, err = c.Write([]byte(svar))
					}
					if err != nil {
						if err != io.EOF {
							log.Println(err)
						}
						return // from for loop
					}
				}
			}(conn)
		}
	}()
	wg.Wait()
}
