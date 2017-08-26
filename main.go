package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"

	"gopkg.in/PuerkitoBio/goquery.v1"
)

func main() {
	var modemURL, password, wifiPassword string

	flag.StringVar(&modemURL, "url", "http://192.168.1.1", "The url of your modem")
	flag.StringVar(&password, "modem-password", "", "The password of your modem")
	flag.StringVar(&wifiPassword, "wifi-password", "", "The wifi password")

	flag.Parse()

	if modemURL == "" {
		log.Fatal("Missing required flag --url")
	}

	if password == "" {
		log.Fatal("Missing required flag --modem-password")
	}

	if wifiPassword == "" {
		log.Fatal("Missing required flag --wifi-password")
	}

	cookieJar, _ := cookiejar.New(nil)

	client := &http.Client{
		Jar: cookieJar,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	log.Println("Opening the login page...")

	resp, err := client.Get(modemURL + "/ui/ti/login")
	if err != nil {
		log.Fatalf("Failed to request the login page %s", err)
	}

	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		log.Fatalf("Failed to parse the html of login page %s", err)
	}

	secret, exists := doc.Find("input[type=hidden][name=nonce]").First().Attr("value")
	if !exists {
		log.Fatal("Failed to find the secret in the login page")
	}

	actionKey, exists := doc.Find("input[type=hidden][name=action__key]").First().Attr("value")
	if !exists {
		log.Fatal("Failed to find the actionkey in the login page")
	}

	// Prepare the encrypted password
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(password))

	encPass := hex.EncodeToString(mac.Sum(nil))

	log.Println("Logging in...")

	_, err = client.PostForm(modemURL+"/ui/ti/login", url.Values{
		"nonce":       {secret},
		"userPwd":     {encPass},
		"action":      {"login"},
		"action__key": {actionKey},
		"login":       {"Accedi"},
	})

	if err != nil {
		log.Fatalf("Failed to do the login %s", err)
	}

	for i := 0; i <= 1; i++ {
		if i == 1 {
			log.Println("Waiting 5 seconds for the wifi to turn off...")
			time.Sleep(5000 * time.Millisecond)
		}

		log.Println("Opening the welcome page...")

		resp, err = client.Get(modemURL + "/ui/ti/welcomePage")
		if err != nil {
			log.Fatalf("Failed to open the welcome page %s", err)
		}

		doc, err = goquery.NewDocumentFromResponse(resp)
		if err != nil {
			log.Fatalf("Failed to parse the welcome page %s", err)
		}

		actionKey, exists = doc.Find("input[type=hidden][name=action__key]").First().Attr("value")
		if !exists {
			log.Fatal("Failed to find the actionkey in the welcome page")
		}

		log.Println("Turning the wifi on/off...")

		_, err = client.PostForm(modemURL+"/ui/ti/welcomePage?radio=1", url.Values{
			"encryptionKey": {wifiPassword},
			"securityMode":  {"WPA2-Personal"},
			"action":        {"wifiSwitch"},
			"action__key":   {actionKey},
		})

		if err != nil {
			log.Fatalf("Failed to send the wifi enable/disable request %s", err)
		}
	}
}
