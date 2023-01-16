package cmd

import (
	"fmt"
	"github.com/byzk-org/bypt/tools"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"os/user"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "bypt",
	Short: "java应用管理程序(客户端)",
	Long: `
 ________      ___    ___ ________  _________   
|\   __  \    |\  \  /  /|\   __  \|\___   ___\ 
\ \  \|\ /_   \ \  \/  / | \  \|\  \|___ \  \_| 
 \ \   __  \   \ \    / / \ \   ____\   \ \  \  
  \ \  \|\  \   \/  /  /   \ \  \___|    \ \  \ 
   \ \_______\__/  / /      \ \__\        \ \__\
    \|_______|\___/ /        \|__|         \|__|
             \|___|/                            
                                     版本: 2.0.0
                                     作者: 无&痕
   Java应用部署管理客户端`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child cechoommands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	//cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	//rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.newApp.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	current, err := user.Current()
	if err != nil {
		tools.ErrOut("获取当前用户失败")
		os.Exit(-1)
	}
	logDir := tools.PathJoin(current.HomeDir, ".bypt", ".data", "client")
	os.MkdirAll(logDir, os.ModePerm)
	cfgFile = tools.PathJoin(logDir, ".bypt.yaml")
	tools.CreateFile(cfgFile)
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".newApp" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".bypt.yaml")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		//fmt.Println("Using config file:", viper.ConfigFileUsed())

	} else {
		fmt.Println(err)
	}
}
