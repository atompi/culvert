/*
Copyright © 2022 Atom Pi <coder.atompi@gmail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"gitee.com/autom-studio/culvert/internal/config"
	"gitee.com/autom-studio/culvert/internal/forward"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "culvert",
	Short:   "A tool which create SSH tunnel for port forward.",
	Long:    `Creat a port forward tunnel through SSH connection.`,
	Version: config.Version,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		var culvertConfig config.CulvertConfig
		err := viper.Unmarshal(&culvertConfig)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Unmarshal config failed: ", err)
			os.Exit(1)
		}

		var sshAuth []ssh.AuthMethod

		if culvertConfig.Server.Password != "" {
			sshAuth = []ssh.AuthMethod{
				ssh.Password(culvertConfig.Server.Password),
			}
		} else if culvertConfig.Server.KeyFile != "" {
			key, err := ioutil.ReadFile(culvertConfig.Server.KeyFile)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Read private key file failed:", err)
				os.Exit(1)
			}
			signer, err := ssh.ParsePrivateKey(key)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Parse private key file failed:", err)
				os.Exit(1)
			}
			sshAuth = []ssh.AuthMethod{
				ssh.PublicKeys(signer),
			}
		} else {
			fmt.Fprintln(os.Stderr, "Have no SSH authorization mathod")
			os.Exit(1)
		}

		sshClientConfig := &ssh.ClientConfig{
			User:            culvertConfig.Server.Username,
			Auth:            sshAuth,
			Timeout:         30 * time.Second,
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		}
		sshServerConfig := &culvertConfig.Server

		forward.Tunnel(sshClientConfig, sshServerConfig, culvertConfig)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ./culvert.yaml)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Search config in current directory with name "culvert" (without extension).
		viper.AddConfigPath("./")
		viper.SetConfigType("yaml")
		viper.SetConfigName("culvert")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	} else {
		fmt.Fprintln(os.Stderr, "Init config file failed:", err)
	}
}
