package cmd

import (
	"fmt"
	"net"

	"github.com/d-leme/tradew-inventory-write/pkg/core"
	"github.com/d-leme/tradew-inventory-write/pkg/inventory"
	"github.com/d-leme/tradew-inventory-write/pkg/inventory/proto"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

// ServerGRPC ...
func ServerGRPC(command *cobra.Command, args []string) {

	settings := new(core.Settings)

	if err := core.FromYAML(command.Flag("settings").Value.String(), settings); err != nil {
		logrus.
			WithError(err).
			Fatal("unable to parse settings, shutting down...")
		return
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", settings.GRPCPort))
	if err != nil {
		logrus.WithError(err).Fatal("failed to start grpc server")
		return
	}

	container := NewContainer(settings)

	grpcServer := grpc.NewServer()

	s := inventory.NewGRPCService(container.InventoryService)

	proto.RegisterInventoryServiceServer(grpcServer, s)

	logrus.Infof("starting grpc service at port %v", settings.GRPCPort)

	if err := grpcServer.Serve(lis); err != nil {
		logrus.WithError(err).Fatal("failed to start grpc server")
	}
}
