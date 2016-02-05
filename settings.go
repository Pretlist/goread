package goapp

import (
	"appengine"
	"time"
)

var (
	ENABLE_PUBSUBHUBBUB bool = !appengine.IsDevAppServer()
	/*STRIPE_PLANS             = []Plan{}*/
)

const (
	GOOGLE_ANALYTICS_ID   = ""
	GOOGLE_ANALYTICS_HOST = ""
	PUBSUBHUBBUB_HOST     = "http://localhost:8080" // e.g., "www.goread.io"
	STRIPE_KEY            = ""
	STRIPE_SECRET         = ""
	STRIPE_PLAN           = "123"
)

const (
	UpdateMin         = time.Minute * 20
	UpdateMax         = time.Hour * 12
	UpdateDefault     = time.Hour * 3
	UpdateFraction    = 0.5
	UpdateJitter      = time.Minute * 3
	UpdateLongFactor  = 20
	NewIntervalWeight = 0.2
)
