package main

import (
	"fmt"
	"log"
    "os"
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

func SetPathFile() (string, error) {
    homeDir, err := os.UserHomeDir()
    if err != nil {
        log.Fatal("error finding home")
        return "", err
    }

    path := fmt.Sprintf("%v/.config/beacon/", homeDir)

    if _, err = os.Stat(path); os.IsNotExist(err) {
        err = os.MkdirAll(path, 0750)
        if err != nil {
            return "", err
        }
    } 

    return path + "beacon.yml", nil
}

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
