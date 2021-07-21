package cmd

import (
	"context"

	"github.com/Tra-Dew/inventory-write/pkg/core"
	"github.com/Tra-Dew/inventory-write/pkg/inventory"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// DispatchItemCreated ...
func DispatchItemCreated(command *cobra.Command, args []string) {
	settings := new(core.Settings)

	err := core.FromYAML(command.Flag("settings").Value.String(), settings)
	if err != nil {
		logrus.
			WithError(err).
			Fatal("unable to parse settings, shutting down...")
	}

	ctx := context.Background()
	container := NewContainer(settings)

	fields := logrus.Fields{"event": settings.Events.ItemsCreated}

	items, err := container.InventoryRepository.GetByStatus(ctx, inventory.ItemPendingCreateDispatch)
	if err != nil {
		logrus.
			WithError(err).
			WithFields(fields).
			Error("error while getting items by statys")
		return
	}

	event := inventory.ParseItemsToItemsCreatedEvent(items)
	messageID, err := container.Producer.Publish(settings.Events.ItemsCreated, event)

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
