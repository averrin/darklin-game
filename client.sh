cd ./client
export GOPATH=$GOPATH:$(pwd)
go build -o ./darklin-client ./*.go;
./darklin-client -host "localhost:8089"
cd -
