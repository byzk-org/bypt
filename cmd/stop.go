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
	"fmt"
	"github.com/byzk-org/bypt/socket"
	"github.com/byzk-org/bypt/tools"
	"github.com/spf13/cobra"
	"github.com/ttacon/chalk"
	"os"
	"path/filepath"
	"strings"
)

// stopCmd represents the stop command
var stopCmd = &cobra.Command{
	Use:   "stop appName",
	Short: "停止应用",
	Long:  `停止一个正在进行中的应用`,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		withStartConfig, _ := cmd.LocalFlags().GetString("withStartConfig")
		if withStartConfig != "" {
			stopWithConfig(withStartConfig)
			return
		}

		if len(args) == 0 {
			cmd.Help()
			return
		}
		appName := args[0]
		conn := socket.GetClientConn()
		defer conn.SendEndMsg()
		conn.WriteDataStr(cmdStop).Wait().WriteDataStr(appName).Wait()
		tools.SuccessOut("应用停止成功")
	},
}

func stopWithConfig(configPath string) {
	abs, err := filepath.Abs(configPath)
	if err != nil {
		tools.ErrOut("获取配置文件绝对路径失败")
	}
	_, err = os.Stat(abs)
	if err != nil {
		tools.ErrOut("配置文件不存在")
	}

	conn := socket.GetClientConn()
	defer conn.SendEndMsg()

	conn.WriteDataStr(cmdStopWithConfig).Wait().WriteDataStr(abs)
	for {
		str := conn.ReadMsgStr()
		if str == "!!!!!!" {
			break
		}

		if strings.HasPrefix(str, "error:") {
			fmt.Println(chalk.Red, str[6:], chalk.Reset)
		} else {
			tools.SuccessOut(str)
		}

	}

	conn.Wait()
}

func init() {
	rootCmd.AddCommand(stopCmd)

	stopCmd.Flags().StringP("withStartConfig", "w", "", "根据启动配置文件停止应用")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// stopCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// stopCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
