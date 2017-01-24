#!/sh

export GOOS=$1
export GOARCH=$2
echo "Building docker-stats to /artifacts/$GOOS/$GOOARCH..."
mkdir -p /artifacts/$GOOS/$GOARCH
go build -o /artifacts/$GOOS/$GOARCH/docker-stats
