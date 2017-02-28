declare -A exts
exts=(["windows"]=".exe")

if [ -z "$1" ]; then
    echo Please, specify commit hash or tag. && exit
fi

VERSION=$(git rev-parse "$1")
TAG=$(git tag --points-at "$VERSION")
if [ -n "$TAG" ]; then 
  VERSION=$TAG
fi

for OS in "linux" "windows" "darwin" "freebsd"; do
    for ARCH in "386" "amd64"; do
        CGO_ENABLED=0 GOOS=$OS GOARCH=$ARCH go build -ldflags "-X main.version=$VERSION" -o bombardier-$OS-$ARCH"${exts[$OS]}"
    done
done