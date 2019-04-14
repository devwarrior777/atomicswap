gRPC support
============

The functionality here can be used directly by golang clients just by
importing it:

import (
    "github.com/devwarrior777/atomicswap/libs"
    "github.com/devwarrior777/atomicswap/libs/ltc"
    "github.com/devwarrior777/atomicswap/libs/xzc"
    //...
)

Other languages
---------------

When the golang gRPC SwapServer is running it will be the server side of grpc.

Client side you can generate into any language supported by grpc such as python,
nodejs, etc.

The `gen` tool generates golang<->golang protobuf code so useful only for
the golang gRPC server side.

A `pygen` tool will make python source binding in the `./python/` folder

A `jsgen` tool will make javascript source bindings in the `./python/` folder

Other language generation protoc parameters are shown in the `example_gen` file

You will need the `protoc-XXX.zip` compiler at: https://github.com/protocolbuffers/protobuf/releases

