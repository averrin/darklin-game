rev=$(git log --pretty=format:'' | wc -l)
go run -ldflags "-X main.VERSION 0.1.$rev#dev" ./main.go -port 8089 -interactive
