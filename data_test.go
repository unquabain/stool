package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type DataTestCase struct {
	Filename string
	Format   Format
}

var DataTestCases = []DataTestCase{
	{
		Filename: `test_data/test.yaml`,
		Format:   FormatYAML,
	},
	{
		Filename: `test_data/test.json`,
		Format:   FormatJSON,
	},
	{
		Filename: `test_data/test.yaml`,
		Format:   FormatUnknown,
	},
	{
		Filename: `test_data/test.json`,
		Format:   FormatUnknown,
	},
}

func (tc DataTestCase) Test(t *testing.T) {
	t.Helper()
	data, err := Parse(tc.Filename, tc.Format)
	assert.NoError(t, err)
	for _, key := range []string{`animals`, `vegetables`, `minerals`} {
		assert.Contains(t, data, key)
	}
	for _, key := range []string{`invertebrates`, `vertebrates`} {
		assert.Contains(t, data[`animals`], key)
	}
	var cursor any = data
	for _, key := range []string{`vegetables`, `trees`, `evergreen`} {
		assert.Contains(t, cursor, key)
		cursor = cursor.(map[string]any)[key]
	}
	assert.Contains(t, cursor, `spruce`)
	cursor = data
	for _, key := range []string{`minerals`, `igneous`} {
		assert.Contains(t, cursor, key)
		cursor = cursor.(map[string]any)[key]
	}
	assert.Contains(t, cursor, `basalt`)
}

func TestData(t *testing.T) {
	for _, tc := range DataTestCases {
		tc.Test(t)
	}
}
