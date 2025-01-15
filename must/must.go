package must

import (
	"os"
)

func must(possibleError error) {
	if possibleError != nil {
		panic(possibleError)
	}
}

func Void(possibleError error) {
	must(possibleError)
}

func AnyType[T any](result T, possibleError error) T {
	must(possibleError)
	return result
}

func Any(result any, possibleError error) any {
	must(possibleError)
	return result
}

// String is more explicitly telling than Any to expect a string
func String(result string, possibleError error) string {
	must(possibleError)
	return result
}

func AnySlice[T any](result []T, possibleError error) []T {
	must(possibleError)
	return result
}

func EnvVarValue(envVarName string) string {
	envVarValue, envVarPresent := os.LookupEnv(envVarName)
	if !envVarPresent {
		errorMessage := "Required environment variable '" + envVarName + "' not set"
		//fmt.Printf("An error occurred: %s\n", errorMessage)
		panic(errorMessage)
	}
	if envVarValue == "" {
		errorMessage := "Required environment variable '" + envVarName + "' is empty"
		//fmt.Printf("An error occurred: %s\n", errorMessage)
		panic(errorMessage)
	}
	return envVarValue
}
