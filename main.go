package termuxlaunchfix

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/toxyl/flo"
)

const (
	EXIT_NO_PROOT             = 100
	EXIT_NO_HOME              = 101
	EXIT_SCRIPT_CREATE_FAILED = 102
)

func getHomeDir() string {
	homePath, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("Failed to get home dir path: %v\n", err)
		os.Exit(EXIT_NO_HOME)
	}
	return homePath
}

func getProotPath() string {
	prootPath, err := exec.LookPath("proot")
	if err != nil {
		fmt.Printf("`proot` not found - install with `pkg install proot`\n")
		os.Exit(EXIT_NO_PROOT)
	}
	return prootPath
}

func saveStartScript(f *flo.FileObj, execPath, prootPath string) {
	if err := f.StoreString(makeStartScript(execPath, prootPath)); err != nil {
		fmt.Printf("Could not create start script: %s.\n", err.Error())
		os.Exit(EXIT_SCRIPT_CREATE_FAILED)
	}
	f.PermExec(true, false, false)

	fmt.Printf("\nI created a script for you that you will need to run once to initialize everything:\n %s\n\nOnce you have done that you can simply run the app with `%s`", strings.Join(append([]string{"./$HOME/" + f.Name()}, os.Args[1:]...), " "), f.Name())
}

func init() {
	if runtime.GOOS == "android" && os.Getenv("TERMUX_VERSION") != "" {
		// An Android, using Termux, the os.Args property contains an additional argument at the front.
		// It looks like it's the path the user used in the CLI, so we remove it to adhere to convention.
		if len(os.Args) > 1 {
			os.Args = os.Args[1:]
		}

		execPath := os.Args[0]
		homePath := getHomeDir()
		scriptPath := filepath.Join(homePath, filepath.Base(execPath)+".proot")
		startScript := flo.File(scriptPath)

		if !startScript.Exists() {
			saveStartScript(startScript, execPath, getProotPath())
			os.Exit(0)
		}
	}
}
