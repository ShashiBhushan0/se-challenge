### Edgecom Code Challenge

#### Problem
Your task is to build a time-series data service that will:

- **Fetch time-series data** from a provided API.
- **Store the data** in a database.
- **Update the data** every 5 minutes using a scheduler.
- Expose the data via a **gRPC** API.
- At startup, **bootstrap the service** with the last 2 years' worth of data.
- Use **Docker** and **docker-compose** to run the service.

#### Key Requirements
1. **Fetching Data**:

   - Scrape time-series data from a provided API.
   - Store the data in a database. (Using a time-series database is a bonus.)
   - The API should collect time-series data at regular intervals.

2. **Database**:
   - Choose a suitable database for storing time-series data. Options include: PostgreSQL, InfluxDB, QuestDB and TDengine.

3. **Scheduler**:
   - The service should have a scheduler to fetch and store new data every 5 minutes.

4. **Data Bootstrapping**:
   - On startup, the service should fetch the last 2 years' worth of data from the API to populate the database with historical data.

5. **gRPC API**:

   - Create a gRPC service that exposes an API to query the time-series data.
   - The API should accept the following parameters:
      - **start**: the start timestamp for the data query.
      - **end**: the end timestamp for the data query.
      - **window**: time interval for grouping the data (e.g., 1m, 5m, 1h, 1d).
      - **aggregation**: type of aggregation to apply on the data.

6. Aggregations:
   - Implement the following aggregations for data queries:
      - **MIN**: Minimum value within each time window.
      - **MAX**: Maximum value within each time window.
      - **AVG**: Average value within each time window.
      - **SUM**: Sum of all values within each time window.

7. Docker and Docker Compose:
   - Package the service using Docker.
   - Provide a docker-compose.yml file to run the service with all dependencies (e.g., database).

#### Technologies
The accepted programming languages are: Go, Python or Rust. You can use any libraries or frameworks that you think are necessary to complete the task.

#### API Description
You will be using the following API to fetch time-series data:
```bash
curl --location 'https://api.edgecomenergy.net/core/asset/3662953a-1396-4996-a1b6-99a0c5e7a5de/series?start=2024-09-13T00:00:00&end=2024-09-17T00:00:00' 
```

Response:
```json
{
    "result": [
        {
            "time": 1694564100,
            "value": 533.076
        },
        {
            "time": 1694565000,
            "value": 619.056
        },
        {
            "time": 1694565900,
            "value": 653.4480000000001
        },
        {
            "time": 1694566800,
            "value": 619.056
        }
    ]
}
```

This API provides time-series data for a specific asset. The API returns data points within a given date range. The query parameters are as follows:

- **start**: The start timestamp for the data range (ISO 8601 format).
- **end**: The end timestamp for the data range (ISO 8601 format).

You can use this API to retrieve historical data during the initial bootstrapping process as well as during regular intervals (every 5 minutes) to keep your database updated with the latest time-series data.

