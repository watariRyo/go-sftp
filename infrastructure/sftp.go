package infrastructure

import (
	"fmt"
	"github.com/pkg/sftp"
	"github.com/watariRyo/go-sftp/config"
	"golang.org/x/crypto/ssh"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func GetSftpConnectionString() (*ssh.Client, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}
	host := cfg.Sftp.Host
	port := cfg.Sftp.Port
	user := cfg.Sftp.User
	pass := cfg.Sftp.Password

	fmt.Println("Create sshClientConfig")
	sshConfig := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(pass),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	fmt.Println("SSH connect")
	addr := fmt.Sprintf("%s:%s", host, port)
	fmt.Println(addr)

	conn, err := ssh.Dial("tcp", addr, sshConfig)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func UploadSftp(conn *ssh.Client, uploadPath, localPath string) error {
	fmt.Println("open an SFTP session over an existing ssh connection")
	client, err := sftp.NewClient(conn)
	if err != nil {
		return err
	}
	defer client.Close()

	fmt.Sprintf("Uploading [%s] to [%s] ...\n\n", localPath, uploadPath)

	srcFile, err := os.Open(localPath)
	if err != nil {
		fmt.Sprintf("Unable to open local file: %v\n", err)
		return err
	}
	defer srcFile.Close()

	// Make remote directories recursion
	parent := filepath.Dir(uploadPath)
	path := string(filepath.Separator)
	dirs := strings.Split(parent, path)
	for _, dir := range dirs {
		path = filepath.Join(path, dir)
		client.Mkdir(path)
	}

	// Note: SFTP To Go doesn't support O_RDWR mode
	dstFile, err := client.OpenFile(uploadPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC)
	if err != nil {
		fmt.Sprintf("Unable to open remote file: %v\n", err)
		return err
	}
	defer dstFile.Close()

	bytes, err := io.Copy(dstFile, srcFile)
	if err != nil {
		fmt.Sprintf("Unable to upload local file: %v\n", err)
		os.Exit(1)
	}
	fmt.Sprintf("%d bytes copied\n", bytes)

	return nil
}

func DownloadFile(conn *ssh.Client, remotePath, localPath string) error {
	fmt.Println("open an SFTP session over an existing ssh connection")
	client, err := sftp.NewClient(conn)
	if err != nil {
		return err
	}
	defer client.Close()

	fmt.Sprintf("Downloading [%s] to [%s] ...\n", remotePath, localPath)
	// Note: SFTP To Go doesn't support O_RDWR mode
	srcFile, err := client.OpenFile(remotePath, os.O_RDONLY)
	if err != nil {
		fmt.Sprintf("Unable to open remote file: %v\n", err)
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(localPath)
	if err != nil {
		fmt.Sprintf("Unable to open local file: %v\n", err)
		return err
	}
	defer dstFile.Close()

	bytes, err := io.Copy(dstFile, srcFile)
	if err != nil {
		fmt.Sprintf("Unable to download remote file: %v\n", err)
		return err
	}
	fmt.Sprintf("%d bytes copied\n", bytes)

	return nil
}
