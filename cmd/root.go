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
	"context"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"

	"gitee.com/autom-studio/culvert/internal/tunnel"
	logkit "gitee.com/autom-studio/go-kits/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "culvert",
	Short:   "A tool which create SSH tunnel for port forward.",
	Long:    `Creat a port forward tunnel through SSH connection.`,
	Version: tunnel.Version,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		var culvertConfig tunnel.CulvertConfig
		err := yaml.Unmarshal([]byte(tunnel.ConfigYaml), &culvertConfig)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Unmarshal config failed: ", err)
			os.Exit(1)
		}

		localPortVar := viper.GetString("port")
		localPort, err := strconv.Atoi(localPortVar)
		if err != nil {
			localPort = 0
		}

		logPath := culvertConfig.Log.Path
		logLevel := culvertConfig.Log.Level
		logger := logkit.InitLogger(logPath, logLevel)
		defer logger.Sync()
		undo := zap.ReplaceGlobals(logger)
		defer undo()

		ctx, cancel := context.WithCancel(context.Background())
		go func() {
			sigc := make(chan os.Signal, 1)
			signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM)
			logger.Sugar().Infof("received %v - initiating shutdown", <-sigc)
			cancel()
		}()

		var wg sync.WaitGroup
		logger.Info("tunnel created")
		defer logger.Info("tunnel closed")
		for _, t := range culvertConfig.Tunnels {
			wg.Add(1)
			if localPort != 0 {
				t.Local.Port = localPort
			}
			go t.BindTunnel(ctx, &wg)
		}
		wg.Wait()
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

	rootCmd.PersistentFlags().StringP("port", "p", "", "local bind port, default is nil use the value bound in programer(view via '-v')")
	viper.BindPFlag("port", rootCmd.PersistentFlags().Lookup("port"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {}
