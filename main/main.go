package main

import (
	"embed"
	"encoding/json"
	"flag"
	"io/fs"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"

	"github.com/SherClockHolmes/webpush-go"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

// go:embed static
var embedStaticFs embed.FS

type GotifyMessage struct {
	Id       int
	AppId    int
	Date     string
	Priority int
	Title    string
	Message  string
	Extras   struct {
		Url string
	}
}

type Config struct {
	Http struct {
		ListenAddress string
	}
	Gotify struct {
		Address string
	}
	Vapid struct {
		PublicKey  string
		PrivateKey string
	}
	Subscriber []*webpush.Subscription
}

type WebPushMessage struct {
	Title string
	Body  string
	URL   string
}

var (
	logger        = logrus.New()
	listenAddress = "127.0.0.1:3000"
	configPath    = "config.json"
	config        = &Config{}
)

func init() {
	flag.StringVar(&configPath, "config", "config.json", "path of the configuration file")
	logger.SetLevel(logrus.DebugLevel)
}

func generateVapidKeyPair() {
	logger.Info("generating VAPID keypair")
	vapidPrivateKey, vapidPublicKey, err := webpush.GenerateVAPIDKeys()
	config.Vapid.PublicKey = vapidPublicKey
	config.Vapid.PrivateKey = vapidPrivateKey
	if err != nil {
		logger.Fatal("failed to generate VAPID keypair ", err)
	}
}

func loadConfig() (err error) {
	logger.Infof("load config file from \"%s\" ", configPath)
	binary, err := ioutil.ReadFile(configPath)
	if err != nil {
		return
	}
	err = json.Unmarshal(binary, config)
	return
}

func saveConfig() (err error) {
	binary, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return
	}
	err = ioutil.WriteFile(configPath, binary, 0644)
	return
}

func addSubscriber(subscription *webpush.Subscription) {
	for _, s := range config.Subscriber {
		if reflect.DeepEqual(s, subscription) {
			logger.Infof("subscription already exists, skipped. %s", subscription.Keys.Auth[0:5])
			return
		}
	}
	logger.Info("add subscriber: ", subscription.Keys.Auth[0:5])
	config.Subscriber = append(config.Subscriber, subscription)
	saveConfig()
}

func createAndListenWebsocket(addr string, messageChan chan *GotifyMessage) (err error) {
	logger.Debug("try to connect ", addr)
	conn, _, err := websocket.DefaultDialer.Dial(addr, nil)
	if err != nil {
		logger.Fatalf("failed to connect %s - %s", addr, err)
		return err
	}
	logger.Debug("start receive message")
	for {
		_, respBinary, err := conn.ReadMessage()
		if err != nil {
			logger.Fatal("failed to read message from websocket connection", err)
			continue
		}
		msg := &GotifyMessage{}
		json.Unmarshal(respBinary, msg)
		logger.Infof("<-recv message %+v", msg)
		messageChan <- msg
	}
}

func sendWebPush(messageChan chan *GotifyMessage) {
	logger.Debug("start listen messages")
	for {
		msg := <-messageChan
		pushMessage := &WebPushMessage{
			Title: msg.Title,
			Body:  msg.Message,
			URL:   msg.Extras.Url,
		}
		for _, subscription := range config.Subscriber {
			logger.Infof("->send message %+v to %s", pushMessage, subscription.Keys.Auth[0:5])
			messageBinary, _ := json.Marshal(pushMessage)

			_, err := webpush.SendNotification(messageBinary, subscription, &webpush.Options{
				VAPIDPublicKey:  config.Vapid.PublicKey,
				VAPIDPrivateKey: config.Vapid.PrivateKey,
			})
			if err != nil {
				logger.Error("failed to push message to server ", err)
			} else {
				logger.Info("send message successfully")
			}
		}
	}
}

func initWebServer() {
	var staticFs http.FileSystem
	_, err := os.Stat("static")
	if os.IsNotExist(err) {
		fsys, _ := fs.Sub(embedStaticFs, "static")
		staticFs = http.FS(fsys)
	} else {
		logger.Debug("webserver is currently use static folder because it exists")
		staticFs = http.Dir("static")
	}

	http.Handle("/", http.FileServer(staticFs))
	http.HandleFunc("/api/getPublicKey", func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Add("Content-Type", "application/json")
		logger.Debug("GET /api/getPublicKey")
		rw.Write([]byte(config.Vapid.PublicKey))
	})
	http.HandleFunc("/api/subscribe", func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Add("Content-Type", "application/json")
		reqBodyBinary, err := ioutil.ReadAll(r.Body)
		if err != nil {
			logger.Error("failed to read data from request ", err)
			rw.WriteHeader(http.StatusBadRequest)
			return
		}
		subscription := &webpush.Subscription{}
		err = json.Unmarshal(reqBodyBinary, subscription)
		if err != nil {
			logger.Error("failed to parse json from request ", err)
			rw.WriteHeader(http.StatusBadRequest)
			return
		}
		addSubscriber(subscription)
		rw.WriteHeader(http.StatusNoContent)

	})

}

func main() {

	flag.Parse()

	logger.Info("Gotify Bridge v1.0 by techmoe (https://github.com/dev-techmoe/gotify-bridge)")
	initWebServer()

	// check config file exists
	// if not, generate a new VAPID keypair and save it
	_, err := os.Stat(configPath)
	if os.IsNotExist(err) {
		generateVapidKeyPair()
		saveConfig()
	}
	// load config
	loadConfig()

	messageChan := make(chan *GotifyMessage)
	go createAndListenWebsocket(config.Gotify.Address, messageChan)
	go sendWebPush(messageChan)

	logger.Infof("webserver listen on %s", listenAddress)
	err = http.ListenAndServe(config.Http.ListenAddress, nil)
	if err != nil {
		logger.Fatalf("failed to listen %s (%s)", listenAddress, err)
	}
}
