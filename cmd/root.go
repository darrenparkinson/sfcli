package cmd

import (
	"fmt"
	"os"

	"github.com/darrenparkinson/sfcli/pkg/salesforce"
	"github.com/spf13/cobra"

	"github.com/spf13/viper"
)

// Config holds the configuration provided by viper
type Config struct {
	Username     string `mapstructure:"username"`
	Password     string `mapstructure:"password"`
	ClientID     string `mapstructure:"client_id"`
	ClientSecret string `mapstructure:"client_secret"`
	BaseURL      string `mapstructure:"baseurl"`
}

// App represents the running application and holds a reference to our salesforce client
type App struct {
	config Config
	sc     *salesforce.Client
}

var cfgFile string
var config Config
var app App

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "sfcli",
	Short: "Salesforce CLI Utility",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.sfcli.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".sfcli" (without extension).
		viper.AddConfigPath(".")
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".sfcli")

	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}

	viper.Unmarshal(&config)
	if config.Username == "" || config.Password == "" || config.ClientID == "" || config.ClientSecret == "" || config.BaseURL == "" {
		fmt.Fprintln(os.Stderr, "Error executing CLI: Missing required environment variables")
		os.Exit(1)
	}

	sc, err := salesforce.NewClient(config.BaseURL, config.Username, config.Password, config.ClientID, config.ClientSecret, nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error executing CLI: Problem initialising salesforce client")
		os.Exit(1)
	}
	app = App{config, sc}

}
