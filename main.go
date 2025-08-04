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

func init() {
	if runtime.GOOS == "android" && os.Getenv("TERMUX_VERSION") != "" {
		// An Android, using Termux, the os.Args property contains an additional argument at the front.
		// It looks like it's the path the user used in the CLI, so we remove it to adhere to convention.
		if len(os.Args) > 1 {
			os.Args = os.Args[1:]
		}

		prefix := os.Getenv("PREFIX")
		if prefix == "" {
			prefix = "/data/data/com.termux/files/usr"
		}

		execPath := os.Args[0]

		homePath, err := os.UserHomeDir()
		if err != nil {
			fmt.Printf("Failed to get home dir path: %v\n", err)
			os.Exit(EXIT_NO_HOME)
		}
		scriptPath := filepath.Join(homePath, filepath.Base(execPath)+".proot")
		startScript := flo.File(scriptPath)

		if !startScript.Exists() {
			prootPath, err := exec.LookPath("proot")
			if err != nil {
				fmt.Printf("`proot` not found - install with `pkg install proot`\n")
				os.Exit(EXIT_NO_PROOT)
			}

			if err := startScript.StoreString(fmt.Sprintf("#!/bin/bash\n%s -b $PREFIX/etc/resolv.conf:/etc/resolv.conf %s $@\n", prootPath, execPath)); err != nil {
				fmt.Printf("Could not create start script: %s.\n", err.Error())
				os.Exit(EXIT_SCRIPT_CREATE_FAILED)
			}

			startScript.PermExec(true, false, false)

			fmt.Printf("\nCreated a start script for you. Run the app like this:\n %s\n", strings.Join(append([]string{startScript.Name()}, os.Args[1:]...), " "))
			os.Exit(0)
		}
	}
}
