package main

import (
	"workshop-builder/build"
	"workshop-builder/clean"
	"workshop-builder/initialize"
	"workshop-builder/serve"
	"workshop-builder/version"

	"github.com/spf13/cobra"
)

func main() {
	var cmdBuild = &cobra.Command{
		Use:   "build",
		Short: "Build the DSCDA Workshop",
		Long:  `build is for building a workshop based off the base DSCDA template, and the configuration provided.`,
		Run: func(cmd *cobra.Command, args []string) {
			build.BuildCmd()
		},
	}
	var cmdServe = &cobra.Command{
		Use:   "serve",
		Short: "Serve the DSCDA Workshop http://localhost:1313",
		Long:  `serve uses Hugo to serve the content.  By default Hugo uses http://localhost:1313.`,
		Run: func(cmd *cobra.Command, args []string) {
			serve.ServeCmd()
		},
	}
	var cmdInit = &cobra.Command{
		Use:   "init",
		Short: "Initialize a sample config.json, and manifest.yml",
		Long:  `init bootstraps a configuration for dscda to build a workshop from, extend the config.json based on your needs. init also creates a basic cf manifest.yml for cf pushing.`,
		Run: func(cmd *cobra.Command, args []string) {
			initialize.InitCmd()
		},
	}
	var cmdClean = &cobra.Command{
		Use:   "clean",
		Short: "Clean up all dscda-builder metadata and generated folders",
		Long:  `Clean the workshop of all excess content that is not required for a dscda push. Technically this will delete both paceWorkshopContent/ and workshopGen/ folders, as well as all git metadata.`,
		Run: func(cmd *cobra.Command, args []string) {
			clean.CleanCmd()
		},
	}
	var cmdVersion = &cobra.Command{
		Use:   "version",
		Short: "dscda-builder version info",
		Long:  `Helpful when troubleshooting issues with users.`,
		Run: func(cmd *cobra.Command, args []string) {
			version.VersionCmd()
		},
	}
	var rootCmd = &cobra.Command{Use: "dscda"}
	rootCmd.AddCommand(cmdBuild)
	rootCmd.AddCommand(cmdServe)
	rootCmd.AddCommand(cmdInit)
	rootCmd.AddCommand(cmdClean)
	rootCmd.AddCommand(cmdVersion)
	rootCmd.Execute()
}
