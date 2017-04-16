package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

const (
	buildsDir    = "builds"
	reposDir     = "repos"
	settingsFile = ".ipfspages.yml"
)

type Settings struct {
	Commands []string
	Output   string
	Dnslinks map[string]string
}

func main() {
	if len(os.Args) != 3 {
		log.Fatal("usage: pagebuild github.com/ORG/REPO COMMIT")
	}
	rparts := strings.Split(os.Args[1], "/")
	if len(rparts) != 3 {
		log.Fatal("usage: pagebuild github.com/ORG/REPO COMMIT")
	}

	owner := rparts[1]
	repo := rparts[2]
	commit := os.Args[2]

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	repoDir := filepath.Join(cwd, reposDir, owner, repo)
	repoURL := fmt.Sprintf("git@github.com:%s/%s.git", owner, repo)
	buildDir := filepath.Join(cwd, buildsDir, commit)
	treeDir := filepath.Join(cwd, buildsDir, commit, "tree")

	if err := cloneRepo(repoDir, repoURL); err != nil {
		log.Fatal(err)
	}

	if err := cloneBuild(buildDir, repoDir, repoURL, commit); err != nil {
		log.Fatal(err)
	}

	commitRelevant, err := exists(filepath.Join(treeDir, settingsFile))
	if err != nil {
		log.Fatal(err)
	}
	if !commitRelevant {
		log.Printf("ignoring github.com/%s/%s %s", owner, repo, commit)
		if err := os.RemoveAll(buildDir); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}

	settings, err := parseSettingsFile(filepath.Join(treeDir, settingsFile))
	if err != nil {
		log.Fatal(err)
	}

	h, err := build(treeDir, settings)
	if err != nil {
		log.Fatal(err)
	}

	// TODO:
	// - ipfs files cp $h /github.com/ipfs/website
	// - ipfs files publish
	log.Printf("https://ipfs.io/ipfs/%s", h)
}

func build(treeDir string, settings *Settings) (string, error) {
	for _, cmd := range settings.Commands {
		if err := command(treeDir, "bash", "-c", cmd); err != nil {
			return "", err
		}
	}

	args := []string{"add", "-Q", "-r", filepath.Join(treeDir, settings.Output)}
	log.Printf("> ipfs %s", strings.Join(args, " "))
	c := exec.Command("ipfs", args...)
	buf := new(bytes.Buffer)
	c.Stdout = buf
	c.Stderr = os.Stderr
	if err := c.Run(); err != nil {
		return "", nil
	}

	return buf.String(), nil
}

func parseSettingsFile(path string) (*Settings, error) {
	yml, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var settings Settings
	return &settings, yaml.Unmarshal(yml, &settings)
}

func cloneRepo(repoDir, repoURL string) error {
	if err := os.MkdirAll(filepath.Join(repoDir, ".."), 0755); err != nil {
		return err
	}

	repoexists, err := exists(repoDir)
	if err != nil {
		return err
	}
	if !repoexists {
		if err = command("./", "git", "clone", "--bare", repoURL, repoDir); err != nil {
			return err
		}
	}

	if err = command(repoDir, "git", "fetch", "--all"); err != nil {
		return err
	}

	return nil
}

func cloneBuild(buildDir, repoDir, repoURL, commit string) error {
	if err := os.RemoveAll(buildDir); err != nil {
		return err
	}
	if err := os.MkdirAll(buildDir, 0755); err != nil {
		return err
	}

	treeDir := filepath.Join(buildDir, "tree")

	if err := command(buildDir, "git", "clone", "--reference", repoDir, repoURL, treeDir); err != nil {
		return err
	}

	if err := command(treeDir, "git", "checkout", "-q", commit); err != nil {
		return err
	}

	return nil
}

func command(dir string, cmd string, args ...string) error {
	log.Printf("> %s %s (cwd=%s)", cmd, strings.Join(args, " "), dir)
	c := exec.Command(cmd, args...)
	c.Dir = dir
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Run()
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
