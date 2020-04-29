package main

import (
	"syscall"

	wxccserver "github.com/liuzhanpeng/go-wx-ccserver"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "停止服务",
	Run: func(cmd *cobra.Command, args []string) {
		pidFile := wxccserver.NewPIDFile(pidFilename)
		pid, err := pidFile.GetPID()
		if err != nil {
			logrus.Fatalln("程序未运行，不能停止")
		}

		logrus.Infof("[%v]stop...", pid)
		syscall.Kill(pid, syscall.SIGTERM)
	},
}

func init() {
	rootCmd.AddCommand(stopCmd)
}
