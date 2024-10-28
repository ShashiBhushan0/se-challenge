package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	"database/sql"
	"log"

	_ "github.com/lib/pq" // PostgreSQL driver
	"github.com/shashibhushan06/aggregator/aggregator"
	"google.golang.org/grpc"
)

type DataPoint struct {
	Time  int64   `json:"time"`
	Value float64 `json:"value"`
}

type SeriesResponse struct {
	Result []DataPoint `json:"result"`
}

func main() {
	StartService()
}

func StartService() {

	fmt.Println("Initiating DB connection...")

	db := getDbConnection()
	defer db.Close()

	lastFetchTime := bootstrapData(db)
	go keepUpdatingData(db, lastFetchTime)
	grpc_port := os.Getenv("GRPC_PORT")

	lis, err := net.Listen("tcp", ":"+grpc_port)
	if err != nil {
		log.Fatal("Listen failed", err)
	}
	grpcServer := grpc.NewServer()

	service := &aggregator.Server{DB: db}
	aggregator.RegisterTimeAggregatorServiceServer(grpcServer, service)
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatal("Server failed", err)
	}
	fmt.Printf("Done!")
}

func keepUpdatingData(db *sql.DB, lastFetchTime time.Time) {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()
	run := false
	if run {
		fmt.Println(db)
	}

	for range ticker.C {
		fmt.Println("Executing function every 5 minutes...")
		currentTime := time.Now()
		// Call your desired function here
		data := getTimeData(lastFetchTime, currentTime)
		fmt.Println("Total responses = ", len(data.Result))
		insertIntoDB(db, &data)
		lastFetchTime = currentTime
		fmt.Println("Data points inserted successfully", lastFetchTime)
	}
}

func bootstrapData(db *sql.DB) time.Time {

	currentTime := time.Now()

	queryString := "select max(time) from data_points dp"
	fmt.Println(queryString)
	resultRows, err := db.Query(queryString)
	if err != nil {
		log.Fatal("Query failed", err)
	}
	defer resultRows.Close()
	lastTimestamp := currentTime.AddDate(-2, 0, 0)
	for resultRows.Next() {
		var val sql.NullInt64
		if err := resultRows.Scan(&val); err != nil {
			fmt.Println("SQL rows Scan Failed", err)
			break
		}
		if val.Int64 != 0 {
			lastTimestamp = time.Unix(val.Int64, 0)
			fmt.Println("Data found till ", lastTimestamp.Format("2006-01-02T15:04:05"))
		} else {
			fmt.Println("Database is empty. Boostrapping")
		}
	}

	fmt.Println("Updating data from ", lastTimestamp.Format("2006-01-02T15:04:05"), "to", currentTime.Format("2006-01-02T15:04:05"))

	data := getTimeData(lastTimestamp, currentTime)
	fmt.Println("Total responses = ", len(data.Result))

	insertIntoDB(db, &data)
	fmt.Println("Data Updated")
	return currentTime

}

func getTimeData(startDate time.Time, endDate time.Time) SeriesResponse {

	var data SeriesResponse

	// startDate := time.Date(2024, 9, 13, 9, 0, 0, 0, time.UTC)
	// endDate := time.Date(2024, 9, 13, 10, 0, 0, 0, time.UTC)
	// data := getTimeData(startDate, endDate)

	url := fmt.Sprintf("https://api.edgecomenergy.net/core/asset/3662953a-1396-4996-a1b6-99a0c5e7a5de/series?start=%s&end=%s",
		startDate.Format("2006-01-02T15:04:05"),
		endDate.Format("2006-01-02T15:04:05"))

	client := &http.Client{}

	// Create a new request with appropriate headers
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return data
	}

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return data
	}
	defer resp.Body.Close()

	// Check for successful response
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error: unexpected status code %d\n", resp.StatusCode)
		return data
	}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		fmt.Println("Error: Failed", err)
	}
	return data
}

func getDbConnection() *sql.DB {
	// Create a database connection
	time.Sleep(10 * time.Second)
	connStr := "postgres://" + os.Getenv("POSTGRES_USER") + ":" + os.Getenv("POSTGRES_PASSWORD") + "@postgres:" + os.Getenv("POSTGRES_PORT") + "/" + os.Getenv("POSTGRES_DB") + "?sslmode=disable"
	var db *sql.DB
	var err error

	db, err = sql.Open("postgres", connStr)

	if err != nil {
		log.Fatal("Error connecting to DB", err)
	}

	// Create a table to store the data points (if it doesn't exist)
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS data_points (time BIGINT, value NUMERIC)")
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func insertIntoDB(db *sql.DB, data *SeriesResponse) {
	// Insert data points into the database
	for _, point := range data.Result {
		_, err := db.Exec("INSERT INTO data_points (time, value) VALUES ($1, $2)", point.Time, point.Value)
		if err != nil {
			log.Println("Error inserting data point:", err)
		}
	}
}

func DemoTest() {

	epoch := int64(1667011200)

	// Convert epoch to time.Time object
	t := time.Unix(epoch, 0) // nanoseconds can be set to a specific value if needed

	// Format the time using desired layout string
	formattedTime := t.Format("2006-01-02T15:04:05") // Replace with your desired layout

	currentTime := time.Now()
	// Calculate 2 years ago
	twoYearsAgo := currentTime.AddDate(-2, 0, 0).Format("2006-01-02T15:04:05")

	fmt.Println("Original epoch:", epoch)
	fmt.Println("Converted time:", formattedTime)
	fmt.Println(twoYearsAgo)
}
