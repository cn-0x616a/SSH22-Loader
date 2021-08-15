package main

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"sync"
	"time"
)


var (
	SSHTimeout = 	5
	SSHPort    = 	22
)

var (
	Errors 	int
	Success int
)

func ClearArg(arg []string) string {
	var String string
	for _, x := range arg {
		String+=x
		String+=" "

	}
	return String
}

func SendSSH(wg *sync.WaitGroup, host string, user string, pw string, cmd string) {
	defer wg.Done()




	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(pw),
		},
		Timeout: time.Second*time.Duration(SSHTimeout),
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	t := time.Now()
	s, err := ssh.Dial("tcp", host+":"+string(SSHPort), config)
	if err != nil {
		fmt.Printf("[-] Failed to connect to %s\n", host)
		Errors++
		return
	}
	defer s.Close()

	session, err := s.NewSession()
	if err != nil {
		fmt.Println(err)
	}
	defer session.Close()

	StdOut, err := session.StdoutPipe()
	if err != nil {
		fmt.Println(err)
	}
	go io.Copy(os.Stdout, StdOut)
	SessStr, err := session.StderrPipe()
	if err != nil {
		fmt.Println(err)
	}
	go io.Copy(os.Stderr, SessStr)

	err = session.Run(cmd)
	if err != nil {
		fmt.Println(err)
	}
	Success++
	fmt.Printf("[+] Command sent to %s within %s\n", host, time.Since(t))
}

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Syntax Error")
		os.Exit(0)
	}
	File := os.Args[1]


	f, err := os.Open(File)
	if err != nil {
		fmt.Println("Could not find File.")
		os.Exit(0)
	}
	defer f.Close()


	Content, err := ioutil.ReadAll(f)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	Servers := strings.Split(string(Content), "\n")

	Command := os.Args[2:]
	cmd := ClearArg(Command)

	var wg sync.WaitGroup

	Start := time.Now()

	for _, x := range Servers {
		x = strings.TrimSpace(x)
		Args := strings.Split(x, ":")
		wg.Add(1)
		go SendSSH(&wg, Args[0], Args[1], Args[2], cmd)
	}
	fmt.Println("Started all Threads!")
	wg.Wait()
	ts := time.Since(Start)
	fmt.Printf("\n\nProccess took %s to finish!\nSuccess: %d | Errors: %d", ts, Success, Errors)
	os.Exit(0)
}
