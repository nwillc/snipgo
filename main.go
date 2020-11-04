/*
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

package main

import (
	"fmt"
	"github.com/goava/di"
	"github.com/nwillc/snipgo/model"
	"github.com/nwillc/snipgo/services"
	"github.com/nwillc/snipgo/ui"
	"log"
)

//go:generate go run gorelease

func main() {
	c, err := di.New(
		di.Provide(newContext),
		di.Provide(newUI),
	)
	if err != nil {
		log.Fatal(err)
	}
	var ctx = services.NewDefaultContext()
	preferences, err := model.ReadPreferences(ctx.JSON, ctx.OS, ctx.IOUTIL, "")
	if err != nil {
		panic("Could not get preferences")
	}
	snippets, err := model.ReadSnippets(ctx.JSON, ctx.OS, ctx.IOUTIL, preferences.DefaultFile)
	if err != nil {
		panic(fmt.Sprintf("failed %v", err))
	}

	categories := snippets.ByCategory()

	var ui *ui.UI
	if err = c.Resolve(&ui); err != nil {
		log.Fatal("Unable to create ui ", err)
	}
	ui.Categories(&categories)
	ui.Run()
}

func newContext() *services.Context {
	return services.NewDefaultContext()
}

func newUI(ctx *services.Context) *ui.UI {
	return ui.NewUI(ctx)
}
