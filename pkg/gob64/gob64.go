package gob64

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
)

func FromGOB64(p interface{}, str string) error {
	by, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return err
	}

	b := bytes.Buffer{}
	b.Write(by)
	d := gob.NewDecoder(&b)

	err = d.Decode(p)
	if err != nil {
		return err
	}

	return nil
}

func ToGOB64(p interface{}) (string, error) {
	b := bytes.Buffer{}
	e := gob.NewEncoder(&b)
	err := e.Encode(p)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b.Bytes()), nil
}
