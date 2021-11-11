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
	config          string
	sourceDir       string
	targetDir       string
	sudo            string
	forcedInstaller string
	verbose         bool
	dryRun          bool
)

func init() {
	//TODO (@morgan): config should have some default sane value if missing, or some kind of detection
	rootCmd.PersistentFlags().StringVar(&config, "config", "", "the path to the configuration")
	rootCmd.PersistentFlags().StringVar(&sourceDir, "source", "", "the base location of the source files to symlink against")
	rootCmd.PersistentFlags().StringVar(&targetDir, "target", "", "the base location to link the source files into")
	rootCmd.PersistentFlags().StringVar(&sudo, "sudo", "", "force enable/disable sudo, overwriting installer/package defaults")
	rootCmd.PersistentFlags().StringVar(&forcedInstaller, "force_installer", "", "force the specified installer to be used, if found")
	rootCmd.PersistentFlags().BoolVarP(&dryRun, "dry", "d", false, "spit out install commands, don't actually run them")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "enable verbose mode")

	// required flags
	_ = rootCmd.MarkPersistentFlagRequired("config")
}

var rootCmd = &cobra.Command{
	Use:   "autostart [task]",
	Short: "Autostart autos your starts",
	Long: `Autostart.sh is a meant as a bootstrapper for *nix like environments, specifically installation of packages
and configuration/dotfile management. It's main goals are ease-of-use when configuring and running.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		task := args[0]
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
		parsedConfig.Verbose = verbose
		if sudo != "" {
			parsedConfig.Sudo = sudo
		}
		if forcedInstaller != "" {
			parsedConfig.ForceInstaller = forcedInstaller
		}
		err = autostart.Start(ctx, *parsedConfig, task)
		if err != nil {
			fmt.Printf("FATAL: %+v\n", err)
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
