package server

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"math/rand"
	"net"
	"time"

	"github.com/alessiosavi/GoGPUtils/helper"
)

func handleConnection(c net.Conn, semaphore <-chan struct{}) {
	fmt.Printf("Serving %s\n", c.RemoteAddr().String())
	data, err := bufio.NewReader(c).ReadBytes(0x04)
	if err != nil {
		fmt.Println(err)
		return
	}
	data = data[:len(data)-1]
	log.Println("Data received!")
	log.Println(helper.MarshalIndent(string(data)))
	c.Write([]byte(fmt.Sprintf("%d", helper.RandomInt(0, math.MaxInt))))
	c.Write([]byte{0x04})
	c.Close()
	time.Sleep(300 * time.Millisecond)
	<-semaphore
}

func Server(port string) {
	l, err := net.Listen("tcp4", ":"+port)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer l.Close()
	rand.NewSource(time.Now().Unix())
	semaphore := make(chan struct{}, 10000)
	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		semaphore <- struct{}{}
		go handleConnection(c, semaphore)
	}
}

func Client(port int) {

	tcpAddr, err := net.ResolveTCPAddr("tcp4", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		panic(err)
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	conn.Write([]byte("test data"))
	if _, err = conn.Write([]byte{0x04}); err != nil {
		panic(err)
	}

	if _, err = bufio.NewReader(conn).ReadString(0x04); err != nil {
		panic(err)
	}
	// log.Println(data[:len(data)-1])

}
