# `REST over gRPC with grpc-gateway for Go`
This are the source codes for the medium.com blog [post](https://medium.com/@arkadybalaba/rest-over-grpc-with-grpc-gateway-for-go-9584bfcbb835).

Follow the blog post for instructions and details.

## To run
Run server from `cmd` folder with: `GO111MODULE=on go run main.go`

Run gRPC client from `cmd/client` folder with: `GO111MODULE=on go run client.go`

Call the HTTP/REST endpoint with a client of you choice.
Example curl request:
`curl --request PUT 
  --url https://localhost:8080/v1/reminder/schedule 
  --header 'content-type: application/json'
  --data '{
	"when": "2019-09-20T02:50:20Z"
}'`
