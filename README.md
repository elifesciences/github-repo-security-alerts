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
    
## Example

```json
{
  "api-dummy": [
    {
      "AgeDays": 19,
      "Summary": "Improper header name validation in guzzlehttp/psr7",
      "URL": "https://github.com/elifesciences/api-dummy/security/dependabot/13",
      "CVE_ID": "CVE-2023-29197",
      "GHSA_ID": "GHSA-wxmh-65f7-jcvw"
    }
  ],
  "data-science-dags": [
    {
      "AgeDays": 6,
      "Summary": "Flask vulnerable to possible disclosure of permanent session cookie due to missing Vary: Cookie header",
      "URL": "https://github.com/elifesciences/data-science-dags/security/dependabot/28",
      "CVE_ID": "CVE-2023-30861",
      "GHSA_ID": "GHSA-m2qf-hxjv-5gpq"
    },
  ],
  "hypothesis-dummy": [
    {
      "AgeDays": 19,
      "Summary": "Improper header name validation in guzzlehttp/psr7",
      "URL": "https://github.com/elifesciences/hypothesis-dummy/security/dependabot/4",
      "CVE_ID": "CVE-2023-29197",
      "GHSA_ID": "GHSA-wxmh-65f7-jcvw"
    }
  ],
  "lens-s3": [
    {
      "AgeDays": 37,
      "Summary": "Potential XSS vulnerability in jQuery",
      "URL": "https://github.com/elifesciences/lens-s3/security/dependabot/1",
      "CVE_ID": "CVE-2020-11022",
      "GHSA_ID": "GHSA-gxr4-xjj5-5px2"
    }
  ]
}
```
    
## Licence

Copyright © 2023 eLife Sciences

Distributed under the GNU Affero General Public Licence, version 3.
