package grpcClient

import (
	"api/grpcModels"
	"api/models"
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

func LogIssue(tx *grpcModels.LogRequest) (*grpcModels.LogResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	grpcCon, err := grpc.DialContext(ctx, models.GRPC_IP, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		//utils.WrapErrorLog(fmt.Sprintf("did not connect: %s", err))
		cancel()
		return nil, err
	}
	defer grpcCon.Close()
	defer cancel()

	c := grpcModels.NewLogClient(grpcCon)
	return c.Log(ctx, tx)
}

func PingMainNode(tx *grpcModels.PingRequest) (*grpcModels.PingResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	grpcCon, err := grpc.DialContext(ctx, models.GRPC_IP, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		//utils.WrapErrorLog(fmt.Sprintf("did not connect: %s", err))
		cancel()
		return nil, err
	}
	defer grpcCon.Close()
	defer cancel()

	c := grpcModels.NewLogClient(grpcCon)
	return c.Ping(ctx, tx)
}

func NSFWReq(tx *grpcModels.NSFWAnnRequest) (*grpcModels.NSFWAnnResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	grpcCon, err := grpc.DialContext(ctx, models.GRPC_IP, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		//utils.WrapErrorLog(fmt.Sprintf("did not connect: %s", err))
		cancel()
		return nil, err
	}
	defer grpcCon.Close()
	defer cancel()

	c := grpcModels.NewLogClient(grpcCon)
	return c.NSFWAnn(ctx, tx)
}
