package mailsender

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/cors"
)

const (
	wrongMethod         = "only POST method is supported"
	badRequest          = "bad request"
	internalServerError = "internal server error"
)

type Server struct {
	http *http.Server
}

func NewServer(listenAddr, smtpAddr, fromMail string) *Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/contact", MailHandler(NewSmtp(smtpAddr), fromMail))

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"POST"},
		AllowedHeaders:   []string{"*"},
		MaxAge:           10,
		AllowCredentials: true,
	}).Handler(mux)

	return &Server{
		http: &http.Server{
			Addr:    listenAddr,
			Handler: c,
		},
	}
}

func (s *Server) Serve() {
	fmt.Printf("listen and serve %s\n", s.http.Addr)

	go s.http.ListenAndServe()

	// Gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGKILL)
	<-quit

	s.http.Close()
	fmt.Println("Stop listening server")
}
