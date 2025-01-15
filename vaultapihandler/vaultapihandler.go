package vaultapihandler

import (
	"fmt"
	vaultApi "github.com/hashicorp/vault/api"
	"github.com/investify-tech/go-utils/log"
	"github.com/investify-tech/go-utils/must"
	"golang.org/x/crypto/ssh/terminal"
	"os"
)

type SecretEngineType string

const EnvVarNameApiToken = "VAULT_API_TOKEN"

const (
	KV1 SecretEngineType = "kv1"
	KV2                  = "kv-v2"
)

// Yet we support only one single vault endpoint
var globalVaultApiClient *vaultApi.Client = nil

func RetrieveSecretValue(
	secretEngineType SecretEngineType, secretEngineName, secretSubPath, secretName string) (string, error) {

	return RetrieveSecretValueFromEndpoint(
		"localhost", "8200", true, // default values
		secretEngineType, secretEngineName, secretSubPath, secretName)
}

func RetrieveSecretValueFromEndpoint(
	host, port string, useHttps bool,
	secretEngineType SecretEngineType, secretEngineName, secretSubPath, secretName string) (string, error) {

	vaultApiClient := getOrCreateVaultApiClient(host, port, useHttps)
	var secretFromApi *vaultApi.Secret
	var err error = nil

	secretApiPath := getSecretPathForApiRequest(secretEngineType, secretEngineName, secretSubPath)

	vaultApiClient.Logical().List(secretApiPath)

	secretFromApi, err = vaultApiClient.Logical().Read(secretApiPath)
	if err != nil {
		return "", err
	}
	if secretFromApi == nil {
		return "", fmt.Errorf("Secret '%s' seems to be not available/existing", secretApiPath)
	}

	return extractSecretFromApiResponse(secretEngineType, secretFromApi, secretName)
}

func getOrCreateVaultApiClient(host, port string, useHttps bool) *vaultApi.Client {
	if globalVaultApiClient != nil {
		return globalVaultApiClient
	}

	vaultApiTlsConfig := new(vaultApi.TLSConfig)
	vaultApiTlsConfig.Insecure = true

	vaultApiConfig := vaultApi.DefaultConfig()
	vaultApiConfig.ConfigureTLS(vaultApiTlsConfig)

	protocol := "http"
	if useHttps {
		protocol += "s"
	}

	vaultApiConfig.Address = fmt.Sprintf("%s://%s:%s", protocol, host, port)

	vaultApiClient, err := vaultApi.NewClient(vaultApiConfig)
	if err != nil {
		log.LogFatalAndQuit(err, "unable to initialize Vault client")
	}

	// Authenticate
	vaultApiClient.SetToken(retrieveVaultApiToken())
	// TODO: Establish a proper token check

	globalVaultApiClient = vaultApiClient
	log.LogInfo("Vault api client object created and cached")

	return vaultApiClient
}

func retrieveVaultApiToken() string {
	var vaultApiToken string

	vaultApiToken = os.Getenv(EnvVarNameApiToken)
	if vaultApiToken != "" {
		log.LogInfo("Retrieved vault api token from env var '%s'", EnvVarNameApiToken)
		return vaultApiToken
	}

	fmt.Print("Vault api token (silent input): ")
	input := must.AnySlice(terminal.ReadPassword(0))
	vaultApiToken = string(input)
	if vaultApiToken != "" {
		log.LogInfo("Retrieved vault api token from command line input")
		return vaultApiToken
	}

	panic("No vault api token could be retrieved")
}

func getSecretPathForApiRequest(secretEngineType SecretEngineType, secretEngineName, secretSubPath string) string {
	if secretEngineType == KV1 {
		return fmt.Sprintf("%s/%s", secretEngineName, secretSubPath)
	} else if secretEngineType == KV2 {
		return fmt.Sprintf("%s/data/%s", secretEngineName, secretSubPath)
	} else {
		panic(cannotDealWithSecretEngineMsg(secretEngineType))
	}
}

func cannotDealWithSecretEngineMsg(secretEngineType SecretEngineType) string {
	return fmt.Sprintf("Cannot deal with secret engine type '%s'", secretEngineType)
}

func extractSecretFromApiResponse(secretEngineType SecretEngineType, secret *vaultApi.Secret, secretName string) (string, error) {
	var value interface{}
	var exists bool
	var err error = nil

	if secretEngineType == KV1 {
		value, exists = secret.Data[secretName]
	} else if secretEngineType == KV2 {
		value, exists = secret.Data["data"].(map[string]interface{})[secretName]
	} else {
		panic(cannotDealWithSecretEngineMsg(secretEngineType))
	}

	if !exists {
		err = fmt.Errorf("secret '%s' could not be found in secret data map", secretName)
	}

	return value.(string), err
}
