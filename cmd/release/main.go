package main

import (
	"bufio"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"strings"
	"time"
)

var output = "version/version.go"
var dotVersionFile = ".version"
var gitUser = "git"

func main() {
	repo := GetRepository("")
	w, err := repo.Worktree()
	CheckIfError(err)
	status, err := w.Status()
	CheckIfError(err)

	/*
	 * Check that we are ready for release.
	 */
	if len(status) != 1 {
		panic(fmt.Errorf("incorrrect file commit status, %d files", len(status)))
	}

	vs := status.File(dotVersionFile)
	if vs.Staging == '?' && vs.Worktree == '?' {
		panic(fmt.Errorf("%s should be only uncommitted file", dotVersionFile))
	}
	/*
	 * Get new version.
	 */
	content, err := ioutil.ReadFile(dotVersionFile)
	CheckIfError(err)
	version := strings.Replace(string(content), "\n", "", -1)

	/*
	 * Create the new version.go file.
	 */
	CreateVersionGo(output, version)

	/*
	 * Add the .version and version.go files.
	 */
	_, err = w.Add(output)
	CheckIfError(err)

	_, err = w.Add(dotVersionFile)
	CheckIfError(err)

	/*
	 * Commit the files.
	 */
	_, err = w.Commit("Generated for "+version, &git.CommitOptions{
		Author: NewSignature(),
	})
	CheckIfError(err)

	/*
	 * Set the new tag.
	 */
	ok, err := SetTag(repo, version)
	CheckIfError(err)

	if !ok {
		panic(fmt.Errorf("unable to set tag %s\n", version))
	}

	/*
	 * Push the tag
	 */
	err = repo.Push(&git.PushOptions{
		RemoteName: "origin",
		RefSpecs:   []config.RefSpec{config.RefSpec("refs/tags/*:refs/tags/*")},
	})
	if err != nil {
		sshKey, _ := PublicKeys()
		err = repo.Push(&git.PushOptions{
			RemoteName: "origin",
			RefSpecs:   []config.RefSpec{config.RefSpec("refs/tags/*:refs/tags/*")},
			Auth:       sshKey,
		})
		if err != nil {
			log.Printf("Push failed, please: git push origin %s; git push", version)
		}
	} else {
		/*
		 * Push the entire repo
		 */
		err = repo.Push(&git.PushOptions{})
		CheckIfError(err)
	}
}

func PublicKeys() (*ssh.PublicKeys, error) {
	path, err := os.UserHomeDir()
	CheckIfError(err)
	path += "/.ssh/id_rsa"

	publicKey, err := ssh.NewPublicKeysFromFile(gitUser, path, "")
	if err != nil {
		return nil, err
	}
	return publicKey, nil
}

func NewSignature() *object.Signature {
	userInfo, err := user.Current()
	CheckIfError(err)
	sig := object.Signature{
		Name: userInfo.Name,
		When: time.Now(),
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
