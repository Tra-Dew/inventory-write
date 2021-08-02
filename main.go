package main

import (
	"github.com/d-leme/tradew-inventory-write/cmd"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			logrus.WithField("error", err).Error("error main")
		}
	}()

	root := &cobra.Command{}

	api := &cobra.Command{
		Use:   "api",
		Short: "Starts api handlers",
		Run:   cmd.ServerHTTP,
	}

	grpc := &cobra.Command{
		Use:   "grpc",
		Short: "Starts api handlers",
		Run:   cmd.ServerGRPC,
	}

	itemUpdatedWorker := &cobra.Command{
		Use:   "dispatch-item-updated-worker",
		Short: "Starts dispatch-item-updated-worker",
		Run:   cmd.DispatchItemUpdated,
	}

	root.PersistentFlags().String("settings", "./settings.yml", "path to settings.yaml config file")
	root.AddCommand(api, grpc, itemUpdatedWorker)

	root.Execute()
}
