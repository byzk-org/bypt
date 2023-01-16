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
	"errors"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/byzk-org/bypt/socket"
	"github.com/byzk-org/bypt/tools"
	"github.com/byzk-org/bypt/vos"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

const (
	pageSize       = 100
	timeFLayoutStr = "20060102150405"
)

// logsCmd represents the logs command
var logsCmd = &cobra.Command{
	Use:   "logs [flags] appName [appVersion]",
	Short: "日志操作",
	Long:  `日志查看、导出、清空等操作`,
	Args:  cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {

		conn := socket.GetClientConn()
		defer conn.SendEndMsg()
		defer func() {
			e := recover()
			if e != nil {
				switch err := e.(type) {
				case error:
					tools.ErrOut(err.Error())
				case string:
					tools.ErrOut(err)
				default:
					tools.ErrOut("未知异常")
				}
			}
		}()

		var (
			startT time.Time
			endT   time.Time
			err    error
		)

		clear, _ := cmd.LocalFlags().GetBool("clear")
		follow, _ := cmd.LocalFlags().GetBool("follow")
		exportFile, _ := cmd.LocalFlags().GetString("export")
		startTime, _ := cmd.LocalFlags().GetString("startTime")
		endTime, _ := cmd.LocalFlags().GetString("endTime")

		logDir := ""
		if len(args) == 1 {
			logDir = conn.WriteDataStr(cmdLogByName).Wait().WriteDataStr(args[0]).ReadMsgStr()
		} else {
			logDir = conn.WriteDataStr(cmdLogByNameAndVersion).Wait().WriteDataStr(args[0]).WriteDataStr(args[1]).ReadMsgStr()
		}

		conn.Wait()
		if logDir == "" {
			tools.ErrOut("获取日志路径失败")
			return
		}

		stat, err := os.Stat(logDir)
		if err != nil {
			tools.ErrOut("获取日志目录状态失败")
			return
		}

		if !stat.IsDir() {
			tools.ErrOut("打开日志目录失败")
			return
		}

		dir, err := ioutil.ReadDir(logDir)
		if err != nil {
			tools.ErrOut("获取日志目录失败")
		}

		logFileName := filepath.Join(logDir, dir[0].Name())
		logFileStat, err := os.Stat(logFileName)
		if err != nil {
			tools.ErrOut("获取日志文件状态失败")
			return
		}

		if logFileStat.IsDir() {
			tools.ErrOut("获取日志文件失败")
			return
		}

		logDb, err := gorm.Open("sqlite3", fmt.Sprintf("file:%s?auto_vacuum=1", logFileName))
		if err != nil {
			tools.ErrOut("打开日志文件失败")
			return
		}
		logDb.LogMode(false)
		//logDb = logDb.Debug()

		if clear {
			deleteLogMode := logDb.Model(&vos.DbLog{})
			if endTime != "" {
				endT, err = time.ParseInLocation(timeFLayoutStr, endTime, time.Local)
				if err != nil {
					tools.ErrOut("转换结束时间个是失败, 请您检查您的输入")
				}
				deleteLogMode = deleteLogMode.Where("at_date <= ?", endT.UnixNano())
			}

			if endTime != "" {
				endT, err = time.ParseInLocation(timeFLayoutStr, endTime, time.Local)
				if err != nil {
					tools.ErrOut("转换结束时间个是失败, 请您检查您的输入")
				}
				deleteLogMode = deleteLogMode.Where("at_date <= ?", endT.UnixNano())
			}

			if !tools.ScanTerminalConfirm("确认要删除日志吗, 此操作不可逆, 请谨慎操作!!!!!") {
				return
			}

			if err = deleteLogMode.Delete(&vos.DbLog{}).Error; err != nil {
				tools.ErrOut("删除日志失败")
			}
			tools.SuccessOut("日志删除成功")
			return
		}

		logModel := logDb.Model(&vos.DbLog{}).Order("at_date asc").Limit(pageSize)

		if startTime != "" {
			startT, err = time.ParseInLocation(timeFLayoutStr, startTime, time.Local)
			if err != nil {
				tools.ErrOut("转换开始时间格式失败, 请您检查您的输入")
			}
			logModel = logModel.Where("at_date >= ?", startT.UnixNano())
		}

		if endTime != "" {
			endT, err = time.ParseInLocation(timeFLayoutStr, endTime, time.Local)
			if err != nil {
				tools.ErrOut("转换结束时间个是失败, 请您检查您的输入")
			}
			logModel = logModel.Where("at_date <= ?", endT.UnixNano())
		}

		if exportFile != "" {
			logExportFile(exportFile, logModel)
			return
		}

		page := 1
		logResult := make([]vos.DbLog, 0, pageSize)

		// 剩余
		var surplus int64 = 0

		for {
			if err = logModel.Offset((int64((page - 1) * pageSize)) - surplus).Find(&logResult).Error; err != nil && err != gorm.ErrRecordNotFound {
				if e, ok := err.(sqlite3.Error); ok && e.Code == sqlite3.ErrBusy {
					continue
				}
				tools.ErrOut("日志文件解析读取失败 => " + err.Error())
				return
			}

			if len(logResult) == 0 {
				if follow {
					if err = watcherFile(logFileName); err != nil {
						tools.ErrOut(err.Error())
						return
					}
					continue
				}
				return
			}

			if len(logResult) < pageSize {
				surplus += int64(pageSize - len(logResult))
			}

			for _, log := range logResult {
				fmt.Println(string(log.Content))
			}

			page += 1
		}

	},
}

func logExportFile(exportFileName string, logModel *gorm.DB) {
	abs, err := filepath.Abs(exportFileName)
	if err != nil {
		tools.ErrOut("获取导出文件绝对路径失败")
	}
	_ = os.MkdirAll(filepath.Dir(abs), 0777)
	exportFile, err := os.Create(abs)
	if err != nil {
		tools.ErrOut("创建导出文件失败")
	}

	page := 1
	logResult := make([]vos.DbLog, 0, pageSize)
	tools.SuccessOut("将导出日志文件到 [" + abs + "] 文件中")
	for {
		if err = logModel.Offset((page - 1) * pageSize).Find(&logResult).Error; err != nil && err != gorm.ErrRecordNotFound {
			tools.ErrOut("日志文件解析读取失败 => " + err.Error())
			return
		}
		if len(logResult) == 0 {
			break
		}

		for _, log := range logResult {
			_, _ = exportFile.Write(log.Content)
			_, _ = exportFile.WriteString("\n")
		}
		page += 1
	}
	tools.SuccessOut("日志导出完成")
}

func watcherFile(filePath string) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return errors.New("创建文件监听失败")
	}
	defer watcher.Close()
	err = watcher.Add(filePath)
	if err != nil {
		return errors.New("创建文件监听器失败")
	}

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return errors.New("检测文件监听失败")
			}

			if event.Op&fsnotify.Write == fsnotify.Write {
				return nil
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return errors.New("获取文件错误监听失败")
			}
			return errors.New("文件监听失败 => " + err.Error())
		}

	}

}

func init() {
	rootCmd.AddCommand(logsCmd)

	logsCmd.Flags().BoolP("follow", "f", false, "持续监听")
	logsCmd.Flags().StringP("export", "e", "", "导出到文件, 带有此参数 [follow] 参数将失效")
	logsCmd.Flags().StringP("startTime", "s", "", "要查询的开始时间, 格式: yyyyMMddHHmmss, 例: 20210329000000")
	logsCmd.Flags().StringP("endTime", "t", "", "要查询的结束时间, 格式: yyyyMMddHHmmss, 例: 20210329000000")
	logsCmd.Flags().BoolP("clear", "c", false, "清除日志")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// logsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// logsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
