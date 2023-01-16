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

// jdkRmCmd represents the jdkRm command
var jdkRmCmd = &cobra.Command{
	Use:   "rm [jdkName]",
	Short: "删除jdk",
	Long:  `删除现已存在的一个或所有jdk信息`,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		all, _ := cmd.LocalFlags().GetBool("all")
		if all && len(args) > 0 {
			tools.ErrOut("带有参数与名称冲突")
		}
		conn := socket.GetClientConn()
		defer conn.SendEndMsg()

		if all {
			if !tools.ScanTerminalConfirm("请您确认是否要删除所有的jdk信息,此项操作不可逆, 请您一定谨慎操作！！！！！") {
				return
			}
			conn.WriteDataStr(cmdJdkRmAll).Wait()
		} else if len(args) == 0 {
			cmd.Help()
			return
		} else {
			if !tools.ScanTerminalConfirm("请您确认是否删除 " + args[0] + " 此项操作不可逆, 请您一定谨慎操作！！！！！") {
				return
			}
			conn.WriteDataStr(cmdJdkRm).Wait().WriteDataStr(args[0]).Wait()
		}

		tools.SuccessOut("删除成功")
	},
}

func init() {
	jdkCmd.AddCommand(jdkRmCmd)
	jdkRmCmd.Flags().BoolP("all", "a", false, "删除全部")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// jdkRmCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// jdkRmCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
