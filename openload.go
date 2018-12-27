package openload

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

const API = "https://api.openload.co/1"

type Ticket struct {
	Ticket  string
	Captcha struct {
		Url    string
		Width  int
		Height int
	}
	WaitTime   int
	ValidUntil time.Time
}

type Download struct {
	Name        string
	Size        int
	Sha1        string
	ContentType string
	UploadAt    time.Time
	Url         string
	Token       string
}

type Info struct {
	Id          string
	Status      int
	Name        string
	Size        string
	Sha1        string
	ContentType string
	CStatus     string
}

func GenerateTicket(cli *http.Client, fileId, login, key string) (*Ticket, error) {
	var u string
	if login != "" && key != "" {
		u = fmt.Sprintf(
			API+"/file/dlticket?file=%s&login=%s&key=%s",
			fileId, login, key)
	} else {
		u = fmt.Sprintf(API+"/file/dlticket?file=%s",
			fileId)
	}
	res, err := cli.Get(u)
	if err != nil {
		return nil, err
	}
	var response struct {
		Status int    `json:"status"`
		Msg    string `json:"msg"`
		Result struct {
			Ticket        string `json:"ticket"`
			CaptchaUrl    string `json:"captcha_url"`
			CaptchaWidth  int    `json:"captcha_w"`
			CaptchaHeight int    `json:"captcha_h"`
			WaitTime      int    `json:"wait_time"`
			ValidUntil    string `json:"valid_until"`
		} `json:"result"`
	}

	err = json.NewDecoder(res.Body).Decode(&response)
	res.Body.Close()
	if err != nil {
		return nil, err
	}
	if response.Status != 200 {
		return nil, errors.New(response.Msg)
	}
	r := response.Result
	validUntil, err := time.Parse("2006-01-02 15:04:05", r.ValidUntil)
	if err != nil {
		return nil, err
	}
	t := &Ticket{
		Ticket:     r.Ticket,
		WaitTime:   r.WaitTime,
		ValidUntil: validUntil,
	}
	t.Captcha.Url = r.CaptchaUrl
	t.Captcha.Width = r.CaptchaWidth
	t.Captcha.Height = r.CaptchaHeight
	return t, nil
}

func RequestDownload(cli *http.Client, fileId, ticket, captcha string) (*Download, error) {
	u := fmt.Sprintf(
		API+"/file/dl?file=%s&ticket=%s&captcha_response=%s",
		fileId, ticket, captcha)
	res, err := cli.Get(u)
	if err != nil {
		return nil, err
	}
	var response struct {
		Status int    `json:"status"`
		Msg    string `json:"msg"`
		Result struct {
			Name        string `json:"name"`
			Size        int    `json:"size"`
			Sha1        string `json:"sha1"`
			ContentType string `json:"content_type"`
			UploadAt    string `json:"upload_at"`
			Url         string `json:"url"`
			Token       string `json:"token"`
		} `json:"result"`
	}
	err = json.NewDecoder(res.Body).Decode(&response)
	res.Body.Close()
	if err != nil {
		return nil, err
	}
	if response.Status != 200 {
		return nil, errors.New(response.Msg)
	}
	r := response.Result
	uploadAt, err := time.Parse("2006-01-02 15:04:05", r.UploadAt)
	if err != nil {
		return nil, err
	}
	d := &Download{
		Name:        r.Name,
		Size:        r.Size,
		Sha1:        r.Sha1,
		ContentType: r.ContentType,
		UploadAt:    uploadAt,
		Url:         r.Url,
		Token:       r.Token,
	}
	return d, nil
}

func FileInfo(cli *http.Client, fileId, login, key string) ([]Info, error) {
	var u string
	if login != "" && key != "" {
		u = fmt.Sprintf(
			API+"/file/info?file=%s&login=%s&key=%s",
			fileId, login, key)
	} else {
		u = fmt.Sprintf(API+"/file/info?file=%s",
			fileId)
	}
	res, err := cli.Get(u)
	if err != nil {
		return nil, err
	}
	var response struct {
		Status int    `json:"status"`
		Msg    string `json:"msg"`
		Result map[string]struct {
			Id          string `json:"id"`
			Status      int    `json:"status"`
			Name        string `json:"name"`
			Size        string `json:"size"`
			Sha1        string `json:"sha1"`
			ContentType string `json:"content_type"`
			CStatus     string `json:"cstatus"`
		} `json:"result"`
	}

	err = json.NewDecoder(res.Body).Decode(&response)
	res.Body.Close()
	if err != nil {
		return nil, err
	}
	if response.Status != 200 {
		return nil, errors.New(response.Msg)
	}
	r := response.Result
	infos := make([]Info, len(r))
	o := 0
	for _, i := range r {
		infos[o] = Info{
			Id:          i.Id,
			Status:      i.Status,
			Name:        i.Name,
			Size:        i.Size,
			Sha1:        i.Sha1,
			ContentType: i.ContentType,
			CStatus:     i.CStatus,
		}
		o++
	}
	return infos, nil
}
