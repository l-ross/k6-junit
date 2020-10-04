# k6-junit

k6-junit provides a utility to convert a [k6](https://k6.io) JSON summary in to a JUnit result file.

## Install

Install via Go: `go get github.com/l-ross/k6-junit/cmd/k6-junit`

Download the latest release from GitHub [here](https://github.com/l-ross/k6-junit/releases/latest)

## Example

When running k6 ensure that the `--summary-export` flag is provided to write the k6 summary to a JSON file.

Example k6 JSON Summary (truncated for brevity):
```json
{
    "metrics": {
        "http_req_duration": {
            "avg": 112.3124,
            "max": 112.3124,
            "med": 112.3124,
            "min": 112.3124,
            "p(90)": 112.3124,
            "p(95)": 112.3124,
            "thresholds": {
                "p(90) < 100": true,
                "p(95) < 120": false
            }
        }
    },
    "root_group": {
        "name": "",
        "path": "",
        "id": "d41d8cd98f00b204e9800998ecf8427e",
        "groups": {},
        "checks": {
            "is status 200": {
                "name": "is status 200",
                "path": "::is status 200",
                "id": "548d37ca5f33793206f7832e7cea54fb",
                "passes": 5,
                "fails": 70
            }
        }
    }
}
```

Then to generate the JUnit output run `k6-junit --in summary.json`, example output:
```xml
<?xml version="1.0" encoding="UTF-8"?>
<testsuites tests="3" failures="2">
  <testsuite name="Checks" tests="1" failures="1">
    <testcase name="is status 200">
      <failure message="5 / 75 (6.67%) checks passed"></failure>
    </testcase>
  </testsuite>
  <testsuite name="Thresholds" tests="2" failures="1">
    <testcase name="p(90) &lt; 100">
      <failure message="threshold exceeded"></failure>
    </testcase>
    <testcase name="p(95) &lt; 120"></testcase>
  </testsuite>
</testsuites>
```
