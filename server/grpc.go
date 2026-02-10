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


func (s *riskServer) NewClient(c *WsClient) {
	for {
		(*c).connect(c.symbol)

		log.Println("[NET] Connection died, restarting...")
	}

}

func (s *riskServer) StreamRisk(req *pb.RiskRequest, stream pb.RiskService_StreamRiskServer) error {
	ctx := stream.Context()

	log.Println("[gRPC] Connecting to new client")

	wsclient := WsClient{
		orderbook: OrderBook{
			Bids: make(map[float64]float64),
			Asks: make(map[float64]float64),
			OBIChan: make(chan float64, 100),
			SpreadChan: make(chan float64, 100),
		},
		symbol: req.Symbol,
	}

	go s.NewClient(&wsclient)

	for {
		var obi float64
		var spread float64

		select {
		case <-ctx.Done():
			log.Println("[gRPC] Connection closed")
			return nil
		case o, ok := <-wsclient.orderbook.OBIChan:
			if !ok {
				log.Println("OBI channel closed")
				return nil
			}
			obi = o
		}

		select {
		case s, ok := <-wsclient.orderbook.SpreadChan:
			if !ok {
				log.Println("Spread channel closed")
				return nil
			}
			spread = s
		default:
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
