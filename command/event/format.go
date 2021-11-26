package event

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func (e event) format(prefix ...string) string {

	var ln []string
	var before string
	var weekly string
	var lasting string
	deltas := e.deltas()

	if len(prefix) > 0 {
		before = strings.Join(prefix, "\n") + "\n\n"
	}

	if e.Weekly {
		if e.Finish != nil {
			weekly = ", "
		}
		weekly += "Weekly"
	}

	if deltas.Lasting != 0 {
		lasting = deltas.lasting(false)
	}

	ico := ":calendar_spiral:"
	until := e.until(deltas)
	h := fmt.Sprintf("%s **__%s__** (%s)", ico, e.Name, until)
	ln = append(ln, h)
	ln = append(ln, e.Summary)
	ln = append(ln, " ")
	ln = append(ln, "**Starts**")
	ln = append(ln, date(e.Start, e.Local))
	ln = append(ln, " ")
	if e.Finish != nil {
		ln = append(ln, "**Finishes**")
		ln = append(ln, date(*e.Finish, e.Local))
		ln = append(ln, " ")
	}
	ln = append(ln, "**Duration**")
	ln = append(ln, lasting+weekly)
	ln = append(ln, " ")
	ln = append(ln, "**Tagged**")
	ln = append(ln, strings.Join(append(e.Category, e.Setting...), ", "))

	for i := range ln {
		ln[i] = "> " + ln[i]
	}

	return before + strings.Join(ln, "\n")
}

func date(n int64, local bool) string {
	t := time.Unix(n, 0).In(time.UTC)
	wd := t.Weekday().String()
	mm := t.Month().String()
	date := ordinal(t.Day())
	tz := "UTC"
	if local {
		tz = "Local Time"
	}
	clock := t.Format("3:04 pm")
	return fmt.Sprintf("%s, %s %s %d, %s %s", wd, date, mm, t.Year(), clock, tz)
}

func ordinal(n int) string {
	if n >= 10 && n < 19 {
		return fmt.Sprint(n, "th")
	}
	switch n % 10 {
	case 1:
		return fmt.Sprint(n, "st")
	case 2:
		return fmt.Sprint(n, "nd")
	case 3:
		return fmt.Sprint(n, "rd")
	}
	return fmt.Sprint(n, "th")
}

func (e event) deltas() *eventDeltas {

	var lasting time.Duration
	if e.Finish != nil {
		lasting = time.Second * time.Duration(*e.Finish-e.Start)
	}

	/*
		This can be negative if we're within the period the
		event occurs so we make sure to clamp it to zero.
	*/
	until := time.Until(time.Unix(e.Start, 0))
	if until < 0 {
		until = 0
	}

	h := int64(until.Hours())
	d := h / 24
	w := d / 7
	m := int64(until.Minutes()) - h*60
	s := int64(until.Seconds()) - h*60*60 - m*60

	return &eventDeltas{
		Months:  d / 30,
		Weeks:   w,
		Days:    d,
		Hours:   h % 24,
		Minutes: m,
		Seconds: s,
		Lasting: lasting,
	}
}

type eventDeltas struct {
	Months  int64
	Weeks   int64
	Days    int64
	Hours   int64
	Minutes int64
	Seconds int64
	Lasting time.Duration
}

func (e event) until(ed *eventDeltas) string {
	now := time.Now()
	if e.Finish == nil {
		if s := time.Unix(e.Start, 0); now.After(s) {
			return "Has Occurred"
		}
	} else {
		if f := time.Unix(*e.Finish, 0); now.After(f) {
			return "Is over"
		}
		if s := time.Unix(e.Start, 0); now.After(s) {
			return "Has begun"
		}
	}
	u := ""
	switch {
	case ed.Months > 1:
		u = fmt.Sprintf("%d months", ed.Months)
	case ed.Months == 1:
		u = "1 month"
	case ed.Weeks > 1:
		u = fmt.Sprintf("%d weeks", ed.Weeks)
	case ed.Weeks == 1:
		u = "1 week"
	case ed.Days > 1:
		u = fmt.Sprintf("%d days", ed.Days)
	case ed.Days == 1:
		u = "1 day"
	case ed.Hours > 1:
		u = fmt.Sprintf("%d hours", ed.Hours)
	case ed.Hours == 1:
		u = "1 hour"
	case ed.Minutes > 1:
		u = fmt.Sprintf("%d mins", ed.Minutes)
	case ed.Minutes == 1:
		u = "1 min"
	case ed.Seconds > 1:
		u = fmt.Sprintf("%d secs", ed.Seconds)
	case ed.Seconds == 1:
		u = "1 sec"
	}
	return "In " + u
}

func (ed eventDeltas) lasting(compact bool) string {

	h := int(ed.Lasting.Hours())
	d := h / 24
	h %= 24
	m := int(ed.Lasting.Minutes()) % 60

	args := []struct {
		val     int
		unit    string
		compact string
	}{
		{val: d, unit: "day", compact: "day"},
		{val: h, unit: "hour", compact: "hr"},
		{val: m, unit: "minute", compact: "min"},
	}

	var ss []string
	for _, a := range args {
		if len(ss) == 2 {
			break
		}
		if a.val == 0 {
			continue
		}
		unit := a.unit
		if compact {
			unit = a.compact
		}
		plural := ""
		if a.val > 1 {
			plural = "s"
		}
		ss = append(ss, strconv.Itoa(a.val)+" "+unit+plural)
	}

	comma := ","
	if compact {
		comma = ""
	}
	return strings.Join(ss, comma+" ")
}
