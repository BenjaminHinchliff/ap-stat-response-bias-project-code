package main

import (
	"encoding/base64"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/mail"
	"os"
	"path/filepath"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
)

func getClient(config *oauth2.Config) *http.Client {
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func sendMessage(srv *gmail.Service, sender mail.Address, receiver mail.Address, subject string, body string) error {
	headers := map[string]string{}
	headers["From"] = sender.String()
	headers["To"] = receiver.String()
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/plain; charset=\"utf-8\""
	headers["Content-Transfer-Encoding"] = "base64"

	var msg string
	for k, v := range headers {
		msg += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	msg += "\r\n" + body

	gmsg := gmail.Message{
		Raw: base64.RawURLEncoding.EncodeToString([]byte(msg)),
	}

	fmt.Println("Sent Message:\n", gmsg)
	// _, err := srv.Users.Messages.Send("me", &gmsg).Do()
	return nil
}

const sourceDir string = "../../data/treatments/"

func main() {
	// next time I'm using SMTP.
	b, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	config, err := google.ConfigFromJSON(b, gmail.GmailSendScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	srv, err := gmail.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Gmail client: %v", err)
	}

	from := mail.Address{Name: "Benjamin Hinchliff", Address: "benjamin.hinchliff21@auhsdschools.org"}

	treatments, err := ioutil.ReadDir(sourceDir)
	if err != nil {
		log.Fatalf("Failed to open treatments dir: %v", err)
	}

	for _, file := range treatments {
		if filepath.Ext(file.Name()) == ".csv" {
			path := sourceDir + file.Name()
			f, err := os.Open(path)
			if err != nil {
				log.Fatalf("failed to open a treatment csv file %v\n", err)
			}
			r := csv.NewReader(f)
			for {
				record, err := r.Read()
				if err == io.EOF {
					break
				}
				if err != nil {
					log.Fatalf("failed to read a row from csv: %v\n", err)
				}
				to := mail.Address{
					Name:    record[1] + " " + record[0],
					Address: record[3],
				}

				sendMessage(srv, from, to, "test2", "hey other ben this is more testing")
				if err != nil {
					log.Fatalln("Failed to send email: ", err)
				}
			}
		}
	}
}
