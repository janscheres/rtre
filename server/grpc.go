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

	log.Println("[gRPC] Connecting to new client")

	for {
		var obi float64
		var spread float64

		select {
		case <-ctx.Done():
			log.Println("[gRPC] Connection closed")
			return nil
		case o, ok := <-s.orderBook.OBIChan:
			if !ok {
				log.Println("OBI channel closed")
				return nil
			}
			obi = o
		}

		select {
		case s, ok := <-s.orderBook.SpreadChan:
			if !ok {
				log.Println("Spread channel closed")
				return nil
			}
			spread = s
		}

		err := stream.Send(&pb.RiskResponse{
			Timestamp: time.Now().UnixNano(),
			Obi: obi,
			Spread: spread,

		})
		//log.Println("sent obi!:)")
		if err != nil {
			return err
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

	log.Println("Ready to start receiving on gRPC :50051")

	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatal("[gRPC] Failed to serve (not slay)", err)
	}
}
