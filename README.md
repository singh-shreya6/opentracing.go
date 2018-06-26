# opentracing.go
Implementation of distributed open-tracing using Golang and Jaeger UI

In microservices architectures of modern application development, there are more applications communicating with each other than ever before. While application performance monitoring is great for debugging inside a single app, as a system expands into multiple services, how can you understand how much time each service is taking, where the exception happens, and the overall health of your system? In particular, how do you measure network latency between services—such as how long a request takes between one app to another?

Enter distributed tracing instrumentation. With the higher-level distribution of services that takes place in a cloud-based environment, tracing will become a key part of the cloud infrastructure supporting those services. In OpenTracing, a trace tells the story of a transaction or workflow as it propagates through a distributed system. The concept of the trace borrows a tool from the scientific community called a directed acyclic graph (DAG), which stages the parts of a process from a clear start to a clear end.

Here I have provided the code for opentracing in a distributed environment where a user is making a request on a browser at some port. This request to server 1 prints out Hello <name>! time is <current time>. After completing this operation, server 1 makes a request to server 2 on another port,by sending the name in request, the second server takes the name in its response body and displays Hello <name>!

So overall the following operations took place:
1) A browser makes a request to http://localhost:8082/name
2) The browser displays Hello name! time is Thu Jun  7 17:18:54 2018
3) Server 1 sends a POST request to Sever 2 with the name at http://localhost:8081
4) Server 2 captures this request in response body and displays Hello name!

Thus in all we get 3 spans, 1st span is of server-1, 2nd span is from server-1 to server-2 and 3rd span of server 2.
The opentracing package used is Jaeger.

Jaeger is run via Docker using the following command:
docker run -d -e COLLECTOR_ZIPKIN_HTTP_PORT=9411 -p5775:5775/udp -p6831:6831/udp -p6832:6832/udp   -p5778:5778 -p16686:16686 -p14268:14268 -p9411:9411 jaegertracing/all-in-one:latest

or it can also be run using:
docker run -d -p 5775:5775/udp -p 16686:16686 jaegertracing/all-in-one:latest

In Linux/Ubuntu, do remember to prefix sudo with these commands.
Now the Jaeger UI runs at http://localhost:16686 
All the codes in Golang.

