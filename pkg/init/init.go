package init

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

func Run(domain string, ssh bool) error {
	path := askString("Please input your domain path")
	homeCmd := askString("Please input your home command")
	searchCmd := askString("Please input your search command")
	gitUser := askString("Please input your git user name")
	gitEmail := askString("Please input your git email")

	err := mkdir(path)
	if err != nil {
		return err
	}

	err = mkdir(filepath.Join(path, "src"))
	if err != nil {
		return err
	}
	err = mkdir(filepath.Join(path, "doc"))
	if err != nil {
		return err
	}

	tmpl := template.New("init.sh")
	tmpl, err = tmpl.Parse(shellTempl)
	if err != nil {
		return err
	}

	shellPath := filepath.Join(path, "init.sh")
	file, err := os.OpenFile(shellPath, os.O_CREATE|os.O_RDWR, 0677)
	if err != nil {
		return err
	}
	defer file.Close()

	var sshOpt string
	if ssh {
		sshOpt = "--ssh"
	}

	err = tmpl.Execute(file, map[string]string{
		"Name":     homeCmd,
		"VarName":  strings.ToUpper(homeCmd),
		"Domain":   domain,
		"SSH":      sshOpt,
		"Search":   searchCmd,
		"GitUser":  gitUser,
		"GitEmail": gitEmail,
	})
	if err != nil {
		return err
	}

	fmt.Printf("Init done, add `source %s` to your .zshrc file.\n", shellPath)
	return nil
}

func askString(msg string) string {
	fmt.Printf("%s: ", msg)
	var ret string
	fmt.Scanf("%s", &ret)
	if ret == "" {
		fmt.Println("cannot be empty")
		os.Exit(1)
	}
	return ret
}

func mkdir(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return os.MkdirAll(path, 0777)
		}
		return err
	}
	if !info.IsDir() {
		return fmt.Errorf("%s is not a directory", path)
	}
	return nil
}

const shellTempl = `
export {{.VarName}}_PATH=$(dirname $0 | awk '{printf "%s/src",$1}')

function _{{.Name}}() {
	jump-src complete --home ${{.VarName}}_PATH
}

function {{.Name}}() {
	if [ $# = 0 ]; then
		cd ${{.VarName}}_PATH
		return
	fi
	jump_path=$(jump-src print --home ${{.VarName}}_PATH $1 | tail -1 | tr -d '\n')
	if [ -n "$jump_path" ] && [ -d "$jump_path" ]; then
		cd $jump_path
		return
	fi
	jump-src create --home ${{.VarName}}_PATH {{.SSH}} --domain {{.Domain}} $1
	jump_path=${{.VarName}}_PATH/$1
	if [ -d "$jump_path" ]; then
		cd $jump_path
		cd .git
		git config user.name "{{.GitUser}}"
		git config user.email "{{.GitEmail}}"
		cd ..
		return
	fi
}
complete -F _{{.Name}} -o filenames {{.Name}}

function {{.Search}}() {
	target=$(jump-src complete --home ${{.VarName}}_PATH | fzf | tr -d '\n')
	if [ -z "$target" ]; then
		return
	fi
	jump_path=${{.VarName}}_PATH/$target
	if [ -d "$jump_path" ]; then
		cd $jump_path
		return
	fi
}
`
