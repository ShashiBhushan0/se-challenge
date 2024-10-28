package main

import (
	"context"
	"fmt"
	"log"

	"github.com/shashibhushan06/aggregator/aggregator"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	var conn *grpc.ClientConn
	conn, err := grpc.NewClient(":8009", grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatal("Cannot connect", err)
	}
	fmt.Println("May be connected")
	defer conn.Close()
	agg := aggregator.NewTimeAggregatorServiceClient(conn)
	queryRquest := aggregator.QueryRequest{Start: 1727443800, End: 1730034000, Window: "1d", Aggregation: "SUM"}
	resp, err := agg.QueryData(context.Background(), &queryRquest)
	if err != nil {
		log.Fatal("Query Failed")
	}
	for _, data := range resp.Data {
		fmt.Println(data)
	}
}
