package main

import (
	"testing"

	"github.com/SherClockHolmes/webpush-go"
	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	configPath = "testdata/config.json"
	err := loadConfig()
	assert.Equal(t, config.Vapid.PublicKey, "PublicKeyPublicKeyPublicKeyPublicKeyPublicKey")
	assert.Equal(t, config.Vapid.PrivateKey, "PrivateKeyPrivateKeyPrivateKeyPrivateKeyPrivateKey")
	assert.Nil(t, err)
}

func TestSaveConfig(t *testing.T) {
	configPath = "testdata/config_write.json"
	err := saveConfig()
	assert.Nil(t, err)
}

func TestAddSubscriber(t *testing.T) {
	addSubscriber(&webpush.Subscription{
		Endpoint: "endpoint",
		Keys: webpush.Keys{
			P256dh: "abcaaaaaaaaaa",
			Auth:   "autaaaaaaaaaaaaaaah",
		},
	})

	assert.Equal(t, config.Subscriber[0].Endpoint, "endpoint")
}
