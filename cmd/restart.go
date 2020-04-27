package main

import (
	"lzpeng/wxccserver"
	"syscall"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var restartCmd = &cobra.Command{
	Use:   "restart",
	Short: "重启服务",
	Run: func(cmd *cobra.Command, args []string) {
		pidFile := wxccserver.NewPIDFile(pidFilename)
		pid, err := pidFile.GetPID()
		if err != nil {
			logrus.Fatalln("程序未运行，不能重启")
			return
		}

		logrus.Infof("[%v]restart...", pid)
		syscall.Kill(pid, syscall.SIGUSR1)
	},
}

func init() {
	rootCmd.AddCommand(restartCmd)
}
