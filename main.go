package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func main () {
	switch os.Args[1] {
	case "run":
		run()
	case "child":
		child()
	default:
		panic("I no undastand")
	}
}

func run () {

	cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[2:]...)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr {
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID,
	}

	must(cmd.Run())
}

func child () {
	fmt.Printf("running %v as PID %d\n", os.Args[2:], os.Getpid())

	cg()

	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	must(syscall.Sethostname([]byte("gocontainer")))
	must(syscall.Chroot("/home/dang/containerfs"))
	must(os.Chdir("/"))
	must(syscall.Mount("proc", "proc", "proc", 0, ""))
	must(cmd.Run())
	must(syscall.Unmount("proc", 0))
}

func cg() {
	cgroups := "/sys/fs/cgroup/"
	pids := filepath.Join(cgroups, "pids")

	must(os.Mkdir(filepath.Join(pids, "dangcgroup"), 0755))
	must(ioutil.WriteFile(filepath.Join(pids, "dangcgroup/pids.max"), []byte("20"), 0700))
	//Removes the new cgroup in plact after the container exits
	must(ioutil.WriteFile(filepath.Join(pids, "dangcgroup/notify_on_release"), []byte("1"), 0700))
	must(ioutil.WriteFile(filepath.Join(pids, "dangcgroup/cgroup.procs"), []byte(strconv.Itoa(os.Getpid())), 0700))
}

func must (err error) {
	if err != nil {
	panic (err)
	}
}
