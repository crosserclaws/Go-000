package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

const (
	listenAddr = ":8080"
)

var (
	programControlCh <-chan void
	wg               sync.WaitGroup
)

type void struct{}

func main() {
	log.Println("Server starts.")

	signalCh := RegisterSignals(syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	programCh := make(chan void)
	programControlCh = programCh
	go listenAndServe()

	callBack := func() { close(programCh) }
	waitForSignals(signalCh, callBack)

	log.Println("Server waits for closing connections.")
	wg.Wait()

	log.Println("Server ends.")
}

func listenAndServe() {
	ln, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	var id int

LS:
	for {
		select {
		case <-programControlCh:
			break LS
		default:
			conn, err := ln.Accept()
			if err != nil {
				log.Printf("Failed to accept: %v", err)
				continue
			}
			id++
			go handleConnection(conn, id)
		}
	}
}

func handleConnection(conn net.Conn, id int) {
	controlCh := make(chan void, 1)
	dataCh := make(chan string, 1)
	defer conn.Close()
	go handleRead(conn, id, controlCh, dataCh)
	handleWrite(conn, id, controlCh, dataCh)
}

func handleRead(conn net.Conn, id int, controlCh chan void, dataCh chan<- string) {
	reader := fmt.Sprintf("Reader[%v]", id)
	log.Println(reader, "starts.")
	wg.Add(1)
	defer wg.Done()
	defer log.Println(reader, "ends.")
	rd := bufio.NewReader(conn)

	for {
		select {
		case <-programControlCh:
			log.Println(reader, "gets the program is going to shutdown.")
			return
		case <-controlCh:
			log.Println(reader, "gets the connection is going to close.")
			return
		default:
			msg, err := readMessage(rd)
			if err != nil {
				log.Printf("%s fails to read: %v", reader, err)
				notify(controlCh)
				return
			}

			log.Printf("%s receives a msg: %s", reader, msg)
			dataCh <- msg
		}
	}
}

func readMessage(rd *bufio.Reader) (string, error) {
	line, _, err := rd.ReadLine()
	if err != nil {
		return "", err
	}
	msg := string(line)
	return msg, nil
}

func handleWrite(conn net.Conn, id int, controlCh chan void, dataCh <-chan string) {
	writer := fmt.Sprintf("Writer[%v]", id)
	log.Println(writer, "starts.")
	wg.Add(1)
	defer wg.Done()
	defer log.Println(writer, "ends.")
	wr := bufio.NewWriter(conn)

	for {
		select {
		case <-programControlCh:
			log.Println(writer, "gets the program is going to shutdown.")
			writeMessage(wr, "[Notification] Server is going to shutdown.")
			return
		case <-controlCh:
			log.Println(writer, "gets the connection is going to close.")
			return
		case msg := <-dataCh:
			log.Printf("%s echos a response: %v", writer, msg)
			if err := writeMessage(wr, msg); err != nil {
				log.Printf("%s fails to write: %v", writer, err)
				notify(controlCh)
				return
			}
		}
	}
}

func writeMessage(wr *bufio.Writer, msg string) error {
	if _, err := wr.WriteString(msg); err != nil {
		return err
	}
	if err := wr.Flush(); err != nil {
		return err
	}
	return nil
}

func notify(ch chan<- void) {
	ch <- void{}
}

// RegisterSignals is a utility function registers given signals.
func RegisterSignals(sig ...os.Signal) <-chan os.Signal {
	log.Println("Register signals:", sig)

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, sig...)
	return ch
}

func waitForSignals(ch <-chan os.Signal, callBack func()) {
L:
	for {
		select {
		case s := <-ch:
			switch s {
			case syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
				log.Println("Ask server to shutdown when capturing a registered signal:", s)
				callBack()
				break L
			default:
				log.Println("Capture an unknown signal:", s)
			}
		}
	}

}
