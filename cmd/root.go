/*
Copyright © 2025 Benjamin Jørgensen <me@benmi.me>
*/
package cmd

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "jobbtid",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("serve called")
		serve()
	},
}

var cfgFilePath string

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFilePath, "config", "", "config file (default is $HOME/.jobbtid.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func greet(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World! %s", time.Now())
}

type serverConfig struct {
	Host string
}

func (conf serverConfig) host() string {
	return conf.Host
}

func readConfigFile() (serverConfig, error) {
	myServerConfig := serverConfig{
		Host: ":8080",
	}
	viper.AddConfigPath(cfgFilePath)
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			fmt.Println("Using default!")
			return myServerConfig, nil
		} else {
			// Config file was found but another error was produced
		}
	} else {
		host := viper.GetString("host")
		if len(host) == 0 {
			host = ":8080"
			myServerConfig.Host = host
		}
	}
	return myServerConfig, err
}

func serve() {
	myConf, err := readConfigFile()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	http.HandleFunc("/", greet)
	http.ListenAndServe(myConf.host(), nil)
}
