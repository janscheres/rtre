import grpc

import rtre_pb2
import rtre_pb2_grpc


def main():
    channel = grpc.insecure_channel("localhost:50051")
    stub = rtre_pb2_grpc.RiskServiceStub(channel)

    request = rtre_pb2.RiskRequest(symbol="btcusdt")

    obiArr = []
    obiSum = 0
    numObi = 0
    mean = 0

    try:
        for response in stub.StreamRisk(request):
            obiArr.append(response.obi)
            obiSum+=response.obi
            numObi+=1

            mean = obiSum/numObi

            print("Received:", response.timestamp, response.obi, mean)
            print("SPREAD", response.spread)

    except grpc.RpcError as e:
        print("gRPC error:", e)

    sumOfSquares=0
    for o in obiArr:
        sumOfSquares+= pow((o-mean), 2)

    print("sample standard deviation:", pow(sumOfSquares/(numObi-1), 0.5))


if __name__ == "__main__":
    main()

