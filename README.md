# Github Repository Security Alerts

Prints a list of summarised Dependabot security alerts that are open, grouped by repository name.

## Dependencies

* Go 1.20+
* A Github token with the `repo` scope or `security_events` scope. 
    - See: https://docs.github.com/en/rest/dependabot/alerts#list-dependabot-alerts-for-an-organization

## Installation

    git clone github-repo-security-alerts
    cd github-repo-security-alerts
    go build .

## Usage

    GITHUB_TOKEN=your-github-token ./github-repo-security-alerts
    
## Licence

Copyright © 2023 eLife Sciences

Distributed under the GNU Affero General Public Licence, version 3.
