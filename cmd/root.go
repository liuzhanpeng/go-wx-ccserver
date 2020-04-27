/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

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
package main

import (
	"log"
	"lzpeng/wxccserver"
	"os"
	"path/filepath"
	"time"

	rotatelogs "github.com/lestrrat/go-file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/spf13/viper"
)

var pidFilename string = os.Args[0] + ".pid"
var cfgFile string
var config wxccserver.Config
var daemon bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use: "wxccserver",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig, initLogger)
}

func initConfig() {
	viper.SetConfigFile(cfgFile)
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("加载配置错误: %v", err)

	}

	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("加载解析错误: %v", err)
	}
}

// 初始化日志
func initLogger() {
	logPath := filepath.Join(config.Log.Path, "log")
	_, err := os.Stat(logPath)
	if err != nil {
		os.Mkdir(filepath.Dir(logPath), os.ModePerm)
	}

	rotationTime := time.Duration(config.Log.RotationTime)
	maxAge := time.Duration(config.Log.MaxAge)
	writer, err := rotatelogs.New(
		logPath+".%Y%m%d",
		rotatelogs.WithLinkName(logPath),
		rotatelogs.WithRotationTime(rotationTime*time.Hour),
		rotatelogs.WithMaxAge(maxAge*time.Hour),
	)
	if err != nil {
		log.Fatalf("logger error %v", err)
	}

	level, err := logrus.ParseLevel(config.Log.Level)
	if err != nil {
		log.Fatalf("logger error %v", err)
	}
	logrus.SetLevel(level)
	hook := lfshook.NewHook(lfshook.WriterMap{
		logrus.DebugLevel: writer,
		logrus.InfoLevel:  writer,
		logrus.WarnLevel:  writer,
		logrus.ErrorLevel: writer,
		logrus.PanicLevel: writer,
	}, &logrus.JSONFormatter{})
	logrus.AddHook(hook)
}
