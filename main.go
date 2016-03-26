package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/garyburd/go-oauth/oauth"
)

func Usage() {
	fmt.Println(`twstorage usage:

# Encrypt Usage:

twstorage <filename>

    * If not authorized will open a browser and
      prompt to authorize to Twitter requiring a PIN
      once authorized, will save tokens to ~/.twstorage

    Returns:
        Tweet : <url>
        Key   : <key>
        Nonce : <nonce>

# Decrypt Usage:

twstorage --key <key> --nonce <nonce> <url>

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

		// chunkify
		chunks := chunkify(enc, 132)

		// upload to twitter
		var tweet, first Tweet
		for i, chunk := range chunks {

			// initial tweet
			fmt.Println("Tweeting chunk: ", i)
			if i != 0 {
				chunk = fmt.Sprintf("@mkaz %s", chunk) // TODO: use screenname, shrink chunksize
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

func TwitterAuthorization() {
	var err error
	authorized := false

	file, config := getConfig()
	token, authorized, err = getAccessToken(config)
	if err != nil {
		log.Fatalf("Failed to get access token: %s, %s, %s", token, authorized, err)
	}

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
