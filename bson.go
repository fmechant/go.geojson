package geojson

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func convertAToArray(obj *map[string]interface{}) {
	for k, v := range *obj {
		(*obj)[k] = arr(v)
	}
}

func arr(v interface{}) interface{} {
	if a, ok := v.(primitive.A); ok {
		var aa []interface{}
		for _, el := range a {
			aa = append(aa, arr(el))
		}
		return aa
	}
	if asMap, ok := v.(map[string]interface{}); ok {
		for key := range asMap {
			asMap[key] = arr(asMap[key])
		}
	}
	return v
}
