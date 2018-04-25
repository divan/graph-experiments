package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"path/filepath"
	"time"

	"github.com/mitchellh/go-homedir"

	"golang.org/x/crypto/ssh"
)

const sshPort = "22"

// SSH represents SSH connection to cluster node,
// which runs statusd/geth in docker instances.
type SSH struct {
	Host    string
	session *ssh.Session
}

// NewSSH opens new SSH connection and prepares Session to work with.
// It's in responsibility of called to close the Session after using it.
func NewSSH(host, user string) (*SSH, error) {
	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			publicKeyFile(),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         10 * time.Second,
	}

	client, err := ssh.Dial("tcp", host+":"+sshPort, config)
	if err != nil {
		return nil, err
	}

	session, err := client.NewSession()
	if err != nil {
		return nil, fmt.Errorf("Failed to create session: %s", err)
	}

	modes := ssh.TerminalModes{
		// ssh.ECHO:          0,     // disable echoing
		ssh.TTY_OP_ISPEED: 9600, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 9600, // output speed = 14.4kbaud
	}

	if err := session.RequestPty("xterm", 80, 40, modes); err != nil {
		session.Close()
		return nil, fmt.Errorf("request for pseudo terminal failed: %s", err)
	}

	return &SSH{
		Host:    host,
		session: session,
	}, nil
}

// Close closes the connection and active session.
func (s *SSH) Close() error {
	if s == nil {
		return nil
	}
	return s.session.Close()
}

func publicKeyFile() ssh.AuthMethod {
	home, err := homedir.Dir()
	if err != nil {
		log.Fatalln(fmt.Sprintf("Cannot find home dir for SSH public key file %v", err))
		return nil
	}
	file := filepath.Join(home, ".ssh", "dec.key")
	buffer, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatalln(fmt.Sprintf("Cannot read SSH public key file %s: %v", file, err))
		return nil
	}

	key, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		log.Fatalln(fmt.Sprintf("Cannot parse SSH public key file %s: %v", file, err))
		return nil
	}
	return ssh.PublicKeys(key)
}

type SSHCommand struct {
	Path   string
	Env    []string
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

func (s *SSH) run(cmd *SSHCommand) error {
	if cmd.Stdin != nil {
		stdin, err := s.session.StdinPipe()
		if err != nil {
			return fmt.Errorf("Unable to setup stdin for session: %v", err)
		}
		go io.Copy(stdin, cmd.Stdin)
	}

	if cmd.Stdout != nil {
		stdout, err := s.session.StdoutPipe()
		if err != nil {
			return fmt.Errorf("Unable to setup stdout for session: %v", err)
		}
		go io.Copy(cmd.Stdout, stdout)
	}

	if cmd.Stderr != nil {
		stderr, err := s.session.StderrPipe()
		if err != nil {
			return fmt.Errorf("Unable to setup stderr for session: %v", err)
		}
		go io.Copy(cmd.Stderr, stderr)
	}

	return s.session.Run(cmd.Path)
}

// Exec executes bash command on remote machine, and returns stdout and error.
func (s *SSH) Exec(path string) (string, error) {
	buf := new(bytes.Buffer)
	errBuf := new(bytes.Buffer)
	cmd := &SSHCommand{
		Path:   path,
		Stdout: buf,
		Stderr: errBuf,
	}

	if err := s.run(cmd); err != nil {
		return "", fmt.Errorf("ssh command: %s: %s", err, errBuf.String())
	}

	// TODO(divan): WTF: without delay somewhere between runs, it sometimes returns truncated output :/
	time.Sleep(100 * time.Millisecond)
	return buf.String(), nil
}
