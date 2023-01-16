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
	"github.com/byzk-org/bypt/consts"
	"github.com/byzk-org/bypt/socket"
	"github.com/byzk-org/bypt/tools"
	"github.com/byzk-org/bypt/vos"
	"github.com/spf13/cobra"
	"github.com/ttacon/chalk"
	"strconv"
	"time"
)

// infoCmd represents the info command
var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "查看当前平台信息",
	Long:  `查看当前平台信息`,
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(infoBanner())
		fmt.Println(infoConfigList())
		infoLogClearBanner()
	},
}

func infoConfigList() *uitable.Table {
	defer func() {
		e := recover()
		if e != nil {
			tools.ErrOut("获取程序配置信息失败")
		}
	}()
	conn := socket.GetClientConn()
	defer conn.SendEndMsg()

	listJsonInfo := conn.WriteDataStr(cmdConfigList).Wait().ReadMsg()

	configList := make([]*vos.Setting, 0)
	if err := json.Unmarshal(listJsonInfo, &configList); err != nil {
		tools.ErrOut("解析响应数据失败")
	}

	if len(configList) == 0 {
		tools.ErrOut("获取配置信息列表失败")
	}

	configMap := make(map[string]*vos.Setting)
	for _, config := range configList {
		configMap[config.Name] = config
	}

	cert, err := tools.ParseCertByPem(consts.CaCert)
	if err != nil {
		tools.ErrOut("解析服务信息失败")
		return nil
	}

	userCert, err := tools.ParseCertByPem(consts.UserCert)
	if err != nil {
		tools.ErrOut("解析客户端信息失败")
		return nil
	}

	table := uitable.New()
	table.Wrap = true
	table.AddRow("服务分组: ", cert.Issuer.CommonName)
	table.AddRow("服务器到期时间: ", cert.NotAfter.Format("2006-01-02 15:03:04"))
	table.AddRow("客户端到期时间: ", userCert.NotAfter.Format("2006-01-02 15:03:04"))
	table.AddRow("程序运行目录: ", configMap["runDir"].Val)
	table.AddRow("日志存放目录: ", configMap["logDir"].Val)
	table.AddRow("程序文件存放目录: ", configMap["appSaveDir"].Val)
	table.AddRow("jdk文件存放目录: ", configMap["jdkSaveDir"].Val)

	return table
}

func infoLogClearBanner() {
	conn := socket.GetClientConn()
	defer conn.SendEndMsg()

	msg := conn.WriteDataStr(cmdInfoLogClear).Wait().ReadMsg()
	tmpInfo := &struct {
		NextClearTime time.Time
		PrevClearTime time.Time
		ClearLogMsg   []byte
		TimeSpace     int64
		TimeUnit      time.Duration
	}{}
	if err := json.Unmarshal(msg, &tmpInfo); err != nil {
		tools.ErrOut("获取日志清除信息失败")
	}

	table := uitable.New()
	table.Wrap = true

	timeUnitStr := "小时"
	switch tmpInfo.TimeUnit {
	case time.Millisecond:
		timeUnitStr = "毫秒"
	case time.Second:
		timeUnitStr = "秒"
	case time.Minute:
		timeUnitStr = "分钟"
	case time.Hour:
		timeUnitStr = "小时"
	default:
		tools.ErrOut("非法的时间单位")
	}

	table.AddRow("日志清除时间间隔: ", strconv.FormatInt(tmpInfo.TimeSpace, 10)+timeUnitStr)

	if !tmpInfo.PrevClearTime.IsZero() {
		table.AddRow("日志最近一次清除时间: ", tmpInfo.PrevClearTime.Format("2006-01-02 15:04:05"))
	}

	if !tmpInfo.NextClearTime.IsZero() {
		table.AddRow("日志下一次清除时间: ", tmpInfo.NextClearTime.Format("2006-01-02 15:04:05"))
	}
	if len(tmpInfo.ClearLogMsg) > 0 {
		table.AddRow(fmt.Sprint(chalk.Yellow, "日志最近一次输出信息: "))
		fmt.Println(table)
		fmt.Println(string(tmpInfo.ClearLogMsg), chalk.Reset)
	} else {
		fmt.Println(table)
	}

}

func infoBanner() string {
	conn := socket.GetClientConn()
	defer conn.Close()

	return conn.WriteDataStr(cmdInfoBanner).Wait().ReadMsgStr()
}

func init() {
	rootCmd.AddCommand(infoCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// infoCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// infoCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
