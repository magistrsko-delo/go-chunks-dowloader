package main

import (
	"fmt"
	"github.com/heptiolabs/healthcheck"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/joho/godotenv"
	"go-chunks-dowloader/Models"
	"go-chunks-dowloader/router"
	"log"
	"net/http"
	"time"

	"github.com/opentracing-contrib/go-stdlib/nethttp"
	opentracing "github.com/opentracing/opentracing-go"
	jaeger "github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/zipkin"
)

func init()  {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
	Models.InitEnv()
}

func main()  {

	health := healthcheck.NewHandler()
	health.AddLivenessCheck("aws-service-check: " + Models.GetEnvStruct().AwsStorageUrl + "health", healthcheck.HTTPGetCheck(Models.GetEnvStruct().AwsStorageUrl + "health", 5*time.Second))

	zipkinPropagator := zipkin.NewZipkinB3HTTPHeaderPropagator()
	injector := jaeger.TracerOptions.Injector(opentracing.HTTPHeaders, zipkinPropagator)
	extractor := jaeger.TracerOptions.Extractor(opentracing.HTTPHeaders, zipkinPropagator)

	zipkinSharedRPCSpan := jaeger.TracerOptions.ZipkinSharedRPCSpan(true)

	sender, err := jaeger.NewUDPTransport(Models.GetEnvStruct().TracingConnection, 0)

	r := mux.NewRouter()

	api := r.PathPrefix("/v1").Subrouter()
	chunkDownloadRouter  := &router.ChunkDownloaderRouter{Router: api}
	chunkDownloadRouter.RegisterHandlers()

	r.NotFoundHandler = http.HandlerFunc(NotFound)

	corsOpts := cors.New(cors.Options{
		AllowedOrigins: []string{Models.GetEnvStruct().OriginAllowed},
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
			http.MethodOptions,
			http.MethodHead,
		},
		AllowedHeaders: []string{
			"X-Requested-With",
		},
	})

	go http.ListenAndServe("0.0.0.0:8888", health)

	if err == nil {
		fmt.Println("success: TRACING")
		tracer, closer := jaeger.NewTracer(
			"chunk-downloader",
			jaeger.NewConstSampler(true),
			jaeger.NewRemoteReporter(
				sender,
				jaeger.ReporterOptions.BufferFlushInterval(1*time.Second)),
			injector,
			extractor,
			zipkinSharedRPCSpan,
		)
		defer closer.Close()
		log.Fatal(http.ListenAndServe(":" + Models.GetEnvStruct().Port, nethttp.Middleware(tracer, corsOpts.Handler(r)))  )
	} else {
		fmt.Println( "err: ", err)
		log.Fatal(http.ListenAndServe(":" + Models.GetEnvStruct().Port, corsOpts.Handler(r)) )
	}

}


func NotFound(w http.ResponseWriter, r *http.Request) {
	rsp := "route not found: " + r.URL.Path
	w.Write([]byte(rsp))
}

