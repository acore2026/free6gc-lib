package aper

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func spansByPath(spans []DecodeSpan) map[string]DecodeSpan {
	byPath := make(map[string]DecodeSpan, len(spans))
	for _, span := range spans {
		byPath[span.Path] = span
	}
	return byPath
}

func TestUnmarshalWithParamsAndTraceRecordsBitOffsets(t *testing.T) {
	var decoded BitStringStructTest1
	spans, err := UnmarshalWithParamsAndTrace([]byte{0xb4}, &decoded, "")

	assert.NoError(t, err)
	assert.Equal(t, BitStringStructTest1Data[0], decoded)
	assert.Equal(t, map[string]DecodeSpan{
		"$.BitStringStructTest1":            {Path: "$.BitStringStructTest1", StartBit: 0, BitLength: 6},
		"$.BitStringStructTest1.BitString1": {Path: "$.BitStringStructTest1.BitString1", StartBit: 0, BitLength: 3},
		"$.BitStringStructTest1.BitString2": {Path: "$.BitStringStructTest1.BitString2", StartBit: 3, BitLength: 3},
	}, spansByPath(spans))
}

func TestUnmarshalWithParamsAndTraceRecordsChoiceAndSequenceElements(t *testing.T) {
	test := choiceTestData[0]
	decoded := reflect.New(reflect.TypeOf(test.Out))
	spans, err := UnmarshalWithParamsAndTrace(test.in, decoded.Interface(), "")

	assert.NoError(t, err)
	assert.Equal(t, test.Out, decoded.Elem().Interface())
	byPath := spansByPath(spans)
	assert.Equal(t, DecodeSpan{Path: "$.choiceTest1.Choice.Present", StartBit: 0, BitLength: 2}, byPath["$.choiceTest1.Choice.Present"])
	assert.Contains(t, byPath, "$.choiceTest1.Choice.List1")
	assert.Contains(t, byPath, "$.choiceTest1.Choice.List1.[0]")
	assert.Contains(t, byPath, "$.choiceTest1.Choice.List1.[1]")
	assert.Contains(t, byPath, "$.choiceTest1.Choice.List1.[2]")
	assert.Less(t, byPath["$.choiceTest1.Choice.List1.[0]"].StartBit, byPath["$.choiceTest1.Choice.List1.[1]"].StartBit)
	assert.Less(t, byPath["$.choiceTest1.Choice.List1.[1]"].StartBit, byPath["$.choiceTest1.Choice.List1.[2]"].StartBit)
}

func TestUnmarshalWithParamsAndTracePreservesOpenTypeAbsoluteOffsets(t *testing.T) {
	test := openTypeTestData[0]
	decoded := reflect.New(reflect.TypeOf(test.Out))
	spans, err := UnmarshalWithParamsAndTrace(test.in, decoded.Interface(), "")

	assert.NoError(t, err)
	assert.Equal(t, test.Out, decoded.Elem().Interface())
	byPath := spansByPath(spans)
	assert.Equal(t, DecodeSpan{Path: "$.openTypeTest1.ID", StartBit: 0, BitLength: 8}, byPath["$.openTypeTest1.ID"])
	assert.Equal(t, DecodeSpan{Path: "$.openTypeTest1.Value", StartBit: 8, BitLength: 96}, byPath["$.openTypeTest1.Value"])
	assert.Equal(t, uint64(16), byPath["$.openTypeTest1.Value.List1"].StartBit)
	assert.Contains(t, byPath, "$.openTypeTest1.Value.List1.[0]")
	assert.GreaterOrEqual(t, byPath["$.openTypeTest1.Value.List1.[0]"].StartBit, uint64(16))
	assert.Less(t, byPath["$.openTypeTest1.Value.List1.[0]"].StartBit, uint64(len(test.in)*8))
}
