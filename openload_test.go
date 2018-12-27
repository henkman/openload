package openload_test

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/henkman/openload"
)

const (
	FILEID = ""
	LOGIN  = ""
	PASS   = ""
)

func TestDownload(t *testing.T) {
	cli := http.Client{
		Timeout: time.Second * 10,
	}
	ticket, err := openload.GenerateTicket(&cli, FILEID, LOGIN, PASS)
	if err != nil {
		fmt.Println(err)
		t.Fail()
		return
	}
	fmt.Println(ticket)
}

func TestFileInfo(t *testing.T) {
	cli := http.Client{
		Timeout: time.Second * 10,
	}
	info, err := openload.FileInfo(&cli, FILEID, "", "")
	if err != nil {
		fmt.Println(err)
		t.Fail()
		return
	}
	fmt.Println(info)
}
