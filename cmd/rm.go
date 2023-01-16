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

// rmCmd represents the rm command
var rmCmd = &cobra.Command{
	Use:   "rm appName appVersion",
	Short: "删除应用",
	Long:  `删除一个或所有服务, 命令之后可以加应用名称和应用版本`,
	Args:  cobra.MaximumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		all, _ := cmd.LocalFlags().GetBool("all")
		if all {
			rmAll()
			return
		}

		if len(args) == 0 {
			_ = cmd.Help()
			return
		}

		if len(args) == 1 {
			rmByName(args[0])
			return
		}

		if !tools.ScanTerminalConfirm("是否确认删除[" + args[0] + "-" + args[1] + "]应用?") {
			return
		}

		conn := socket.GetClientConn()
		defer conn.SendEndMsg()

		conn.WriteDataStr(cmdRmByNameAndVersion).Wait().WriteDataStr(args[0]).WriteDataStr(args[1]).Wait()
		tools.SuccessOut("删除成功")

	},
}

func rmByName(name string) {
	if !tools.ScanTerminalConfirm("是否确认删除[" + name + "]应用?") {
		return
	}

	conn := socket.GetClientConn()
	defer conn.SendEndMsg()

	conn.WriteDataStr(cmdRmByName).Wait().WriteDataStr(name).Wait()
	tools.SuccessOut("删除成功")

}

func rmAll() {
	if !tools.ScanTerminalConfirm("是否确认删除所有应用?") {
		return
	}

	conn := socket.GetClientConn()
	defer conn.SendEndMsg()

	conn.WriteDataStr(cmdRmAll).Wait().Wait()
	tools.SuccessOut("删除成功")
}

func init() {
	rootCmd.AddCommand(rmCmd)
	rmCmd.Flags().BoolP("all", "a", false, "删除全部")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// rmCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// rmCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
