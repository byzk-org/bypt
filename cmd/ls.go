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
	"github.com/ttacon/chalk"
)

// lsCmd represents the ls command
var lsCmd = &cobra.Command{
	Use:   "ls [appName] [appVersion]",
	Short: "查看应用",
	Long:  `查看所有应用或某个应用的详情`,
	Args:  cobra.MaximumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 1 {
			lsQueryByAppName(args[0])
			return
		}

		if len(args) == 2 {
			lsQueryByAppNameAndVersion(args[0], args[1])
			return
		}

		conn := socket.GetClientConn()
		defer conn.SendEndMsg()
		jsonBytes := conn.WriteDataStr(cmdAppList).Wait().ReadMsg()
		appInfoList := make([]vos.AppInfo, 0)
		if err := json.Unmarshal(jsonBytes, &appInfoList); err != nil {
			tools.ErrOut("转换数据结构失败")
		}
		table := uitable.New()
		table.Wrap = false
		table.AddRow("名称", "描述", "导入时间", "最后更新时间", "当前版本", "存在版本数量")
		for _, appInfo := range appInfoList {
			table.AddRow(appInfo.Name, appInfo.Desc, appInfo.CreateTime.Format("2006-01-02 15:04:05"), appInfo.EndUpdateTime.Format("2006-01-02 15:04:05"), appInfo.CurrentVersion, len(appInfo.Versions))
		}

		fmt.Println(table)
	},
}

func lsQueryByAppName(appName string) {
	conn := socket.GetClientConn()
	defer conn.SendEndMsg()

	msg := conn.WriteDataStr(cmdAppListByAppName).Wait().WriteDataStr(appName).ReadMsg()
	appInfo := &vos.AppInfo{}
	if err := json.Unmarshal(msg, &appInfo); err != nil {
		tools.ErrOut("转换数据结构失败")
	}

	table := uitable.New()
	table.Wrap = true
	table.AddRow("名称: ", appInfo.Name)
	table.AddRow("描述: ", appInfo.Desc)
	table.AddRow("导入时间: ", appInfo.CreateTime.Format("2006-01-02 15:04:05"))
	table.AddRow("最后更新时间: ", appInfo.EndUpdateTime.Format("2006-01-02 15:04:05"))
	table.AddRow("当前使用的版本: ", appInfo.CurrentVersion)
	table.AddRow()
	table.AddRow("所有版本列表:")

	versionList := uitable.New()
	versionList.Wrap = false
	versionList.AddRow("版本名称", "MD5", "SHA1", "描述", "导入时间", "最后更新时间")
	for _, versionInfo := range appInfo.Versions {
		versionList.AddRow(versionInfo.Name, hex.EncodeToString(versionInfo.ContentMd5), hex.EncodeToString(versionInfo.ContentSha1), versionInfo.Desc, versionInfo.CreateTime.Format("2006-01-02 15:04:05"), versionInfo.EndUpdateTime.Format("2006-01-02 15:04:05"))
	}

	fmt.Println(table)

	fmt.Println(versionList)

}

func lsQueryByAppNameAndVersion(appName string, appVersion string) {
	conn := socket.GetClientConn()
	defer conn.SendEndMsg()

	msg := conn.WriteDataStr(cmdAppListByAppNameAndVersion).Wait().WriteDataStr(appName).WriteDataStr(appVersion).ReadMsg()
	appInfo := &vos.AppInfo{}
	if err := json.Unmarshal(msg, appInfo); err != nil {
		tools.ErrOut("转换数据结构失败")
	}

	versionInfo := appInfo.CurrentVersionInfo

	table := uitable.New()
	table.AddRow("应用名称: ", appInfo.Name)
	table.AddRow("应用描述: ", appInfo.Desc)
	table.AddRow("应用导入时间: ", appInfo.CreateTime.Format("2006-01-02 15:04:05"))
	fmt.Println(table)
	fmt.Println()
	fmt.Println()
	versionTable := uitable.New()
	versionTable.AddRow("版本 " + appInfo.CurrentVersion + " 信息:")
	versionTable.AddRow("|--版本名称: ", versionInfo.Name)
	versionTable.AddRow("|--版本描述: ", versionInfo.Desc)
	versionTable.AddRow("|--MD5: ", hex.EncodeToString(versionInfo.ContentMd5))
	versionTable.AddRow("|--SHA1: ", hex.EncodeToString(versionInfo.ContentSha1))
	versionTable.AddRow("|--版本导入时间: ", versionInfo.CreateTime.Format("2006-01-02 15:04:05"))
	versionTable.AddRow("|--版本最后更新时间: ", versionInfo.CreateTime.Format("2006-01-02 15:04:05"))

	fmt.Println(versionTable)

	if len(versionInfo.EnvConfigInfos) > 0 {
		fmt.Println()
		versionConfigTable := uitable.New()
		versionConfigTable.AddRow("|--可用的配置信息: ")
		versionConfigTable.AddRow("配置名称", "配置描述", "默认值")
		for _, c := range versionInfo.EnvConfigInfos {
			versionConfigTable.AddRow(c.Name, c.Desc, c.DefaultVal)
		}
		fmt.Println(versionConfigTable)
	}

	if len(versionInfo.PluginInfo) > 0 {
		fmt.Println()
		fmt.Println()
		pluginTable := uitable.New()
		pluginTable.AddRow("|--插件信息:")
		fmt.Println(pluginTable)

		for _, plugin := range versionInfo.PluginInfo {
			pluginTable = uitable.New()
			pluginTable.AddRow(fmt.Sprint(chalk.Cyan, "插件名称: "+plugin.Name, chalk.Reset))
			pluginTable.AddRow("插件类别: " + plugin.Type)
			pluginTable.AddRow("插件描述: " + plugin.Desc)
			pluginTable.AddRow("插件MD5摘要: " + hex.EncodeToString(plugin.Md5))
			pluginTable.AddRow("插件SHA1摘要: " + hex.EncodeToString(plugin.Sha1))
			if len(plugin.EnvConfig) > 0 {
				pluginTable.AddRow("插件配置:")
				fmt.Println(pluginTable)
				pluginConfigTable := uitable.New()
				pluginConfigTable.AddRow("名称", "描述", "默认值")
				for _, c := range plugin.EnvConfig {
					pluginConfigTable.AddRow(c.Name, c.Desc, c.DefaultVal)
				}
				fmt.Println(pluginConfigTable)
			} else {
				fmt.Println(pluginTable)
			}
			fmt.Println()
		}
	}

}

func init() {
	rootCmd.AddCommand(lsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// lsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// lsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
