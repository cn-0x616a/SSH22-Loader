# Golang-SSHLoader
SSHLoader written in golang and made for raw speed. Can be used as command executer for multiple servers at once or is also used for Botnet "SSH Loading".


[![Go Reference](https://pkg.go.dev/badge/golang.org/x/image.svg)](golang.org/x/crypto/ssh)


**Configuration**



```go
var (
	SSHTimeout = 	5 // Timeout from SSH connection in seconds.
	SSHPort    = 	22 // SSH Port
)
```

Also works with private keys

```go
pem, err := os.ReadFile("key.pem")
	if err != nil {
		fmt.Println(err)
	}
	signer, err := ssh.ParsePrivateKey(pem)
	if err != nil {
		fmt.Println(err)
	}

	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		Timeout: time.Second*time.Duration(SSHTimeout),
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
```

**Setup**

Firstly we are going to clone the repository.
```
apt install git
git clone https://github.com/Yariya/Golang-SSHLoader/
```

Then we cd into the repository
```
cd Golang-SSHLoader
```

Downloading dependencies
```linux
apt install golang-go
go get golang.org/x/crypto/ssh # May not be needed
```

Build
```
go build loader.go
```

# Syntax


Server(s).txt Syntax

```
1.1.1.1:root:Yariya
...
```

Execute it
```
./loader servers.txt COMMAND
```
