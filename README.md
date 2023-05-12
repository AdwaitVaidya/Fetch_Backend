# Receipt Processor

Run following commands in order.
```shell script
docker build -t image .
```

To build docker image. Then to start docker container with binding to port 8080

```shell script
docker run --rm -p 8080:8080 image
```
Open another terminal and 
cd into folder with unit tests .
Run a basic unit test that adds a couple of receipts and checks if the recipets are added correctly and whther the points are correct.

```shell script
cd cmd/server/
go test -v
```

