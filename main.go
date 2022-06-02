package main

import (
	"fmt"
	"os"

	initf "github.com/fioncat/jump-src/pkg/init"
	"github.com/fioncat/jump-src/pkg/jump"
	"github.com/fioncat/jump-src/pkg/list"
	"github.com/spf13/cobra"
)

var (
	home   string
	domain string
	ssh    bool
)

var cmd = &cobra.Command{
	Use:   "jump-src --home <src-homepath> [--complete] [--domain <domain>] [--ssh] [<target>]",
	Short: "jump to source code",
}

var completeCmd = &cobra.Command{
	Use: "complete --home <src-home>",

	Run: func(cmd *cobra.Command, args []string) {
		projs, err := list.Run(home)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		for _, proj := range projs {
			fmt.Println(proj)
		}
	},
}

var printCmd = &cobra.Command{
	Use: "print --home <src-home> <target>",

	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			cmd.Usage()
			os.Exit(1)
		}
		err := jump.Print(home, args[0])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

var createCmd = &cobra.Command{
	Use: "create --home <src-home> --domain <domain> [--ssh] <target>",

	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			cmd.Usage()
			os.Exit(1)
		}
		err := jump.Create(home, domain, ssh, args[0])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

var initCmd = &cobra.Command{
	Use: "init --domain <domain> [--ssh]",

	Run: func(cmd *cobra.Command, args []string) {
		err := initf.Run(domain, ssh)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func main() {
	cmd.PersistentFlags().StringVarP(&home, "home", "", ".", "src home path")
	cmd.PersistentFlags().StringVarP(&domain, "domain", "", "github.com", "Git domain")
	cmd.PersistentFlags().BoolVarP(&ssh, "ssh", "", false, "use ssh when cloning")

	cmd.AddCommand(completeCmd, printCmd, createCmd, initCmd)
	cmd.Execute()
}
