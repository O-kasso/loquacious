// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"errors"
	"log"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

// sayCmd represents the say command
var sayCmd = &cobra.Command{
	Use:   "say",
	Short: "Synthesize speech for text provided directly via CLI args",
	Long:  "E.g. `loq say Hello World",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) > 0 {
			log.Println("Loq Say activated")
			text := strings.Join(args, " ")
			macSay(text)
			return nil
		}
		return errors.New("No input text provided for speech synthesis")
	},
}

func macSay(text string) {
	cmd := exec.Command("say", "-v", "samantha", text)
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Loq Say complete.")
}

func init() {
	rootCmd.AddCommand(sayCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// sayCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// sayCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
