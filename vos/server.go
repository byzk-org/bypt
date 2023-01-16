package vos

import (
	"time"
)

type AppPluginType string

const (
	AppPluginTypeListener AppPluginType = "listener"
	AppPluginTypeNormal   AppPluginType = "normal"
)

type AppRestartType string

const (
	// AppRestartTypeErrorAuto 开机自启
	AppRestartTypeErrorAuto   AppRestartType = "auto"
	AppRestartTypeAlways      AppRestartType = "always"
	AppRestartTypeErrorAny    AppRestartType = "error-any"
	AppRestartTypeErrorApp    AppRestartType = "error-app"
	AppRestartTypeErrorPlugin AppRestartType = "error-plugin"
)

type appRunStatus string

const (
	appRunStatusRunner      appRunStatus = "正在运行"
	appRunStatusRunError    appRunStatus = "运行异常"
	appRunStatusWaitRestart appRunStatus = "等待重启"
	appRunStatusRunRestart  appRunStatus = "正在重启"
)

type Setting struct {
	Name    string `gorm:"primary_key" json:"name,omitempty"`
	Desc    string `json:"desc,omitempty"`
	Val     string `json:"val,omitempty"`
	StopApp bool   `json:"stopApp,omitempty"`
}

type AppInfo struct {
	// Name 应用名称
	Name               string           `gorm:"primary_key" json:"name,omitempty"`
	Desc               string           `json:"desc,omitempty"`
	CreateTime         time.Time        `json:"createTime,omitempty"`
	EndUpdateTime      time.Time        `json:"endUpdateTime,omitempty"`
	CurrentVersion     string           `json:"currentVersion,omitempty"`
	CurrentVersionInfo *AppVersionInfo  `gorm:"-" json:"currentVersionInfo,omitempty"`
	Versions           []AppVersionInfo `gorm:"-" json:"versions,omitempty"`
	Restart            AppRestartType   `gorm:"-" json:"restart,omitempty"`
}

type AppVersionInfo struct {
	AppName        string       `json:"-"`
	Name           string       `json:"name,omitempty"`
	Desc           string       `json:"desc,omitempty"`
	CreateTime     time.Time    `json:"createTime,omitempty"`
	EndUpdateTime  time.Time    `json:"endUpdateTime,omitempty"`
	ContentMd5     []byte       `json:"md5,omitempty"`
	ContentSha1    []byte       `json:"sha1,omitempty"`
	EnvConfigInfos []*AppConfig `json:"envConfigInfos,omitempty"`
	PluginInfo     []*AppPlugin `json:"pluginInfo,omitempty"`
}
type AppPlugin struct {
	AppName    string
	AppVersion string
	Name       string
	Desc       string
	Content    []byte
	Md5        []byte
	Sha1       []byte
	Sign       []byte
	Type       AppPluginType
	EnvConfig  []*AppConfig
}

// AppStartInfo app启动信息
type AppStartInfo struct {
	// Name app信息
	Name string `json:"name,omitempty"`
	// Version 启动的版本号
	Version string `json:"version,omitempty"`
	// EnvConfig 环境配置
	EnvConfig []AppConfig `json:"envConfig,omitempty"`
	// JdkPath 本地jdk路径
	JdkPath string `json:"jdkPath,omitempty"`
	// JdkPackName 管理器内jdk的名称
	JdkPackName string `json:"jdkPackName,omitempty"`
	// JdkArgs jar包之前的jdk命令
	JdkArgs []string `json:"jdkFlags,omitempty"`
	// Args 应用参数
	Args []string `json:"args,omitempty"`
	// Restart 是否跟随服务重启
	Restart AppRestartType `json:"restart,omitempty"`
	// CopyFiles 要拷贝的文件
	CopyFiles []string `json:"copyFiles,omitempty"`
	// SaveAppSuffix 保留应用文件后缀
	SaveAppSuffix   bool                         `json:"saveAppSuffix,omitempty"`
	Xmx             string                       `json:"xmx,omitempty"`
	Xms             string                       `json:"xms,omitempty"`
	Xmn             string                       `json:"xmn,omitempty"`
	PermSize        string                       `json:"permSize,omitempty"`
	MaxPermSize     string                       `json:"maxPermSize,omitempty"`
	PluginEnvConfig map[string]map[string]string `gorm:"-" json:"pluginEnvConfig,omitempty"`
}

type AppConfig struct {
	Name       string `json:"name,omitempty"`
	Val        string `json:"val,omitempty"`
	Desc       string `json:"desc,omitempty"`
	DefaultVal string `json:"defaultVal,omitempty"`
}

type AppStatusInfo struct {
	StartArgs          *AppStartInfo     `json:"startArgs,omitempty"`
	Name               string            `json:"name,omitempty"`
	Desc               string            `json:"desc,omitempty"`
	AppInfo            *AppInfo          `json:"appInfo,omitempty"`
	VersionStr         string            `json:"versionStr,omitempty"`
	VersionInfo        *AppVersionInfo   `json:"versionInfo,omitempty"`
	StartTime          time.Time         `json:"startTime,omitempty"`
	HaveErr            bool              `json:"haveErr,omitempty"`
	ErrMsg             string            `json:"errMsg,omitempty"`
	JavaCmd            string            `json:"javaCmd,omitempty"`
	PluginOutPutBuffer map[string][]byte `json:"pluginOutPutBuffer,omitempty"`
	Status             appRunStatus      `json:"status,omitempty"`
}

type DbLog struct {
	AppName    string `json:"appName,omitempty"`
	AppVersion string `json:"appVersion,omitempty"`
	Content    []byte `json:"content,omitempty"`
	AtDate     int64  `json:"atDate,omitempty"`
}

type SyncInfo struct {
	All bool `json:"all,omitempty"`
	//StartInfo  bool   `json:"startInfo,omitempty"`
	Jdk        bool   `json:"jdk,omitempty"`
	App        bool   `json:"app,omitempty"`
	Version    bool   `json:"version,omitempty"`
	AppName    string `json:"appName,omitempty"`
	AppVersion string `json:"appVersion,omitempty"`
}

type JdkInfo struct {
	Name          string    `json:"name,omitempty"`
	Desc          string    `json:"desc,omitempty"`
	CreateTime    time.Time `json:"createTime,omitempty"`
	EndUpdateTime time.Time `json:"endUpdateTime,omitempty"`
	MD5           []byte    `json:"md5,omitempty"`
	SHA1          []byte    `json:"sha1,omitempty"`
	Content       []byte    `json:"content,omitempty"`
	Sign          []byte    `json:"sign,omitempty"`
}
