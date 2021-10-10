package reexec

import (
	"os"
	"os/exec"
	"path/filepath"
)

func Command(args ...string) *exec.Cmd {
	name, err := filepath.EvalSymlinks(os.Args[0])
	if err != nil {
		panic(err)
	}

	return exec.Command(name, args...)
}
