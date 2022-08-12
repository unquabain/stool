package main

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TemplateTestCase struct {
	Template string
	Expected string
}

func (ttc TemplateTestCase) Test(t *testing.T) {
	t.Helper()

	tmplt, err := GetTemplate(ttc.Template, ``)
	assert.NoError(t, err)

	buff := new(bytes.Buffer)
	assert.NoError(t, tmplt.Execute(buff, TemplateTestData))
	assert.Equal(t, ttc.Expected, string(buff.Bytes()))
}

var TemplateTestData = []any{map[string]any{
	`boring`: []any{
		`walls`,
		map[string]any{
			`textbooks`: map[string]any{
				`subject`: []any{
					`money`,
					`personal growth`,
				},
			},
			`movies`: map[string]any{
				`subject`: []any{
					`emotions`,
					`training`,
				},
			},
		},
	},
	`interesting`: []any{
		`skies`,
		map[string]any{
			`textbooks`: map[string]any{
				`subject`: []any{
					`math`,
					`science`,
					`language`,
				},
			},
			`movies`: map[string]any{
				`subject`: []any{
					`space wars`,
					`robots`,
					`murder`,
					`dinosaurs`,
				},
			},
		},
	},
}}

var TemplateTestCases = []TemplateTestCase{
	{
		Template: `{{ index (index .interesting 1).movies.subject 2 }}`,
		Expected: `murder`,
	},
	{
		Template: `{{ index (index .interesting 1).movies.subject | yaml }}`,
		Expected: "- space wars\n- robots\n- murder\n- dinosaurs\n",
	},
	{
		Template: `{{ index (index .interesting 1).movies.subject | json }}`,
		Expected: `["space wars","robots","murder","dinosaurs"]`,
	},
	{
		Template: `{{ index (index .interesting 1).movies.subject | jspretty }}`,
		Expected: `[
    "space wars",
    "robots",
    "murder",
    "dinosaurs"
]`,
	},
}

func TestTemplate(t *testing.T) {
	for _, tc := range TemplateTestCases {
		tc.Test(t)
	}
}
