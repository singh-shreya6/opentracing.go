package go_opentracing

import (
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"net/http"
	"context"
	"github.com/opentracing/opentracing-go/ext"
	"log"
	"time"
)


func Init(s string) {
	// init open tracing
	udpTransport, _ := jaeger.NewUDPTransport("localhost:5775", 0)
	reporter:= jaeger.NewRemoteReporter(udpTransport)
	sampler := jaeger.NewConstSampler(true)
	tracer, _ := jaeger.NewTracer(s, sampler, reporter)
	opentracing.SetGlobalTracer(tracer)
}

func Introduce_span(ctx context.Context,spanName string)(opentracing.Span,context.Context) {
	span,ctx:=opentracing.StartSpanFromContext(ctx, spanName)
	return span,ctx
}


func Serialise(ctx context.Context,req *http.Request){
	req = req.WithContext(ctx)
	if span := opentracing.SpanFromContext(ctx); span != nil {
		opentracing.GlobalTracer().Inject(
			span.Context(),
			opentracing.HTTPHeaders,
			opentracing.HTTPHeadersCarrier(req.Header))
	}
}

func Deserialize(r *http.Request,spanName string) (opentracing.Span,*http.Request){
	time.Sleep(250 * time.Millisecond)
	var serverSpan opentracing.Span
	appSpecificOperationName := spanName
	wireContext, err := opentracing.GlobalTracer().Extract(
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(r.Header))
	if err != nil {
		log.Fatal("Error")
	}
	serverSpan = opentracing.StartSpan(
		appSpecificOperationName,
		ext.RPCServerOption(wireContext))
	ctx := opentracing.ContextWithSpan(context.Background(), serverSpan)
	r=r.WithContext(ctx)
	return serverSpan,r
}

func HttpMiddleware(serverName string, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		span, newCtx := opentracing.StartSpanFromContext(r.Context(), serverName)
		defer func() {
			span.Finish()
		}()
		r = r.WithContext(newCtx)
		h.ServeHTTP(w, r)
	})
}

