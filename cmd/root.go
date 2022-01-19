/*
Copyright © 2021 Morgan Hein <work@morganhe.in>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

var (
	dryRun  bool
	cfgFile string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "shoelace",
	Short: "shoelace autos your starts",
	Long: `shoelace.sh is a meant as a bootstrapper for *nix like environments, specifically installation of packages
and configuration/dotfile management. It's main goals are ease-of-use when configuring and running.`,
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().BoolVarP(&dryRun, "dry-run", "d", false, "echo commands only")
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/shoelace/config.toml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
		fmt.Printf("Using config file from %v\n", cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".shoelace.sh" (without extension).
		viper.AddConfigPath(filepath.Join(home, ".config/shoelace/"))
		viper.AddConfigPath(filepath.Join(home, ".shoelace/"))
		viper.SetConfigType("toml")
		viper.SetConfigName("config")
		fmt.Println("config file not specified, will search default locations")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		_, err = fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
		if err != nil {
			panic(err)
		}
	}
}

func loadConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
		tryLoadConfig()
		return
	}
	//home dir first
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)

	viper.AddConfigPath(filepath.Join(home, ".config/shoelace/"))
	viper.AddConfigPath(filepath.Join(home, ".shoelace/"))
	viper.SetConfigType("toml")
	viper.SetConfigName("config")

	if loaded := tryLoadConfig(); loaded {
		return
	}

	//then check defaults
	viper.AddConfigPath(filepath.Join(home, "/usr/share/shoelace/"))
	viper.SetConfigType("toml")
	viper.SetConfigName("default")

	if loaded := tryLoadConfig(); loaded {
		return
	}

	//do we need to check that the appropriate information wasn't provided via Environment variables, instead of erroring out here?
	cobra.CheckErr(errors.New("could not load a config file"))
}

func tryLoadConfig() bool {
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		return false
	}
	_, err = fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	if err != nil {
		cobra.CheckErr(err)
	}
	return true
}
