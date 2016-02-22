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
	"net/url"
	"strconv"

	"bytes"

	//"io/ioutil"
 	//"regexp"

	mpg "github.com/mjibson/goread/_third_party/github.com/MiniProfiler/go/miniprofiler_gae"
	"github.com/mjibson/goread/_third_party/github.com/mjibson/goon"

	"github.com/pusher/pusher-http-go"

	"github.com/stripe/stripe-go"
    "github.com/stripe/stripe-go/client"

	"appengine"
	"appengine/datastore"
	"appengine/urlfetch"
)

type APResponseEnvelope struct {
	Timestamp     string `json:"timestamp"`
	Ack           string `json:"ack"`
	CorrelationId string `json:"correlationId"`
	Build         string `json:"build"`
}

type APRsp struct {
	ResponseEnvelope  APResponseEnvelope `json:"responseEnvelope"`
	Token     		  string 			 `json:"token"`
	TokenSecret		  string    		 `json:"tokenSecret"`
	Scope			[]string             `json:"scope"`                   
}

/*type StripeWebHook struct {
    Id 			string `json:"id"`
    Created  	int    `json:"created"`
    Type 		string `json:"type"`
}*/

type Event struct {
	ID       string     `json:"id"`
	Live     bool       `json:"livemode"`
	Created  int      `json:"created"`
	Data     *EventData `json:"data"`
	Webhooks uint64     `json:"pending_webhooks"`
	Type     string     `json:"type"`
	Req      string     `json:"request"`
	UserID   string     `json:"user_id"`
}

// EventData is the unmarshalled object as a map.
type EventData struct {
	Raw  *RawData        `json:"object"`
}

type RawData struct {
	ID       string       `json:"id"`
	Amount   int 		  `json:"amount"`
	Status 	 string       `json:"status"`
	Currency string 	  `json:"currency"`
	Source   *SourceData  `json:"source"`
}

type SourceData struct {
	Name    string       `json:"name"`
	Funding string       `json:"funding"`
}

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


func SubscribeStripe(c mpg.Context, w http.ResponseWriter, r *http.Request) {

	ar := strings.Split(r.FormValue("state"), "~")

	gn := goon.FromContext(c)

	sc := StripeSubscription{
		Active:		 "true",
		CreatedBy:   ar[0],
		Channel:     ar[1],
		DateCreated: ar[2],
		Code:        r.FormValue("code"),
		Scope:       r.FormValue("scope"),
	}
	gn.Put(&sc)

	cn := appengine.NewContext(r)
	urlfetchClient := urlfetch.Client(cn)

	client := pusher.Client{
		AppId:      "178872",
		Key:        "2aad67c195708eaa0e5f",
		Secret:     "048f50b1be4faa0aa64b",
		HttpClient: urlfetchClient,
	}

	client.Trigger("test_channel", "my_event", sc)
	return
}

func SubscriptionList(c mpg.Context, w http.ResponseWriter, r *http.Request) {
	gn := goon.FromContext(c)
	q := datastore.NewQuery(gn.Kind(&StripeSubscription{}))
	var sc []StripeSubscription
	_, err1 := gn.GetAll(q, &sc)
	if err1 != nil {
		return
	}
	b, err2 := json.Marshal(sc)
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

func PaymentSuccess(c mpg.Context, w http.ResponseWriter, r *http.Request) {

	cn := appengine.NewContext(r)
    urlfetchClient := urlfetch.Client(cn)

    client := pusher.Client{
        AppId:  "178872",
        Key:    "2aad67c195708eaa0e5f",
        Secret: "048f50b1be4faa0aa64b",
        HttpClient: urlfetchClient,
    }

    mapH := map[string]string{"status": "Successfully Paid !"}
    mapB, err := json.Marshal(mapH)

    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    client.Trigger("test_channel", "my_event", mapH)
	
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write(mapB)
}

func PaymentStripe(c mpg.Context, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	cn := appengine.NewContext(r)
    httpClient := urlfetch.Client(cn)

    sc := client.New("sk_test_eDxoEXHMVzXkHS2Vjvxndjz3", stripe.NewBackends(httpClient))

    params := stripe.ChargeParams{
	    Desc:     "Pretlist Subscription",
	    Amount:   20,
	    Currency: "usd",
	}
	params.SetSource(&stripe.CardParams{
	    Name:   r.FormValue("stripeEmail"),
	    Number: "378282246310005",
	    Month:  "06",
	    Year:   "17",
	})

	_, err := sc.Charges.New(&params)

	if err == nil {
		//fmt.Fprintf(w, "Successful test payment!")
	} else {
		http.Error(w, err.Error(), http.StatusInternalServerError)
        return
	}

}

func WebhookStripe(c mpg.Context, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	cn := appengine.NewContext(r)
    urlfetchClient := urlfetch.Client(cn)

    client := pusher.Client{
        AppId:  "178872",
        Key:    "2aad67c195708eaa0e5f",
        Secret: "048f50b1be4faa0aa64b",
        HttpClient: urlfetchClient,
    }

	gn := goon.FromContext(c)

	py := Payment{Id: r.FormValue("paymentId")}
	if err := gn.Get(&py);

	err != nil {
		decoder := json.NewDecoder(r.Body)
	    var swh Event   
	    err := decoder.Decode(&swh)
	    
	    if err != nil {
	        http.Error(w, err.Error(), http.StatusInternalServerError)
	        return
	    }
		py = Payment{
				Id: swh.Data.Raw.ID,
			    Active: "true",
			    PaidBy: swh.Data.Raw.Source.Name,
			    Status: swh.Data.Raw.Status,
			    DateCreated: strconv.Itoa(swh.Created),
			    Type: swh.Type,
			    Funding: swh.Data.Raw.Source.Funding,
			    Amount: strconv.Itoa(swh.Data.Raw.Amount),
			    Currency:swh.Data.Raw.Currency,
			    Source: "stripe",
			}
		gn.Put(&py)
		//mapH := map[string]string{"source": "stripe", "status": swh.Data.Raw.Status}
		//mapB, err := json.Marshal(py)

	    /*if err != nil {
	        http.Error(w, err.Error(), http.StatusInternalServerError)
	        return
	    }*/

    	client.Trigger("test_channel", "my_event", py)
		return
	}

	py.Active = r.FormValue("active")
    gn.Put(&py)
	client.Trigger("test_channel", "my_event", py)
	//w.Write(mapB)
}

func PaymentList(c mpg.Context, w http.ResponseWriter, r *http.Request) {
	gn := goon.FromContext(c)
	q := datastore.NewQuery(gn.Kind(&Payment{}))
	var ho []Payment
	
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


func GetAccessToken(c mpg.Context, w http.ResponseWriter, r *http.Request) {
	//c := appengine.NewContext(r)

	ap := map[string]string{"token": r.FormValue("request_token"), "verifier": r.FormValue("verification_code"), "requestEnvelope": "{'errorLanguage':'en_US'}"}

	b, err := json.Marshal(ap)

	if err != nil {
		return
	}
	buf := bytes.NewReader(b)

	req, err := http.NewRequest("POST", "https://svcs.sandbox.paypal.com/Permissions/GetAccessToken", buf)
	if err != nil {
		return
	}


	h := &req.Header
	h.Set("X-PAYPAL-SECURITY-USERID", "payme_api1.pretlist.com")
	h.Set("X-PAYPAL-SECURITY-PASSWORD", "DXRLZABPS3VEX44W")
	h.Set("X-PAYPAL-SECURITY-SIGNATURE", "AV56f8z5u6-nc0hOEpPsCNZgh-WeAgzupj2SjH4bg5xOc2SXgmLp3XRK")
	h.Set("X-PAYPAL-REQUEST-DATA-FORMAT", "JSON")
	h.Set("X-PAYPAL-RESPONSE-DATA-FORMAT", "JSON")
	h.Set("X-PAYPAL-APPLICATION-ID", "APP-80W284485P519543T")

	cl := &http.Client{}
	rsp, err := cl.Do(req)
	if err != nil {
		return
	}
	defer rsp.Body.Close()

	b2 := make([]byte, 1024)
	n, err := rsp.Body.Read(b2)
	aprsp := &APRsp{}
	err = json.Unmarshal(b2[0:n], aprsp)
	if err != nil {
		return
	}


	/*cn := appengine.NewContext(r)
    urlfetchClient := urlfetch.Client(cn)

    client := pusher.Client{
        AppId:  "178872",
        Key:    "2aad67c195708eaa0e5f",
        Secret: "048f50b1be4faa0aa64b",
        HttpClient: urlfetchClient,
    }

    mapH := map[string]string{"request_token": r.FormValue("request_token"), "verification_code": r.FormValue("verification_code")}
    mapB, err := json.Marshal(mapH)

    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    client.Trigger("test_channel", "my_event", mapH)*/
	
	w.Write(b)
}

func IPN(c mpg.Context, w http.ResponseWriter, r *http.Request) {
 	err := r.ParseForm() // need this to get PayPal's HTTP POST of IPN data

 	if err != nil {
 		//fmt.Println(err)
 		return
 	}

 	if r.Method == "POST" {

 		var postStr string = ""

 		for k, v := range r.Form {
 			//fmt.Println("key :", k)
 			//fmt.Println("value :", strings.Join(v, ""))

           // NOTE : Store the IPN data k,v into a slice. It will be useful for database entry later.

 			postStr = postStr + k + "=" + url.QueryEscape(strings.Join(v, "")) + " "
 		}

 		cn := appengine.NewContext(r)
	    urlfetchClient := urlfetch.Client(cn)

	    client := pusher.Client{
	        AppId:  "178872",
	        Key:    "2aad67c195708eaa0e5f",
	        Secret: "048f50b1be4faa0aa64b",
	        HttpClient: urlfetchClient,
	    }

	    mapH:= map[string]string{"status": postStr}
 		client.Trigger("test_channel", "my_event", mapH)
	
 	}
}

func WebhookPaypal(c mpg.Context, w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")

		cn := appengine.NewContext(r)
	    urlfetchClient := urlfetch.Client(cn)

	    client := pusher.Client{
	        AppId:  "178872",
	        Key:    "2aad67c195708eaa0e5f",
	        Secret: "048f50b1be4faa0aa64b",
	        HttpClient: urlfetchClient,
	    }

		gn := goon.FromContext(c)

		py := Payment{Id: r.FormValue("paymentId")}
		if err := gn.Get(&py);

		err != nil {
		    
			py = Payment{
					Id: r.FormValue("txn_id"),
				    Active: "true",
				    PaidBy: r.FormValue("payer_email"),
				    Status: r.FormValue("payment_status"),
				    DateCreated: r.FormValue("payment_date"),
				    Type: r.FormValue("transaction_subject"),
				    Funding: r.FormValue("transaction_subject"),
				    Amount: r.FormValue("mc_gross"),
				    Currency:r.FormValue("mc_currency"),
				    Source: "paypal",
				}
			gn.Put(&py)
			//mapH := map[string]string{"source": "stripe", "status": swh.Data.Raw.Status}
			//mapB, err := json.Marshal(py)

		   /* if err != nil {
		        http.Error(w, err.Error(), http.StatusInternalServerError)
		        return
		    }*/

	    	client.Trigger("test_channel", "my_event", py)
			return
		}

		py.Active = r.FormValue("active")
	    gn.Put(&py)
		client.Trigger("test_channel", "my_event", py)
		//w.Write(mapB)
}



func CheckoutStripe(c mpg.Context, w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "checkout-stripe.html", nil)
}