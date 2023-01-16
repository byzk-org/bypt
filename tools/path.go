package tools

import (
	"github.com/byzk-org/bypt/vos"
	"path"
	"runtime"
	"strings"
)

func PathJoin(ele ...string) string {
	str := path.Join(ele...)
	if vos.SystemType(runtime.GOOS) == vos.WINDOWS {
		str = strings.ReplaceAll(str, "/", "\\")
		str = strings.ReplaceAll(str, "\\\\", "\\")
	}
	return str
}
