package vaultclient

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/hashicorp/vault/api"
)

// KV2Secret is the structure for a KV2 secret
type KV2Secret struct {
	RequestID     string                 `json:"request_id"`
	LeaseID       string                 `json:"lease_id"`
	LeaseDuration int                    `json:"lease_duration"`
	Renewable     bool                   `json:"renewable"`
	Data          map[string]interface{} `json:"data"`
	Warnings      []string               `json:"warnings"`
}

// GetKV2Secret returns the contents of a KV2 secret as a JSON string
func GetKV2Secret(kv2Path, kv2Name string) (string, error) {
	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		return "", err
	}

	client.SetToken(os.Getenv("VAULT_TOKEN"))

	secret, err := client.Logical().Read(fmt.Sprintf("%s/data/%s", kv2Path, kv2Name))
	if err != nil {
		return "", err
	}

	kv2Secret := KV2Secret{
		RequestID:     secret.RequestID,
		LeaseID:       secret.LeaseID,
		LeaseDuration: secret.LeaseDuration,
		Renewable:     secret.Renewable,
		Data:          secret.Data,
		Warnings:      secret.Warnings,
	}

	kv2SecretJSON, err := json.Marshal(kv2Secret)
	if err != nil {
		return "Error", err
	}
	println(json.Marshal(kv2Secret))

	return string(kv2SecretJSON), nil
}

// RenewToken renews the Vault token every 10 seconds and logs the time remaining for each renewal
func RenewToken() {
	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		log.Fatalf("Error initializing Vault client: %s", err)
	}

	client.SetToken(os.Getenv("VAULT_TOKEN"))

	for {
		secret, err := client.Auth().Token().RenewSelf(0)
		if err != nil {
			log.Fatalf("Error renewing Vault token: %s", err)
		}

		remaining := time.Duration(secret.Auth.LeaseDuration) * time.Second
		log.Printf("Vault token renewed. Time remaining: %s", remaining)
		time.Sleep(10 * time.Second)
	}
}
