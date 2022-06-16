package w3m

import (
	"crypto/md5"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func scrape(url string) string {
	return fmt.Sprintf(`import cfscrape
scraper = cfscrape.create_scraper()
print (scraper.get("%s").content.decode("utf-8"))
`, url)
}

func BypassCloudflare(url string) (string, error) {
	path, err := exec.LookPath("python")
	if err != nil {
		return "", nil
	}

	filename := "/tmp/" + fmt.Sprintf("%x", md5.Sum([]byte(url)))
	f, err := os.Create(filename)
	if err != nil {
		return "", err
	}
	defer f.Close()

	cmd := exec.Command(path, "-")
	cmd.Stdout = f

	in, err := cmd.StdinPipe()
	if err != nil {
		return "", err
	}

	err = cmd.Start()
	if err != nil {
		return "", err
	}

	_, err = in.Write([]byte(scrape(url)))
	if err != nil {
		return "", err
	}

	err = in.Close()
	if err != nil {
		return "", err
	}

	err = cmd.Wait()
	if err != nil {
		return "", err
	}

	file, err := filepath.Abs(filename)
	if err != nil {
		return "", err
	}

	return "file://" + file, nil
}

func Download(url string) (string, error) {
	path, err := exec.LookPath("w3m")
	if err != nil {
		return "", nil
	}

	filename := "/tmp/" + fmt.Sprintf("%x", md5.Sum([]byte(url)))
	f, err := os.Create(filename)
	if err != nil {
		return "", err
	}
	defer f.Close()

	cmd := exec.Command(path, "-dump_source", url)
	cmd.Stdout = f

	err = cmd.Run()
	if err != nil {
		return "", err
	}

	file, err := filepath.Abs(filename)
	if err != nil {
		return "", err
	}

	return "file://" + file, nil
}
