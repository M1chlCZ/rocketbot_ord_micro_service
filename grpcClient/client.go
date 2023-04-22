package grpcClient

import (
	"api/grpcModels"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

func DetectNSFW(tx *grpcModels.NSFWRequest) (*grpcModels.NSFWResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	grpcCon, err := grpc.DialContext(ctx, fmt.Sprintf(":%d", 4000), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		//utils.WrapErrorLog(fmt.Sprintf("did not connect: %s", err))
		cancel()
		return nil, err
	}
	defer grpcCon.Close()
	defer cancel()

	c := grpcModels.NewNSFWClient(grpcCon)
	return c.Detect(ctx, tx)
}
