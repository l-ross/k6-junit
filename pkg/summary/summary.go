package summary

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"sort"
)

// A Summary is a struct representing the output of a parsed K6 summary file
type Summary struct {
	Checks           []Check
	checkFailedCount int
	Thresholds       []Threshold
	thrFailedCount   int
}

// A K6 Threshold
type Threshold struct {
	// Name of the threshold
	Name string
	// Result of the threshold
	Passed bool
}

// A K6 Check
type Check struct {
	// Name of the check
	Name string `json:"name"`
	// Number of instances of this check that passed
	Passes int `json:"passes"`
	// Number of instances of this check that failed
	Fails int `json:"fails"`
}

// NewSummary creates a Summary from the provided byte slice.
func NewSummary(b []byte) (*Summary, error) {
	return NewSummaryFromReader(bytes.NewBuffer(b))
}

// NewSummaryFromReader creates a Summary from the provided io.Reader.
func NewSummaryFromReader(r io.Reader) (*Summary, error) {
	s := &Summary{}

	// Unmarshal first level of keys
	var d map[string]json.RawMessage
	if err := json.NewDecoder(r).Decode(&d); err != nil {
		return nil, err
	}

	//
	// Get Checks
	//

	rootGroupRaw := d["root_group"]
	rg := rootGroup{}
	if err := json.Unmarshal(rootGroupRaw, &rg); err != nil {
		return nil, err
	}

	var checks map[string]Check
	if err := json.Unmarshal(rg.Checks, &checks); err != nil {
		return nil, err
	}

	for _, v := range checks {
		if v.Fails != 0 {
			s.checkFailedCount++
		}
		s.Checks = append(s.Checks, v)
	}

	sort.Slice(s.Checks, func(i, j int) bool {
		return s.Checks[i].Name < s.Checks[j].Name
	})

	//
	// Get Thresholds
	//

	metrics := d["metrics"]
	if err := json.Unmarshal(metrics, &d); err != nil {
		return nil, err
	}

	for _, v := range d {
		m := metric{}

		if err := json.Unmarshal(v, &m); err != nil {
			return nil, err
		}

		if m.Thresholds != nil {
			t := threshold{}
			if err := json.Unmarshal(m.Thresholds, &t); err != nil {
				return nil, err
			}

			for k, v := range t {
				if v {
					s.thrFailedCount++
				}

				s.Thresholds = append(s.Thresholds, Threshold{
					Name:   k,
					Passed: !v,
				})
			}
		}
	}

	sort.Slice(s.Thresholds, func(i, j int) bool {
		return s.Thresholds[i].Name < s.Thresholds[j].Name
	})

	return s, nil
}

type rootGroup struct {
	Checks json.RawMessage `json:"checks"`
}

type metric struct {
	Thresholds json.RawMessage `json:"thresholds"`
}

type threshold map[string]bool

// JUnit returns the JUnit XML output from a Summary.
func (s *Summary) JUnit() ([]byte, error) {
	j := &junit{}

	j.Tests = len(s.Checks) + len(s.Thresholds)
	j.Failures = s.checkFailedCount + s.thrFailedCount

	// Add Checks
	cts := &testSuite{
		Name:     "Checks",
		Tests:    len(s.Checks),
		Failures: s.checkFailedCount,
	}

	for _, c := range s.Checks {
		tc := &testCase{
			Name: c.Name,
		}

		if c.Fails > 0 {
			total := c.Fails + c.Passes
			per := float64(c.Passes) / float64(total) * 100
			tc.Failure = &failure{
				Message: fmt.Sprintf("%d / %d (%0.2f%%) checks passed", c.Passes, total, per),
			}
		}

		cts.TestCases = append(cts.TestCases, tc)
	}

	j.TestSuites = append(j.TestSuites, cts)

	// Add Thresholds
	tts := &testSuite{
		Name:     "Thresholds",
		Tests:    len(s.Thresholds),
		Failures: s.thrFailedCount,
	}

	for _, t := range s.Thresholds {
		tc := &testCase{
			Name: t.Name,
		}

		if !t.Passed {
			tc.Failure = &failure{
				Message: "threshold exceeded",
			}
		}

		tts.TestCases = append(tts.TestCases, tc)
	}

	j.TestSuites = append(j.TestSuites, tts)

	b, err := xml.MarshalIndent(j, "", "  ")
	if err != nil {
		return nil, err
	}

	b = append([]byte(xml.Header), b...)

	return b, nil
}

type junit struct {
	XMLName    xml.Name     `xml:"testsuites"`
	Tests      int          `xml:"tests,attr"`
	Failures   int          `xml:"failures,attr"`
	TestSuites []*testSuite `xml:"testsuite"`
}

type testSuite struct {
	Name      string      `xml:"name,attr"`
	Tests     int         `xml:"tests,attr"`
	Failures  int         `xml:"failures,attr"`
	TestCases []*testCase `xml:"testcase"`
}

type testCase struct {
	Name    string   `xml:"name,attr"`
	Failure *failure `xml:"failure,omitempty"`
}

type failure struct {
	Message string `xml:"message,attr"`
}
