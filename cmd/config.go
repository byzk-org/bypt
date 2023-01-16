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
	"strings"

	"github.com/spf13/cobra"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "修改与查看配置",
	Long:  `可修改项的查看与配置`,
	Example: `查看配置: bypt config ls
修改配置: bypt config -w key=val`,
	Run: func(cmd *cobra.Command, args []string) {
		str, err := cmd.LocalFlags().GetString("write")
		if err != nil || str == "" {
			_ = cmd.Help()
			return
		}

		if !strings.Contains(str, "=") {
			tools.ErrOut("要修改的格式必须为: key=value请检查您输入的格式是否正确!!!")
			return
		}

		split := strings.Split(str, "=")
		key := split[0]
		val := split[1]

		conn := socket.GetClientConn()
		conn.WriteDataStr(cmdConfigSetting).Wait().WriteDataStr(key).WriteDataStr(val).Wait()
		conn.SendEndMsg()
		tools.SuccessOut("配置修改成功")

	},
}

func init() {
	rootCmd.AddCommand(configCmd)

	configCmd.Flags().StringP("write", "w", "", "修改默认的配置格式: key=value")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// configCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// configCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
