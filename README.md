# Receipt Processor

Run following commands in order.
'''
docker build -t image .
'''

To build docker image. Then to start docker container with binding to port 8080

'''
docker run --rm -p 8080:8080 image
'''

Run a basic unit test that adds a couple of receipts and checks if the recipets are added correctly and whther the points are correct.

'''
go test -v
'''
