package sftpserver

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"log"
	"net"
	"os"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

type SFTPServer struct {
	Addr       string
	Username   string
	Password   string
	PrivateKey ssh.Signer
	listener   net.Listener
}

func New(addr, username, password string) *SFTPServer {
	hostKey, err := generateTempHostKey()
	if err != nil {
		return nil
	}

	return &SFTPServer{
		Addr:       addr,
		Username:   username,
		Password:   password,
		PrivateKey: hostKey,
	}
}

func (s *SFTPServer) Start() error {
	config := &ssh.ServerConfig{
		PasswordCallback: s.passwordCallback,
	}
	config.AddHostKey(s.PrivateKey)

	var err error
	s.listener, err = net.Listen("tcp", s.Addr)
	if err != nil {
		return fmt.Errorf("failed to listen for connection: %w", err)
	}

	go s.acceptConnections(config)
	return nil
}

func (s *SFTPServer) Stop() error {
	return s.listener.Close()
}

// generateTempHostKey generates a new RSA private key and saves it to a temporary file.
// It returns the signer for use with ssh.ServerConfig and the path to the temporary file.
func generateTempHostKey() (ssh.Signer, error) {
	// Generate a new RSA key.
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, fmt.Errorf("failed to generate RSA key: %w", err)
	}

	// Marshal the RSA private key to ASN.1 DER encoded form.
	privateKeyDER := x509.MarshalPKCS1PrivateKey(key)

	// Create a PEM block for the private key.
	privateKeyBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyDER,
	}

	// Create a temporary file for the private key.
	tempFile, err := os.CreateTemp("", "ssh_host_key")
	if err != nil {
		return nil, fmt.Errorf("failed to create temporary file for host key: %w", err)
	}
	defer tempFile.Close()

	// Write the PEM block to the temporary file.
	if err := pem.Encode(tempFile, privateKeyBlock); err != nil {
		return nil, fmt.Errorf("failed to write PEM block to temporary file: %w", err)
	}

	// Create a signer from the private key for use with ssh.ServerConfig.
	signer, err := ssh.NewSignerFromSigner(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create signer from private key: %w", err)
	}

	return signer, nil
}

func (s *SFTPServer) passwordCallback(c ssh.ConnMetadata, pass []byte) (*ssh.Permissions, error) {
	if c.User() == s.Username && string(pass) == s.Password {
		return nil, nil
	}
	return nil, fmt.Errorf("password rejected for %q", c.User())
}

func (s *SFTPServer) acceptConnections(config *ssh.ServerConfig) {
	for {
		nConn, err := s.listener.Accept()
		if err != nil {
			log.Printf("failed to accept incoming connection: %v", err)
			break
		}
		go s.handleConnection(nConn, config)
	}
}

func (s *SFTPServer) handleConnection(nConn net.Conn, config *ssh.ServerConfig) {
	sshConn, chans, reqs, err := ssh.NewServerConn(nConn, config)
	if err != nil {
		log.Printf("failed to handshake: %v", err)
		return
	}
	defer sshConn.Close()

	// The incoming requests must be serviced.
	go ssh.DiscardRequests(reqs)

	s.handleChannels(chans)
}

func (s *SFTPServer) handleChannels(chans <-chan ssh.NewChannel) {
	for newChannel := range chans {
		if newChannel.ChannelType() != "session" {
			newChannel.Reject(ssh.UnknownChannelType, "unknown channel type")
			continue
		}

		channel, requests, err := newChannel.Accept()
		if err != nil {
			continue
		}

		go s.handleRequests(requests, channel)
	}
}

func (s *SFTPServer) handleRequests(requests <-chan *ssh.Request, channel ssh.Channel) {
	for req := range requests {
		ok := false
		switch req.Type {
		case "subsystem":
			if string(req.Payload[4:]) == "sftp" {
				ok = true
			}
		}

		req.Reply(ok, nil)
		if ok {
			s.startSFTPSession(channel)
		}
	}
}

func (s *SFTPServer) startSFTPSession(channel ssh.Channel) {
	server, err := sftp.NewServer(
		channel,
	)
	if err != nil {
		log.Printf("failed to start sftp server: %v", err)
		return
	}
	if err := server.Serve(); err == io.EOF {
		server.Close()
	} else if err != nil {
		log.Printf("sftp server completed with error: %v", err)
	}
}
