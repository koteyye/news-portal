package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	testRestAddress        = "127.0.0.1:8085"
	testUserServiceAddress = "localhost:3200"
	testLogLevel           = 0
	testDBDSN              = "postgresql://postgres:postgres@localhost:5433/news?sslmode=disable"
	testS3Address          = "127.0.0.1:9001"
	testS3KeyID            = "BTviNoT5vb3qiXEP"
	testS3SecretKey        = "b2jl46Iqc9FAH4vDoJlr0y6HlHSkLPnc"
	testSecretKey          = "supersecretkey"
)

var testCorsAllowed = []string{"http://localhost:8083"}

func Test_GetConfig(t *testing.T) {
	t.Run("get config", func(t *testing.T) {
		t.Setenv("CONFIG_PATH", "./config.json")

		cfg, err := GetConfig()
		wantCfg := &Config{
			RESTAddress:       testRestAddress,
			UserServerAddress: testUserServiceAddress,
			LogLevel:          0,
			DBDSN:             testDBDSN,
			S3Address:         testS3Address,
			S3KeyID:           testS3KeyID,
			S3SecretKey:       testS3SecretKey,
			CorsAllowed:       testCorsAllowed,
			SecretKey:         testSecretKey,
		}
		assert.NoError(t, err)
		assert.Equal(t, wantCfg, cfg)
	})
}
