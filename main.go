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
	"github.com/nwillc/snipgo/model"
	"github.com/nwillc/snipgo/services"
	"github.com/nwillc/snipgo/ui"
)

//go:generate go run gorelease

func main() {
	var os = services.NewOs()
	var json = services.NewJson()
	var ioUtils = services.NewIoUtil()

	var categories = defaultCategories(json, os, ioUtils)
	ui.NewUI(json, os, ioUtils, categories).Run()
}

func defaultCategories(json services.Json, os services.Os, ioUtil services.IoUtil) *model.Categories {
	preferences, err := model.ReadPreferences(json, os, ioUtil, "")
	if err != nil {
		panic("Could not get preferences")
	}
	snippets, err := model.ReadSnippets(json, os, ioUtil, preferences.DefaultFile)
	if err != nil {
		panic(fmt.Sprintf("failed %v", err))
	}

	categories := snippets.ByCategory()
	return &categories
}
