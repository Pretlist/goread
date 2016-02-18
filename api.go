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
	"strings"

	mpg "github.com/mjibson/goread/_third_party/github.com/MiniProfiler/go/miniprofiler_gae"
	"github.com/mjibson/goread/_third_party/github.com/mjibson/goon"

	"github.com/pusher/pusher-http-go"

	"appengine"
	"appengine/datastore"
	"appengine/urlfetch"
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

func AddHangouts(c mpg.Context, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	gn := goon.FromContext(c)

	ho := Hangout{Id: r.FormValue("hangoutId")}
	if err := gn.Get(&ho); err != nil {
		ho = Hangout{
			Id:            r.FormValue("hangoutId"),
			Active:        r.FormValue("active"),
			ParentChannel: r.FormValue("channel"),
			CreatedBy:     r.FormValue("createdBy"),
			DateCreated:   r.FormValue("dateCreated"),
			Hook:          r.FormValue("hook"),
			Token:         r.FormValue("token"),
		}
		gn.Put(&ho)
		return
	}

	cn := appengine.NewContext(r)
	urlfetchClient := urlfetch.Client(cn)

	client := pusher.Client{
		AppId:      "178872",
		Key:        "2aad67c195708eaa0e5f",
		Secret:     "048f50b1be4faa0aa64b",
		HttpClient: urlfetchClient,
	}

	if r.FormValue("channelName") != "" {
		if strings.Contains(ho.ParentChannel, r.FormValue("channelName")) {
			//fmt.Printf("Found subStr in str \n")
		} else {
			channels := ho.ParentChannel + ";" + r.FormValue("channelName")
			ho.ParentChannel = channels
		}
	}

	ho.Hook = r.FormValue("url")
	ho.Active = r.FormValue("active")
	gn.Put(&ho)

	mapH := map[string]string{"url": ho.Hook, "displayName": ho.CreatedBy, "active": ho.Active, "channel": ho.ParentChannel}
	mapB, err := json.Marshal(mapH)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	client.Trigger("test_channel", "my_event", mapH)

	//Set the cross origin resource sharing header to allow AJAX

	w.Write(mapB)
}

func HangoutList(c mpg.Context, w http.ResponseWriter, r *http.Request) {
	gn := goon.FromContext(c)
	q := datastore.NewQuery(gn.Kind(&Hangout{}))
	var ho []Hangout

	_, err1 := gn.GetAll(q, &ho)
	if err1 != nil {
		return
	}

	b, err2 := json.Marshal(ho)
	if err2 != nil {
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write(b)
}
