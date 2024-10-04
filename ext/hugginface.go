package ext

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
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

	err = hf.saveImage(body)
	if err != nil {
		fmt.Println("couldn't save image: ", err)
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

func (hf HuggingFace) saveImage(resp []byte) error {
	img, _, err := image.Decode(bytes.NewReader(resp))
	if err != nil {
		return err
	}

	out, _ := os.Create(
		fmt.Sprintf(".temp/images/%v.jpg", time.Now().Unix()))
	defer out.Close()

	err = jpeg.Encode(out, img, &jpeg.Options{Quality: 100})
	if err != nil {
		return err
	}

	return nil
}
