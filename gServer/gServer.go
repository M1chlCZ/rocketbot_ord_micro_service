package gServer

import (
	grpcModels "api/gModels"
	"api/utils"
	"crypto/tls"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
	"net"
)

func NewGServer() {
	tlsCredentials, err := loadTLSCredentials()
	if err != nil {
		log.Fatal("cannot load TLS credentials: ", err)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", utils.STAKEPORT))
	if err != nil {
		utils.WrapErrorLog(err.Error())
	}
	utils.ReportMessage(fmt.Sprintf("gRPC Online on port %d!", utils.STAKEPORT))
	grpcServer := grpc.NewServer(grpc.Creds(tlsCredentials))
	s := RocketbotORDServer{}
	grpcModels.RegisterRocketbotORDServer(grpcServer, &s)
	if err := grpcServer.Serve(lis); err != nil {
		utils.WrapErrorLog(err.Error())
	}
}

func loadTLSCredentials() (credentials.TransportCredentials, error) {
	serverCert, err := tls.LoadX509KeyPair("./.cert/server-cert.pem", "./.cert/server-key.pem")
	if err != nil {
		utils.ReportMessage("Failed to load server's certificate and private key")
		return nil, err
	}

	config := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientAuth:   tls.NoClientCert,
	}

	return credentials.NewTLS(config), nil
}
