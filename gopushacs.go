// By: Indra BW

package main

import (
	"encoding/json"

	"github.com/andelf/go-curl"
)

// REST API Push structure
type JsonPush struct {
	Alert   string `json:"alert"`
	Title   string `json:"title"`
	Vibrate bool   `json:"vibrate"`
	Sound   string `json:"sound"`
}

// Callback function for curl data, write to chan
func write_data(ptr []byte, userdata interface{}) bool {
	ch, ok := userdata.(chan string)
	if ok {
		ch <- string(ptr)
		return true // ok
	} else {
		println("ERROR!")
		return false
	}
	return false
}

func main() {
	//--- SETUP
	key := "APP_KEY"
	username := "ACS_USER"
	password := "ACS_PASSWORD"
	channel := "PUSH_CHANNEL"
	cookief := "cookie.txt"
	to_ids := "everyone"
	message := "MY_MESSAGE"
	title := "MY_TITLE"
	vibrate := true
	sound := "default"

	// Inject payload post push data
	jsonD := JsonPush{
		message,
		title,
		vibrate,
		sound}
	jsonB, _ := json.Marshal(jsonD)

	// Login post params
	postLoginData := "login=" + username + "&password=" + password

	// Initialize curl
	easy := curl.EasyInit()
	defer easy.Cleanup()

	// Silently retreive cookie please...
	silentTransfer := func(buf []byte, userdata interface{}) bool {
		return true
	}

	//--- LOGIN
	easy.Setopt(curl.OPT_URL, "https://api.cloud.appcelerator.com/v1/users/login.json?key="+key)
	easy.Setopt(curl.OPT_VERBOSE, false)
	easy.Setopt(curl.OPT_COOKIEJAR, cookief)
	easy.Setopt(curl.OPT_COOKIEFILE, cookief)
	easy.Setopt(curl.OPT_POST, 1)
	easy.Setopt(curl.OPT_FOLLOWLOCATION, 1)
	easy.Setopt(curl.OPT_TIMEOUT, 60)
	easy.Setopt(curl.OPT_POSTFIELDS, postLoginData)
	easy.Setopt(curl.OPT_WRITEFUNCTION, silentTransfer)

	if err := easy.Perform(); err != nil {
		println("ERROR: ", err.Error())
	}

	postPushData := "channel=" + channel + "&to_ids=" + to_ids + "&payload=" + string(jsonB)

	//-- SEND PUSH
	easy.Setopt(curl.OPT_URL, "https://api.cloud.appcelerator.com/v1/push_notification/notify.json?key="+key)
	easy.Setopt(curl.OPT_POSTFIELDS, postPushData)
	easy.Setopt(curl.OPT_WRITEFUNCTION, write_data)

	// make channel, print status
	ch := make(chan string)
	go func(ch chan string) {
		for {
			data := <-ch
			println(data)
		}
	}(ch)

	easy.Setopt(curl.OPT_WRITEDATA, ch)

	if err := easy.Perform(); err != nil {
		println("ERROR: ", err.Error())
	}
}
