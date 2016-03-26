package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/garyburd/go-oauth/oauth"
)

// ** CHANGE THIS
// a quick hack for username, since used in reply
const TwitterUsername string = "mkaz"

func Usage() {
	fmt.Println(`twstorage usage:

# Encrypt Usage:

twstorage <filename>

    * If not authorized will open a browser and
      prompt to authorize to Twitter requiring a PIN
      once authorized, will save tokens to ~/.twstorage

    Returns:
        Key   : <key>
        Tweet : <url>

# Decrypt Usage:

twstorage -k <key> <url>

`)
	os.Exit(1)
}

var key string
var token *oauth.Credentials

func main() {

	flag.StringVar(&key, "k", "", "encryption key")
	flag.Parse()

	if flag.NArg() == 0 {
		Usage()
	}

	// do twitter authorization
	TwitterAuthorization()

	if key == "" {

		// encrypting
		filename := flag.Arg(0)
		fmt.Println("Encrypting: ", filename)

		// generate keys
		key = randomString(32)
		fmt.Println("Key  : ", key)

		// read in file
		content, err := ioutil.ReadFile(filename)
		if err != nil {
			log.Fatalln("Error reading file", filename)
		}

		// encrypt
		enc, err := encrypt(content, key)
		if err != nil {
			log.Fatalln("Error encrypting file", err)
		}

		// chunkify tweet limits
		chunksize := 137 - len(TwitterUsername)
		fmt.Println("Chunk Size: ", chunksize)
		chunks := chunkify(enc, chunksize)

		// upload to twitter
		var tweet, first Tweet
		for i, chunk := range chunks {

			// initial tweet
			fmt.Println("Tweeting chunk: ", i)
			if i != 0 {
				chunk = fmt.Sprintf("@%s %s", TwitterUsername, chunk)
			}

			tweet, err = postTweet(token, chunk, tweet.Identifier)
			if err != nil {
				log.Fatalf("Error posting tweet: %s %v", chunk, err)
			}

			if i == 0 {
				first = tweet
			}

			time.Sleep(250 * time.Millisecond) // take it easy
		}
		fmt.Printf("Tweet: https://twitter.com/%s/status/%s\n", tweet.User.ScreenName, first.Identifier)
		os.Exit(1)
	}

	// decrypting

	url := flag.Arg(0)
	fmt.Println("Fetching: ", url)

	// parse tweet id from url
	li := strings.LastIndex(url, "/") + 1
	tweetId := url[li:len(url)]

	// fetch from twitter
	encText, err := getTweetChain(token, tweetId)
	if err != nil {
		fmt.Println("Error fetching Tweet chain: ", err)
	}

	// use key to decrypt
	dec, err := decrypt(encText, key)
	if err != nil {
		fmt.Println("Error decrypting text")
	}

	fmt.Println(dec)

}

// URL starts at bottom of chain and works way to top
func getTweetChain(token *oauth.Credentials, tweetId string) (str string, err error) {
	var tweet Tweet

	for tweetId != "" {
		tweet, err = getTweet(token, tweetId)
		if err != nil {
			return "", err
		}

		// TODO: use screenname, strip RT proper
		rt := fmt.Sprintf("@%s ", TwitterUsername)
		str = strings.Replace(tweet.Text, rt, "", 1) + str
		tweetId = tweet.ReplyId
	}

	return str, nil
}

// Reads Config and Setups Auth
func TwitterAuthorization() {
	var err error
	authorized := false

	file, config := getConfig()
	token, authorized, err = getAccessToken(config)
	if err != nil {
		log.Fatalf("Failed to get access token: %s, %s, %s", token, authorized, err)
	}

	// save config file
	if authorized {
		b, err := json.MarshalIndent(config, "", "  ")
		if err != nil {
			log.Fatal("failed to store file:", err)
		}
		err = ioutil.WriteFile(file, b, 0700)
		if err != nil {
			log.Fatal("failed to store file:", err)
		}
	}
}

// config stores at: ~/.twstorage
// Requires client secret/token from https://dev.twitter.com/
// If authorized, includes authorization tokens
func getConfig() (string, map[string]string) {

	dir := os.Getenv("HOME")
	settings := filepath.Join(dir, ".twstorage")
	config := map[string]string{}

	b, err := ioutil.ReadFile(settings)
	if err != nil {
		fmt.Println("Error reading settings.")
		fmt.Println("Create ~/.twstorage with your app settings from dev.twitter.com")
		fmt.Println("{")
		fmt.Println(`  "ClientSecret": "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789",`)
		fmt.Println(`  "ClientToken": "abcdefghijklmnopqrstuvwxyz"`)
		fmt.Println("}")
		os.Exit(1)
	} else {
		err = json.Unmarshal(b, &config)
		if err != nil {
			log.Fatalf("Error in JSON? Could not unmarshal %v: %v", settings, err)
		}
	}

	return settings, config
}
