package images

import (
	"bytes"
	"fmt"
	"github.com/DictumMortuum/servus/pkg/models"
	"github.com/Kagami/go-avif"
	"github.com/jmoiron/sqlx"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"net/http"
	"os"
	"sort"
)

type ImageConfig struct {
	URL    string
	Config image.Config
	Body   io.ReadCloser
}

func imagePath(id int64) string {
	return fmt.Sprintf("/home/dimitris/Code/DictumMortuum/json-api/rest/v1/boardgames/%d/image.avif", id)
}

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

func Boardgame(db *sqlx.DB, args *models.QueryBuilder) (interface{}, error) {
	rs := []ImageConfig{}

	path := imagePath(args.Id)

	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return nil, nil
	}

	urls, err := getBoardgameUrls(db, args.Id)
	if err != nil {
		return nil, err
	}

	urls2, err := getPriceUrls(db, args.Id)
	if err != nil {
		return nil, err
	}
	urls = append(urls, urls2...)

	fmt.Println(urls)

	if len(urls) == 0 {
		return rs, nil
	}

	for _, url := range urls {
		tmp, err := checkImage(args.Id, url)
		fmt.Println(err)
		if err == nil {
			fmt.Println(tmp.Config.Height * tmp.Config.Width)
			// if tmp.Config.Height*tmp.Config.Width < 1638400 {
			rs = append(rs, *tmp)
			// }
		}
	}

	sort.Slice(rs, func(i, j int) bool {
		return (rs[i].Config.Width * rs[i].Config.Height) > (rs[j].Config.Width * rs[j].Config.Height)
	})

	decodeImage(args.Id, rs[0])

	return rs, nil
}
