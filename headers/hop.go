package headers

import (
	"fmt"
	"net/mail"
	"strings"
	"time"
)

type Hop struct {
	Num      int
	Delay    int
	From     string
	By       string
	With     string
	Time     string
	DateTime time.Time
}

func (h Hop) String() string {
	return fmt.Sprintf("from %s by %s with %s; %s", h.From, h.By, h.With, h.DateTime.UTC())
}

// Hops takes the Received headers and returns the hops in order
func Hops(received []string) []Hop {
	hops := make([]Hop, len(received))

	for i, rec := range received {
		dateStr := strings.Split(rec, ";")
		t, err := mail.ParseDate(strings.TrimSpace(dateStr[1]))

		if err != nil {
			panic(err) //TODO
		}

		var h Hop
		h.DateTime = t.UTC()
		h.Time = h.DateTime.Format("02 Jan 2006 15:04:05 MST")

		//re := regexp.MustCompile("from[[:space:]]+ ([[:alpha:]]*)")m
		//s := re.Split(parts[0], 2)
		//fmt.Printf("%+v\n", s)
		//fmt.Println(len(s))

		//parts := strings.Fields(dateStr[0])
		//fmt.Printf("from: %s\n", parts[1])

		rec := dateStr[0]
		//fmt.Println(rec)

		// This seems all a bit fragile

		//fromIdx := strings.Index(rec, "from")
		byIdx := strings.Index(rec, "by")
		withIdx := strings.Index(rec, "with")

		h.From = strings.TrimSpace(strings.Replace(rec[:byIdx], "from ", "", -1))
		h.By = strings.TrimSpace(strings.Replace(rec[byIdx:withIdx], "by ", "", -1))
		h.With = strings.TrimSpace(strings.Replace(rec[withIdx:], "with ", "", -1))

		hops[i] = h

	}

	//Reverse hops
	for i := len(hops)/2 - 1; i >= 0; i-- {
		opp := len(hops) - 1 - i
		hops[i], hops[opp] = hops[opp], hops[i]
	}

	for i, h := range hops {
		h.Num = i + 1 // Set the Num field for humans who count from 1 rather than zero
		if i > 0 {
			h.Delay = int(h.DateTime.Sub(hops[i-1].DateTime).Seconds())
		}
		hops[i] = h
	}
	return hops

}
