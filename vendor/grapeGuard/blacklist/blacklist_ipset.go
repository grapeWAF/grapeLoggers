package blacklist

import (
	"fmt"
	"os/exec"
	"strings"

	guard "grapeGuard"
)

const (
	execIpset  = `ipset`
	checkIpset = `--v`
	createHash = `create guardHash hash:ip hashsize 5120 maxelem 3200000 timeout 0`
	addIp      = `add guardHash %v timeout %v`
	removeIp   = `del guardHash %v`
)

var (
	useIPSet = false
)

func RunCommand(command string) error {
	cmd := exec.Command(execIpset, strings.Fields(command)...)
	if cmd != nil {
		return cmd.Run()
	}

	return fmt.Errorf("miss command...")
}

func checkIPSet() error {
	if guard.IsIPSet() == false {
		return nil
	}

	err := RunCommand(checkIpset)
	if err != nil {
		useIPSet = false
	}

	return err
}

func createIPSet() {
	if guard.IsIPSet() == false {
		return
	}

	RunCommand(createHash)
}

func addIpset(ip string, ttl int) error {
	if guard.IsIPSet() == false {
		return nil
	}

	return RunCommand(fmt.Sprintf(addIp, ip, ttl))
}

func removeIpset(ip string) error {
	if guard.IsIPSet() == false {
		return nil
	}

	return RunCommand(fmt.Sprintf(removeIp, ip))
}
