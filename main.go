package main

import (
	"github.com/Tra-Dew/inventory-write/cmd"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			logrus.WithField("error", err).Error("error main")
		}
	}()

	root := &cobra.Command{
		Use:   "api",
		Short: "Starts api handlers",
		Run:   cmd.Server,
	}

	root.PersistentFlags().String("settings", "./settings.yml", "path to settings.yaml config file")
	root.Execute()
}
