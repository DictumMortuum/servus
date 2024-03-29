package router

import (
	"github.com/DictumMortuum/servus/pkg/util"
	"github.com/gin-gonic/gin"
	"github.com/ziutek/telnet"
	"time"
)

const timeout = 10 * time.Second

func expect(t *telnet.Conn, d ...string) error {
	err := t.SetReadDeadline(time.Now().Add(timeout))
	if err != nil {
		return err
	}

	err = t.SkipUntil(d...)
	if err != nil {
		return err
	}

	return nil
}

func sendln(t *telnet.Conn, s string) error {
	err := t.SetWriteDeadline(time.Now().Add(timeout))
	if err != nil {
		return err
	}

	buf := make([]byte, len(s)+1)
	copy(buf, s)
	buf[len(s)] = '\n'

	_, err = t.Write(buf)
	if err != nil {
		return err
	}

	return nil
}

func reset(host, user, password string) error {
	t, err := telnet.Dial("tcp", host)
	if err != nil {
		return err
	}

	t.SetUnixWriteMode(true)

	err = expect(t, "Login:")
	if err != nil {
		return err
	}

	err = sendln(t, user)
	if err != nil {
		return err
	}

	err = expect(t, "Password:")
	if err != nil {
		return err
	}

	err = sendln(t, password)
	if err != nil {
		return err
	}

	err = expect(t, "WAP>")
	if err != nil {
		return err
	}

	err = sendln(t, "reset")
	if err != nil {
		return err
	}

	return nil
}

func Reset(c *gin.Context) {
	err := reset("192.168.5.254:23", "Forthnet", "F0rth@c$n3t#")
	if err != nil {
		util.Error(c, err)
	}

	util.Success(c, nil)
}
