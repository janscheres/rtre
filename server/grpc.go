package main

import (
	"log"
	"net"
	"time"

	pb "github.com/janscheres/rtre/pb"
	"google.golang.org/grpc"
)

type riskServer struct {
	pb.UnimplementedRiskServiceServer
}


func (s *riskServer) StreamRisk(req *pb.RiskRequest, stream pb.RiskService_StreamRiskServer) error {
	ctx := stream.Context()

	log.Println("[gRPC] Connecting to new client")

	wsclient := &WsClient{
		ctx: ctx,
		symbol: req.Symbol,
		orderbook: OrderBook{
			Bids: make(map[float64]float64),
			Asks: make(map[float64]float64),
			OBIChan: make(chan float64, 100),
			SpreadChan: make(chan float64, 100),
		},
	}

	go wsclient.run()

	var lastOBI float64
	var lastSpread float64

	for {
		select {
		case <-ctx.Done():
			log.Println("[gRPC] Connection closed")
			return nil

		case obi := <-wsclient.orderbook.OBIChan:
			lastOBI = obi
			err := stream.Send(&pb.RiskResponse{
				Timestamp: time.Now().UnixNano(),
				Obi: lastOBI,
				Spread: lastSpread,

			})

			if err != nil {
				return err
			}

		case spread := <-wsclient.orderbook.SpreadChan:
			lastSpread = spread
			err := stream.Send(&pb.RiskResponse{
				Timestamp: time.Now().UnixNano(),
				Obi: lastOBI,
				Spread: lastSpread,

			})

			if err != nil {
				return err
			}
		}
	}
}

func startgRPCServer(s *riskServer) {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal("[NET] Failed to start gRPC server", err)
	}

	grpcServer := grpc.NewServer()
	//pb.RegisterRiskServiceServer(grpcServer, &riskServer{orderBook: o})
	pb.RegisterRiskServiceServer(grpcServer, s)

	log.Println("Ready to start receiving on gRPC :50051")

	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatal("[gRPC] Failed to serve (not slay)", err)
	}
}
