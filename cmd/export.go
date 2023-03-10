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
	"path/filepath"
)

// exportCmd represents the export command
var exportCmd = &cobra.Command{
	Use:   "export path",
	Short: "导出启动配置",
	Long:  `导出现服务中所有存在的启动配置信息`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		s := args[0]
		dir, err := filepath.Abs(s)
		if err != nil {
			tools.ErrOut("获取保存目录失败")
		}
		conn := socket.GetClientConn()
		defer conn.SendEndMsg()

		str := conn.WriteDataStr(cmdExport).Wait().WriteDataStr(dir).ReadMsgStr()
		tools.SuccessOut("文件导出至 => " + str)
		conn.Wait()
	},
}

func init() {
	rootCmd.AddCommand(exportCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// exportCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// exportCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
