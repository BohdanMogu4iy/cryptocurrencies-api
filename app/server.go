package app

import (
	"context"
	"cryptocurrencies-api/controllers"
	"cryptocurrencies-api/middlewares"
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"
)

type server struct {
	router *mux.Router
	srv *http.Server
}

func (s *server) configureRouter() {
	s.router = mux.NewRouter()

	userRouter := s.router.PathPrefix("/v1/user/").Methods("POST").Subrouter()
	userRouter.HandleFunc("/create", controllers.CreateAccount)
	userRouter.HandleFunc("/auth", controllers.Authenticate)

	privateRouter := s.router.PathPrefix("/private/").Subrouter()
	privateRouter.Use(middlewares.JwtAuthentication)
	//privateRouter.HandleFunc("/logout", s.handlerLogoutRequest()).Methods("GET")
	privateRouter.HandleFunc("/test", controllers.TestController).Methods("GET")

	err := s.router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		pathTemplate, err := route.GetPathTemplate()
		if err == nil {
			fmt.Println("ROUTE:", pathTemplate)
		}
		pathRegexp, err := route.GetPathRegexp()
		if err == nil {
			fmt.Println("Path regexp:", pathRegexp)
		}
		queriesTemplates, err := route.GetQueriesTemplates()
		if err == nil {
			fmt.Println("Queries templates:", strings.Join(queriesTemplates, ","))
		}
		queriesRegexps, err := route.GetQueriesRegexp()
		if err == nil {
			fmt.Println("Queries regexps:", strings.Join(queriesRegexps, ","))
		}
		methods, err := route.GetMethods()
		if err == nil {
			fmt.Println("Methods:", strings.Join(methods, ","))
		}
		fmt.Println()
		return nil
	})

	if err != nil {
		fmt.Println(err)
	}
}

func (s *server) init() {
	s.configureRouter()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000" //localhost
	}

	s.srv = &http.Server{
		Addr: ":" + port,
		WriteTimeout: time.Second * 5,
		ReadTimeout:  time.Second * 5,
		//IdleTimeout:  time.Second * 60,
		Handler: s.router,
	}
}

func createServer() (s *server){
	s = &server{
		router:       mux.NewRouter(),
	}
	s.init()
	return
}

func RunServer(){
	s:=createServer()

	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second * 15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()

	go func() {
		log.Printf("Server is listening on %v", s.srv.Addr)
		if err := s.srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	<-c

	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()

	s.srv.Shutdown(ctx)

	log.Println("shutting down")
	os.Exit(0)
}


