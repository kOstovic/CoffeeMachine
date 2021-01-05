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

// drinkCmd represents the drink command
var drinkCmd = &cobra.Command{
	Use:   "drink",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("drink called")
	},
}
// getAllAvailableDrinksCmd represents the getAllAvailableDrinks command
var getAllAvailableDrinksCmd = &cobra.Command{
	Use:   "getAllAvailableDrinks",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("getAllAvailableDrinks called")
	},
}
// getConsumeDrinkCmd represents the getConsumeDrink command
var getConsumeDrinkCmd = &cobra.Command{
	Use:   "getConsumeDrink",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("getConsumeDrink called")
	},
}
// postAddDrinkCmd represents the postAddDrink command
var postAddDrinkCmd = &cobra.Command{
	Use:   "postAddDrink",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("postAddDrink called")
	},
}
// postRemoveDrinkCmd represents the postRemoveDrink command
var postRemoveDrinkCmd = &cobra.Command{
	Use:   "postRemoveDrink",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("postRemoveDrink called")
	},
}
func init() {
	rootCmd.AddCommand(drinkCmd)
	drinkCmd.AddCommand(getAllAvailableDrinksCmd)
	drinkCmd.AddCommand(getConsumeDrinkCmd)
	drinkCmd.AddCommand(postAddDrinkCmd)
	drinkCmd.AddCommand(postRemoveDrinkCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// drinkCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// drinkCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
