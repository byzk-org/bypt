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
	"encoding/json"
	"fmt"
	"github.com/gosuri/uitable"
	"github.com/byzk-org/bypt/socket"
	"github.com/byzk-org/bypt/tools"
	"github.com/byzk-org/bypt/vos"
	"github.com/spf13/cobra"
)

// configListCmd represents the configList command
var configListCmd = &cobra.Command{
	Use:   "ls",
	Short: "查看所有可配置项",
	Long:  `查看配置列表`,
	Run: func(cmd *cobra.Command, args []string) {
		conn := socket.GetClientConn()
		conn.WriteDataStr(cmdConfigList)
		conn.ReadMsgStr()
		listJsonInfo := conn.ReadMsg()

		configList := make([]vos.Setting, 0)
		if err := json.Unmarshal(listJsonInfo, &configList); err != nil {
			tools.ErrOut("解析响应数据失败")
		}

		table := uitable.New()
		table.MaxColWidth = 100
		table.Wrap = false
		table.AddRow("名称", "值", "是否需要停止所有应用", "描述")

		for _, config := range configList {
			stopApp := "是"
			if !config.StopApp {
				stopApp = "否"
			}
			table.AddRow(config.Name, config.Val, stopApp, config.Desc)
		}

		fmt.Println(table)

		conn.SendEndMsg()
	},
}

func init() {
	configCmd.AddCommand(configListCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// configListCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// configListCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
