package main

import (
	"cloud.google.com/go/compute/metadata"
	"context"
	"fmt"
	pb "github.com/kazshinohara/pb/grpc-echo"
	"google.golang.org/grpc"
	meta "google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

var (
	port    = os.Getenv("PORT")
	version = os.Getenv("VERSION")
	kind    = os.Getenv("KIND")
)

type EchoServiceServer struct {
	pb.UnimplementedEchoServiceServer
}

func (s *EchoServiceServer) GetAll(ctx context.Context, empty *emptypb.Empty) (*pb.All, error) {
	region := resolveRegion()
	cluster := resolveCluster()
	hostname := resolveHostname()
	sourceIp := resolveSourceIp(ctx)

	return &pb.All{
		Kind:     kind,
		Version:  version,
		Region:   region,
		Cluster:  cluster,
		Hostname: hostname,
		SourceIp: sourceIp,
	}, nil
}

func (s *EchoServiceServer) GetKind(ctx context.Context, empty *emptypb.Empty) (*pb.Kind, error) {
	return &pb.Kind{
		Kind: kind,
	}, nil
}

func (s *EchoServiceServer) GetVersion(ctx context.Context, empty *emptypb.Empty) (*pb.Version, error) {
	return &pb.Version{
		Version: version,
	}, nil
}

func (s *EchoServiceServer) GetRegion(ctx context.Context, empty *emptypb.Empty) (*pb.Region, error) {
	region := resolveRegion()
	return &pb.Region{
		Region: region,
	}, nil
}

func (s *EchoServiceServer) GetCluster(ctx context.Context, empty *emptypb.Empty) (*pb.Cluster, error) {
	cluster := resolveCluster()
	return &pb.Cluster{
		Cluster: cluster,
	}, nil
}

func (s *EchoServiceServer) GetHostname(ctx context.Context, empty *emptypb.Empty) (*pb.Hostname, error) {
	hostname := resolveHostname()
	return &pb.Hostname{
		Hostname: hostname,
	}, nil
}

func (s *EchoServiceServer) GetSourceIp(ctx context.Context, empty *emptypb.Empty) (*pb.SourceIp, error) {
	sourceIp := resolveSourceIp(ctx)
	return &pb.SourceIp{
		SourceIp: sourceIp,
	}, nil
}

func (s *EchoServiceServer) GetHeader(ctx context.Context, hn *pb.HeaderName) (*pb.HeaderValue, error) {
	var values []string
	if md, ok := meta.FromIncomingContext(ctx); ok {
		values = md.Get(hn.RequestHeaderName)
	}
	if len(values) > 0 {
		return &pb.HeaderValue{
			RequestHeaderValue: values[0],
		}, nil
	}
	return &pb.HeaderValue{
		RequestHeaderValue: "unknown",
	}, nil

}

func (s *EchoServiceServer) GetHostnameServerStream(conf *pb.ServerStreamConfig, stream pb.EchoService_GetHostnameServerStreamServer) error {
	hostname := resolveHostname()
	for i := 0; i < int(conf.NumberOfResponse); i++ {
		if err := stream.Send(&pb.Hostname{
			Hostname: hostname,
		}); err != nil {
			return  err
		}
		time.Sleep(time.Second * time.Duration(conf.Interval))
	}
	return nil
}

//TODO: Implement Client streaming RPC & Bi-directional streaming RPC

func resolveRegion() string {
	if !metadata.OnGCE() {
		log.Println("This app is not running on GCE")
	} else {
		zone, err := metadata.Zone()
		if err != nil {
			log.Printf("could not get zone info: %v", err)
			return "unknown"
		}
		region := zone[:strings.LastIndex(zone, "-")]
		return region
	}
	return "unknown"
}

func resolveCluster() string {
	if !metadata.OnGCE() {
		log.Println("This app is not running on GCE")
	} else {
		cluster, err := metadata.Get("/instance/attributes/cluster-name")
		if err != nil {
			log.Printf("could not get cluster name: %v", err)
			return "unknown"
		}
		return cluster
	}
	return "unknown"
}

func resolveHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		fmt.Printf("could not get hostname: %v", err)
		return "unknown"
	}
	return hostname
}

func resolveSourceIp(ctx context.Context) string {
	var values []string
	if md, ok := meta.FromIncomingContext(ctx); ok {
		values = md.Get("X-Forwarded-For")
		if len(values) > 0 {
			return values[0]
		}
	}

	var addr string
	if pr, ok := peer.FromContext(ctx); ok {
		if tcpAddr, ok := pr.Addr.(*net.TCPAddr); ok {
			addr = tcpAddr.IP.String()
			return addr
		} else {
			addr = pr.Addr.String()
			return addr
		}
	}
	return "unknown"
}

func main() {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterEchoServiceServer(s, &EchoServiceServer{})
	reflection.Register(s)
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
