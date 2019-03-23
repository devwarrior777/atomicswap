The functionality here can be used directly by golang clients just by
importing it.

When a server is set up this code will be on the server side of grpc.

Client side you can generate into any language supported by grpc.

The gen tool generates golang<->golang protobuf code so useful only for
the grpc server side.

Above this directory is one for each supported coin with all the code to
make atomicswap transactions by connecting to a coin-specific full-node(s)
containing the desired wallet(s)

This is adapted from decred/atomicswap and heavily refactred and modified 

At the moment paths are hard coded but later I will make everything dependent
on configuration files

All is currently using go12.1
