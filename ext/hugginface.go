package ext

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const VITGPT2_APIURL = "https://api-inference.huggingface.co/models/ZB-Tech/Text-to-Image"

type HuggingFace struct {
	apikey string
}

func NewHuggingFace(apikey string) *HuggingFace {
	return &HuggingFace{
		apikey: apikey,
	}
}

func (hf HuggingFace) GenerateImg(text string) (*bytes.Reader, error) {
	resp, err := hf.sendHFHttpRequest(text)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return nil, err
	}

	return bytes.NewReader(body), nil
}

func (hf HuggingFace) sendHFHttpRequest(text string) (*http.Response, error) {
	req, err := http.NewRequest("POST", VITGPT2_APIURL, strings.NewReader(text))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", hf.apikey))

	client := &http.Client{}
	return client.Do(req)
}
