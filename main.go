package main

import (
	"github.com/datastax-cda/workshop-builder//serve"
	"github.com/datastax-cda/workshop-builder/build"
	"github.com/datastax-cda/workshop-builder/clean"
	"github.com/datastax-cda/workshop-builder/initialize"
	"github.com/datastax-cda/workshop-builder/version"
	"github.com/spf13/cobra"
)

func main() {
	var cmdBuild = &cobra.Command{
		Use:   "build",
		Short: "Build the PACE Workshop",
		Long:  `build is for building a workshop based off the base PACE template, and the configuration provided.`,
		Run: func(cmd *cobra.Command, args []string) {
			build.BuildCmd()
		},
	}
	var cmdServe = &cobra.Command{
		Use:   "serve",
		Short: "Serve the PACE Workshop http://localhost:1313",
		Long:  `serve uses Hugo to serve the content.  By default Hugo uses http://localhost:1313.`,
		Run: func(cmd *cobra.Command, args []string) {
			serve.ServeCmd()
		},
	}
	var cmdInit = &cobra.Command{
		Use:   "init",
		Short: "Initialize a sample config.json, and manifest.yml",
		Long:  `init bootstraps a configuration for pace to build a workshop from, extend the config.json based on your needs. init also creates a basic cf manifest.yml for cf pushing.`,
		Run: func(cmd *cobra.Command, args []string) {
			initialize.InitCmd()
		},
	}
	var cmdClean = &cobra.Command{
		Use:   "clean",
		Short: "Clean up all pace-builder metadata and generated folders",
		Long:  `Clean the workshop of all excess content that is not required for a pace push. Technically this will delete both paceWorkshopContent/ and workshopGen/ folders, as well as all git metadata.`,
		Run: func(cmd *cobra.Command, args []string) {
			clean.CleanCmd()
		},
	}
	var cmdVersion = &cobra.Command{
		Use:   "version",
		Short: "pace-builder version info",
		Long:  `Helpful when troubleshooting issues with users.`,
		Run: func(cmd *cobra.Command, args []string) {
			version.VersionCmd()
		},
	}
	var rootCmd = &cobra.Command{Use: "pace"}
	rootCmd.AddCommand(cmdBuild)
	rootCmd.AddCommand(cmdServe)
	rootCmd.AddCommand(cmdInit)
	rootCmd.AddCommand(cmdClean)
	rootCmd.AddCommand(cmdVersion)
	rootCmd.Execute()
}
