package main

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)
var (
	Errors 	int32
	Success int32
)
const (
	SSHTimeout = 	5
	SSHPort    = 	"22"
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
	s, err := ssh.Dial("tcp", host+":"+SSHPort, config)
	if err != nil {
		fmt.Printf("[-] Failed to connect to %s\n", host)
		atomic.AddInt32(&Errors, 1)
		return
	}
	defer func(s *ssh.Client) {
		err := s.Close()
		if err != nil {
			fmt.Println("Close Err!")
		}
	}(s)

	session, err := s.NewSession()
	if err != nil {
		fmt.Printf("\033[31m[-] Failed to open session to %s\n\033[39m", host)
	}
	defer session.Close()

	err = session.Run(cmd)
	if err != nil {
		fmt.Printf("\033[31m[-] Failed to run command on %s\n\033[39m", host)
	}
	ts := time.Since(t)
	fmt.Printf("\033[32m[+] Command sent to %s within %.2f\n\033[39m", host, float64(ts.Milliseconds())/1000.0)
	atomic.AddInt32(&Success, 1)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("%s <serverlist> <command>", os.Args[0])
		os.Exit(0)
	}
	File := os.Args[1]


	f, err := os.Open(File)
	if err != nil {
		fmt.Println("Could not find File.")
		os.Exit(0)
	}

	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			fmt.Println("Close Err!")
		}
	}(f)
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
		if x == "" {
			continue
		}
		x = strings.TrimSpace(x)
		Args := strings.Split(x, ":")
		if len(Args) < 2 {
			continue
		}
		wg.Add(1)
		go SendSSH(&wg, Args[0], Args[1], Args[2], cmd)
	}
	wg.Wait()
	ts := time.Since(Start)
	fmt.Printf("\033[32m\nProccess took %.2fs to finish!\nSuccess: %d | Errors: %d\033[39m", float64(ts.Milliseconds())/1000.0, Success, Errors)
	os.Exit(0)
}
