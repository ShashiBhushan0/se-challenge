package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"database/sql"
	"log"

	_ "github.com/lib/pq" // PostgreSQL driver
)

// type SeriesResponse struct {
// 	// Add any specific fields you expect in the response here
// 	// based on the API documentation. For example:
// 	// Timestamps []time.Time `json:"timestamps"`
// 	// Values     []float64  `json:"values"`
// }

type DataPoint struct {
	Time  int64   `json:"time"`
	Value float64 `json:"value"`
}

type SeriesResponse struct {
	Result []DataPoint `json:"result"`
}

func main() {

	// url := "https://api.edgecomenergy.net/core/asset/3662953a-1396-4996-a1b6-99a0c5e7a5de/series?start=2024-09-13T00:00:00&end=2024-09-17T00:00:00"

	// Define the desired date range

	bootstrapData()
	// Code()

}

func Code() {
	db := getDbConnection()
	defer db.Close()

	startDate := time.Date(2024, 9, 13, 9, 0, 0, 0, time.UTC)
	endDate := time.Date(2024, 9, 13, 10, 0, 0, 0, time.UTC)
	data := getTimeData(startDate, endDate)

	insertIntoDB(db, &data)

	fmt.Println("Data points inserted successfully!")
}

func bootstrapData() time.Time {

	currentTime := time.Now()
	// Calculate 2 years ago
	twoYearsAgo := currentTime.AddDate(-2, 0, 0)
	// twoYearsAgo := currentTime.AddDate(0, 0, -2)

	fmt.Println(twoYearsAgo.Format("2006-01-02T15:04:05"))
	twoYearsAgo = twoYearsAgo.Add(5 * time.Minute)
	fmt.Println(twoYearsAgo.Format("2006-01-02T15:04:05"))
	data := getTimeData(twoYearsAgo, currentTime)
	fmt.Println("Total responses = ", len(data.Result))
	return currentTime

}

func getTimeData(startDate time.Time, endDate time.Time) SeriesResponse {
	var data SeriesResponse

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
		log.Fatal(err)
	}
	return data
}

func getDbConnection() *sql.DB {
	// Create a database connection
	connStr := "postgres://shashi:mysecretpassword@localhost:5432/postgres?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	// defer db.Close()

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
