package shell

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"testing"
)

func TestCD(t *testing.T) {
	testCases := []struct {
		name     string
		argument string
	}{
		{
			name:     "Parent directory",
			argument: "..",
		}, {
			name:     "Same directory",
			argument: ".",
		}, {
			name:     "Parent directory of parent directory",
			argument: "../..",
		}, {
			name:     "Home directory",
			argument: "",
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			startDirectory, err := os.Getwd()
			if err != nil {
				t.Error(err)
				return
			}
			if err := Cd(testCase.argument); err != nil {
				t.Error(err)
				return
			}
			nowDirectory, err := os.Getwd()
			if err != nil {
				t.Error(err)
				return
			}
			var expectedDirectory string
			if testCase.argument == "" {
				expectedDirectory, err = os.UserHomeDir()
				if err != nil {
					t.Error(err)
					return
				}
			} else {
				expectedDirectory, err = filepath.Abs(filepath.Clean(startDirectory + "/" + testCase.argument))
				if err != nil {
					t.Error(err)
					return
				}
			}

			if expectedDirectory != nowDirectory {
				t.Errorf("got %s, want %s", nowDirectory, expectedDirectory)
			}
		})
	}
}

func TestEcho(t *testing.T) {
	testCases := []struct {
		name           string
		arguments      []string
		expectedOutput string
	}{
		{
			name:           "Default",
			arguments:      []string{"hello", "world"},
			expectedOutput: "hello world",
		}, {
			name:           "One argument",
			arguments:      []string{"hello"},
			expectedOutput: "hello",
		}, {
			name:           "No arguments",
			arguments:      []string{},
			expectedOutput: "",
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			got := Echo(testCase.arguments...)

			if got != testCase.expectedOutput {
				t.Errorf("got %s, want %s", got, testCase.expectedOutput)
			}
		})
	}
}

func TestPWD(t *testing.T) {
	got, err := Pwd()
	if err != nil {
		t.Error(err)
		return
	}
	expected, err := os.Getwd()
	if err != nil {
		t.Error(err)
		return
	}
	if got != expected {
		t.Errorf("got %s, want %s", got, expected)
	}
}

func TestPS(t *testing.T) {
	got, err := Ps()
	if err != nil {
		t.Error(err)
		return
	}
	if len(got) == 0 {
		t.Errorf("cannot be 0 processes")
	}
}

func TestKill(t *testing.T) {
	if err := Kill("9999999999999"); err.Error() != "OpenProcess: The parameter is incorrect." {
		t.Error(err)
	}
	cmd := exec.Command("go", "run", "../task.go")
	if err := cmd.Start(); err != nil {
		t.Error(err)
	}
	if err := Kill(strconv.Itoa(cmd.Process.Pid)); err != nil {
		t.Error(err)
	}
}

func TestExec(t *testing.T) {
	got, err := Exec("echo", strings.NewReader(""), "hello", "world")
	if err != nil {
		t.Error(err)
	}
	if got != "hello world\n" {
		t.Errorf("got %s, want hello world", got)
	}

	got, err = Exec("the_most_inaccessible_team_with_the_most_difficult_name", strings.NewReader(""), "hello", "world")
	if err == nil {
		t.Error(err)
	}
	if got != "" {
		t.Errorf("got %s when the output was not expected", got)
	}
}

func TestShell(t *testing.T) {
	testCases := []struct {
		name           string
		inputText      string
		expectedOutput string
		expectOutput   bool
		errorInOutput  bool
		expectedError  error
	}{
		{
			name:           "cd",
			inputText:      "cd ..\n\\quit\n",
			expectedOutput: "",
			expectOutput:   false,
			errorInOutput:  false,
			expectedError:  nil,
		}, {
			name:           "echo",
			inputText:      "echo test1\necho test2 test3\n\\quit\n",
			expectedOutput: "test1\ntest2 test3\n",
			expectOutput:   true,
			errorInOutput:  false,
			expectedError:  nil,
		}, {
			name:          "ps",
			inputText:     "ps\n\\quit\n",
			expectOutput:  false,
			errorInOutput: false,
			expectedError: nil,
		}, {
			name:          "pwd",
			inputText:     "pwd\n\\quit\n",
			expectOutput:  false,
			errorInOutput: false,
			expectedError: nil,
		}, {
			name:           "kill",
			inputText:      "kill -99999999\n\\quit\n",
			expectedOutput: "error: OpenProcess: The parameter is incorrect.\n",
			expectOutput:   true,
			errorInOutput:  true,
			expectedError:  nil,
		}, {
			name:           "UNIX: not enough arguments",
			inputText:      "unix\n\\quit\n",
			expectedOutput: "error: " + ErrNotEnoughArguments.Error() + "\n",
			expectOutput:   true,
			errorInOutput:  true,
			expectedError:  nil,
		}, {
			name:           "UNIX: go version",
			inputText:      "unix go version\n\\quit\n",
			expectedOutput: fmt.Sprintf("go version %s %s/%s\n", runtime.Version(), runtime.GOOS, runtime.GOARCH),
			expectOutput:   true,
			errorInOutput:  false,
			expectedError:  nil,
		}, {
			name:           "Pipe",
			inputText:      "echo test1 | grep test1\n\\quit\n",
			expectedOutput: "test1\n",
			expectOutput:   true,
			errorInOutput:  false,
			expectedError:  nil,
		}, {
			name:           "Pipe wrong string",
			inputText:      " | \n\\quit\n",
			expectedOutput: "error: " + ErrNotEnoughCommands.Error() + "\n",
			expectOutput:   true,
			errorInOutput:  true,
			expectedError:  nil,
		}, {
			name:          "Not exist command",
			inputText:     "ls\n\\quit\n",
			expectOutput:  false,
			errorInOutput: false,
			expectedError: nil,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			reader := strings.NewReader(testCase.inputText)
			var buffer bytes.Buffer
			err := Shell(reader, &buffer)
			if testCase.expectedError != nil {
				if err == nil || err.Error() != testCase.expectedError.Error() {
					t.Errorf("error: got %v, want %v", err, testCase.expectedError)
				}
			} else {
				if err != testCase.expectedError {
					t.Errorf("error: got %v, want %v", err, testCase.expectedError)
				}
			}
			got := buffer.String()
			if testCase.errorInOutput && !strings.Contains(got, "error: ") {
				t.Errorf("expecting error in shell, got %s", got)
			}
			if !testCase.errorInOutput && strings.Contains(got, "error: ") {
				t.Errorf("not expecting error in shell, got %s", got)
			}
			if testCase.expectOutput {
				if got != testCase.expectedOutput {
					t.Errorf("got %s, want %s", got, testCase.expectedOutput)
				}
			}
		})
	}
}
