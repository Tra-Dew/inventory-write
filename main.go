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

	dispatchItemCreatedWorker := &cobra.Command{
		Use:   "dispatch-item-created-worker",
		Short: "Starts dispatch-item-created-worker",
		Run:   cmd.DispatchItemCreated,
	}

	dispatchItemUpdatedWorker := &cobra.Command{
		Use:   "dispatch-item-updated-worker",
		Short: "Starts dispatch-item-updated-worker",
		Run:   cmd.DispatchItemUpdated,
	}

	root.PersistentFlags().String("settings", "./settings.yml", "path to settings.yaml config file")
	root.AddCommand(
		api,
		itemsLockRequestedConsumer,
		dispatchItemLockedWorker,
		dispatchItemCreatedWorker,
		dispatchItemUpdatedWorker,
	)

	root.Execute()
}
