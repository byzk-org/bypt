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
	"github.com/gosuri/uiprogress"
	"github.com/byzk-org/bypt/socket"
	"github.com/byzk-org/bypt/tools"
	"github.com/byzk-org/bypt/vos"
	"github.com/spf13/cobra"
	"github.com/ttacon/chalk"
	"strconv"
	"strings"
	"time"
)

// syncCmd represents the sync command
var syncCmd = &cobra.Command{
	Use:   "sync [appName] [version]",
	Args:  cobra.MaximumNArgs(2),
	Short: "同步远程应用到本地",
	Long:  `同步远程应用、版本、插件等信息到本地`,
	Run: func(cmd *cobra.Command, args []string) {
		syncAll, _ := cmd.LocalFlags().GetBool("all")
		syncApp, _ := cmd.LocalFlags().GetBool("app")
		syncJdk, _ := cmd.LocalFlags().GetBool("jdk")
		syncVersion, _ := cmd.LocalFlags().GetBool("version")

		syncInfo := &vos.SyncInfo{
			All:     syncAll,
			App:     syncApp,
			Jdk:     syncJdk,
			Version: syncVersion,
		}

		if len(args) > 0 {
			syncInfo.AppName = args[0]
			if len(args) == 2 {
				syncInfo.AppVersion = args[1]
			}
		}

		haveNum := 0

		if syncInfo.All {
			haveNum += 1
		}

		if syncInfo.App {
			haveNum += 1
		}

		if syncInfo.Version {
			haveNum += 1
		}

		if syncInfo.AppName != "" && haveNum == 0 {
			tools.ErrOut("缺少应用同步项")
		}

		if haveNum == 0 {
			cmd.Help()
			return
		}

		marshal, _ := json.Marshal(syncInfo)
		conn := socket.GetClientConn()
		defer conn.Close()

		conn.WriteDataStr(cmdSyncInfo).Wait().WriteData(marshal)

		progress := uiprogress.New()
		progress.Start()
		for {
			if syncShowProcess(progress, conn) {
				tools.SuccessOut("同步成功")
				return
			}
		}

	},
}

func syncShowProcess(progress *uiprogress.Progress, conn *socket.Conn) bool {
	var (
		receiveSize int64 = 0
	)
	//progress := uiprogress.New()
	//progress.Start()
	//defer progress.Stop()

	msg, ok := syncMsgHeaderHandler(conn)
	if ok {
		return true
	}
	msgSizeStr := conn.ReadMsgStr()
	msgSize, err := strconv.ParseInt(msgSizeStr, 10, 64)
	if err != nil {
		tools.ErrOut("转换 " + msg + " 消息长度失败")
	}

	conn.WriteDataStr("ok")

	bar := progress.AddBar(100).AppendCompleted()
	bar.PrependFunc(func(b *uiprogress.Bar) string {
		return msg
	})

	progress.SetRefreshInterval(time.Second * 1)

	for receiveSize < msgSize {
		tmpMsg := conn.ReadMsgStr()
		receiveSize, err = strconv.ParseInt(tmpMsg, 10, 64)
		if err != nil {
			tools.ErrOut("转换数据长度失败")
		}
		conn.WriteDataStr("ok")
		_ = bar.Set(int(float64(receiveSize) / float64(msgSize) * 100))
	}
	return false
}

func syncMsgHeaderHandler(conn *socket.Conn) (string, bool) {
	str := conn.ReadMsgStr()
	if str == "!!!!!!" {
		return "", true
	}
	if strings.Contains(str, ":") {
		split := strings.Split(str, ":")
		if len(split) < 2 {
			return str, false
		}

		level := split[0]
		msg := strings.Join(split[1:], "")
		switch level {
		case "info":
			fmt.Println(chalk.Cyan, msg, chalk.Reset)
		case "warn":
			fmt.Println(chalk.Yellow, msg, chalk.Reset)
		case "error":
			fmt.Println(chalk.Red, msg, chalk.Reset)
		default:
			fmt.Println(msg)
		}
		return syncMsgHeaderHandler(conn)
	}
	return str, false
}

func init() {
	rootCmd.AddCommand(syncCmd)
	syncCmd.Flags().BoolP("all", "a", false, "同步所有信息, 包括应用、版本以及jdk信息")
	syncCmd.Flags().BoolP("app", "m", false, "同步应用和版本信息")
	syncCmd.Flags().BoolP("version", "v", false, "同步版本信息")
	syncCmd.Flags().BoolP("jdk", "j", false, "同步jdk信息")
}
