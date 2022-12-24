package main

import (
	"github.com/watariRyo/go-sftp/infrastructure"
)

func main() {
	conn, err := infrastructure.GetSftpConnectionString()
	if err != nil {
		panic("Can not get sftp connection")
	}
	infrastructure.UploadSftp(conn, "./upload/test.txt", "./sftp-local-test/upload/hoge.txt")
	infrastructure.DownloadFile(conn, "./upload/test.txt", "./sftp-local-test/download/hoge.txt")
}
