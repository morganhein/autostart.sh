package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/morganhein/autostart-sh/pkg"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "autostart",
	Short: "Autostart autos your starts",
	Long: `Autostart.sh is a meant as a bootstrapper for *nix like environments, specifically installation of packages
and configuration/dotfile management. It's main goals are ease-of-use when configuring and running.`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		err := autostart.RunTask(ctx, "default")
		if err != nil {
			panic(err)
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
