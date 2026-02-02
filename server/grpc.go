package main

import (
	"net"
	"log"
		"time"

	"google.golang.org/grpc"
	pb "github.com/janscheres/rtre/pb"
)

type riskServer struct {
	pb.UnimplementedRiskServiceServer

	orderBook *OrderBook
}

func (s *riskServer) StreamRisk(req *pb.RiskRequest, stream pb.RiskService_StreamRiskServer) error {
	ctx := stream.Context()

	for {
		select {
		case <-ctx.Done():
			log.Println("[gRPC] Connection closed")
			return nil
		case obi, ok := <-s.orderBook.OBIChan:
			if !ok {
				log.Println("OBI channel closed")
				return nil
			}

            err := stream.Send(&pb.RiskResponse{
                Timestamp: int64(t.Unix()),
				Obi: obi,

            })
			log.Println("sent obi!:)")
            if err != nil {
                return err
            }
		}
	}
}

func startgRPCServer(o *OrderBook) {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal("[NET] Failed to start gRPC server", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterRiskServiceServer(grpcServer, &riskServer{orderBook: o})

	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatal("[gRPC] Failed to serve (not slay)", err)
	}
}
