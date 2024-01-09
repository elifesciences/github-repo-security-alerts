# Github Repository Security Alerts

Prints a list of summarised Dependabot security alerts that are open, grouped by repository name.

Accepts a map of `project-name => list-of-maintainers` as a JSON file to further group by maintainer. 
See [maintainers-txt](https://github.com/elifesciences/maintainers-txt)

## Dependencies

* Go 1.20+
* A Github token with the `repo` scope or `security_events` scope. 
    - See: https://docs.github.com/en/rest/dependabot/alerts#list-dependabot-alerts-for-an-organization

## Installation

    git clone https://github.com/elifesciences/github-repo-security-alerts
    cd github-repo-security-alerts
    go build .

## Usage

    GITHUB_TOKEN=your-github-token ./github-repo-security-alerts

or

    GITHUB_TOKEN=your-github-token ./github-repo-security-alerts project-maintainers.json

with `project-maintainers.json` looking like:

```json
{
  "project-foo": [
    "jane"
  ],
  "project-baz": [
    "jack",
    "john"
  ],
}
```

## Example

```json
{
  "project-foo": [
    {
      "AgeDays": 6,
      "Summary": "Flask vulnerable to possible disclosure of permanent session cookie due to missing Vary: Cookie header",
      "URL": "https://github.com/elifesciences/project-foo/security/dependabot/28",
      "CVE_ID": "CVE-2023-30861",
      "GHSA_ID": "GHSA-m2qf-hxjv-5gpq"
    },
  ],
  "project-bar": [
    {
      "AgeDays": 19,
      "Summary": "Improper header name validation in guzzlehttp/psr7",
      "URL": "https://github.com/elifesciences/project-bar/security/dependabot/4",
      "CVE_ID": "CVE-2023-29197",
      "GHSA_ID": "GHSA-wxmh-65f7-jcvw"
    }
  ],
  "project-baz": [
    {
      "AgeDays": 37,
      "Summary": "Potential XSS vulnerability in jQuery",
      "URL": "https://github.com/elifesciences/project-baz/security/dependabot/1",
      "CVE_ID": "CVE-2020-11022",
      "GHSA_ID": "GHSA-gxr4-xjj5-5px2"
    }
  ]
}
```

and with a `project-maintainers.json` file:

```json
{
  "jack": {
    "project-baz": [
      {
        "AgeDays": 59,
        "Summary": "Potential XSS vulnerability in jQuery",
        "URL": "https://github.com/elifesciences/project-baz/security/dependabot/1",
        "CVE_ID": "CVE-2020-11022",
        "GHSA_ID": "GHSA-gxr4-xjj5-5px2"
      }
    ]
  },
  "john": {
    "project-baz": [
      {
        "AgeDays": 59,
        "Summary": "Potential XSS vulnerability in jQuery",
        "URL": "https://github.com/elifesciences/project-baz/security/dependabot/1",
        "CVE_ID": "CVE-2020-11022",
        "GHSA_ID": "GHSA-gxr4-xjj5-5px2"
      }
    ]
  },
  "jane": {
    "project-foo": [
      {
        "AgeDays": 18,
        "Summary": "Apache Airflow vulnerable to Privilege Context Switching Error",
        "URL": "https://github.com/elifesciences/project-foo/security/dependabot/13",
        "CVE_ID": "CVE-2023-25754",
        "GHSA_ID": "GHSA-jchm-fm4q-c2fp"
      },
      {
        "AgeDays": 21,
        "Summary": "Apache Airflow vulnerable to stored Cross-site Scripting",
        "URL": "https://github.com/elifesciences/project-foo/security/dependabot/12",
        "CVE_ID": "CVE-2023-29247",
        "GHSA_ID": "GHSA-vcf6-3wv2-5vcr"
      }
    ]
}
```

## Licence

Copyright © 2024 eLife Sciences

Distributed under the GNU Affero General Public Licence, version 3.
