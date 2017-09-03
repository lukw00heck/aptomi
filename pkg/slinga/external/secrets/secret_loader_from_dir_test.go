package secrets

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLoadSecrets(t *testing.T) {
	secretLoader := NewSecretLoaderFromDir("../../testdata/unittests")

	{
		secrets := secretLoader.LoadSecretsByUserID("1").Labels
		assert.Equal(t, 5, len(secrets))
		assert.Equal(t, "aliceappkey", secrets["twitterAppKey"])
		assert.Equal(t, "aliceappsecret", secrets["twitterAppSecret"])
		assert.Equal(t, "alicetokenkey", secrets["twitterTokenKey"])
		assert.Equal(t, "alicetokensecret", secrets["twitterTokenSecret"])
		assert.Equal(t, "bigsecretvalue", secrets["bigsecret"])
	}

	{
		secrets := secretLoader.LoadSecretsByUserID("2").Labels
		assert.Equal(t, 5, len(secrets))
		assert.Equal(t, "bobappkey", secrets["twitterAppKey"])
		assert.Equal(t, "bobappsecret", secrets["twitterAppSecret"])
		assert.Equal(t, "bobtokenkey", secrets["twitterTokenKey"])
		assert.Equal(t, "bobtokensecret", secrets["twitterTokenSecret"])
		assert.Equal(t, "bigsecretvalue", secrets["bigsecret"])
	}

	{
		secrets := secretLoader.LoadSecretsByUserID("3").Labels
		assert.Equal(t, 1, len(secrets))
		assert.Equal(t, "bigsecretvalue", secrets["bigsecret"])
	}
}