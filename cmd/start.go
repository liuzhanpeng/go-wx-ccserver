package main

import (
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"syscall"

	wxccserver "github.com/liuzhanpeng/go-wx-ccserver"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "启动服务",
	Run: func(cmd *cobra.Command, args []string) {
		if daemon {
			command := exec.Command(os.Args[0], "start")
			command.Start()
			return
		}

		pidFile := wxccserver.NewPIDFile(pidFilename)
		pidFile.SetPID(os.Getpid())
		defer pidFile.Remove()

		srv := wxccserver.NewServer(&config)

		go func() {
			// 监听信号
			sigChan := make(chan os.Signal, 1)
			signal.Notify(sigChan, syscall.SIGHUP, syscall.SIGUSR1, syscall.SIGUSR2, syscall.SIGINT, syscall.SIGTERM)

			for {
				sig := <-sigChan

				switch sig {
				case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM:
					srv.Stop()
					logrus.Info("stopped")
					return
				case syscall.SIGUSR1:
					srv.Stop()

					command := exec.Command("./"+filepath.Base(os.Args[0]), "start")
					command.Start()
					return
				case syscall.SIGUSR2:
					if err := srv.Reload(); err != nil {
						logrus.Errorf("reload:", err)
					}

					logrus.Info("reloaded")

					return
				}
			}
		}()

		logrus.Infof("[%v]start...", os.Getpid())
		srv.Start()
	},
}

func init() {
	startCmd.Flags().StringVar(&cfgFile, "config", "./config.toml", "配置文件路径")
	startCmd.Flags().BoolVar(&daemon, "daemon", false, "是否以服务形式运行")

	rootCmd.AddCommand(startCmd)
}
