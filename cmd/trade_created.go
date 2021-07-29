package cmd

import (
	"context"
	"reflect"

	"github.com/d-leme/tradew-inventory-write/pkg/core"
	"github.com/d-leme/tradew-inventory-write/pkg/inventory"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// TradeCreated ...
func TradeCreated(command *cobra.Command, args []string) {
	settings := new(core.Settings)

	err := core.FromYAML(command.Flag("settings").Value.String(), settings)
	if err != nil {
		logrus.
			WithError(err).
			Fatal("unable to parse settings, shutting down...")
	}

	container := NewContainer(settings)

	consumer := core.NewMessageBrokerSubscriber(
		core.WithSessionSNS(container.SNS),
		core.WithSessionSQS(container.SQS),
		core.WithSubscriberID(settings.Events.TradeCreated),
		core.WithMaxRetries(3),
		core.WithType(reflect.TypeOf(inventory.TradeCreatedEvent{})),
		core.WithTopicID(settings.Events.TradeCreated),
		core.WithHandler(func(payload interface{}) error {
			message := payload.(*inventory.TradeCreatedEvent)

			fields := logrus.Fields{
				"trade_id":              message.ID,
				"owner_id":              message.OwnerID,
				"wanted_items_owner_id": message.WantedItemsOwnerID,
			}

			ctx := context.Background()

			logrus.
				WithFields(fields).
				Info("processing received event")

			wantedIDs := make([]string, len(message.WantedItems))
			for i, id := range wantedIDs {
				wantedIDs[i] = id
			}

			wantedItems, err := container.InventoryRepository.Get(ctx, &message.WantedItemsOwnerID, wantedIDs)

			if err != nil {
				logrus.
					WithError(err).
					WithFields(fields).
					Info("error while getting wanted items")
				return err
			}

			// validating if all wanted items exist and belong to the same person
			if len(wantedItems) != len(message.WantedItems) {
				// publish error
			}

			items := make([]*inventory.LockItemModel, len(message.OfferedItems))

			for i, item := range message.OfferedItems {
				items[i] = &inventory.LockItemModel{
					ID:       item.ID,
					Quantity: item.Quantity,
				}
			}

			req := &inventory.LockItemsRequest{
				LockedBy: message.ID,
				Items:    items,
			}

			err = container.InventoryService.LockItems(ctx, message.OwnerID, req)
			if err != nil {
				logrus.
					WithError(err).
					WithFields(fields).
					Info("error while locking items")
				return err
			}

			logrus.
				WithFields(fields).
				Info("items locked successfully")

			event := inventory.ParseItemsToItemsLockCompletedEvent(message.ID, message.OfferedItems)
			messageID, err := container.Producer.Publish(settings.Events.ItemsLockCompleted, event)

			if err != nil {
				logrus.
					WithError(err).
					WithFields(fields).
					Error("error while dispatching message")
				return err
			}

			fields["message_id"] = messageID
			fields["event"] = settings.Events.ItemsLockCompleted

			logrus.
				WithFields(fields).
				Info("dipached event")

			return nil
		}))

	if err := consumer.Run(); err != nil {
		logrus.
			WithError(err).
			Error("shutting down with error")
	}
}
