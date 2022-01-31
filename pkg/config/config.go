package config

import (
	"errors"
	"fmt"
	"github.com/morganhein/shoelace/pkg/manager"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

//Try to load the cfg file from the various locations using viper
//FileConfig should contain the parameters from cobra as well!
//Unmarshal config

// LoadConfig reads in config file and ENV variables if set.
func LoadConfig(cfgFile string) (manager.FileConfig, error) {
	cfg := manager.FileConfig{}
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
		err := tryLoadConfig()
		if err != nil {
			return
		}
	}
	//home dir first
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)

	viper.AddConfigPath(filepath.Join(home, ".config/shoelace/"))
	viper.AddConfigPath(filepath.Join(home, ".shoelace/"))
	viper.SetConfigType("toml")
	viper.SetConfigName("config")

	if loaded := tryLoadConfig(); loaded {
		return nil
	}

	//then check defaults
	viper.AddConfigPath(filepath.Join(home, "/usr/share/shoelace/"))
	viper.SetConfigType("toml")
	viper.SetConfigName("default")

	if loaded := tryLoadConfig(); loaded {
		return nil
	}

	//do we need to check that the appropriate information wasn't provided via Environment variables, instead of erroring out here?
	return errors.New("could not load a config file")
}

func tryLoadConfig(rc manager.RunningConfig) error {
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}
	c.
		_, err = fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	return nil
}
