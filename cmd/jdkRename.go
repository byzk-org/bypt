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
	"github.com/byzk-org/bypt/socket"
	"github.com/byzk-org/bypt/tools"
	"github.com/spf13/cobra"
)

// jdkRenameCmd represents the jdkRename command
var jdkRenameCmd = &cobra.Command{
	Use:   "rename srcName distName",
	Short: "重命名",
	Long:  `修改现已存在的jdk的名称`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		conn := socket.GetClientConn()
		defer conn.SendEndMsg()

		conn.WriteDataStr(cmdJdkRename).Wait().WriteDataStr(args[0]).WriteDataStr(args[1]).Wait()

		tools.SuccessOut("修改成功")

	},
}

func init() {
	jdkCmd.AddCommand(jdkRenameCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// jdkRenameCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// jdkRenameCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
