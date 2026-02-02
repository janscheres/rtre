# Real-Time Risk Engine

## Running
Start the server with `go run .` in the server directory

Run the Python client with `uv run client/client.py`

## Protobuf file generation commands
`protoc --proto_path=. \
  --go_out=./server --go_opt=paths=source_relative \
  --go-grpc_out=./server --go-grpc_opt=paths=source_relative \
  rtre.proto`

`uv add grpcio grpcio-tools protobuf`

`uv run python -m grpc_tools.protoc \
  -I. \
  --python_out=client \
  --grpc_python_out=client \
  rtre.proto`
