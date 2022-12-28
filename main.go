package main

import (
	"fmt"
	"log"
	"net/smtp"

	"github.com/spf13/viper"
)

type Email struct {
	Sender     string
	Pass       string
	Recipients []string
	Subject    string
	Body       string
}

// handles all of the sending logic
func (e *Email) Send(server, port string) {

	// Auth via Sender app password
	auth := smtp.PlainAuth("", e.Sender, e.Pass, server,)

	smtp.SendMail(
		fmt.Sprintf("%v:%v", server, port),
		auth,
		e.Sender,
		e.Recipients,
		[]byte(fmt.Sprintf("Subject: %v \r\n\r\n%v", e.Subject, e.Body)),
	)
}

func LoadConfig() (config Email, conf_err error) {

	//config_file.AddConfigPath("/etc/email.yml")

	config_file := viper.New()
	config_file.AddConfigPath(".")
	config_file.SetConfigFile("email.yml")

	err := config_file.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}

	conf_err = config_file.Unmarshal(&config)

	return
}

func main() {

	email, err := LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	email.Send("smtp.gmail.com", "587")
}
