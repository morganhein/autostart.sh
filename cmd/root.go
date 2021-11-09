package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	autostart "github.com/morganhein/autostart.sh/pkg"
	"github.com/spf13/cobra"
)

var (
	task   string
	config string
	dryRun bool
)

func init() {
	rootCmd.PersistentFlags().StringVar(&task, "task", "", "the task to install")
	//TODO (@morgan): config should have some default sane value if missing, or some kind of detection
	rootCmd.PersistentFlags().StringVar(&config, "config", "", "the path to the configuration")
	rootCmd.PersistentFlags().BoolVarP(&dryRun, "dry", "d", false, "spit out install commands, don't actually run them")
}

var rootCmd = &cobra.Command{
	Use:   "autostart",
	Short: "Autostart autos your starts",
	Long: `Autostart.sh is a meant as a bootstrapper for *nix like environments, specifically installation of packages
and configuration/dotfile management. It's main goals are ease-of-use when configuring and running.`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute*1)
		defer cancel()
		config, err := autostart.LoadPackageConfig(ctx, config)
		if err != nil {
			panic(err)
		}
		config.DryRun = dryRun
		err = autostart.Start(ctx, *config, task)
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
