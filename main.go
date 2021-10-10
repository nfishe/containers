package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"

	"github.com/nfishe/containers-from-scratch/reexec"
	utilruntime "github.com/nfishe/containers-from-scratch/util/runtime"
	"github.com/opencontainers/runc/libcontainer/user"
	"golang.org/x/sys/unix"
)

func main() {
	switch os.Args[1] {
	case "run":
		run()
	case "child":
		child()
	default:
		panic("help")
	}
}

func run() {
	cmd := reexec.Command(append([]string{"child"}, os.Args[2:]...)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS |
			syscall.CLONE_NEWUSER |
			syscall.CLONE_NEWPID |
			syscall.CLONE_NEWNS,
		Unshareflags: syscall.CLONE_NEWNS,
		UidMappings: []syscall.SysProcIDMap{{
			ContainerID: 0,
			HostID:      1000,
			Size:        1,
		}},
		GidMappings: []syscall.SysProcIDMap{{
			ContainerID: 0,
			HostID:      1000,
			Size:        1,
		}},
	}

	utilruntime.Must(cmd.Run())
}

func child() {
	fmt.Printf("Running %v \n", os.Args[2:])

	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{}

	path, err := filepath.Abs("rootfs")
	if err != nil {
		log.Fatal(err)
	}

	utilruntime.Must(unix.Chdir(path))
	if uid := unix.Geteuid(); uid == 0 {
		if err := chroot(path); err != nil {
			panic(err)
		}

		userent, err := user.LookupUid(65532)
		if err != nil {
			panic(err)
		}

		utilruntime.Must(syscall.Setuid(userent.Uid))
	}
	utilruntime.Must(cmd.Run())
}

func chroot(path string) error {
	if err := unix.Chroot(path); err != nil {
		return err
	}

	return unix.Chdir("/")
}
