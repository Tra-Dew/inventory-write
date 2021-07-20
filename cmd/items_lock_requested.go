package cmd

import (
	"context"
	"reflect"

	"github.com/Tra-Dew/inventory-write/pkg/core"
	"github.com/Tra-Dew/inventory-write/pkg/inventory"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// ItemsLockRequested ...
func ItemsLockRequested(command *cobra.Command, args []string) {
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
		core.WithSubscriberID(settings.Events.ItemsLockRequested),
		core.WithMaxRetries(3),
		core.WithType(reflect.TypeOf(inventory.ItemsLockRequestedEvent{})),
		core.WithTopicID(settings.Events.ItemsLockRequested),
		core.WithHandler(func(payload interface{}) error {
			message := payload.(*inventory.ItemsLockRequestedEvent)

			fields := logrus.Fields{
				"owner_id":       message.OwnerID,
				"correlation_id": message.CorrelationID,
				"event":          settings.Events.ItemsLockRequested,
			}

			logrus.
				WithFields(fields).
				Info("processing received event")

			ctx := context.Background()

			items := make([]*inventory.LockItemModel, len(message.Items))

			for i, item := range message.Items {
				items[i] = &inventory.LockItemModel{
					ID:       item.ID,
					Quantity: item.Quantity,
				}
			}

			req := &inventory.LockItemsRequest{Items: items}

			err := container.InventoryService.LockItems(ctx, message.OwnerID, message.CorrelationID, req)

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

			return nil
		}))

	if err := consumer.Run(); err != nil {
		logrus.
			WithError(err).
			Error("shutting down with error")
	}
}
