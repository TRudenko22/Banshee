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
func (e *Email) Send(server, port string) (err error) {

	// Auth via Sender app password
    auth := smtp.PlainAuth("", e.Sender, e.Pass, server,)

    err = smtp.SendMail(
	    fmt.Sprintf("%v:%v", server, port),
		auth,
		e.Sender,
		e.Recipients,
		[]byte(fmt.Sprintf("Subject: %v \r\n\r\n%v", e.Subject, e.Body)),
	) 
  
    return
}


func LoadConfig() (config Email, err error) {

	config_file := viper.New()
	config_file.AddConfigPath(".")
	config_file.SetConfigFile("email.yml")

	err = config_file.ReadInConfig()
	if err != nil {
        return 
	}

	err = config_file.Unmarshal(&config)

	return
}

func main() {

	email, err := LoadConfig()
	if err != nil {
    	log.Fatal("Error reading configuration: ", err)
	}

	err = email.Send("smtp.gmail.com", "587")
  	if err != nil {
    	log.Fatal("Error Sending Email: Check Configuration Variables")
  	}
}
