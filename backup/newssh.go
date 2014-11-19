package backup

import (
	"fmt"
	"io"

	"golang.org/x/crypto/ssh"
)

func Connect(userName string, password string, host string, port int, verbose bool, errPipe io.Writer) (*ssh.Session, error) {
	// Define the Client Config as :
	clientConfig := &ssh.ClientConfig{
		User: userName,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
	}

	target := fmt.Sprintf("%s:%d", host, port)
	client, err := ssh.Dial("tcp", target, clientConfig)
	if err != nil {
		if verbose {
			fmt.Fprintln(errPipe, "Failed to dial: "+err.Error())
		}
		return nil, err
	}

	session, err := client.NewSession()
	if err != nil {
		if verbose {
			fmt.Fprintln(errPipe, "Failed to create session: "+err.Error())
		}
	}

	return session, err

}
