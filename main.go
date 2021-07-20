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

	root := &cobra.Command{}

	api := &cobra.Command{
		Use:   "api",
		Short: "Starts api handlers",
		Run:   cmd.Server,
	}

	itemsLockRequestedConsumer := &cobra.Command{
		Use:   "items-lock-requested-consumer",
		Short: "Starts items-lock-requested-consumer",
		Run:   cmd.ItemsLockRequested,
	}

	dispatchItemLockedWorker := &cobra.Command{
		Use:   "dispatch-item-locked-worker",
		Short: "Starts dispatch-item-locked-worker",
		Run:   cmd.DispatchItemLocked,
	}

	root.PersistentFlags().String("settings", "./settings.yml", "path to settings.yaml config file")
	root.AddCommand(
		api,
		itemsLockRequestedConsumer,
		dispatchItemLockedWorker,
	)

	root.Execute()
}
