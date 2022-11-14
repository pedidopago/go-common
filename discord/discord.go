package discord

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"strings"
)

type NewWebhookMessageInput struct {
	URL      string
	Content  string
	Username string
	Logs     string
	Filename string
}

func NewWebhookMessage(input NewWebhookMessageInput) error {
	if input.URL == "" {
		return nil
	}
	if input.Logs == "" {
		ddata := struct {
			Content  string `json:"content"`
			Username string `json:"username"`
		}{
			Content:  input.Content,
			Username: input.Username,
		}
		jdd, _ := json.Marshal(ddata)
		req, _ := http.NewRequest(http.MethodPost, input.URL, bytes.NewReader(jdd))
		req.Header.Set("Content-Type", "application/json")
		req.Close = true
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		io.Copy(ioutil.Discard, resp.Body)
		if resp.StatusCode != 200 && resp.StatusCode != 204 {
			return fmt.Errorf("%d", resp.StatusCode)
		}
		return nil
	}
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	filename := "output.txt"
	if input.Filename != "" {
		filename = input.Filename
	}
	ww, _ := writer.CreateFormFile("file1", filename)
	io.Copy(ww, strings.NewReader(input.Logs))
	ddata := struct {
		Content  string `json:"content"`
		Username string `json:"username"`
	}{
		Content:  input.Content,
		Username: input.Username,
	}
	jdd, _ := json.Marshal(ddata)
	writer.WriteField("payload_json", string(jdd))
	writer.Close()
	req, _ := http.NewRequest(http.MethodPost, input.URL, &body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Close = true
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	io.Copy(ioutil.Discard, resp.Body)
	if resp.StatusCode != 200 && resp.StatusCode != 204 {
		return fmt.Errorf("%d", resp.StatusCode)
	}
	return nil
}
