package geojson

import (
	"fmt"
	"reflect"
)

func decodeBoundingBoxValue(bb reflect.Value) ([]float64, error) {
	if bb.IsZero() {
		return nil, nil
	}
	if bb.Kind() != reflect.Slice {
		return nil, fmt.Errorf("invalid bounding box %#v", bb.Interface())
	}
	result := make([]float64, bb.Len())
	for i := range result {
		c := indexValue(bb, i)
		if c.Kind() == reflect.Int {
			result[i] = float64(c.Int())
			continue
		}
		if c.Kind() == reflect.Float32 || c.Kind() == reflect.Float64 {
			result[i] = c.Float()
			continue
		}
		return nil, fmt.Errorf("invalid bounding box %#v", bb.Interface())
	}
	return result, nil
}
