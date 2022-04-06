package w3m

import (
	"crypto/md5"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

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
