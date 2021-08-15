package app

import (
	"context"
	"cryptocurrencies-api/config"
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
	srv    *http.Server
}

func (s *server) configureRouter() {
	s.router = mux.NewRouter()
	userRouter := s.router.PathPrefix("/" + config.ServerConfig.Version + "/user/").Methods("POST").Subrouter()
	userRouter.Use(middlewares.Cors)
	userRouter.HandleFunc("/create", controllers.CreateAccount)
	userRouter.HandleFunc("/login", controllers.LoginAccount)

	refreshRouter := s.router.PathPrefix("/" + config.ServerConfig.Version + "/user/").Methods("GET").Subrouter()
	refreshRouter.Use(middlewares.Cors)
	refreshRouter.Use(middlewares.JwtValidation)
	refreshRouter.Use(middlewares.JwtRefreshValidation)
	refreshRouter.HandleFunc("/refreshToken", controllers.RefreshToken)

	privateRouter := s.router.PathPrefix("/" + config.ServerConfig.Version + "/cryptocurrency/").Methods("GET").Subrouter()
	privateRouter.Use(middlewares.Cors)
	privateRouter.Use(middlewares.JwtValidation)
	privateRouter.HandleFunc("/btcRate", controllers.BtcRate)

	s.walkRouters()
}

func (s *server) walkRouters() {
	if err := s.router.Walk(
		func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
			if pathTemplate, err := route.GetPathTemplate(); err == nil {
				fmt.Println("ROUTE:", pathTemplate)
			}else {return err}
			if pathRegexp, err := route.GetPathRegexp(); err == nil {
				fmt.Println("Path regexp:", pathRegexp)
			}else {return err}
			if queriesTemplates, err := route.GetQueriesTemplates(); err == nil {
				fmt.Println("Queries templates:", strings.Join(queriesTemplates, ","))
			}else {return err}
			if queriesRegexps, err := route.GetQueriesRegexp(); err == nil {
				fmt.Println("Queries regexps:", strings.Join(queriesRegexps, ","))
			}else {return err}
			if methods, err := route.GetMethods(); err == nil {
				fmt.Println("Methods:", strings.Join(methods, ","))
			}else {return err}
			fmt.Println()
			return nil
		}); err != nil {
		fmt.Println(err)
	}
}

func createServer() (s *server) {
	s = &server{
		router: mux.NewRouter(),
	}

	s.configureRouter()

	s.srv = &http.Server{
		Addr:         ":" + config.ServerConfig.Port,
		WriteTimeout: time.Second * 5,
		ReadTimeout:  time.Second * 5,
		//IdleTimeout:  time.Second * 60,
		Handler: s.router,
	}
	return
}

func RunServer() {
	s := createServer()

	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
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

	err := s.srv.Shutdown(ctx)
	if err != nil {
		return
	}

	log.Println("shutting down")
	os.Exit(0)
}
