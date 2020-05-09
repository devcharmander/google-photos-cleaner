package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type Conf struct {
	creadentialsPath string //file path where the credentials json file is stored
	scope            []string
}

type CLI struct {
	config *Conf
	client *http.Client
}

// Retrieve a token, saves the token, then returns the generated client.
func getClient(cfg *Conf) *http.Client {

	b, err := ioutil.ReadFile(cfg.creadentialsPath)
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, cfg.scope...)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
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

// Retrieves a token from a local file.
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

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0700)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func main() {

	cfg := &Conf{
		creadentialsPath: "credentials.json",
		scope:            []string{"https://www.googleapis.com/auth/photoslibrary", "https://www.googleapis.com/auth/photoslibrary.readonly", "https://www.googleapis.com/auth/photoslibrary.readonly.appcreateddata"},
	}

	client := getClient(cfg)

	cli := &CLI{
		client: client,
		config: cfg,
	}

	cli.prerequisites()

	fmt.Println("calculating number of junk items ..")
	cli.getJunk(listMediaItems)
	info := cli.getMediaIds()
	fmt.Println("There are ", len(info), "files that can be deleted")

	fmt.Println("creating backup ...")

	cli.createBackup(info)

	fmt.Println("backyp created at ", backupPath)

	fmt.Println("Fetching your albums from Google Photos...")
	n, kv := cli.ListAlbums(listAlbums)
	fmt.Println("There are ", n, "albums in this account")
	fmt.Println("Select the newly created album from the list (Enter album number)")
	var an int
	_, err := fmt.Scan(&an)
	if err != nil {
		log.Fatalf("Could not read the user input")
	}
	if an > n || an < 0 {
		log.Fatalf("Invalid album id")
	}

	fmt.Println("Deleting ", len(info), " items")
	albumId := kv[an]
	cli.RemoveMediaItems(albumId, info)

}

func (cli *CLI) prerequisites() {

	fmt.Println("Google doesn't provide API endpoints to delete photos or move them to a new album, so there are some pre-requisites")
	fmt.Println("Please make sure that you have copied all your photos to a single album")
	fmt.Println("Don't worry, copying your photos to a new album doesn't consume space")

	fmt.Println("Cleanup would happen in three steps:")
	fmt.Println("1- Select an album by providing the album number")
	fmt.Println("2- a copy of the files that would be deleted will be saved into your local disk, as a backup")
	fmt.Println("3- actual deletion happens from the album (the album where you have all your photos)")

	err := os.MkdirAll(backupPath, 0700)
	if err != nil {
		log.Fatalf("Could not create the backup directory. Error %v", err)
	}

}
