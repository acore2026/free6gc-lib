package aper

import "reflect"

// DecodeSpan identifies the portion of an APER payload consumed while
// decoding a field. Offsets and lengths are expressed in bits because APER
// fields are not necessarily byte-aligned.
type DecodeSpan struct {
	Path      string
	StartBit  uint64
	BitLength uint64
}

type decodeTrace struct {
	spans []DecodeSpan
}

func (trace *decodeTrace) record(path string, startBit, endBit uint64) {
	if trace == nil || path == "" || endBit <= startBit {
		return
	}
	trace.spans = append(trace.spans, DecodeSpan{
		Path:      path,
		StartBit:  startBit,
		BitLength: endBit - startBit,
	})
}

func decodeRootPath(value reflect.Value) string {
	for value.IsValid() && value.Kind() == reflect.Pointer {
		if value.IsNil() {
			value = reflect.New(value.Type().Elem()).Elem()
			break
		}
		value = value.Elem()
	}
	if !value.IsValid() {
		return "$"
	}
	if name := value.Type().Name(); name != "" {
		return "$." + name
	}
	return "$"
}

func decodeChildPath(parent, child string) string {
	if parent == "" {
		return ""
	}
	return parent + "." + child
}
