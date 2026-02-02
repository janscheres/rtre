import grpc

import rtre_pb2
import rtre_pb2_grpc


def main():
    channel = grpc.insecure_channel("localhost:50051")
    stub = rtre_pb2_grpc.RiskServiceStub(channel)

    request = rtre_pb2.RiskRequest(symbol="btcusdt")

    try:
        for response in stub.StreamRisk(request):
            print("Received:", response.timestamp, response.obi)
    except grpc.RpcError as e:
        print("gRPC error:", e)


if __name__ == "__main__":
    main()

