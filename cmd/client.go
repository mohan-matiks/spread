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

var releaseCmd = &cobra.Command{
	Use:   "release",
	Short: "Release a new version of the app",
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
			cli.BundleConfig{
				RemoteURL:           remoteURL,
				AuthKey:             authKey,
				AppName:             appName,
				Environment:         environment,
				Description:         description,
				TargetVersion:       targetVersion,
				ProjectDir:          projectDir,
				OSName:              osName,
				IsTypescriptProject: isTypescriptProject,
				DisableMinify:       disableMinify,
				Hermes:              hermes,
			},
		)
	},
}

func init() {
	rootCmd.AddCommand(releaseCmd)
	releaseCmd.Flags().StringVarP(&remoteURL, "remote", "r", "", "API base URL (required)")
	releaseCmd.Flags().StringVarP(&authKey, "auth-key", "a", "", "API auth key (required)")
	releaseCmd.Flags().StringVarP(&appName, "app-name", "n", "", "App name (required)")
	releaseCmd.Flags().StringVarP(&environment, "environment", "e", "", "Environment (required)")
	releaseCmd.Flags().StringVarP(&targetVersion, "target-version", "t", "", "Target version (required)")
	releaseCmd.Flags().StringVarP(&osName, "os-name", "o", "", "OS name (required)")

	releaseCmd.Flags().StringVarP(&projectDir, "project-dir", "p", "", "Project directory (optional)")
	releaseCmd.Flags().BoolVarP(&isTypescriptProject, "is-typescript", "i", false, "Is typescript project (optional)")
	releaseCmd.Flags().BoolVarP(&disableMinify, "disable-minify", "m", false, "Disable minify (optional)")
	releaseCmd.Flags().BoolVarP(&hermes, "hermes", "z", false, "Hermes (optional)")
	releaseCmd.Flags().StringVarP(&description, "description", "d", "", "Description (optional)")

	releaseCmd.MarkFlagRequired("remote")         // Mark as required
	releaseCmd.MarkFlagRequired("auth-key")       // Mark as required
	releaseCmd.MarkFlagRequired("app-name")       // Mark as required
	releaseCmd.MarkFlagRequired("environment")    // Mark as required
	releaseCmd.MarkFlagRequired("target-version") // Mark as required
	releaseCmd.MarkFlagRequired("os-name")        // Mark as required
}
