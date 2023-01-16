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
	"github.com/byzk-org/bypt/tools"
	"github.com/spf13/cobra"
	"github.com/ttacon/chalk"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

// importCmd represents the import command
var importCmd = &cobra.Command{
	Use:   "import dir",
	Short: "导入一个文件内的所有程序包",
	Long:  `导入一个文件内的所有程序包`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		importPath := args[0]
		abs, err := filepath.Abs(importPath)
		if err != nil {
			tools.ErrOut("获取文件夹的绝对路径失败")
		}
		stat, err := os.Stat(abs)
		if err != nil {
			tools.ErrOut("获取文件夹信息失败")
			return
		}

		if !stat.IsDir() {
			tools.ErrOut("需要传入文件夹地址而不是文件")
		}

		dirList, err := ioutil.ReadDir(abs)
		if err != nil {
			tools.ErrOut("获取文件夹内信息失败")
			return
		}

		successNum := 0
		errNum := 0
		for _, f := range dirList {
			if f.IsDir() {
				tools.WarningOut(fmt.Sprintf("跳过目录[%s]", f.Name()))
				continue
			}

			execFile := filepath.Join(abs, f.Name())
			fmt.Printf("正在导入文件[%s]...\n", f.Name())

			command := exec.Command(execFile, "install")
			command.Env = os.Environ()
			command.Dir = abs
			output, e := command.CombinedOutput()
			if e != nil {
				fmt.Println(chalk.Red, fmt.Sprintf("导入[%s]失败 => %s", f.Name(), output), chalk.Reset)
				errNum += 1
				continue
			}
			tools.SuccessOut(fmt.Sprintf("导入[%s],成功", f.Name()))
			fmt.Println()
			successNum += 1
		}

		fmt.Printf("文件夹内总计数量: %d\n", len(dirList))
		fmt.Printf("成功导入数量: %d\n", successNum)
		fmt.Printf("失败导入数量: %d\n", errNum)

	},
}

func init() {
	rootCmd.AddCommand(importCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// importCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// importCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
