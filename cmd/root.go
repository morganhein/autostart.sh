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
	task      string
	config    string
	dryRun    bool
	sourceDir string
	targetDir string
)

func init() {
	rootCmd.PersistentFlags().StringVar(&task, "task", "", "the task to install")
	//TODO (@morgan): config should have some default sane value if missing, or some kind of detection
	rootCmd.PersistentFlags().StringVar(&config, "config", "", "the path to the configuration")
	rootCmd.PersistentFlags().StringVar(&sourceDir, "source", "", "the base location of the source files to symlink against")
	rootCmd.PersistentFlags().StringVar(&targetDir, "target", "", "the base location to link the source files into")
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
		parsedConfig, err := autostart.LoadPackageConfig(ctx, config)
		if err != nil {
			panic(err)
		}
		parsedConfig.ConfigLocation = config
		parsedConfig.DryRun = dryRun
		parsedConfig.SourceDir = sourceDir
		parsedConfig.TargetDir = targetDir
		err = autostart.Start(ctx, *parsedConfig, task)
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
