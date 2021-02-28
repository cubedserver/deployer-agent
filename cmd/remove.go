package cmd

import (
	"os/user"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// removeCmd represents the remove command
var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove a system account from Deployer",
	Long:  `Remove a system account from Deployer's automatic SSH Key management.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Check the user exists on the server
		u, err := user.Lookup(username)
		if err != nil || u == nil {
			color.Red("Unable to find user `%s`. Please check the username, and re-create the user on Deployer.", username)
			return
		}

		// Load existing accounts from config
		var accounts []Account
		configErr := viper.UnmarshalKey("accounts", &accounts)

		if configErr != nil {
			panic("There was a problem setting up the user account. Please try again or contact Deployer for assistance")
		}

		// Remove the account and update the config
		// Loop over accounts and search for the username
		var updatedAccounts []Account
		for _, data := range accounts {
			if data.Username != username {
				updatedAccounts = append(updatedAccounts, Account{data.Username, data.ApiKey})
			}
		}

		// Save the updated accounts list
		viper.Set("accounts", updatedAccounts)
		viper.WriteConfig()
		color.Green("\nThe selected account has been removed from Deployer.")
		color.Green("\nThe authorized_keys file has been left in tact to allow you to manually update it.")
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)

	// User flag
	removeCmd.Flags().StringVarP(&username, "username", "u", "", "The username of the system account to add to Deployer")
	removeCmd.MarkFlagRequired("user")
}
