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
	"io/ioutil"
	"math"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/google/go-github/v52/github"
	"golang.org/x/oauth2"
)

// type aliases, just to help readability
type Project = string
type Maintainer = string

// subset of interesting fields from a github.DependabotAlert struct
type Alert struct {
	AgeDays int
	Summary string
	URL     string
	CVE_ID  string
	GHSA_ID string
}

func panicOnErr(err error, action string) {
	if err != nil {
		panic(fmt.Sprintf("failed with '%s' while '%s'", err.Error(), action))
	}
}

func as_json(thing interface{}) string {
	json_blob_bytes, err := json.Marshal(thing)
	panicOnErr(err, "marshalling JSON data into a byte array")
	var out bytes.Buffer
	json.Indent(&out, json_blob_bytes, "", "  ")
	return out.String()
}

// ---

func github_token() string {
	token, present := os.LookupEnv("GITHUB_TOKEN")
	if !present {
		panic("envvar GITHUB_TOKEN not set.")
	}
	return token
}

// "https://github.com/elifesciences/journal-cms/security/dependabot/19" => "journal-cms"
func extract_project_from_url(github_url string) Project {
	u, err := url.Parse(github_url)
	panicOnErr(err, "parsing a URL")
	p := u.Path
	bits := strings.Split(p, "/")
	return bits[2]
}

func parse_maintainer_alias_map(args []string) map[Project][]Maintainer {
	maintainer_alias_map := map[Project][]Maintainer{}
	if len(args) > 0 {
		path := args[0]
		txt, _ := ioutil.ReadFile(path)
		json.Unmarshal(txt, &maintainer_alias_map)
	}
	return maintainer_alias_map
}

func fetch_alert_list(org_name, token string) map[Project][]Alert {
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
	dependabot_alert_list, _, err := client.Dependabot.ListOrgAlerts(ctx, org_name, opts)
	panicOnErr(err, "listing org security alerts")

	now := time.Now()
	idx := map[Project][]Alert{}
	for _, dependabot_alert := range dependabot_alert_list {
		pname := extract_project_from_url(dependabot_alert.GetHTMLURL())

		project_alert_list, present := idx[pname]
		if !present {
			project_alert_list = []Alert{}
		}

		age := now.Sub(dependabot_alert.GetCreatedAt().Time)
		age_days := int(math.Ceil(age.Hours() / 24))
		alert := Alert{
			CVE_ID:  dependabot_alert.SecurityAdvisory.GetCVEID(),
			GHSA_ID: dependabot_alert.SecurityAdvisory.GetGHSAID(),
			Summary: dependabot_alert.SecurityAdvisory.GetSummary(),
			URL:     dependabot_alert.GetHTMLURL(),
			AgeDays: age_days,
		}
		project_alert_list = append(project_alert_list, alert)
		idx[pname] = project_alert_list
	}

	return idx
}

// returns true if `str` is not a slack channel (#foo) and contains an '@'
func is_email_address(str string) bool {
	return str != "" && str[0] != '#' && strings.Contains(str, "@")
}

func main() {
	args := os.Args[1:]

	token := github_token()
	org_name := "elifesciences"

	maintainer_alias_map := parse_maintainer_alias_map(args)
	project_alert_map := fetch_alert_list(org_name, token)

	if len(project_alert_map) > 0 && len(maintainer_alias_map) > 0 {
		// we have project alerts and we have project maintainers.
		// group the projects by maintainers.
		// maintainer=>project=>[]Alert
		maintainer_project_map := map[Maintainer]map[Project][]Alert{}
		for project, alert_list := range project_alert_map {
			project_maintainer_list, present := maintainer_alias_map[project]
			if !present {
				// project has no maintainers!
				// it's possible the repository is new but using vulnerable deps.
				// projects with no maintainers are handled in `maintainers-txt` project.
				fmt.Fprintf(os.Stderr, "skipping project '%s' with %d alert(s): no maintainers found\n", project, len(alert_list))
				continue
			}
			for _, maintainer := range project_maintainer_list {
				if !is_email_address(maintainer) {
					fmt.Fprintf(os.Stderr, "skipping maintainer, doesn't look like an email address: %s\n", maintainer)
					continue
				}
				project_map, present := maintainer_project_map[maintainer]
				if !present {
					project_map = map[Project][]Alert{}
				}
				project_map[project] = alert_list
				maintainer_project_map[maintainer] = project_map
			}
		}
		fmt.Println(as_json(maintainer_project_map))

	} else if len(project_alert_map) > 0 {
		// we have project alerts but no list of project maintainers.
		// output everything as-is
		fmt.Println(as_json(project_alert_map))
	}
}
