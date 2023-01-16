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
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/gosuri/uitable"
	"github.com/byzk-org/bypt/socket"
	"github.com/byzk-org/bypt/tools"
	"github.com/byzk-org/bypt/vos"
	"github.com/spf13/cobra"
)

// jdkLSCmd represents the jdkLS command
var jdkLSCmd = &cobra.Command{
	Use:   "ls [jdkName]",
	Short: "查看jdk相关信息",
	Long:  `查看现有jdk列表或jdk单个详情`,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 1 {
			jdkLsName(args[0])
			return
		}
		conn := socket.GetClientConn()
		defer conn.SendEndMsg()

		msg := conn.WriteDataStr(cmdJdkLs).Wait().ReadMsg()
		jdkList := make([]*vos.JdkInfo, 0)
		if err := json.Unmarshal(msg, &jdkList); err != nil {
			tools.ErrOut("解析数据结构失败")
		}

		table := uitable.New()
		table.AddRow("jdk名称", "jdk描述", "创建时间")
		for _, jdkInfo := range jdkList {
			table.AddRow(jdkInfo.Name, jdkInfo.Desc, jdkInfo.CreateTime.Format("2006-01-02 15:04:05"))
		}
		fmt.Println(table)
	},
}

func jdkLsName(name string) {
	conn := socket.GetClientConn()
	defer conn.SendEndMsg()
	msg := conn.WriteDataStr(cmdJdkLsName).Wait().WriteDataStr(name).ReadMsg()
	jdkInfo := &vos.JdkInfo{}
	if err := json.Unmarshal(msg, &jdkInfo); err != nil {
		tools.ErrOut("解析jdk信息失败")
	}

	table := uitable.New()
	table.Wrap = true
	table.AddRow("jdk名称: ", jdkInfo.Name)
	table.AddRow("jdk描述: ", jdkInfo.Desc)
	table.AddRow("MD5摘要: ", hex.EncodeToString(jdkInfo.MD5))
	table.AddRow("SHA1摘要: ", hex.EncodeToString(jdkInfo.SHA1))
	table.AddRow("创建时间: ", jdkInfo.CreateTime.Format("2006-01-02 15:04:05"))
	table.AddRow("最后更新时间: ", jdkInfo.EndUpdateTime.Format("2006-01-02 15:04:05"))
	fmt.Println(table)
}

func init() {
	jdkCmd.AddCommand(jdkLSCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// jdkLSCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// jdkLSCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
