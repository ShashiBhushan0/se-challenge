package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/shashibhushan06/aggregator/aggregator"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {

	grpc_port := "8011"

	reader := bufio.NewReader(os.Stdin)
	startTime := int64(0)
	var err error
	for startTime == 0 {
		fmt.Print("Enter start time : ")
		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)
		startTime, err = strconv.ParseInt(text, 10, 64)
		if err != nil {
			fmt.Println("Parsing error. Please retry")
		}
	}
	endTime := int64(0)
	for endTime == 0 {
		fmt.Print("Enter end time : ")
		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)
		endTime, err = strconv.ParseInt(text, 10, 64)
		if err != nil {
			fmt.Println("Parsing error. Please retry")
		}
	}

	window := ""
	for window == "" {
		fmt.Print("Enter window time: (eg. 1m, 1h, 1d) : ")
		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)
		if len(text) < 2 {
			fmt.Println("Invalid. Please retry")
			continue
		}
		timeValue := text[:len(text)-1]
		timeType := text[len(text)-1:]
		_, err := strconv.ParseInt(timeValue, 10, 64)
		if err != nil {
			fmt.Println("Invalid. Please retry")
			continue
		}

		switch timeType {
		case "m":
		case "h":
		case "d":
		default:
			fmt.Println("Invalid. Please retry")
			continue
		}
		window = text
	}

	aggregation := ""
	for aggregation == "" {
		fmt.Print("Enter aggregation: (eg. MIN, MAX, AVG, SUM) : ")
		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)
		switch strings.ToUpper(text) {
		case "MIN":
		case "MAX":
		case "AVG":
		case "SUM":
		default:
			fmt.Println("Invalid. Please retry")
			continue
		}
		aggregation = text
	}

	var conn *grpc.ClientConn
	conn, err = grpc.NewClient(":"+grpc_port, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatal("Cannot connect", err)
	}
	defer conn.Close()
	agg := aggregator.NewTimeAggregatorServiceClient(conn)
	queryRquest := aggregator.QueryRequest{Start: startTime, End: endTime, Window: window, Aggregation: aggregation}
	resp, err := agg.QueryData(context.Background(), &queryRquest)
	if err != nil {
		log.Fatal("Query Failed", err)
	}
	for _, data := range resp.Data {
		fmt.Println(data)
	}
}
