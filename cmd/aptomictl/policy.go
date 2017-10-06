package main

import (
	"fmt"

	"github.com/Aptomi/aptomi/pkg/slinga/client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	policyCmd = &cobra.Command{
		Use:   "policy",
		Short: "policy subcommand",
		Long:  "policy subcommand long",
	}
	policyShowCmd = &cobra.Command{
		Use:   "show",
		Short: "show policy",
		Long:  "show policy long",

		Run: func(cmd *cobra.Command, args []string) {
			err := client.Show(viper.GetViper())
			if err != nil {
				panic(fmt.Sprintf("Error while showing policy: %s", err))
			}
		},
	}
	policyApplyCmd = &cobra.Command{
		Use:   "apply",
		Short: "apply policy files",
		Long:  "apply policy files long",

		Run: func(cmd *cobra.Command, args []string) {
			err := client.Apply(viper.GetViper())
			if err != nil {
				panic(fmt.Sprintf("Error while applying specified policy: %s", err))
			}
		},
	}
)

func init() {
	policyApplyCmd.Flags().StringSliceP("policyPaths", "f", make([]string, 0), "Paths to files, dirs with policy to apply")
	err := viper.BindPFlag("apply.policyPaths", policyApplyCmd.Flags().Lookup("policyPaths"))
	if err != nil {
		panic(err) // todo is it ok to panic here?
	}

	policyCmd.AddCommand(policyApplyCmd, policyShowCmd)
	aptomiCtlCmd.AddCommand(policyCmd)
}