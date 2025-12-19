package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	connect_go "github.com/bufbuild/connect-go"
	"github.com/kelseyhightower/envconfig"
	apiv1 "github.com/ruhika006/funding-analytics/gen/api/v1"
	apiv1connect "github.com/ruhika006/funding-analytics/gen/api/v1/apiv1connect"
	"github.com/ruhika006/funding-analytics/handlers"
)

// Wrapper to adapt bare handlers to Connect interface
type ConnectServer struct {
	server *handlers.Server
}

func (cs *ConnectServer) GetRecords(ctx context.Context, req *connect_go.Request[apiv1.GetRecordsRequest]) (*connect_go.Response[apiv1.GetRecordsResponse], error) {
	resp, err := cs.server.GetRecords(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect_go.NewResponse(resp), nil
}

func (cs *ConnectServer) GetTopIndustries(ctx context.Context, req *connect_go.Request[apiv1.GetTopIndustriesRequest]) (*connect_go.Response[apiv1.GetTopIndustriesResponse], error) {
	resp, err := cs.server.GetTopIndustries(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect_go.NewResponse(resp), nil
}

func (cs *ConnectServer) GetTopIndustriesByYear(ctx context.Context, req *connect_go.Request[apiv1.GetTopIndustriesByYearRequest]) (*connect_go.Response[apiv1.GetTopIndustriesByYearResponse], error) {
	resp, err := cs.server.GetTopIndustriesByYear(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect_go.NewResponse(resp), nil
}

func (cs *ConnectServer) GetTopFundedStartups(ctx context.Context, req *connect_go.Request[apiv1.GetTopFundedStartupsRequest]) (*connect_go.Response[apiv1.GetTopFundedStartupsResponse], error) {
	resp, err := cs.server.GetTopFundedStartups(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect_go.NewResponse(resp), nil
}

func (cs *ConnectServer) GetTopCityAndIndustries(ctx context.Context, req *connect_go.Request[apiv1.GetTopCityAndIndustriesRequest]) (*connect_go.Response[apiv1.GetTopCityAndIndustriesResponse], error) {
	resp, err := cs.server.GetTopCityAndIndustries(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect_go.NewResponse(resp), nil
}

type Specification struct {
	User string
	Password string
	Addr string
}

func main() {
	// Get configuration from environment
	var spec Specification
	if err := envconfig.Process("myapp", &spec); err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize ClickHouse connection
	chClient, err := initClickHouseClient(spec)
	if err != nil {
		log.Fatalf("Failed to initialize ClickHouse client: %v", err)
	}
	defer chClient.Close()

	// Verify ClickHouse connectivity
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	if err := chClient.Ping(ctx); err != nil {
		cancel()
		log.Fatalf("ClickHouse ping failed: %v", err)
	}
	cancel()

	log.Println("Connected to ClickHouse")

	mux := http.NewServeMux()

	// Create server with ClickHouse client
	server := handlers.NewServer(chClient)
	connectServer := &ConnectServer{server: server}
	path, handler := apiv1connect.NewQueryServiceHandler(connectServer)
	mux.Handle(path, handler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	addr := ":" + port
	log.Printf("Server listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}


// initClickHouseClient initializes a ClickHouse connection
func initClickHouseClient(spec Specification) (clickhouse.Conn, error) {
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{spec.Addr},
		Auth: clickhouse.Auth{
			Username: spec.User,
			Password: spec.Password,
		},
		Settings: clickhouse.Settings{
			"max_execution_time": 60,
		},
	})
	if err != nil {
		return nil, err
	}
	return conn, nil
}