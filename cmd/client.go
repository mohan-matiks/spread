package cmd

import (
	"fmt"

	"github.com/SwishHQ/spread/cli"
	"github.com/spf13/cobra"
)

var remoteURL string
var authKey string
var appName string
var environment string
var description string
var targetVersion string
var projectDir string
var osName string
var isTypescriptProject bool
var disableMinify bool
var hermes bool

var pushBundleCmd = &cobra.Command{
	Use:   "pushBundle",
	Short: "Push a new bundle to the server",
	Run: func(cmd *cobra.Command, args []string) {
		if remoteURL == "" {
			fmt.Println("Error: --remote flag is required")
			return
		}
		if authKey == "" {
			fmt.Println("Error: --auth-key flag is required")
			return
		}
		if appName == "" {
			fmt.Println("Error: --app-name flag is required")
			return
		}
		if environment == "" {
			fmt.Println("Error: --environment flag is required")
			return
		}
		if targetVersion == "" {
			fmt.Println("Error: --target-version flag is required")
			return
		}
		if osName == "" {
			fmt.Println("Error: --os-name flag is required")
			return
		}
		cli.PushBundle(
			remoteURL,
			authKey,
			appName,
			environment,
			description,
			targetVersion,
			projectDir,
			osName,
			isTypescriptProject,
			disableMinify,
			hermes,
		)
	},
}

func init() {
	rootCmd.AddCommand(pushBundleCmd)
	pushBundleCmd.Flags().StringVarP(&remoteURL, "remote", "r", "", "API base URL (required)")
	pushBundleCmd.Flags().StringVarP(&authKey, "auth-key", "a", "", "API auth key (required)")
	pushBundleCmd.Flags().StringVarP(&appName, "app-name", "n", "", "App name (required)")
	pushBundleCmd.Flags().StringVarP(&environment, "environment", "e", "", "Environment (required)")
	pushBundleCmd.Flags().StringVarP(&targetVersion, "target-version", "t", "", "Target version (required)")
	pushBundleCmd.Flags().StringVarP(&osName, "os-name", "o", "", "OS name (required)")

	pushBundleCmd.Flags().StringVarP(&projectDir, "project-dir", "p", "", "Project directory (optional)")
	pushBundleCmd.Flags().BoolVarP(&isTypescriptProject, "is-typescript", "i", false, "Is typescript project (optional)")
	pushBundleCmd.Flags().BoolVarP(&disableMinify, "disable-minify", "m", false, "Disable minify (optional)")
	pushBundleCmd.Flags().BoolVarP(&hermes, "hermes", "z", false, "Hermes (optional)")
	pushBundleCmd.Flags().StringVarP(&description, "description", "d", "", "Description (optional)")

	pushBundleCmd.MarkFlagRequired("remote")         // Mark as required
	pushBundleCmd.MarkFlagRequired("auth-key")       // Mark as required
	pushBundleCmd.MarkFlagRequired("app-name")       // Mark as required
	pushBundleCmd.MarkFlagRequired("environment")    // Mark as required
	pushBundleCmd.MarkFlagRequired("target-version") // Mark as required
	pushBundleCmd.MarkFlagRequired("os-name")        // Mark as required
}
