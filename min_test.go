package tuna

import "testing"

func TestMin(t *testing.T) {
	ExtractorTestCases{
		{
			stream: NewStream(
				Row{"flux": "3.0"},
				Row{"flux": "4.0"},
				Row{"flux": "2.0"},
			),
			extractor: NewMin("flux"),
			output:    "flux_min\n2\n",
		},
	}.Run(t)
}