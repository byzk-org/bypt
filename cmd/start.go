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
	"github.com/byzk-org/bypt/socket"
	"github.com/byzk-org/bypt/tools"
	"github.com/byzk-org/bypt/vos"
	"github.com/spf13/cobra"
	"github.com/ttacon/chalk"
	"os"
	"path/filepath"
	"strings"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start [flags] appName",
	Short: "启动应用",
	Long:  `启动一个应用`,
	Args:  cobra.MaximumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		withStartConfig, _ := cmd.LocalFlags().GetString("withStartConfig")
		if withStartConfig != "" {
			startWithConfig(withStartConfig)
			return
		}
		if len(args) == 0 {
			cmd.Help()
			return
		}
		configs, _ := cmd.LocalFlags().GetStringSlice("envConfig")
		copyFiles, _ := cmd.LocalFlags().GetStringSlice("copyFile")
		currentVersion, _ := cmd.LocalFlags().GetString("currentVersion")
		javaPackName, _ := cmd.LocalFlags().GetString("javaPackName")
		javaCmdPath, _ := cmd.LocalFlags().GetString("javaCmdPath")
		restart, _ := cmd.LocalFlags().GetString("restart")
		saveAppSuffix, _ := cmd.LocalFlags().GetBool("saveAppSuffix")
		Xmx, _ := cmd.LocalFlags().GetString("Xmx")
		Xms, _ := cmd.LocalFlags().GetString("Xms")
		Xmn, _ := cmd.LocalFlags().GetString("Xmn")
		PermSize, _ := cmd.LocalFlags().GetString("PermSize")
		MaxPermSize, _ := cmd.LocalFlags().GetString("MaxPermSize")
		pluginEnvConfig, _ := cmd.LocalFlags().GetStringSlice("pluginEnvConfig")

		restartType := vos.AppRestartType(restart)
		if restart != "" {
			switch restartType {
			case vos.AppRestartTypeErrorAuto:
			case vos.AppRestartTypeAlways:
			case vos.AppRestartTypeErrorAny:
			case vos.AppRestartTypeErrorApp:
			case vos.AppRestartTypeErrorPlugin:
			default:
				tools.ErrOut("未知的重启类型 => " + restart)

			}
		}

		name := args[0]

		if len(args) > 1 {
			args = args[1:]
		} else {
			args = nil
		}

		configInfos := make([]vos.AppConfig, 0)
		for _, configStr := range configs {
			split := strings.Split(configStr, "=")
			if len(split) != 2 {
				tools.ErrOut("应用配置格式为 key=value 请您检查输入的格式")
			}
			configInfos = append(configInfos, vos.AppConfig{
				Name: split[0],
				Val:  split[1],
			})
		}

		pluginEnvConfigMap := make(map[string]map[string]string)
		if len(pluginEnvConfig) > 0 {
			for _, pluginStr := range pluginEnvConfig {
				split := strings.Split(pluginStr, ":")
				if len(split) < 2 {
					tools.ErrOut("插件配置格式必须满足: 名称:key=value")
					return
				}

				pluginConfigs, ok := pluginEnvConfigMap[split[0]]
				if !ok {
					pluginConfigs = make(map[string]string)
					pluginEnvConfigMap[split[0]] = pluginConfigs
				}

				pluginConfigValueAndKey := strings.Split(strings.Join(split[1:], ":"), "=")
				if len(pluginConfigValueAndKey) < 2 {
					tools.ErrOut("插件配置格式必须满足: 名称:key=value")
					return
				}
				pluginConfigs[pluginConfigValueAndKey[0]] = strings.Join(pluginConfigValueAndKey[1:], "=")
			}
		}

		appStartInfo := &vos.AppStartInfo{
			Name:            name,
			Version:         currentVersion,
			JdkPath:         javaCmdPath,
			Args:            args,
			JdkPackName:     javaPackName,
			Restart:         restartType,
			EnvConfig:       configInfos,
			CopyFiles:       copyFiles,
			SaveAppSuffix:   saveAppSuffix,
			Xmx:             Xmx,
			Xms:             Xms,
			Xmn:             Xmn,
			PermSize:        PermSize,
			MaxPermSize:     MaxPermSize,
			PluginEnvConfig: pluginEnvConfigMap,
		}

		conn := socket.GetClientConn()
		defer conn.SendEndMsg()

		marshal, _ := json.Marshal(appStartInfo)
		conn.WriteDataStr(cmdStart).Wait().WriteData(marshal).Wait()
		tools.SuccessOut("应用启动成功")

	},
}

func startWithConfig(configPath string) {
	abs, err := filepath.Abs(configPath)
	if err != nil {
		tools.ErrOut("获取配置文件绝对路径失败")
	}

	_, err = os.Stat(abs)
	if err != nil {
		tools.ErrOut("启动配置不存在")
	}

	conn := socket.GetClientConn()
	defer conn.SendEndMsg()

	conn.WriteDataStr(cmdStartWithConfig).Wait().WriteDataStr(abs)
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
	rootCmd.AddCommand(startCmd)
	startCmd.Flags().SortFlags = false
	startCmd.Flags().StringP("withStartConfig", "w", "", "根据启动配置文件启动")
	startCmd.Flags().StringP("currentVersion", "v", "", "以指定版本启动,并将指定的版本设置为当前版本")
	startCmd.Flags().StringSliceP("copyFile", "f", []string{}, "拷贝文件到程序运行目录当中")
	startCmd.Flags().StringSliceP("envConfig", "e", []string{}, "覆盖程序配置, 格式: key=value")
	startCmd.Flags().StringP("javaPackName", "p", "", "使用平台内的java包名称")
	startCmd.Flags().StringP("javaCmdPath", "c", "", "java命令路径")
	startCmd.Flags().BoolP("saveAppSuffix", "s", false, "是否保留应用文件后缀")
	startCmd.Flags().StringP("restart", "r", "", `重启模式(每次重启间隔60s)
|-- auto: 开机自动启动
|-- always: 开机、异常等所有意外和正常结束情况都会执行重启操作
|-- error-any: 发生插件错误以及应用错误时重启
|-- error-app: 发生应用错误重启
|-- error-plugin: 发生插件错误重启`)
	startCmd.Flags().StringSlice("pluginEnvConfig", nil, `插件可变参数配置, 格式: 
插件名称(可以是插件名称的前几位):key=value
例如:
8d1:envHost=127.0.0.1:8080`)
	startCmd.Flags().String("Xmx", "", "最大总对内存, 推荐设置为物理内存的1/4")
	startCmd.Flags().String("Xms", "", "初始总堆内存, 推荐和最大堆内存一样大(GC之后就不必调整堆内存大小)")
	startCmd.Flags().String("Xmn", "", "年轻代堆内存, 官方推荐为整个堆的3/8")
	startCmd.Flags().String("PermSize", "", "堆的初始大小, 一般设置为128m即可, 原则为预留30%的空间")
	startCmd.Flags().String("MaxPermSize", "", "堆的最大大小, 一般设置为128m即可, 原则为预留30%的空间")
	//Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// startCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// startCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
