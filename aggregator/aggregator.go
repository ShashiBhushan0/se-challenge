package aggregator

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"
)

type Server struct {
	DB *sql.DB
	UnimplementedTimeAggregatorServiceServer
}

func (s *Server) query(startTime int64, endTime int64, aggregationWindow string, aggregationType string) (*QueryResponse, error) {
	timeValue := aggregationWindow[:len(aggregationWindow)-1]
	timeType := aggregationWindow[len(aggregationWindow)-1:]

	aggregationType = strings.ToUpper(aggregationType)
	fmt.Println(aggregationType)

	timeDiff := int64(0)
	val, err := strconv.ParseInt(timeValue, 10, 64)
	if err != nil {
		fmt.Println("Parsing error")
	}

	switch timeType {
	case "m":
		// fmt.Println(timeValue, "Minutes")
		timeDiff += val * 60
	case "h":
		// fmt.Println(timeValue, "Hours")
		timeDiff += val * 60 * 60
	case "d":
		// fmt.Println(timeValue, "Days")
		timeDiff += val * 24 * 60 * 60
	}
	dps := []*DataPoint{}
	for nextStartTime, nextEndTime := startTime, startTime+timeDiff; nextStartTime < endTime; nextStartTime, nextEndTime = nextStartTime+timeDiff, nextEndTime+timeDiff {
		if nextEndTime > endTime {
			nextEndTime = endTime
		}
		// queryString := "select " + aggregationType + "(time) from data_points dp where time between " + startTime + " and " + endTime
		queryString := fmt.Sprintf("select %v(value) from data_points dp where time between %v and %v", aggregationType, nextStartTime, nextEndTime)
		fmt.Println(queryString)
		resultRows, err := s.DB.Query(queryString)
		if err != nil {
			fmt.Println("Query Failed")
		}
		defer resultRows.Close()
		for resultRows.Next() {
			var val sql.NullFloat64
			if err := resultRows.Scan(&val); err != nil {
				fmt.Println("SQL rows Scan Failed")
			}
			result := val.Float64
			fmt.Println("Result = ", int64(val.Float64))
			s1 := DataPoint{StartTime: nextStartTime, EndTime: nextEndTime, Value: result}
			dps = append(dps, &s1)

		}
	}
	// s1 := DataPoint{StartTime: 21, EndTime: 22, Value: 20}
	// dps = append(dps, &s1)
	qresp := QueryResponse{Data: dps}
	return &qresp, nil
}

func (s *Server) QueryData(ctx context.Context, queryRequest *QueryRequest) (*QueryResponse, error) {
	fmt.Println("Got a query")
	startTime := queryRequest.Start
	endTime := queryRequest.End
	aggregationWindow := queryRequest.Window
	aggregationType := queryRequest.Aggregation

	qresp, err := s.query(startTime, endTime, aggregationWindow, aggregationType)
	if err != nil {
		log.Fatal("Query Failed")
	}

	return qresp, nil
}
