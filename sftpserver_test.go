package sftpserver

import (
	"log"
	"testing"

	"github.com/pkg/sftp"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/ssh"
)

const (
	addr                          = "localhost:2022"
	username                      = "sculley"
	password                      = "password"
	failedServerStartErrorMessage = "Failed to start test SFTP server:"
)

func TestSuccessfulConnection(t *testing.T) {
	server := New(addr, username, password)
	if err := server.Start(); err != nil {
		log.Fatal("Failed to start test SFTP server:", err)
	}
	defer server.Stop()

	conn, err := ssh.Dial("tcp", addr, &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	})

	require.NoError(t, err)
	require.NotNil(t, conn)

	client, err := sftp.NewClient(conn)

	require.NoError(t, err)
	require.NotNil(t, client)
}

func TestFailedConnection(t *testing.T) {
	server := New(addr, username, password)
	if err := server.Start(); err != nil {
		log.Fatal(failedServerStartErrorMessage, err)
	}
	defer server.Stop()

	conn, err := ssh.Dial("tcp", addr, &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password("wrongpassword"),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	})

	require.Error(t, err)
	require.Nil(t, conn)
}
