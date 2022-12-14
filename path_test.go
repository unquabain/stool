package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type PathTestCase struct {
	Path            string
	ExpectedError   string
	ExpectedResults []any
}

func (ptc PathTestCase) Test(t *testing.T) {
	t.Helper()
	data, err := Parse(`test_data/test.yaml`, FormatYAML)
	assert.NoError(t, err, ptc.Path)

	results, err := Evaluate(data, ptc.Path)
	if ptc.ExpectedError != `` {
		assert.ErrorContains(t, err, ptc.ExpectedError, ptc.Path)
		return
	}
	for _, expected := range ptc.ExpectedResults {
		assert.Contains(t, results, expected, ptc.Path)
	}
}

var PathTestCases = []PathTestCase{
	{
		Path: `minerals.sedimentary`,
		ExpectedResults: []any{
			[]any{`sandstone`, `shale`, `chalk`},
		},
	},
	{
		Path: `minerals[*]`,
		ExpectedResults: []any{
			[]any{`sandstone`, `shale`, `chalk`},
			[]any{`obsidian`, `granite`, `basalt`},
			[]any{`slate`, `schist`, `marble`},
		},
	},
	{
		Path: `minerals[*][. < "ci"]`,
		ExpectedResults: []any{
			`basalt`, `chalk`,
		},
	},
	{
		Path: `animals.vertebrates.mammals.length()`,
		ExpectedResults: []any{
			3,
		},
	},
	{
		Path: `animals.invertebrates[length() == 1]`,
		ExpectedResults: []any{
			[]any{`clam`},
		},
	},
	{
		Path: `animals.invertebrates[length() < 2]`,
		ExpectedResults: []any{
			[]any{`clam`},
		},
	},
	{
		Path: `animals.invertebrates[length() != 1]`,
		ExpectedResults: []any{
			[]any{`fly`, `ant`},
		},
	},
	{
		Path: `animals.invertebrates[length() >= 1]`,
		ExpectedResults: []any{
			[]any{`clam`},
			[]any{`fly`, `ant`},
		},
	},
	{
		Path: `animals.invertebrates[length() <= 2]`,
		ExpectedResults: []any{
			[]any{`clam`},
			[]any{`fly`, `ant`},
		},
	},
	{
		Path: `animals.invertebrates[length() == 1][0].json()`,
		ExpectedResults: []any{
			`"clam"`,
		},
	},
	{
		Path: `animals.invertebrates[length() == 1].yaml()`,
		ExpectedResults: []any{
			"- clam\n",
		},
	},
	{
		Path: `animals.invertebrates.yaml()`,
		ExpectedResults: []any{
			"insects:\n    - fly\n    - ant\nmollusks:\n    - clam\n",
		},
	},
	{
		Path: `animals.invertebrates[length() == 1].json()`,
		ExpectedResults: []any{
			`["clam"]`,
		},
	},
	{
		Path: `animals.invertebrates.json()`,
		ExpectedResults: []any{
			`{"insects":["fly","ant"],"mollusks":["clam"]}`,
		},
	},
	{
		Path: `animals.invertebrates.jsonpretty()`,
		ExpectedResults: []any{
			`{
    "insects": [
        "fly",
        "ant"
    ],
    "mollusks": [
        "clam"
    ]
}`,
		},
	},
	{
		Path: `animals.keys()`,
		ExpectedResults: []any{
			[]any{`vertebrates`, `invertebrates`},
		},
	},
	{
		Path: `meta.description.jeval().type`,
		ExpectedResults: []any{
			`yaml document`,
		},
	},
	{
		Path: `animals.vertebrates.yaml().yeval().mammals`,
		ExpectedResults: []any{
			[]any{
				`horse`,
				`shrew`,
				`cat`,
			},
		},
	},
	{
		Path: `animals.vertebrates.json().jeval().mammals`,
		ExpectedResults: []any{
			[]any{
				`horse`,
				`shrew`,
				`cat`,
			},
		},
	},
}

func TestPath(t *testing.T) {
	for _, tc := range PathTestCases {
		tc.Test(t)
	}
}
