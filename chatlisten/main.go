package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	messenger "github.com/maciekmm/messenger-platform-go-sdk"
	"github.com/spf13/viper"
)

var (
	appSecret   = ""
	botToken    = ""
	verify      = ""
	tlsCertFile = ""
	tlsKeyFile  = ""
	bot         *messenger.Messenger
)

func main() {
	viper.SetConfigName("bot")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Printf("Read config error: %s\n", err)
	}
	viper.AutomaticEnv()

	appSecret = viper.GetString("APP_SECRET")
	botToken = viper.GetString("TOKEN")
	verify = viper.GetString("VERIFY")
	tlsCertFile = viper.GetString("TLS_CERT")
	tlsKeyFile = viper.GetString("TLS_KEY")

	bot = &messenger.Messenger{
		VerifyToken: verify,
		AppSecret:   appSecret,
		AccessToken: botToken,
	}
	bot.MessageReceived = onMessageReceived

	http.HandleFunc("/webhook", bot.Handler)
	log.Fatal(http.ListenAndServeTLS(":443", tlsCertFile, tlsKeyFile, nil))
}

func onMessageReceived(e messenger.Event, opts messenger.MessageOpts, msg messenger.ReceivedMessage) {
	profile, err := bot.GetProfile(opts.Sender.ID)
	if err != nil {
		log.Fatal(err)
	}
	out := Out{e, opts, msg, *profile}
	b, err := json.Marshal(out)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(b[:]))
}

type Out struct {
	E       messenger.Event           `json:"event"`
	Opts    messenger.MessageOpts     `json:"opts"`
	Msg     messenger.ReceivedMessage `json:"message"`
	Profile messenger.Profile         `json:"profile"`
}