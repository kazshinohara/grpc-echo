# gRPC Echo
A test gRPC application.  
gRPC Echo is based upon [whereami](https://github.com/kazshinohara/whereami) which is a Rest API
teaching your Pod's hostname, request header value, Google Cloud region etc.  
As of now, only server code is available, please take your favorite gRPC client tool.

Here are examples how to use gRPC Echo with Evans.  
*Note: gRPC Echo enables Server Reflection, don't forget "-r" option.*  

Unary RPC
```shell
❯ echo '{ "request_header_name": "hoge" }' | evans --port 8080 --header hoge="fuga" -r cli call EchoService.GetHeader
{
  "requestHeaderValue": "fuga"
}
```
```shell
❯ echo '{}' | evans --port 8080 -r cli call EchoService.GetVersion
{
  "version": "v0.1"
}
```
```shell
❯ echo '{}' | evans --port 8080 -r cli call EchoService.GetHostname
{
  "hostname": "Kazuus-MacBook-Air.local"
}
```
```shell
❯ echo '{}' | evans --port 8080 -r cli call EchoService.GetKind
{
  "kind": "test"
}
```
```shell
❯ echo '{}' | evans --port 8080 -r cli call EchoService.GetSourceIp
{
  "sourceIp": "127.0.0.1"
}
```
```shell
❯ while true; do echo '{}' | evans --port 8080 -r cli call EchoService.GetHostname; sleep 1; done
{
  "hostname": "Kazuus-MacBook-Air.local"
}
{
  "hostname": "Kazuus-MacBook-Air.local"
}
{
  "hostname": "Kazuus-MacBook-Air.local"
}
^C
```
Server streaming RPC
```shell
❯ echo '{ "interval": 1, "number_of_response": 5 }' | evans --port 8080 -r cli call EchoService.GetHostnameServerStream 
{
  "hostname": "Kazuus-MacBook-Air.local"
}
{
  "hostname": "Kazuus-MacBook-Air.local"
}
{
  "hostname": "Kazuus-MacBook-Air.local"
}
{
  "hostname": "Kazuus-MacBook-Air.local"
}
{
  "hostname": "Kazuus-MacBook-Air.local"
}
```

via ASM
```shell
❯ echo '{}' | evans --host asm.gcpx.org --port 443 -t -r cli call EchoService.GetHostname                                            (asm-cluster-01/default)
{
  "hostname": "grpc-echo-68b8797599-lq826"
}
```
```shell
❯ echo '{ "interval": 1, "number_of_response": 10 }' | evans --host asm.gcpx.org --port 443 -t -r cli call EchoService.GetHostnameServerStream
{
  "hostname": "grpc-echo-68b8797599-wlrkx"
}
```