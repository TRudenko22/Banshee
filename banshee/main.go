package main

import (
	"fmt"
	"log"
	"net/smtp"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Email struct {
	Sender     string
	Pass       string
	Recipients []string
	Subject    string
	Body       string
}

// Handles all of the sending logic
func (e *Email) Send(server, port string) (err error) {

	// Auth via Sender app password
	auth := smtp.PlainAuth("", e.Sender, e.Pass, server)

	err = smtp.SendMail(
		fmt.Sprintf("%v:%v", server, port),
		auth,
		e.Sender,
		e.Recipients,
		[]byte(fmt.Sprintf("Subject: %s \r\n\r\n%s", e.Subject, e.Body)),
	)

	for _, i := range e.Recipients {
		log.Printf("%s  -->  %s\n", e.Sender, i)
	}

	return
}

func (e *Email) Output() {
	fmt.Printf("Sender  : %s\n", e.Sender)
	fmt.Printf("Subject : %s\n", e.Subject)
}

// Takes ~/.config/banshee/banshee.yml and returns an Email struct
func LoadConfig(path string) (config Email, err error) {

	if path == "" {
		return Email{}, fmt.Errorf("No configuration path present")
	}

	lstFiles, err := os.ReadDir(path)
	if err != nil {
		return
	}

	configFile := ""
	for _, file := range lstFiles {
		ext := filepath.Ext(file.Name())
		if ext == ".yml" || ext == ".yaml" {
			configFile = file.Name()
			log.Println("FOUND", file.Name())
		}
	}

	if configFile == "" {
		return Email{}, fmt.Errorf("Couldn't file a YAML configuration")
	}

	config_file := viper.New()
	config_file.SetConfigType("yaml")
	config_file.SetConfigFile(path + configFile)

	err = config_file.ReadInConfig()
	if err != nil {
		return Email{}, fmt.Errorf("error reading in config file")
	}

	err = config_file.Unmarshal(&config)

	return
}

var (
	emailServer = os.Getenv("EMAIL_SERVER")
	emailPort   = os.Getenv("EMAIL_PORT")
)

func main() {
	var email Email
	var err error

	email, err = LoadConfig("/data/")
	if err != nil {
		log.Fatal(err)
	}

	email.Output()

	err = email.Send("smtp.gmail.com", "587")
	if err != nil {
		log.Fatal("Error Sending Email: Check Configuration Variables")
	}
}
