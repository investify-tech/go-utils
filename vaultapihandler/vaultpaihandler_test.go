package vaultapihandler_test

import (
	"context"
	"github.com/investify-tech/go-utils/log"
	"github.com/investify-tech/go-utils/must"
	"github.com/investify-tech/go-utils/vaultapihandler"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"os"
	"strings"
	"testing"
)

const VaultRootToken = "supertoken"

func TestWithVaultContainer(test *testing.T) {
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Name:       "testcontainers-vault",
		Image:      "hashicorp/vault:1.4.3",
		WaitingFor: wait.ForLog("upgrading keys finished"),
		Env:        map[string]string{"VAULT_DEV_ROOT_TOKEN_ID": VaultRootToken},
		CapAdd:     []string{"IPC_LOCK"},
	}
	vaultContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		test.Error(err)
	}
	defer func() {
		if err := vaultContainer.Terminate(ctx); err != nil {
			test.Fatalf("failed to terminate container: %s", err.Error())
		}
	}()

	vaultInitCommands := "" +
		"export VAULT_ADDR=http://127.0.0.1:8200 && " +
		"vault login " + VaultRootToken + " && " +
		"vault secrets enable -path=sec-engine-v1 -version=1 kv && " +
		"vault kv put sec-engine-v1/secret-v1 key=\"value-v1\" &&" +
		"vault secrets enable -path=sec-engine-v2 -version=2 kv &&" +
		"vault kv put sec-engine-v2/secret-v2 key=\"value-v2\" "
	vaultContainer.Exec(ctx, []string{"sh", "-c", vaultInitCommands})

	endpointInfo := strings.Split(must.String(vaultContainer.Endpoint(ctx, "")), ":")
	vaultContainerHost, vaultContainerPort := endpointInfo[0], endpointInfo[1]
	log.LogInfo("vault test container started and initialized: http://%s:%s (reachable only during the test)",
		vaultContainerHost, vaultContainerPort)

	testCases := []struct {
		name             string
		secretEngineType vaultapihandler.SecretEngineType
		secretEngineName string
		secretSubPath    string
		secretName       string
		expectedValue    string
	}{
		{
			name:             "Secret engine v1",
			secretEngineType: vaultapihandler.KV1,
			secretEngineName: "sec-engine-v1",
			secretSubPath:    "secret-v1",
			secretName:       "key",
			expectedValue:    "value-v1",
		},
		{
			name:             "Secret engine v2",
			secretEngineType: vaultapihandler.KV2,
			secretEngineName: "sec-engine-v2",
			secretSubPath:    "secret-v2",
			secretName:       "key",
			expectedValue:    "value-v2",
		},
	}

	os.Setenv(vaultapihandler.EnvVarNameApiToken, VaultRootToken)

	for _, testCase := range testCases {
		test.Run(testCase.name, func(t *testing.T) {
			actual := must.String(vaultapihandler.RetrieveSecretValueFromEndpoint(
				vaultContainerHost, vaultContainerPort, false,
				testCase.secretEngineType, testCase.secretEngineName, testCase.secretSubPath, testCase.secretName))
			expected := testCase.expectedValue

			if !(actual == expected) {
				t.Errorf("Expected '%v' but got '%v'", expected, actual)
			}
		})
	}

}
