gRPC support
============

The functionality here can be used directly by golang clients just by
importing it:

import (
    "github.com/devwarrior777/atomicswap/libs/ltc"
	"github.com/devwarrior777/atomicswap/libs/xzc"
	//...
)

Other languages
---------------

When the gRPC SwapServer is running this code will be on the server side of grpc.

Client side you can generate into any language supported by grpc such as python,
nodejs, etc.

The gen tool generates golang<->golang protobuf code so useful only for
the golang gRPC server side.
