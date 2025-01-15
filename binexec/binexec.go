package binexec

import (
	"fmt"
	"github.com/investify-tech/go-utils/must"
	"github.com/mcuadros/go-version"
	"os/exec"
	"regexp"
	"strings"
)

// CheckRequiredBinaries verifies the availability and versions of the required binaries, panicking if any are missing
// or outdated.
func CheckRequiredBinaries(requiredBinaries [][]string) {

	missingExecutables := make(map[string]string)
	for _, binInfo := range requiredBinaries {
		if !isBinaryInValidVersionAvailable(binInfo) {
			missingExecutables[binInfo[0]] = binInfo[2]
		}
	}
	if len(missingExecutables) > 0 {
		var missingExecutablesList []string
		for key := range missingExecutables {
			missingExecutablesList = append(missingExecutablesList, key+" ("+missingExecutables[key]+")")
		}
		panic("Can't progress as one or more binaries are required, but missing: " + strings.Join(missingExecutablesList, ", "))
	}
}

func isBinaryInValidVersionAvailable(executableInfo []string) bool {
	executable := executableInfo[0]
	executableVersionCommand := executableInfo[1]
	executableExpectedVersion := executableInfo[2]

	path, _ := exec.LookPath(executable)
	if path == "" {
		return false
	}

	var semVerRegEx = regexp.MustCompile(
		// Use raw strings to avoid having to quote the backslashes.
		`([0-9]+)\.([0-9]+)\.([0-9]+)(?:-([0-9A-Za-z-]+(?:\.[0-9A-Za-z-]+)*))?(?:\+[0-9A-Za-z-]+)?`)
	actualVersion := semVerRegEx.FindString(must.String(RunBinary(executableVersionCommand)))
	return version.Compare(actualVersion, executableExpectedVersion, ">=")
}

// RunBashCommand runs a given command line within a bash shell so that it supports piping and so on
func RunBashCommand(command string) (string, error) {
	cmd := exec.Command("bash", "-c", command)

	stdoutStderr, commandError := cmd.CombinedOutput()
	commandOutput := string(stdoutStderr)

	return commandOutput, transformCommandError(commandError, commandOutput)
}

// RunBinary runs an external binary command from the code and provides the output (stdout and stderr combined)
func RunBinary(binaryCommandWithArgs string) (string, error) {
	binary := strings.Split(binaryCommandWithArgs, " ")[0]
	argumentList := []string{}
	for _, arg := range strings.Split(binaryCommandWithArgs, " ")[1:] {
		if arg != "" { // Sort out unusable '' arguments which might have been created by strings.Split()
			argumentList = append(argumentList, arg)
		}
	}

	cmd := exec.Command(binary, argumentList...)
	stdoutStderr, commandError := cmd.CombinedOutput()
	commandOutput := string(stdoutStderr)

	return commandOutput, transformCommandError(commandError, commandOutput)
}

func transformCommandError(commandError error, detailedErrormessage string) error {
	if commandError == nil {
		return commandError
	}
	return fmt.Errorf("Bash command error: '%s'\nBash command output:\n%s", commandError.Error(), detailedErrormessage)
}
