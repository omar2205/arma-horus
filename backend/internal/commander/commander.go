package commander

import (
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
)

// RunCMD executes a command (w/ optional args)
// returns the PID(-1 if failed) and error
func RunCMD(cmd string, args ...string) (int, error) {
	c := exec.Command(cmd, args...)
	err := c.Start()
	if err != nil {
		return -1, err
	}

	return c.Process.Pid, nil
}

// killCMD kills a process by PID
// returns nil if successful
func KillCMD(pid int) error {
	process, err := os.FindProcess(pid)
	if err != nil {
		return err
	}

	return process.Kill()
}

func ReadPIDFromFile(filename string) (int, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return 0, err
	}

	s := string(data)
	s = s[:len(s)-1]

	n, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}

	return n, nil
}
