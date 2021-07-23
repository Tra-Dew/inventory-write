package cmd

import (
	"context"
	"reflect"

	"github.com/Tra-Dew/inventory-write/pkg/core"
	"github.com/Tra-Dew/inventory-write/pkg/inventory"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// TradeAccepted ...
func TradeAccepted(command *cobra.Command, args []string) {
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
		core.WithSubscriberID(settings.Events.TradeAccepted),
		core.WithMaxRetries(3),
		core.WithType(reflect.TypeOf(inventory.TradeOfferAcceptedEvent{})),
		core.WithTopicID(settings.Events.TradeAccepted),
		core.WithHandler(func(payload interface{}) error {
			message := payload.(*inventory.TradeOfferAcceptedEvent)

			fields := logrus.Fields{
				"trade_id": message.ID,
				"owner_id": message.OwnerID,
			}

			logrus.
				WithFields(fields).
				Info("processing received event")

			ctx := context.Background()

			offeredItems := make([]*inventory.TradeItemModel, len(message.OfferedItems))

			for i, item := range message.WantedItems {
				offeredItems[i] = &inventory.TradeItemModel{
					ID:       item.ID,
					Quantity: item.Quantity,
				}
			}

			wantedItems := make([]*inventory.TradeItemModel, len(message.WantedItems))

			for i, item := range message.WantedItems {
				wantedItems[i] = &inventory.TradeItemModel{
					ID:       item.ID,
					Quantity: item.Quantity,
				}
			}

			req := &inventory.TradeItemsRequest{
				OfferedItems: offeredItems,
				WantedItems:  wantedItems,
			}

			if err := container.InventoryService.TradeItems(ctx, req); err != nil {
				logrus.
					WithError(err).
					WithFields(fields).
					Info("error while trading items")
				return err
			}

			logrus.
				WithFields(fields).
				Info("items locked successfully")

			//TODO: dispatch trade completed event

			return nil
		}))

	if err := consumer.Run(); err != nil {
		logrus.
			WithError(err).
			Error("shutting down with error")
	}
}
