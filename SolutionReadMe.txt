This gRPC server will bootstrap the data for last 2 years. In the event of service restart it will fetch data from the last pulled data.

Steps to run the project
1.  Update docker-compose.yml for the following
    a. update the volume as per your file system.
    b. if updating the grpc port please update the expose as well.
2.  $ docker build -t my-grpc-server .
3.  $ docker-compose up -d

In order to test run the sample client
example.
$ go run .\sampleclient\client.go
Enter start time : 1727787000
Enter end time : 1727789400
Enter window time: (eg. 1m, 1h, 1d) : 10m
Enter aggregation: (eg. MIN, MAX, AVG, SUM) : SUM
startTime:1727787000 endTime:1727787600 value:8064.924000000001
startTime:1727787600 endTime:1727788200 value:8064.924000000001
startTime:1727788200 endTime:1727788800 value:8082.120000000001
startTime:1727788800 endTime:1727789400 value:8047.728000000001


NOTE: This is a demo code and can be cleaned up and configured as per requirement.