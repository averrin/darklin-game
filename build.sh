export GOPATH=$GOPATH:$(pwd)
rm ./core
rev=$(git log --pretty=format:'' | wc -l)
go build -ldflags "-X main.VERSION=0.1.$rev" -o ./core ./main.go;
echo "Build completed"
