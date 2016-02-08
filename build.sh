export GOPATH=$GOPATH:$(pwd)
# if [ -f ./core ]; then
#   rm ./core || true
# fi
rev=$(git log --pretty=format:'' | wc -l)
go build -ldflags "-X main.VERSION=0.1.$rev" -o ./core ./main.go;
echo "Build completed"
