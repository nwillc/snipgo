package main

import (
	"bufio"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)

var output = "version/version.go"
var dotVersionFile = ".version"

func main() {
	repo := GetRepository("")
	w, err := repo.Worktree()
	CheckIfError(err)
	status, err := w.Status()
	CheckIfError(err)

	if len(status) != 1 {
		panic(fmt.Errorf("incorrrect file commit status, %d files", len(status)))
	}

	vs := status.File(dotVersionFile)
	if vs.Staging == '?' && vs.Worktree == '?' {
		panic(fmt.Errorf("%s should be only uncommitted file", dotVersionFile))
	}

	// Get the target version
	content, err := ioutil.ReadFile(dotVersionFile)
	CheckIfError(err)
	version := strings.Replace(string(content), "\n", "", -1)

	CreateVersionGo(output, version)

	log.Printf("git status:\n %v", status)

	_, err = w.Add(output)
	CheckIfError(err)

	_, err = w.Add(dotVersionFile)
	CheckIfError(err)

	_, err = w.Commit("Generated for "+version, &git.CommitOptions{
		Author: NewSignature(),
	})
	CheckIfError(err)

	ok, err := SetTag(repo, version)
	CheckIfError(err)

	if !ok {
		log.Printf("unable to set %s\n", version)
	}

	err = repo.Push(&git.PushOptions{
		RemoteName: "origin",
		RefSpecs: []config.RefSpec{config.RefSpec(version)},
	})
	CheckIfError(err)
}

func NewSignature() *object.Signature {
	sig := object.Signature{
		Name:  "nwillc",
		Email: "nwillc@gmail.com",
		When:  time.Now(),
	}
	return &sig
}
func CreateVersionGo(fileName string, version string) {
	// Create new version.go
	versionTemplate, err := template.New("version").Parse(versionTemplateStr)
	CheckIfError(err)
	data := struct {
		Version string
	}{version}

	f, err := os.Create(fileName)
	CheckIfError(err)
	w := bufio.NewWriter(f)
	err = versionTemplate.Execute(w, data)
	CheckIfError(err)
	err = w.Flush()
	CheckIfError(err)
}

func GetRepository(repo string) *git.Repository {
	if repo == "" {
		repo = "."
	}
	r, err := git.PlainOpenWithOptions(repo, &git.PlainOpenOptions{DetectDotGit: true})
	if err != nil {
		panic(err)
	}
	return r
}

func TagExists(r *git.Repository, tag string) bool {
	tagFoundErr := "tag was found"
	tags, err := r.Tags()
	if err != nil {
		log.Printf("get tags error: %s", err)
		return false
	}
	res := false
	err = tags.ForEach(func(t *plumbing.Reference) error {
		if strings.Contains(t.Name().String(), tag) {
			res = true
			return fmt.Errorf(tagFoundErr)
		}
		return nil
	})
	if err != nil && err.Error() != tagFoundErr {
		log.Printf("iterate tags error: %s", err)
		return false
	}
	return res
}

func SetTag(r *git.Repository, tag string) (bool, error) {
	if TagExists(r, tag) {
		log.Printf("tag %s already exists", tag)
		return false, nil
	}
	log.Printf("Set tag %s", tag)
	h, err := r.Head()
	if err != nil {
		log.Printf("get HEAD error: %s", err)
		return false, err
	}
	_, err = r.CreateTag(tag, h.Hash(), &git.CreateTagOptions{
		Tagger:  NewSignature(),
		Message: "Release " + tag,
	})

	if err != nil {
		log.Printf("create tag error: %s", err)
		return false, err
	}

	return true, nil
}

func CheckIfError(err error) {
	if err == nil {
		return
	}

	panic(err)
}

var versionTemplateStr = `/*
 * Copyright (c) 2020, nwillc@gmail.com
 *
 * Permission to use, copy, modify, and/or distribute this software for any
 * purpose with or without fee is hereby granted, provided that the above
 * copyright notice and this permission notice appear in all copies.
 *
 * THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
 * WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
 * MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR
 * ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
 * WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN
 * ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF
 * OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
 */

package version

// Version number for official releases updated with go generate.
var Version = "{{.Version}}"
`
