package education

import (
	"google.golang.org/grpc"
)

type Client struct {
	conn *grpc.ClientConn
}
