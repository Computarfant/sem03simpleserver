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
			var x string

			go func(c net.Conn) {
				for {
					buf := make([]byte, 1024)
					n, err := conn.Read(buf)
					if err != nil {
						if err != io.EOF {
							log.Println(err)
						}
						return // from for loop
					}

					dekryptertMelding := mycrypt.Krypter([]rune(string(buf[:n])), mycrypt.ALF_SEM03, -4)
					log.Println("Dekrypter melding: ", string(dekryptertMelding))

					if strings.HasPrefix(string(dekryptertMelding), "Kjevik") {
						fields := strings.Split(string(dekryptertMelding), ";")
						if len(fields) >= 4 {
							celsius, err := strconv.ParseFloat(fields[3], 64)
							if err != nil {
								log.Println(err)
								continue
							}
							fahrenheit := conv.CelsiusToFahrenheit(celsius)
							x = fmt.Sprintf("%s;%s;%s;%.1f\n", fields[0], fields[1], fields[2], fahrenheit)
							_, err = conn.Write([]byte(x))
							if err != nil {
								log.Println(err)
								return // from for loop
							}
						} else {
							log.Println("Invalid input:", string(dekryptertMelding))
							continue
						}
					} else {
						msg := string(dekryptertMelding)
						switch msg {
						case "ping":
							svar := mycrypt.Krypter([]rune("pong"), mycrypt.ALF_SEM03, 4)
							_, err = conn.Write([]byte(string(svar)))
						case x:
							svar := mycrypt.Krypter([]rune(x), mycrypt.ALF_SEM03, 4)
							_, err = conn.Write([]byte(string(svar)))
						default:
							svar := mycrypt.Krypter([]rune(msg), mycrypt.ALF_SEM03, 4)
							_, err = conn.Write([]byte(string(svar)))
						}
						if err != nil {
							log.Println(err)
							return // from for loop
						}
					}
				}
			}(conn)
		}
	}()
	wg.Wait()
}
