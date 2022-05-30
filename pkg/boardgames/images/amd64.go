//go:build !arm64
// +build !arm64

package images

import (
	"bytes"
	"fmt"
	"github.com/Kagami/go-avif"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"net/http"
	"os"
)

func decodeImage(id int64, cfg ImageConfig) error {
	dst, err := os.Create(imagePath(id))
	if err != nil {
		return err
	}

	resp, err := http.Get(cfg.URL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	avifImage, _, err := image.Decode(resp.Body)
	if err != nil {
		fmt.Println("1", id, err)
		return err
	}

	err = avif.Encode(dst, avifImage, nil)
	if err != nil {
		fmt.Println("2")
		return err
	}

	return nil
}

func checkImage(id int64, url string) (*ImageConfig, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	res := io.NopCloser(bytes.NewReader(bodyBytes))

	img, _, err := image.DecodeConfig(res)
	if err != nil {
		return nil, err
	}

	return &ImageConfig{
		Config: img,
		Body:   res,
		URL:    url,
	}, nil
}
