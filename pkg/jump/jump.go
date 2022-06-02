package jump

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func Create(home, domain string, ssh bool, target string) error {
	path := filepath.Join(home, target)
	_, err := os.Stat(path)
	if err == nil {
		return nil
	}

	tmp := strings.Split(target, "/")
	switch len(tmp) {
	case 1:
		confirm("Do you want to create group `%s`", target)
		err := os.MkdirAll(path, 0777)
		if err != nil {
			return err
		}
		return nil

	case 2:

	default:
		return fmt.Errorf("invalid target `%s`", target)
	}
	group, proj := tmp[0], tmp[1]

	var cloneURL string
	if ssh {
		cloneURL = fmt.Sprintf("git@%s:%s/%s.git", domain, group, proj)
	} else {
		cloneURL = fmt.Sprintf("https://%s/%s/%s.git", domain, group, proj)
	}

	confirm("Do you want to clone `%s`", cloneURL)
	execCmd("git", "clone", cloneURL, path)
	return nil
}

func Print(home, target string) error {
	path := filepath.Join(home, target)
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	if !info.IsDir() {
		return fmt.Errorf("%s is not a directory", path)
	}
	fmt.Println(path)
	return nil
}

func execCmd(name string, args ...string) {
	cmd := exec.Command(name, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		os.Exit(1)
	}
}

func confirm(layer string, args ...interface{}) {
	msg := fmt.Sprintf(layer, args...)
	fmt.Printf("%s ? (y/n) ", msg)
	var ret string
	fmt.Scanf("%s", &ret)
	if ret == "y" {
		return
	}
	os.Exit(0)
}
