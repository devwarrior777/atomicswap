import grpc
import logging

import atomicswap_pb2
import atomicswap_pb2_grpc

''' Test Python swap-lib bindings; 
    Set use_tls=false in server config.ini
'''
def run():
    # NOTE(gRPC Python Team): .close() is possible on a channel and should be
    # used in circumstances in which the with statement does not fit the needs
    # of the code.
    #
    # For more channel options, please see https://grpc.io/grpc/core/group__grpc__arg__keys.html
    with grpc.insecure_channel(
            target='localhost:10010',
            options=[('grpc.lb_policy_name', 'pick_first'),
                     ('grpc.enable_retries', 0), ('grpc.keepalive_timeout_ms',
                                                  3000)]) as channel:
        stub = atomicswap_pb2_grpc.SwapLibStub(channel)
        # Timeout in seconds.
        # Please refer gRPC Python documents for more detail. https://grpc.io/grpc/python/grpc.html
        response = stub.PingWalletRPC(
            atomicswap_pb2.PingWalletRPCRequest(
                coin=atomicswap_pb2.LTC, 
                testnet=True,
                hostport='localhost',
                rpcuser='dev',
                rpcpass='dev',
                wpass='123'
                ), timeout=3)
        print("PingWalletRPC response: ", response.errorno, response.errstr)
        response = stub.NewAddress(
            atomicswap_pb2.NewAddressRequest(
                coin=atomicswap_pb2.LTC, 
                testnet=True,
                hostport='localhost',
                rpcuser='dev',
                rpcpass='dev',
                wpass='123'
                ), timeout=3)
        print("NewAddress response: ", response.errorno, response.errstr)
        if response.errorno == atomicswap_pb2.OK:
            print("New address:",response.address)
    print("Tests Done")


if __name__ == '__main__':
    logging.basicConfig()
    run()