package summary

import (
	"fmt"
	"log"
)

func ExampleSummary() {
	// Truncated k6 JSON summary
	k6JsonSummary := []byte(`
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
	}`)

	// Parse summary
	s, err := NewSummary(k6JsonSummary)
	if err != nil {
		log.Fatal(err)
	}

	// Create JUnit format
	j, err := s.JUnit()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(j))

	// Output:
	// <?xml version="1.0" encoding="UTF-8"?>
	// <testsuites tests="3" failures="2">
	//   <testsuite name="Checks" tests="1" failures="1">
	//     <testcase name="is status 200">
	//       <failure message="5 / 75 (6.67%) checks passed"></failure>
	//     </testcase>
	//   </testsuite>
	//   <testsuite name="Thresholds" tests="2" failures="1">
	//     <testcase name="p(90) &lt; 100">
	//       <failure message="threshold exceeded"></failure>
	//     </testcase>
	//     <testcase name="p(95) &lt; 120"></testcase>
	//   </testsuite>
	// </testsuites>
}
