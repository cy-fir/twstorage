package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"runtime"

	"github.com/fatih/color"
	"github.com/garyburd/go-oauth/oauth"
)

type Tweet struct {
	Text       string `json:"text"`
	Identifier string `json:"id_str"`
	ReplyId    string `json:"in_reply_to_status_id_str"`
	User       struct {
		Name       string `json:"name"`
		ScreenName string `json:"screen_name"`
	} `json:"user"`
}

var oauthClient = oauth.Client{
	TemporaryCredentialRequestURI: "https://api.twitter.com/oauth/request_token",
	ResourceOwnerAuthorizationURI: "https://api.twitter.com/oauth/authenticate",
	TokenRequestURI:               "https://api.twitter.com/oauth/access_token",
}

func clientAuth(requestToken *oauth.Credentials) (*oauth.Credentials, error) {
	var err error
	browser := "xdg-open"
	url_ := oauthClient.AuthorizationURL(requestToken, nil)

	args := []string{url_}
	if runtime.GOOS == "windows" {
		browser = "rundll32.exe"
		args = []string{"url.dll,FileProtocolHandler", url_}
	} else if runtime.GOOS == "darwin" {
		browser = "open"
		args = []string{url_}
	} else if runtime.GOOS == "plan9" {
		browser = "plumb"
	}
	color.Set(color.FgHiRed)
	fmt.Println("Open this URL and enter PIN.")
	color.Set(color.Reset)
	fmt.Println(url_)
	browser, err = exec.LookPath(browser)
	if err == nil {
		cmd := exec.Command(browser, args...)
		cmd.Stderr = os.Stderr
		err = cmd.Start()
		if err != nil {
			log.Fatal("failed to start command:", err)
		}
	}

	fmt.Print("PIN: ")
	stdin := bufio.NewScanner(os.Stdin)
	if !stdin.Scan() {
		log.Fatal("canceled")
	}
	accessToken, _, err := oauthClient.RequestToken(http.DefaultClient, requestToken, stdin.Text())
	if err != nil {
		log.Fatal("failed to request token:", err)
	}
	return accessToken, nil
}

func getAccessToken(config map[string]string) (*oauth.Credentials, bool, error) {

	oauthClient.Credentials.Token = config["ClientToken"]
	oauthClient.Credentials.Secret = config["ClientSecret"]

	authorized := false
	var token *oauth.Credentials
	accessToken, foundToken := config["AccessToken"]
	accessSecert, foundSecret := config["AccessSecret"]
	if foundToken && foundSecret {
		token = &oauth.Credentials{accessToken, accessSecert}
	} else {
		requestToken, err := oauthClient.RequestTemporaryCredentials(http.DefaultClient, "", nil)
		if err != nil {
			log.Print("failed to request temporary credentials:", err)
			return nil, false, err
		}
		token, err = clientAuth(requestToken)
		if err != nil {
			log.Print("failed to request temporary credentials:", err)
			return nil, false, err
		}

		config["AccessToken"] = token.Token
		config["AccessSecret"] = token.Secret
		authorized = true
	}
	return token, authorized, nil
}

// TWITTER ENDPOINTS

func postTweet(token *oauth.Credentials, msg string, replyId string) (tweet Tweet, err error) {
	apiurl := "https://api.twitter.com/1.1/statuses/update.json"

	param := make(url.Values)
	param.Set("status", msg)
	if replyId != "" {
		param.Set("in_reply_to_status_id", replyId)
	}

	oauthClient.SignParam(token, "POST", apiurl, param)
	res, err := http.PostForm(apiurl, url.Values(param))
	if err != nil {
		log.Println("failed to post tweet:", err)
		return tweet, err
	}

	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Println("failed to get timeline:", res.StatusCode)
		log.Println("Body: ", res.Body)
		return tweet, err
	}

	err = json.NewDecoder(res.Body).Decode(&tweet)
	if err != nil {
		log.Println("failed to parse new tweet:", err)
		return tweet, err
	}

	return tweet, nil
}

func getTweet(token *oauth.Credentials, tweetId string) (tweet Tweet, err error) {

	apiurl := "https://api.twitter.com/1.1/statuses/show.json"

	param := make(url.Values)
	param.Set("id", tweetId)

	oauthClient.SignParam(token, "GET", apiurl, param)
	apiurl = apiurl + "?" + param.Encode()
	res, err := http.Get(apiurl)
	if err != nil {
		log.Println("Failed to Get tweet:", err)
		return tweet, err
	}

	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Println("Status Code:", res.StatusCode)
		bytes, _ := ioutil.ReadAll(res.Body)
		log.Println("Body: ", string(bytes))
		return tweet, err
	}

	err = json.NewDecoder(res.Body).Decode(&tweet)
	if err != nil {
		log.Println("failed to parse new tweet:", err)

		return tweet, err
	}

	return tweet, nil

}
