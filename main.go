package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/kr/pty"
	daemon "github.com/sevlyar/go-daemon"
)

func tryPtyOpen() error {
	sigch := make(chan os.Signal, 1)
	defer close(sigch)
	signal.Notify(sigch, syscall.SIGHUP)
	defer signal.Reset(syscall.SIGHUP)

	p, t, err := pty.Open()
	if err != nil {
		return fmt.Errorf("failed to open pty: %s", err)
	}
	if err := p.Close(); err != nil {
		return fmt.Errorf("failed to close pty: %s", err)
	}
	if err := t.Close(); err != nil {
		return fmt.Errorf("failed to close tty: %s", err)
	}
	return nil
}

func main() {
	cntxt := &daemon.Context{
		PidFileName: "pid",
		PidFilePerm: 0644,
		LogFileName: "log",
		LogFilePerm: 0640,
		WorkDir:     "./",
		Umask:       027,
	}

	d, err := cntxt.Reborn()
	if err != nil {
		log.Fatal("Unable to run: ", err)
	}
	if d != nil {
		return
	}
	defer cntxt.Release()

	log.Print("daemon started")

	if err := tryPtyOpen(); err != nil {
		log.Print(err)
		return
	}

	c := exec.Command("ls")
	p, err := pty.Start(c)
	if err != nil {
		log.Print(err)
	}
	b, err := ioutil.ReadAll(p)
	if err != nil {
		log.Print(err)
	}
	log.Print(string(b))
}
