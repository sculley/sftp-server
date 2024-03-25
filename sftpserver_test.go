package sftpserver

import (
	"crypto/rand"
	"crypto/rsa"
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
	cfg := Config{
		Addr:     addr,
		Username: username,
		Password: password,
	}

	server := New(cfg)
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

func TestSuccessfulConnectionWithAuthorizedKey(t *testing.T) {
	// Generate test RSA keys
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	privateKeySigner, err := ssh.NewSignerFromKey(privateKey)
	require.NoError(t, err)

	publicKey := privateKeySigner.PublicKey()

	// Assuming New takes a struct now, update fields accordingly
	server := New(Config{
		Addr:           addr,
		AuthorizedKeys: []ssh.PublicKey{publicKey},
	})
	require.NotNil(t, server)

	if err := server.Start(); err != nil {
		log.Fatal("Failed to start test SFTP server:", err)
	}
	defer server.Stop()

	// Set up the SSH client configuration using the test private key for authentication
	sshConfig := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(privateKeySigner),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	conn, err := ssh.Dial("tcp", addr, sshConfig)
	require.NoError(t, err)
	require.NotNil(t, conn)
	defer conn.Close()

	client, err := sftp.NewClient(conn)
	require.NoError(t, err)
	require.NotNil(t, client)
	defer client.Close()
}

func TestFailedConnection(t *testing.T) {
	cfg := Config{
		Addr:     addr,
		Username: username,
		Password: password,
	}

	server := New(cfg)
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

func TestFailedToListen(t *testing.T) {
	cfg := &Config{}
	cfg.Addr = "localhost:22"
	cfg.Username = username
	cfg.Password = password

	server := New(*cfg)
	err := server.Start()

	require.Error(t, err)
}
