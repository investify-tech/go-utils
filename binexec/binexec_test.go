package binexec

import "testing"

func TestCheckRequiredBinaries(t *testing.T) {
	type args struct {
		requiredBinaries [][]string
	}
	tests := []struct {
		name        string
		args        args
		shouldPanic bool
	}{
		{
			name: "All binaries present and valid",
			args: args{
				requiredBinaries: [][]string{
					{"bash", "bash --version", "4.0.0"},
					{"sh", "sh --version", "3.0.0"},
				},
			},
			shouldPanic: false,
		},
		{
			name: "Binary missing",
			args: args{
				requiredBinaries: [][]string{
					{"baaaash", "baaaash --version", "4.0.0"},
				},
			},
			shouldPanic: true,
		},
		{
			name: "Binary exists but version invalid",
			args: args{
				requiredBinaries: [][]string{
					{"sh", "sh --version", "324.0.0"},
				},
			},
			shouldPanic: true,
		},
		{
			name: "No binaries required",
			args: args{
				requiredBinaries: [][]string{},
			},
			shouldPanic: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.shouldPanic {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("CheckRequiredBinaries() did not panic when it should have")
					}
				}()
			} else {
				defer func() {
					if r := recover(); r != nil {
						t.Errorf("CheckRequiredBinaries() panicked when it should not have")
					}
				}()
			}
			CheckRequiredBinaries(test.args.requiredBinaries)
		})
	}
}

func TestRunBashCommand(t *testing.T) {
	tests := []struct {
		name           string
		command        string
		expectedOutput string
		expectError    bool
	}{
		{
			name:           "Successful command",
			command:        "echo 'Hello, World!'",
			expectedOutput: "Hello, World!\n",
			expectError:    false,
		},
		{
			name:           "Failing command",
			command:        "invalidcommand",
			expectedOutput: "",
			expectError:    true,
		},
		{
			name:           "Empty command",
			command:        "",
			expectedOutput: "",
			expectError:    false,
		},
		{
			name:           "Command with valid error message",
			command:        "bash -c '>&2 echo Error message; exit 1'",
			expectedOutput: "Error message\n",
			expectError:    true,
		},
		{
			name:           "Long-running command",
			command:        "echo 'This is a test'",
			expectedOutput: "This is a test\n",
			expectError:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := RunBashCommand(tt.command)

			if tt.expectError {
				if err == nil {
					t.Errorf("RunBashCommand() expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("RunBashCommand() unexpected error: %v", err)
				}
				if output != tt.expectedOutput {
					t.Errorf("RunBashCommand() expected '%v', got '%v'", tt.expectedOutput, output)
				}
			}
		})
	}
}

// RunBinary is implicitly tested by TestCheckRequiredBinaries() above
