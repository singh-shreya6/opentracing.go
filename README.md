Opentracing using Go

Components:

There are 3 folders in this opentracing package:
1.	server_1 containing main.go
2.	server_2 containing main.go
3.	go-opentracing containing opentracing_helper.go

Jaeger UI: Jaeger UI is run via Docker using the following command: 
docker run -d -e COLLECTOR_ZIPKIN_HTTP_PORT=9411 -p5775:5775/udp -p6831:6831/udp -p6832:6832/udp -p5778:5778 -p16686:16686 -p14268:14268 -p9411:9411 jaegertracing/all-in-one:latest
or it can also be run using:
 docker run -d -p 5775:5775/udp -p 16686:16686 jaegertracing/all-in-one:latest
The Jaeger UI runs at http://localhost:16686

Functionality:

The following operations are taking place.
1.	A browser makes a request to http://localhost:8082/name
2.	The browser displays Hello name! time is Thu Jun 7 17:18:54 2018
3.	Server 1 sends a POST request to Sever 2 with the name at http://localhost:8081
4.	Server 2 captures this request in response body and displays Hello name!

Functions in code:
1)	main.go of server 1
a) func main:
i)	It calls the opentracing initializer function present in opentracing_helper.go
ii)	It creates a server listening at port 8082
iii)	It calls the opentracing middleware function to introduce a span before entering the function handle

b)	func handle:
arguments: http Response Writer ,http  Request
return type: none
i)	It calculates time
ii)	It takes name from URL
iii)	It calls the function makeRequest passing it the contextual information,URL and name
iv)	It prints out name and time

c)	func makeRequest:
arguments: Context, server URL string, name string
return type: none
i)	It calls a function Introduce_span present in opentracing_helper.go which starts a span
ii)	It makes a POST request to port 8081 with name
iii)	It calls the function Serialise present in opentracing_helper.go to serialise the request.
iv)	It creates a client
v)	It reads the response it gets from server at 8081
vi)	It prints the name by taking it out of response body

2)	main.go of server 2
a)	func main:
i)	It initializes opentracing for server 2
ii)	It creates a server listening at port 8081
iii)	It calls the opentracing middleware function to introduce a span before entering the function handle

b)	func handle:
i)	It calls the Deserialize function present in opentracing_helper.go
ii)	It reads the request body coming through POST request of server 1
iii)	It takes the name from the request body 
iv)	It displays the name

3)	Opentracing_helper.go
a)	func Init:
arguments: name of tracer string
return type: none
i)	It initializes the jaeger tracing
ii)	It creates a sampler and reporter required for initializing tracer of Jaeger
iii)	It initializes the tracer using the name of tracer, sampler and reporter
iv)	It set the tracer to Global Tracer

b)	func Introduce_Span:
arguments: .Context, name of span string
return type: Span, Context
i)	It takes the context information and the name of span as its arguments
ii)	It returns the created span and its context information

c)	func Serialise:
arguments: Context, http Request
return type: none
i)	It transmit the span's Trace Context as HTTP headers on the outbound request from server1

d)	func Deserialize:
arguments: http Request, name of span string
return type: Span, http Request
i)	It takes the raw data and reconstructs the byte stream
ii)	It returns the span as well as the http request

e)	func HttpMiddleware:
arguments: name of server string, h http Handler
return type: handler
i)	It starts a span from context
ii)	It returns the handler function

Output:
In all we get 3 spans for opentracing:
1.	 1st span is of server-1
2.	 2nd span is from server-1 to server-2
3.	 3rd span of server-2

