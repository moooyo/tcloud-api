package util

import (
	"time"

	"gopkg.in/gomail.v2"
)

var ch chan *gomail.Message

const codeLength = 6

func GenerateRegisterCode() string {
	return GenerateCaptcha(codeLength)
}

func SendRegisterCode(email string, code string) {
	m := gomail.NewMessage()
	m.SetHeader("From", "email@lengyu.me")
	m.SetHeader("To", email)
	m.SetHeader("Subject", "欢迎注册TCloud, 您还需要继续进行账户确认")
	m.SetBody("text/html", code)
	ch <- m
}

func InitMailServer() {
	config := GetConfig().Mail
	if config.Disable {
		return
	}

	ch = make(chan *gomail.Message, 1024)
	go func() {
		d := gomail.NewDialer(config.Host, config.Port, config.Username, config.Passowrd)

		var s gomail.SendCloser
		var err error
		open := false
		for {
			select {
			case m, ok := <-ch:
				if !ok {
					return
				}
				if !open {
					if s, err = d.Dial(); err != nil {
						ERROR("%e", err)
						break
					}
					open = true
				}
				if err := gomail.Send(s, m); err != nil {
					ERROR("%e", err)
					break
				}
			case <-time.After(30 * time.Second):
				if open {
					if err := s.Close(); err != nil {
						ERROR("%e", err)
						break
					}
					open = false
				}
			}
		}
	}()
}
