package main

import (
	"context"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	gh "github.com/google/go-github/github"
	oauth2 "golang.org/x/oauth2"
)

const pagesYml = ".ipfspages.yml"

type PagesServer struct {
	context       context.Context
	webhookSecret []byte
	github        *gh.Client
	ymlBranch     string
}

func NewPagesServer(ctx context.Context, secret []byte, github *gh.Client) *PagesServer {
	ps := &PagesServer{
		context:       ctx,
		webhookSecret: secret,
		github:        github,
		ymlBranch:     "master",
	}
	return ps
}

func (ps *PagesServer) RefreshTargets(ctx context.Context, org string) error {
	opt := &gh.RepositoryListByOrgOptions{ListOptions: gh.ListOptions{PerPage: 10}}
	for {
		repos, resp, err := ps.github.Repositories.ListByOrg(ctx, org, opt)
		if err != nil {
			return err
		}
		for _, repo := range repos {
			err = ps.refreshTarget(ctx, repo)
			if err != nil {
				return err
			}
		}
		if resp.NextPage == 0 {
			break
		}
		opt.ListOptions.Page = resp.NextPage
	}
	return nil
}

func (ps *PagesServer) refreshTarget(ctx context.Context, repo *gh.Repository) error {
	owner := strings.Split(repo.GetFullName(), "/")[0]
	name := repo.GetName()

	ref, resp, err := ps.github.Git.GetRef(ctx, owner, name, "heads/"+ps.ymlBranch)
	if resp.StatusCode == 404 {
		ps.removeTarget(repo, nil)
		return nil
	}
	if err != nil {
		return err
	}

	tree, _, err := ps.github.Git.GetTree(ctx, owner, name, ref.Object.GetSHA(), false)
	if err != nil {
		return err
	}

	var ymlentry *gh.TreeEntry
	for _, entry := range tree.Entries {
		if entry.GetPath() == pagesYml {
			ymlentry = &entry
			break
		}
	}
	if ymlentry == nil {
		ps.removeTarget(repo, nil)
		return nil
	}

	ymlblob, _, err := ps.github.Git.GetBlob(ctx, owner, name, ymlentry.GetSHA())
	if err != nil {
		return err
	}

	ps.addTarget(repo, ymlblob)

	return nil
}

func (ps *PagesServer) addTarget(repo *gh.Repository, _ *gh.Blob) error {
	log.Printf("addTarget: %s", repo.GetFullName())
	return nil
}

func (ps *PagesServer) removeTarget(repo *gh.Repository, _ *gh.Blob) error {
	log.Printf("removeTarget: %s", repo.GetFullName())
	return nil
}

func (ps *PagesServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_, event, err := ps.parseWebhook(r)
	if err != nil {
		ps.errorResponse(w, r, err)
		return
	}

	// - refreshTarget if on ymlBranch, and changed files includes .ipfspages.yml
	// - build if on one of branches specified in .ipfspages.yml

	switch event := event.(type) {
	case *gh.PushEvent:
		log.Printf("received push to %s", event.Repo.GetFullName())
		err = ps.receivePushEvent(event)
		if err != nil {
			ps.errorResponse(w, r, err)
		}
	default:
	}

	w.WriteHeader(http.StatusAccepted)
}

func (ps *PagesServer) receivePushEvent(event *gh.PushEvent) error {
	remove := false
	add := false
	for _, commit := range event.Commits {
		for _, addee := range append(commit.Added, commit.Modified...) {
			if addee == pagesYml {
				add = true
				remove = false
			}
		}
		for _, addee := range commit.Removed {
			if addee == pagesYml {
				remove = true
				add = false
			}
		}
	}
	if add || remove {
		fullname := strings.Split(event.Repo.GetFullName(), "/")
		repo, resp, err := ps.github.Repositories.Get(ps.context, fullname[0], fullname[1])
		if resp.StatusCode == 404 {
			return nil
		}
		if err != nil {
			return err
		}

		if add {
			return ps.addTarget(repo, nil)
		} else if remove {
			return ps.removeTarget(repo, nil)
		}
	}
	return nil
}

func (ps *PagesServer) errorResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("%s %s -- %s", r.Method, r.RequestURI, err)
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("500 Internal Server Error\n\n" + err.Error()))
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

	ctx := context.Background()
	oat := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: accessToken})
	oac := oauth2.NewClient(ctx, oat)
	github := gh.NewClient(oac)

	ps := NewPagesServer(ctx, webhookSecret, github)

	// for _, org := range []string{"ipfs", "ipld", "orbitdb", "libp2p", "multiformats"} {
	for _, org := range []string{"ipld"} {
		err = ps.RefreshTargets(ctx, org)
		if err != nil {
			log.Fatalf("RefreshTargets: %s", err)
		}
	}

	http.Handle(listen.Path, ps)
	log.Printf("Listening for webhooks on: %s", listen)
	log.Fatal(http.ListenAndServe(listen.Host, nil))
}
