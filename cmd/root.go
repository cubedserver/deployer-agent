package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/spf13/viper"
)

var cfgFile string

type Account struct {
	Username string
	ApiKey   string
}

var keysFileTemplate = []byte(`# This file is managed by Deployer.\n
# Any changes made will be overwritten.\n
# If you do not have an account please contact the server owner for assistance.`)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "deployer",
	Short: "Deployer Server Agent",
	Long:  `The Deployer Server Agent is an easy to use command line application, allowing your server to automatically sync your teams SSH keys.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.AddConfigPath("/etc/deployer")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	// if err := viper.ReadInConfig(); err != nil {
	// 	color.Red("Your server is not configured to use Deployer!")
	// 	fmt.Println("Please follow the instructions on your server details page inside your Deployer account, or contact us for assistance.", viper.ConfigFileUsed())
	// }
}
