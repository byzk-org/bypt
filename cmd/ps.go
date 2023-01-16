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

// psCmd represents the ps command
var psCmd = &cobra.Command{
	Use:   "ps [appName] [pluginName]",
	Short: "查看当前正在运行中的进程",
	Long:  `查看当前所有正在运行中的应用进程`,
	Args:  cobra.MaximumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 1 {
			psOneApp(args[0])
			return
		}

		if len(args) == 2 {
			psOneAppPlugin(args[0], args[1])
			return
		}

		conn := socket.GetClientConn()
		defer conn.SendEndMsg()
		msg := conn.WriteDataStr(cmdPsList).Wait().ReadMsg()

		appList := make([]*vos.AppStatusInfo, 0)
		if err := json.Unmarshal(msg, &appList); err != nil {
			tools.ErrOut("转换结果失败")
		}
		table := uitable.New()
		table.Wrap = false
		table.AddRow("应用名称", "应用版本", "版本描述", "运行状态", "错误信息")
		for _, app := range appList {
			//statusText := fmt.Sprint(chalk.Green, "正在运行", chalk.Reset)
			//if app.HaveErr {
			//	statusText = fmt.Sprint(chalk.Red, "运行异常", chalk.Reset)
			//}
			//table.AddRow(app.Name, app.VersionStr, app.VersionInfo.Desc, statusText, app.ErrMsg)
			table.AddRow(app.Name, app.VersionStr, app.VersionInfo.Desc, app.Status, app.ErrMsg)
		}
		fmt.Println(table)
	},
}

func psOneAppPlugin(appName string, pluginName string) {
	conn := socket.GetClientConn()
	defer conn.SendEndMsg()
	msg := conn.WriteDataStr(cmdPsAppPlugin).Wait().WriteDataStr(appName).WriteDataStr(pluginName).ReadMsg()
	list := make([]*struct {
		Name   string `json:"name,omitempty"`
		Desc   string `json:"desc,omitempty"`
		Md5    []byte `json:"md5,omitempty"`
		Sha1   []byte `json:"sha1,omitempty"`
		Output []byte `json:"output,omitempty"`
	}, 0)
	err := json.Unmarshal(msg, &list)
	if err != nil {
		tools.ErrOut("解析插件数据失败")
	}

	if len(list) == 0 {
		tools.ErrOut("未找到对应插件")
	}

	for _, plugin := range list {
		table := uitable.New()
		table.AddRow("插件名称: ", plugin.Name)
		table.AddRow("插件描述: ", plugin.Desc)
		table.AddRow("插件MD5: ", hex.EncodeToString(plugin.Md5))
		table.AddRow("插件SHA1: ", hex.EncodeToString(plugin.Sha1))
		table.AddRow(fmt.Sprint(chalk.Yellow, "插件输出信息:"))
		fmt.Println(table)
		fmt.Printf("%s\n", plugin.Output)
		fmt.Println(chalk.Reset)
	}

}

func psOneApp(appName string) {
	conn := socket.GetClientConn()
	defer conn.SendEndMsg()
	msg := conn.WriteDataStr(cmdPsApp).Wait().WriteDataStr(appName).ReadMsg()

	appInfo := &vos.AppStatusInfo{}
	if err := json.Unmarshal(msg, &appInfo); err != nil {
		tools.ErrOut("转换数据结构失败")
	}

	//statusText := fmt.Sprint(chalk.Green, "正在运行", chalk.Reset)
	//if appInfo.HaveErr {
	//	statusText = fmt.Sprint(chalk.Red, "运行异常", chalk.Reset)
	//}
	table := uitable.New()
	table.Wrap = true
	table.AddRow("应用名称: ", appInfo.Name)
	table.AddRow("应用描述: ", appInfo.Desc)
	table.AddRow("版本名称: ", appInfo.VersionInfo.Name)
	table.AddRow("版本描述: ", appInfo.VersionInfo.Desc)
	table.AddRow("MD5: ", hex.EncodeToString(appInfo.VersionInfo.ContentMd5))
	table.AddRow("SHA1: ", hex.EncodeToString(appInfo.VersionInfo.ContentSha1))
	if appInfo.StartArgs.Restart != "" {
		table.AddRow("重启配置: ", appInfo.StartArgs.Restart)
	}
	table.AddRow("运行状态: ", appInfo.Status)
	table.AddRow("错误信息: ", appInfo.ErrMsg)
	table.AddRow("应用启动参数:", appInfo.StartArgs.Args)
	table.AddRow("java启动参数: ", appInfo.StartArgs.JdkArgs)
	table.AddRow("java命令路径:", appInfo.JavaCmd)

	configTable := uitable.New()
	configTable.Wrap = false
	configTable.AddRow("配置信息如下: ")
	configTable.AddRow("名称", "现在值", "描述")

	for _, config := range appInfo.StartArgs.EnvConfig {
		configTable.AddRow(config.Name, config.Val, config.Desc)
	}

	fmt.Println(table)
	fmt.Println(chalk.Yellow)
	fmt.Println(configTable)
	fmt.Println(chalk.Reset, chalk.Bold)

	if len(appInfo.VersionInfo.PluginInfo) > 0 {
		for _, p := range appInfo.VersionInfo.PluginInfo {
			pluginName := hex.EncodeToString(p.Md5) + hex.EncodeToString(p.Sha1)
			pluginOutPut := appInfo.PluginOutPutBuffer[pluginName]
			pluginTable := uitable.New()
			pluginTable.AddRow("插件名称: ", pluginName)
			pluginTable.AddRow("插件描述: ", p.Desc)
			if len(p.EnvConfig) > 0 {
				pluginTable.AddRow("插件配置: ")
			}
			fmt.Println(pluginTable)

			if len(p.EnvConfig) > 0 {
				pluginConfigEnvConfig := uitable.New()
				pluginConfigEnvConfig.AddRow("名称", "值", "描述")
				for _, pe := range p.EnvConfig {
					pluginConfigEnvConfig.AddRow(pe.Name, pe.Val, pe.Desc)
				}
				fmt.Println(pluginConfigEnvConfig)
			}

			if len(pluginOutPut) > 0 {
				pluginOutPutTable := uitable.New()
				pluginOutPutTable.AddRow(fmt.Sprint(chalk.Yellow, "插件输出信息: ", chalk.Reset))
				fmt.Println(pluginOutPutTable)
			}

			if len(pluginOutPut) > 0 {
				fmt.Print(chalk.Yellow)
				fmt.Printf("%s", pluginOutPut)
				fmt.Println(chalk.Reset)
			}
			fmt.Println()
		}
	}

}

func init() {
	rootCmd.AddCommand(psCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// psCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// psCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
