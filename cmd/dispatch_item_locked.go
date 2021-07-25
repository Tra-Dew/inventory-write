package cmd

import (
	"context"

	"github.com/d-leme/tradew-inventory-write/pkg/core"
	"github.com/d-leme/tradew-inventory-write/pkg/inventory"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// DispatchItemLocked ...
func DispatchItemLocked(command *cobra.Command, args []string) {
	settings := new(core.Settings)

	err := core.FromYAML(command.Flag("settings").Value.String(), settings)
	if err != nil {
		logrus.
			WithError(err).
			Fatal("unable to parse settings, shutting down...")
	}

	ctx := context.Background()
	container := NewContainer(settings)

	fields := logrus.Fields{"event": settings.Events.ItemsLockRequested}

	items, err := container.InventoryRepository.GetByStatus(ctx, inventory.ItemPendingLockDispatch)
	if err != nil {
		logrus.
			WithError(err).
			WithFields(fields).
			Error("error while getting items by status")
		return
	}

	lenItems := len(items)

	logrus.WithFields(fields).Infof("%d new items to publish", lenItems)

	if lenItems < 1 {
		return
	}

	event := inventory.ParseItemsToItemsLockCompletedEvent(items)
	messageID, err := container.Producer.Publish(settings.Events.ItemsLockRequested, event)

	if err != nil {
		logrus.
			WithError(err).
			WithFields(fields).
			Error("error while dispatching message")
		return
	}

	fields["message_id"] = messageID

	logrus.
		WithFields(fields).
		Info("dipached event")

	for _, item := range items {
		item.UpdateStatus(inventory.ItemAvailable)
	}

	if err := container.InventoryRepository.UpdateBulk(ctx, nil, items); err != nil {
		logrus.
			WithError(err).
			WithFields(fields).
			Error("error while updating items")
		return
	}

	logrus.
		WithFields(fields).
		Info("worker complete")
}
