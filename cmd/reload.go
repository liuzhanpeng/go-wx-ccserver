package main

import (
	"syscall"

	wxccserver "github.com/liuzhanpeng/go-wx-ccserver"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var reloadCmd = &cobra.Command{
	Use:   "reload",
	Short: "重新加载公众号账号信息",
	Run: func(cmd *cobra.Command, args []string) {
		pidFile := wxccserver.NewPIDFile(pidFilename)
		pid, err := pidFile.GetPID()
		if err != nil {
			logrus.Fatalln("程序未运行")
			return
		}

		logrus.Infof("[%v]reload...", pid)
		syscall.Kill(pid, syscall.SIGUSR2)
	},
}

func init() {
	rootCmd.AddCommand(reloadCmd)
}
