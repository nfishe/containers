package main

import (
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strconv"
	"syscall"

	"github.com/nfishe/containers/reexec"
	utilruntime "github.com/nfishe/containers/util/runtime"

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
	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{}

	path, err := filepath.Abs("rootfs")
	if err != nil {
		log.Fatal(err)
	}

	utilruntime.Must(chroot(path))

	log.Printf("UID: %v", unix.Getuid())
	utilruntime.Must(cmd.Run())
}

func chroot(path string) error {
	if err := unix.Chdir(path); err != nil {
		return err
	}

	if err := unix.Chroot(path); err != nil {
		return err
	}

	if err := substituteUser(65532, 65532); err != nil {
		return err
	}

	return unix.Chdir("/")
}

func substituteUser(uid, gid int) error {
	if uid := unix.Geteuid(); uid == 0 {
		userent, err := user.LookupId(strconv.Itoa(uid))
		if err != nil {
			panic(err)
		}

		gid, err := strconv.Atoi(userent.Gid)
		if err != nil {
			log.Panic(err)
		}
		uid, err := strconv.Atoi(userent.Uid)
		if err != nil {
			log.Panic(err)
		}

		if err := syscall.Setgid(gid); err != nil {
			return err
		}

		if err := syscall.Setuid(uid); err != nil {
			return err
		}
	}
	return nil
}
