package main

import (
	"bufio"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"html/template"
	"io/ioutil"
	"os"
	"strings"
)

var output = "version/version.go"

func main() {

	if err := tagging(); err != nil {
		panic(err)
	}

	content, err := ioutil.ReadFile(".version")
	if err != nil {
		panic(err)
	}
	version := strings.Replace(string(content), "\n", "", -1)

	t, err := template.New("version").Parse(versionTemplate)

	if err != nil {
		panic(err)
	}

	data := struct {
		Version string
	}{version}
	f, err := os.Create(output)
	if err != nil {
		panic(err)
	}
	w := bufio.NewWriter(f)
	err = t.Execute(w, data)
	if err != nil {
		panic(err)
	}
	err = w.Flush()
	if err != nil {
		panic(err)
	}
}

func tagging() error {
	r, err := git.PlainOpen(".git")
	if err != nil {
		return err
	}
	tagrefs, err := r.Tags()
	if err != nil {
		return err
	}
	err = tagrefs.ForEach(func(t *plumbing.Reference) error {
		fmt.Println(t.Name())
		return nil
	})
	if err != nil {
		return err
	}
	return fmt.Errorf("not yet implemented")
}

func publicKey() (*ssh.PublicKeys, error) {
	return nil, fmt.Errorf("not yet implemented")
}

var versionTemplate = `/*
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
