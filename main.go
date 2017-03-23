package main

import (
	"log"
	"net/http"
	"net/url"
	"os"

	gh "github.com/google/go-github/github"
)

type PagesServer struct {
	webhookSecret []byte
	accessToken   string
}

func NewPagesServer(secret []byte, token string) *PagesServer {
	ps := &PagesServer{webhookSecret: secret, accessToken: token}
	return ps
}

func (ps *PagesServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	whtype, event, err := ps.parseWebhook(r)
	if err != nil {
		errorResponse(w, r, err)
		return
	}

	switch event := event.(type) {
	case *gh.PushEvent:
		log.Printf("pushed %s to %s", event.HeadCommit.GetID(), event.Repo.GetFullName())
	default:
		log.Printf("ignoring %s", whtype)
	}

	w.WriteHeader(http.StatusAccepted)
}

func (ps *PagesServer) parseWebhook(r *http.Request) (string, interface{}, error) {
	payload, err := gh.ValidatePayload(r, ps.webhookSecret)
	if err != nil {
		return "", nil, err
	}

	whtype := gh.WebHookType(r)

	event, err := gh.ParseWebHook(whtype, payload)
	if err != nil {
		return "", nil, err
	}

	return whtype, event, nil
}

func main() {
	webhookSecret := []byte(os.Getenv("WEBHOOK_SECRET"))
	accessToken := os.Getenv("ACCESS_TOKEN")

	args := os.Args
	if len(args) == 1 {
		args = append(args, "http://127.0.0.1:4242/ipfs-pages/webhook")
	}
	listen, err := url.Parse(args[1])
	if err != nil {
		log.Fatalf("argument must be URL, e.g. http://127.0.0.1:4242/webhook -- %s", err)
	}

	ps := NewPagesServer(webhookSecret, accessToken)
	http.Handle(listen.Path, ps)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	log.Printf("Listening for webhooks on: %s", listen)
	log.Fatal(http.ListenAndServe(listen.Host, nil))
}

func errorResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("%s %s -- %s", r.Method, r.RequestURI, err)
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("500 Internal Server Error\n\n" + err.Error()))
}
