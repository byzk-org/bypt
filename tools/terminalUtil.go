package tools

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/manifoldco/promptui"
	"github.com/manifoldco/promptui/list"
	"github.com/mattn/go-colorable"
	"github.com/ttacon/chalk"
	"os"
)

var stdOut = bufio.NewWriter(colorable.NewColorableStdout())

func OutWithWidth(str string, width int) string {
	runes := []rune(str)
	lens := 0

	for _, v := range runes {
		if v > 256 {
			lens += 2
		} else {
			lens += 1
		}
	}

	if lens >= width {
		return str
	}

	other := width - lens
	isDouble := true
	if other%2 != 0 {
		other += 1
		isDouble = false
	}

	end := other / 2

	blankRune := make([]rune, end, end)
	for i := 0; i < end; i++ {
		blankRune[i] = ' '
	}

	endBlankStr := string(blankRune)
	leftBlank := endBlankStr
	if !isDouble {
		leftBlank = endBlankStr[:end-1]
	}
	return fmt.Sprintf("%s%s%s", leftBlank, str, endBlankStr)

}

func OutWithWidthFunc(width int) func(str string) string {
	return func(str string) string {
		return OutWithWidth(str, width)
	}
}

func ErrOutAndExit(str string) {
	_, _ = fmt.Fprintln(stdOut, chalk.Red, chalk.Bold.TextStyle(str), chalk.Reset)
	_ = stdOut.Flush()
	panic("")
}

func ErrOut(str string) {
	_, _ = fmt.Fprintln(stdOut, chalk.Red, chalk.Bold.TextStyle(str), chalk.Reset)
	_ = stdOut.Flush()
	os.Exit(1)
}

func SuccessOut(str string) {
	_, _ = fmt.Fprintln(stdOut, chalk.Green, str, chalk.Reset)
	_ = stdOut.Flush()
}

func InfoOut(str string) {
	_, _ = fmt.Fprintln(stdOut, chalk.Cyan, str, chalk.Reset)
	_ = stdOut.Flush()
}

func WarningOut(str string) {
	_, _ = fmt.Fprintln(stdOut, chalk.Yellow, str, chalk.Reset)
	_ = stdOut.Flush()
}

func ScanTerminalInput(str string) string {

	input := promptui.Prompt{
		Label: str,
	}

	val, err := input.Run()
	if err != nil {
		ErrOutAndExit("Ctrl+C 退出程序")
	}
	return val
}

func ScanTerminalInputWithValidate(str string, validateFunc promptui.ValidateFunc) string {
	prompt := &promptui.Prompt{
		Label:    str,
		Validate: validateFunc,
	}
	run, err := prompt.Run()
	if err == promptui.ErrInterrupt {
		ErrOutAndExit("Ctrl+C 退出程序")
	}

	return run

}

func ScanTerminalInputWithVerifyFormatFun(str string, errStr string, verify VerifyFormatFun, allowNil bool) string {
	prompt := &promptui.Prompt{
		Label: str,
		Validate: func(s string) error {
			if allowNil && s == "" {
				return nil
			}
			if !verify(s) {
				return errors.New(errStr)
			}
			return nil
		},
	}

	run, err := prompt.Run()
	if err == promptui.ErrInterrupt {
		ErrOutAndExit("Ctrl+C 退出程序")
	}
	return run
}

func ScanTerminalRequiredInput(str string) string {
	prompt := promptui.Prompt{
		Label: str,
		Validate: func(s string) error {
			if len(s) == 0 {
				return errors.New("字段不允许为空!")
			}
			return nil
		},
	}

	val, err := prompt.Run()
	if err != nil {
		ErrOutAndExit("Ctrl+C 退出程序")

	}

	return val
}

func ScanTerminalMultiSelect(str string, keys []string, defaultKeys ...string) []string {
	multiSelect := &survey.MultiSelect{
		Message:  str,
		Options:  keys,
		Help:     "按空格进行选中/取消，上下键移动选项，回车键确认",
		Default:  defaultKeys,
		PageSize: 10,
	}

	var data []string
	err := survey.AskOne(multiSelect, &data)
	if err == terminal.InterruptErr {
		ErrOutAndExit("Ctrl+C 退出程序")
	}
	return data
}

func ScanTerminalSelectWithIndex(str string, keys interface{}, defaultKeys ...string) (int, string) {
	selectInput := promptui.Select{
		Label: str,
		Items: keys,
		Size:  8,
	}
	index, data, err := selectInput.Run()
	if err != nil {
		ErrOutAndExit("Ctrl+C 退出程序")
	}
	return index, data
}

func ScanTerminalSelect(str string, keys interface{}, defaultKeys ...string) string {
	selectInput := promptui.Select{
		Label: str,
		Items: keys,
	}
	_, data, err := selectInput.Run()
	if err != nil {
		ErrOutAndExit("Ctrl+C 退出程序")
	}
	//multiSelect := &survey.Select{
	//	Message:  str,
	//	Options:  keys,
	//	Help:     "按空格进行选中/取消，上下键移动选项，回车键确认",
	//	Default:  defaultKeys,
	//	PageSize: 10,
	//}
	//
	//var data string
	//_ = survey.AskOne(multiSelect, &data)
	//if data == "" && len(defaultKeys) > 0 {
	//	data = defaultKeys[0]
	//}
	return data
}

func ScanTerminalSelectWithTemplate(str string, keys interface{}, template *promptui.SelectTemplates, search list.Searcher) (int, string) {
	selectPromp := &promptui.Select{
		Label:     str,
		Items:     keys,
		Templates: template,
		Size:      8,
		Searcher:  search,
	}
	run, s, err := selectPromp.Run()
	if err == promptui.ErrInterrupt {
		ErrOutAndExit("Ctrl+C 退出程序")
	}
	return run, s
}

func ScanTerminalConfirm(str string) bool {
	confirm := &survey.Confirm{
		Message: str,
	}

	var data bool
	err := survey.AskOne(confirm, &data)
	if err == terminal.InterruptErr {
		ErrOutAndExit("Ctrl+C 退出程序")
	}
	return data
}
