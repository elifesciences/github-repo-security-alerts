/*
   Copyright (C) 2023 eLife Sciences

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU Affero General Public License as
   published by the Free Software Foundation, either version 3 of the
   License, or (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU Affero General Public License for more details.

   You should have received a copy of the GNU Affero General Public License
   along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/google/go-github/v52/github"
	"golang.org/x/oauth2"
)

func panicOnErr(err error, action string) {
	if err != nil {
		panic(fmt.Sprintf("failed with '%s' while '%s'", err.Error(), action))
	}
}

func github_token() string {
	token, present := os.LookupEnv("GITHUB_TOKEN")
	if !present {
		panic("envvar GITHUB_TOKEN not set.")
	}
	return token
}

// "https://github.com/elifesciences/journal-cms/security/dependabot/19" => "journal-cms"
func extract_project_from_url(github_url string) string {
	u, err := url.Parse(github_url)
	panicOnErr(err, "parsing a URL")
	p := u.Path
	bits := strings.Split(p, "/")
	return bits[2]
}

type Alert struct {
	AgeDays int
	Summary string
	URL     string
	CVE_ID  string
	GHSA_ID string
}

func as_json(thing interface{}) string {
	json_blob_bytes, err := json.Marshal(thing)
	panicOnErr(err, "marshalling JSON data into a byte array")
	var out bytes.Buffer
	json.Indent(&out, json_blob_bytes, "", "  ")
	return out.String()
}

func main() {
	token := github_token()
	org_name := "elifesciences"

	// ---

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)

	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	state := "open"
	opts := &github.ListAlertsOptions{
		State: &state,
	}

	alert_list, _, err := client.Dependabot.ListOrgAlerts(ctx, org_name, opts)
	panicOnErr(err, "listing org security alerts")

	now := time.Now()
	idx := map[string][]Alert{}
	for _, alert := range alert_list {
		pname := extract_project_from_url(alert.GetHTMLURL())

		project_alert_list, present := idx[pname]
		if !present {
			project_alert_list = []Alert{}
		}

		age := now.Sub(alert.GetCreatedAt().Time)
		age_days := int(math.Ceil(age.Hours() / 24))

		a := Alert{
			CVE_ID:  alert.SecurityAdvisory.GetCVEID(),
			GHSA_ID: alert.SecurityAdvisory.GetGHSAID(),
			Summary: alert.SecurityAdvisory.GetSummary(),
			URL:     alert.GetHTMLURL(),
			AgeDays: age_days,
		}
		project_alert_list = append(project_alert_list, a)
		idx[pname] = project_alert_list
	}

	fmt.Println(as_json(idx))
}
