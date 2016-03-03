/*
 * Copyright (c) 2013 Matt Jibson <matt.jibson@gmail.com>
 *
 * Permission to use, copy, modify, and distribute this software for any
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

package goapp

import (
	"encoding/json"
	"net/http"

	mpg "github.com/mjibson/goread/_third_party/github.com/MiniProfiler/go/miniprofiler_gae"
	"github.com/mjibson/goread/_third_party/github.com/mjibson/goon"

	"appengine/datastore"
)

func ChannelList(c mpg.Context, w http.ResponseWriter, r *http.Request) {
	gn := goon.FromContext(c)
	q := datastore.NewQuery(gn.Kind(&Channel{}))
	var ch []Channel

	_, err1 := gn.GetAll(q, &ch)
	if err1 != nil {
		return
	}

	b, err2 := json.Marshal(ch)
	if err2 != nil {
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write(b)
}

func UsersList(c mpg.Context, w http.ResponseWriter, r *http.Request) {
	gn := goon.FromContext(c)
	q := datastore.NewQuery(gn.Kind(&User{}))
	var us []User

	_, err1 := gn.GetAll(q, &us)
	if err1 != nil {
		return
	}

	b, err2 := json.Marshal(us)
	if err2 != nil {
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write(b)
}

func CreateMessage(c mpg.Context, w http.ResponseWriter, r *http.Request) {
	gn := goon.FromContext(c)

	msg := Message{
		Title:       r.FormValue("title"),
		CreatedBy:   r.FormValue("createdBy"),
		DateCreated: r.FormValue("dateCreated"),
		Channel:     r.FormValue("channel"),
		DocId:       r.FormValue("docId"),
		Content:     r.FormValue("content"),
	}
	gn.Put(&msg)
}

func GetMessages(c mpg.Context, w http.ResponseWriter, r *http.Request) {
	gn := goon.FromContext(c)
	q := datastore.NewQuery(gn.Kind(&Message{}))
	var msgs []Message
	_, err1 := gn.GetAll(q, &msgs)
	if err1 != nil {
		return
	}
	b, err2 := json.Marshal(msgs)
	if err2 != nil {
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write(b)
}
