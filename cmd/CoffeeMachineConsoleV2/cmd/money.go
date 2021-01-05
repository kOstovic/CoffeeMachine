/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// moneyCmd represents the money command
var moneyCmd = &cobra.Command{
	Use:   "money",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("money called")
	},
}
// getAllAvailableDenominationCmd represents the getAllAvailableDenomination command
var getAllAvailableDenominationCmd = &cobra.Command{
	Use:   "getAllAvailableDenomination",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("getAllAvailableDenomination called")
	},
}
// getDenominationByNameCmd represents the getDenominationByName command
var getDenominationByNameCmd = &cobra.Command{
	Use:   "getDenominationByName",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("getDenominationByName called")
	},
}
// putDenominationCmd represents the putDenomination command
var putDenominationCmd = &cobra.Command{
	Use:   "putDenomination",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("putDenomination called")
	},
}
// putDenominationByNameCmd represents the putDenominationByName command
var putDenominationByNameCmd = &cobra.Command{
	Use:   "putDenominationByName",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("putDenominationByName called")
	},
}
// patchDenominationCmd represents the patchDenomination command
var patchDenominationCmd = &cobra.Command{
	Use:   "patchDenomination",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("patchDenomination called")
	},
}
func init() {
	rootCmd.AddCommand(moneyCmd)
	moneyCmd.AddCommand(getAllAvailableDenominationCmd)
	moneyCmd.AddCommand(getDenominationByNameCmd)
	moneyCmd.AddCommand(putDenominationCmd)
	moneyCmd.AddCommand(putDenominationByNameCmd)
	moneyCmd.AddCommand(patchDenominationCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// moneyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// moneyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
