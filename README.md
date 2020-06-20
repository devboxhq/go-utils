# Go Utils
Utilities and packages containing helpful code!

List of packages:
- Path related (util)
- Serialization (util)
- GRPC Middleware (grpc/middleware)
- JWT Auth (auth/jwt)

## Installation
```
go get -u github.com/uptimize/go-utils
```

## Examples
- get root path based on relative path
```golang
import "github.com/uptimize/go-utils/pkg/util"

func main() {
    path := util.FromRootPath("binaries")
}
```
- json to bytes
```golang
import "github.com/uptimize/go-utils/pkg/util"

func main() {
    data := util.MustJsonToBytes(jsonData)
}
```
- GRPC zap middleware
```golang
import (
    "github.com/uptimize/go-utils/pkg/grpc/middleware"
    "go.uber.org/zap"
    "google.golang.org/grpc"
)

func main() {
    logger, _ := zap.NewProduction()

    // create manager and add zap middleware
    manager := middleware.Manager{}
    _ = manager.AddMiddleware(middleware.NewZapMiddleware())

    // provide middlewares to GRPC server
    grpcServer := grpc.NewServer(manager.BuildServerOptions()...)
}
```
