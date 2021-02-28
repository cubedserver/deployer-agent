package cmd

import (
	"log"
	"os"
	"os/user"
	"strconv"

	"github.com/spf13/viper"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var username string
var apikey string

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a system account to Deployer",
	Long:  `Add a new system account (e.g root) to Deployer to have it's SSH Keys automatically managed.`,
	Run: func(cmd *cobra.Command, args []string) {

		// Check the user exists on the server
		u, err := user.Lookup(username)
		if err != nil {
			color.Red("Unable to find user `%s`. Please check the username, and re-create the user on Deployer.", username)
			os.Exit(1)
		}

		color.Green("Found system user: %s\nSetting up Deployer for the account.", u.Username)

		// Read in the existing accounts and get ready for adding another user
		viper.ReadInConfig()

		var accounts []Account
		configErr := viper.UnmarshalKey("accounts", &accounts)

		if configErr != nil {
			panic("There was a problem setting up the user account. Please try again or contact Deployer for assistance")
		}

		for _, data := range accounts {
			if data.Username == username {
				color.Green("The user %s is already configured - no changes needed.", data.Username)
				os.Exit(1)
			}
		}

		// Append the new account and update the config
		accounts = append(accounts, Account{username, apikey})
		viper.Set("accounts", accounts)
		viper.WriteConfig()

		// Work out the uid and gid for chowning
		uid, _ := strconv.Atoi(u.Uid)
		gid, _ := strconv.Atoi(u.Gid)

		// Set up our path and file vars
		homeDir := u.HomeDir
		keysDir := homeDir + "/.ssh"
		keysFile := keysDir + "/authorized_keys"
		backupKeysFile := keysDir + "/authorized_keys.bak"

		// If the .ssh directory doesnt exist, create it and set it to be owned by the user
		if _, keysDirErr := os.Stat(keysDir); os.IsNotExist(keysDirErr) {
			color.Yellow("It looks like " + keysDir + " does not yet exist. Lets create it now.")
			os.MkdirAll(keysDir, 0700)
			os.Chown(keysDir, uid, gid)
		}

		// If we find an authorized_keys file, move it to a backup location so we dont overwrite it.
		if _, keysFileErr := os.Stat(keysFile); !os.IsNotExist(keysFileErr) {
			// File exists, move it to its backup location
			backupFileErr := os.Rename(keysFile, backupKeysFile)
			if backupFileErr != nil {
				log.Fatal(backupFileErr)
			}
			color.Yellow("An existing authorized_keys file was found. This has been moved to %s", backupKeysFile)

		}

		// Create our authorized_keys file
		file, fileError := os.Create(keysFile)
		if fileError != nil {
			panic("Unable to create authorized_keys file. Please check that the user you are running the agent as has the correct privieges.")
		}

		// Write the template to the file
		file.Write(keysFileTemplate)

		// Ensure the file is owned by the correct user
		os.Chown(keysFile, uid, gid)

		color.Green("The user was successfully configured and is now managed by Deployer.")
	},
}

func init() {
	rootCmd.AddCommand(addCmd)

	// User flag
	addCmd.Flags().StringVarP(&username, "username", "u", "", "The username of the system account to add to Deployer")
	addCmd.MarkFlagRequired("user")

	// API Key Flag
	addCmd.Flags().StringVarP(&apikey, "apikey", "k", "", "The unique API Key for the system account, provided when adding the account via your Deployer control panel.")
	addCmd.MarkFlagRequired("api-key")
}
