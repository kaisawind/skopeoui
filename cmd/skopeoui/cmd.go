package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/kaisawind/skopeoui/pkg/configs"
	"github.com/kaisawind/skopeoui/pkg/configs/cmd"
	"github.com/kaisawind/skopeoui/pkg/stats"
	"github.com/kaisawind/skopeoui/server"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	name = "skopeoui"
	logo = cmd.ArtLogo(name)

	version  bool
	logLevel string
	config   string

	rootCmd = &cobra.Command{
		Use:   name,
		Short: fmt.Sprintf("%s core service of iotx.", name),
		Long:  logo,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("Version   : %s\n", stats.Version)
			fmt.Printf("BuildTime : %s\n", stats.BuildTime)
			if version {
				return nil
			}
			fmt.Print(logo)
			if logLevel != "" {
				lvl, err := logrus.ParseLevel(logLevel)
				if err != nil {
					logrus.SetLevel(logrus.InfoLevel)
				} else {
					logrus.SetLevel(lvl)
				}
				if lvl >= logrus.DebugLevel {
					logrus.Infoln("debug on")
					f, e := os.OpenFile(os.Args[0]+".log", os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
					if e == nil {
						logrus.SetOutput(io.MultiWriter(f, os.Stdout))
						logrus.RegisterExitHandler(func() {
							f.Close()
						})
					}
				}
				logrus.Debugf("log level %s", lvl)
			}
			configs.PrintAllSettings()
			opts := server.NewOptions().
				ApplyDBConfig(configs.GetDB()).
				ApplyHttpConfig(configs.GetHttp())
			s := server.NewServer(opts)
			defer s.Close()
			logrus.Infof("%s is starting...", name)
			return s.Serve()
		},
	}
)

func init() {
	binds := map[string]string{}
	cobra.OnInitialize(func() {
		configs.SetConfigFile(config)
		for k, v := range binds {
			flag := rootCmd.PersistentFlags().Lookup(v)
			if flag == nil {
				flag = rootCmd.Flags().Lookup(v)
			}
			if flag != nil {
				configs.BindPFlag(k, flag)
			}
		}
	})
	rootCmd.Flags().BoolVarP(&version, "version", "v", false, "print version")
	rootCmd.PersistentFlags().StringVar(&logLevel, "log-level", "info", "set log level")
	rootCmd.PersistentFlags().StringVarP(&config, "config", "c", "", "set skopeoui config file path")

	for k, v := range cmd.Envs() {
		flag := strings.ReplaceAll(k, ".", "-")
		rootCmd.PersistentFlags().String(flag, v.Default, v.Usage)
		binds[k] = flag
	}
}
