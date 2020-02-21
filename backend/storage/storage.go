package storage

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/minio/minio-go"
)

type S3Config struct {
	Region    string `json:"region,omitempty"`
	AccessKey string `json:"access_key,omitempty"`
	SecretKey string `json:"secret_key,omitempty"`
	Endpoint  string `json:"endpoint,omitempty"`
}

var (
	DocumentBucket = "documents"
)

var (
	S3 *minio.Client
)

func InitStorage() {
	secretPath := "/config/s3.json"

	if altp, ok := os.LookupEnv("S3_SECRETS"); ok {
		secretPath = altp
	}

	f, err := ioutil.ReadFile(secretPath)

	if err != nil {
		panic(err)
	}

	var config S3Config

	err = json.Unmarshal(f, &config)

	if err != nil {
		panic(err)
	}

	client, err := minio.NewWithRegion(
		config.Endpoint,
		config.AccessKey,
		config.SecretKey,
		true,
		config.Region,
	)

	if err != nil {
		panic(err)
	}

	S3 = client

	if bucket, ok := os.LookupEnv("S3_DOCUMENT_BUCKET"); ok {
		DocumentBucket = bucket
	}
}
