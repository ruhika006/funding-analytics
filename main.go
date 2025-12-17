package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"

	connect_go "github.com/bufbuild/connect-go"
	_ "github.com/duckdb/duckdb-go/v2"
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

func main() {
	db, err := sql.Open("duckdb", "")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	fmt.Println("Opened up DuckDB")

	_, err = db.Exec(`
		CREATE TABLE startup AS
		SELECT * FROM read_csv('./app/startup_funding.csv')
	`)
	if err != nil {
		log.Fatal("Unable to load into table", err)
	}

	mux := http.NewServeMux()

	// Create server with DuckDB client
	server := handlers.NewServer(db)
	connectServer := &ConnectServer{server: server}
	path, handler := apiv1connect.NewQueryServiceHandler(connectServer)
	mux.Handle(path, handler)

	log.Println("Server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
