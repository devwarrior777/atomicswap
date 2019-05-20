Client Test
===========

`client_test.go` is used to test the swap-lib server functionality.

- Make sure configured coins are running and the XXXtestdata.go has correct RPC Info settings
- Make sure the swap-lib server is running

```bash
cd libs/protobind/server
./server
```

You will need your own TLS certs or switch TLS off in the libs/protobind/server/config.ini

The client test uses it's own config.ini just for testing

```bash
cd libs/protobind/server/svrtest
go test
```

