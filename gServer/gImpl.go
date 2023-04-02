package gServer

import (
	grpcModels "api/gModels"
	"context"
)

type RocketbotORDServer struct {
	grpcModels.UnimplementedRocketbotORDServer
}

func (s *RocketbotORDServer) GetORDDetails(ctx context.Context, in *grpcModels.GetOrdDetailsReq) (*grpcModels.GetOrdDetailsReply, error) {
	panic("implement me")
}
