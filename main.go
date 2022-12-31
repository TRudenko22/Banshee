package main

import (
	"fmt"
	"log"
	"net/smtp"
	"os"

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
		[]byte(fmt.Sprintf("Subject: %v \r\n\r\n%v", e.Subject, e.Body)),
	)

	for _, i := range e.Recipients {
		log.Printf("%v  -->  %v\n", e.Sender, i)
	}

	return
}

func (e *Email) Output() {
	fmt.Printf("Sender  : %v\n", e.Sender)
	fmt.Printf("Subject : %v\n", e.Subject)
}

// Creates /home/$USER/.config/beacon directory
func SetPathFile() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("error locating home directory")
	}

	path := fmt.Sprintf("%v/.config/beacon/", homeDir)

	if _, err = os.Stat(path); os.IsNotExist(err) {
		err = os.MkdirAll(path, 0750)
		if err != nil {
			return "", fmt.Errorf("error creating %v", path)
		}
	}

	return path + "beacon.yml", nil
}

// Takes ~/.config/beacon/beacon.yml and returns an Email struct
func LoadConfig() (config Email, err error) {

	pathFile, pathError := SetPathFile()
	if pathError != nil {
		return
	}

	config_file := viper.New()
	config_file.SetConfigType("yaml")
	config_file.SetConfigFile(pathFile)

	err = config_file.ReadInConfig()
	if err != nil {
		return Email{}, fmt.Errorf("error reading in config file")
	}

	err = config_file.Unmarshal(&config)

	return
}

func main() {

	email, err := LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	email.Output()

	err = email.Send("smtp.gmail.com", "587")
	if err != nil {
		log.Fatal("Error Sending Email: Check Configuration Variables")
	}
}
