package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type jsonSettings struct {
	Settings settings `json:"server"`
}

type settings struct {
	Port        int  `json:"port"`
	Development bool `json:"development"`
	LinkedIn    struct {
		ClientID        string `json:"client_id"`
		ClientSecret    string `json:"client_secret"`
		RedirectURLHost string `json:"redirect_url_host"`
	} `json:"linkedin"`
}

// Settings global settings variable.
var Settings settings

func init() {
	raw, err := ioutil.ReadFile("./settings.json")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	var s jsonSettings

	json.Unmarshal(raw, &s)

	Settings = s.Settings
}
