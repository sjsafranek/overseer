#go get -u -v github.com/golang/protobuf/{proto,protoc-gen-go}
#go get -u -v google.golang.org/grpc
#protoc --go_out=plugins=grpc:. *.proto

protoc --go_out=. --go_opt=paths=source_relative     --go-grpc_out=. --go-grpc_opt=paths=source_relative     overseer.proto
