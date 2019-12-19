package geojson

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	"go.mongodb.org/mongo-driver/bson"
)

// A GeometryType serves to enumerate the different GeoJSON geometry types.
type GeometryType string

// The geometry types supported by GeoJSON 1.0
const (
	GeometryPoint           GeometryType = "Point"
	GeometryMultiPoint      GeometryType = "MultiPoint"
	GeometryLineString      GeometryType = "LineString"
	GeometryMultiLineString GeometryType = "MultiLineString"
	GeometryPolygon         GeometryType = "Polygon"
	GeometryMultiPolygon    GeometryType = "MultiPolygon"
	GeometryCollection      GeometryType = "GeometryCollection"
)

var (
	valueOfType        = reflect.ValueOf("type")
	valueOfBbox        = reflect.ValueOf("bbox")
	valueOfCoordinates = reflect.ValueOf("coordinates")
	valueOfGeomeytries = reflect.ValueOf("geometries")
)

// A Geometry correlates to a GeoJSON geometry object.
type Geometry struct {
	Type            GeometryType `json:"type"`
	BoundingBox     []float64    `json:"bbox,omitempty"`
	Point           []float64
	MultiPoint      [][]float64
	LineString      [][]float64
	MultiLineString [][][]float64
	Polygon         [][][]float64
	MultiPolygon    [][][][]float64
	Geometries      []*Geometry
	CRS             map[string]interface{} `json:"crs,omitempty"` // Coordinate Reference System Objects are not currently supported
}

// NewPointGeometry creates and initializes a point geometry with the give coordinate.
func NewPointGeometry(coordinate []float64) *Geometry {
	return &Geometry{
		Type:  GeometryPoint,
		Point: coordinate,
	}
}

// NewMultiPointGeometry creates and initializes a multi-point geometry with the given coordinates.
func NewMultiPointGeometry(coordinates ...[]float64) *Geometry {
	return &Geometry{
		Type:       GeometryMultiPoint,
		MultiPoint: coordinates,
	}
}

// NewLineStringGeometry creates and initializes a line string geometry with the given coordinates.
func NewLineStringGeometry(coordinates [][]float64) *Geometry {
	return &Geometry{
		Type:       GeometryLineString,
		LineString: coordinates,
	}
}

// NewMultiLineStringGeometry creates and initializes a multi-line string geometry with the given lines.
func NewMultiLineStringGeometry(lines ...[][]float64) *Geometry {
	return &Geometry{
		Type:            GeometryMultiLineString,
		MultiLineString: lines,
	}
}

// NewPolygonGeometry creates and initializes a polygon geometry with the given polygon.
func NewPolygonGeometry(polygon [][][]float64) *Geometry {
	return &Geometry{
		Type:    GeometryPolygon,
		Polygon: polygon,
	}
}

// NewMultiPolygonGeometry creates and initializes a multi-polygon geometry with the given polygons.
func NewMultiPolygonGeometry(polygons ...[][][]float64) *Geometry {
	return &Geometry{
		Type:         GeometryMultiPolygon,
		MultiPolygon: polygons,
	}
}

// NewCollectionGeometry creates and initializes a geometry collection geometry with the given geometries.
func NewCollectionGeometry(geometries ...*Geometry) *Geometry {
	return &Geometry{
		Type:       GeometryCollection,
		Geometries: geometries,
	}
}

// MarshalJSON converts the geometry object into the correct JSON.
// This fulfills the json.Marshaler interface.
func (g Geometry) MarshalJSON() ([]byte, error) {
	// defining a struct here lets us define the order of the JSON elements.
	type geometry struct {
		Type        GeometryType           `json:"type"`
		BoundingBox []float64              `json:"bbox,omitempty"`
		Coordinates interface{}            `json:"coordinates,omitempty"`
		Geometries  interface{}            `json:"geometries,omitempty"`
		CRS         map[string]interface{} `json:"crs,omitempty"`
	}

	geo := &geometry{
		Type: g.Type,
	}

	if g.BoundingBox != nil && len(g.BoundingBox) != 0 {
		geo.BoundingBox = g.BoundingBox
	}

	switch g.Type {
	case GeometryPoint:
		geo.Coordinates = g.Point
	case GeometryMultiPoint:
		geo.Coordinates = g.MultiPoint
	case GeometryLineString:
		geo.Coordinates = g.LineString
	case GeometryMultiLineString:
		geo.Coordinates = g.MultiLineString
	case GeometryPolygon:
		geo.Coordinates = g.Polygon
	case GeometryMultiPolygon:
		geo.Coordinates = g.MultiPolygon
	case GeometryCollection:
		geo.Geometries = g.Geometries
	}

	return json.Marshal(geo)
}

// UnmarshalGeometry decodes the data into a GeoJSON geometry.
// Alternately one can call json.Unmarshal(g) directly for the same result.
func UnmarshalGeometry(data []byte) (*Geometry, error) {
	g := &Geometry{}
	err := json.Unmarshal(data, g)
	if err != nil {
		return nil, err
	}

	return g, nil
}

// UnmarshalJSON decodes the data into a GeoJSON geometry.
// This fulfills the json.Unmarshaler interface.
func (g *Geometry) UnmarshalJSON(data []byte) error {
	var object map[string]interface{}
	err := json.Unmarshal(data, &object)
	if err != nil {
		return err
	}
	return decodeGeometry(g, reflect.ValueOf(object))
}

// Scan implements the sql.Scanner interface allowing
// geometry structs to be passed into rows.Scan(...interface{})
// The columns must be received as GeoJSON Geometry.
// When using PostGIS a spatial column would need to be wrapped in ST_AsGeoJSON.
func (g *Geometry) Scan(value interface{}) error {
	var data []byte

	switch value.(type) {
	case string:
		data = []byte(value.(string))
	case []byte:
		data = value.([]byte)
	default:
		return errors.New("unable to parse this type into geojson")
	}

	return g.UnmarshalJSON(data)
}

// MarshalBSON converts the geometry object into the correct JSON.
// This fulfills the bson.Marshaler interface.
func (g Geometry) MarshalBSON() ([]byte, error) {
	type geometry struct {
		Type        GeometryType           `bson:"type"`
		BoundingBox []float64              `bson:"bbox,omitempty"`
		Coordinates interface{}            `bson:"coordinates,omitempty"`
		Geometries  interface{}            `bson:"geometries,omitempty"`
		CRS         map[string]interface{} `bson:"crs,omitempty"`
	}

	geo := &geometry{
		Type: g.Type,
	}

	if g.BoundingBox != nil && len(g.BoundingBox) != 0 {
		geo.BoundingBox = g.BoundingBox
	}

	switch g.Type {
	case GeometryPoint:
		geo.Coordinates = g.Point
	case GeometryMultiPoint:
		geo.Coordinates = g.MultiPoint
	case GeometryLineString:
		geo.Coordinates = g.LineString
	case GeometryMultiLineString:
		geo.Coordinates = g.MultiLineString
	case GeometryPolygon:
		geo.Coordinates = g.Polygon
	case GeometryMultiPolygon:
		geo.Coordinates = g.MultiPolygon
	case GeometryCollection:
		geo.Geometries = g.Geometries
	}

	return bson.Marshal(geo)
}

// UnmarshalBSON decodes the data into a GeoJSON geometry.
// This fulfills the bson.Unmarshaler interface.
func (g *Geometry) UnmarshalBSON(data []byte) error {
	var object map[string]interface{}
	err := bson.Unmarshal(data, &object)
	if err != nil {
		return err
	}

	return decodeGeometry(g, reflect.ValueOf(object))
}

func decodeGeometry(g *Geometry, value reflect.Value) error {
	if value.Kind() != reflect.Map {
		return fmt.Errorf("unable to decode %#v into geometry", value)
	}
	typeProp := mapIndexValue(value, valueOfType)
	if typeProp.Kind() != reflect.String {
		return fmt.Errorf("type property not defined in geometry %#v", value)
	}
	g.Type = GeometryType(typeProp.String())

	bbProp := mapIndexValue(value, valueOfBbox)
	if bbProp.Kind() != reflect.Invalid {
		bb, err := decodeBoundingBoxValue(bbProp)
		if err != nil {
			return err
		}
		g.BoundingBox = bb
	}

	var err error
	switch g.Type {
	case GeometryPoint:
		g.Point, err = decodePosition(mapIndexValue(value, valueOfCoordinates))
	case GeometryMultiPoint:
		g.MultiPoint, err = decodePositionSet(mapIndexValue(value, valueOfCoordinates))
	case GeometryLineString:
		g.LineString, err = decodePositionSet(mapIndexValue(value, valueOfCoordinates))
	case GeometryMultiLineString:
		g.MultiLineString, err = decodePathSet(mapIndexValue(value, valueOfCoordinates))
	case GeometryPolygon:
		g.Polygon, err = decodePathSet(mapIndexValue(value, valueOfCoordinates))
	case GeometryMultiPolygon:
		g.MultiPolygon, err = decodePolygonSet(mapIndexValue(value, valueOfCoordinates))
	case GeometryCollection:
		g.Geometries, err = decodeGeometries(mapIndexValue(value, valueOfGeomeytries))
	}

	return err
}

func decodePosition(pos reflect.Value) ([]float64, error) {
	if pos.Kind() != reflect.Slice {
		return nil, fmt.Errorf("invalid position, got %#v", pos)
	}
	result := make([]float64, pos.Len())
	for i := range result {
		coord := indexValue(pos, i)
		k := coord.Kind()
		if k == reflect.Int || k == reflect.Int8 || k == reflect.Int16 || k == reflect.Int32 || k == reflect.Int64 {
			result[i] = float64(coord.Int())
			continue
		}
		if k == reflect.Float32 || k == reflect.Float64 {
			result[i] = coord.Float()
			continue
		}
		return nil, fmt.Errorf("invalid coordinate in %#v, got %#v", pos, coord)
	}
	return result, nil
}

func decodePositionSet(posSet reflect.Value) ([][]float64, error) {
	if posSet.Kind() != reflect.Slice {
		return nil, fmt.Errorf("invalid set of positions, got %#v", posSet)
	}
	result := make([][]float64, posSet.Len())
	for i := range result {
		pos, err := decodePosition(indexValue(posSet, i))
		if err != nil {
			return nil, err
		}
		result[i] = pos
	}
	return result, nil
}

func decodePathSet(pathSet reflect.Value) ([][][]float64, error) {
	if pathSet.Kind() != reflect.Slice {
		return nil, fmt.Errorf("invalid path set, got %#v", pathSet)
	}
	result := make([][][]float64, pathSet.Len())
	for i := range result {
		posSet, err := decodePositionSet(indexValue(pathSet, i))
		if err != nil {
			return nil, err
		}
		result[i] = posSet
	}
	return result, nil
}

func decodePolygonSet(polygonSet reflect.Value) ([][][][]float64, error) {
	if polygonSet.Kind() != reflect.Slice {
		return nil, fmt.Errorf("invalid polygon, got %#v", polygonSet)
	}
	result := make([][][][]float64, polygonSet.Len())
	for i := range result {
		pathSet, err := decodePathSet(indexValue(polygonSet, i))
		if err != nil {
			return nil, err
		}
		result[i] = pathSet
	}
	return result, nil
}

func decodeGeometries(geoms reflect.Value) ([]*Geometry, error) {
	if geoms.Kind() != reflect.Slice {
		return nil, fmt.Errorf("invalid geometries %#v", geoms)
	}
	geometries := make([]*Geometry, geoms.Len())
	for i := range geometries {
		var g Geometry
		v := indexValue(geoms, i)
		if v.Kind() != reflect.Map {
			return nil, fmt.Errorf("invalid geometry %#v found in geometries %v", v, geoms)
		}
		err := decodeGeometry(&g, v)
		if err != nil {
			return nil, err
		}
		geometries[i] = &g
	}
	return geometries, nil
}

// IsPoint returns true with the geometry object is a Point type.
func (g *Geometry) IsPoint() bool {
	return g.Type == GeometryPoint
}

// IsMultiPoint returns true with the geometry object is a MultiPoint type.
func (g *Geometry) IsMultiPoint() bool {
	return g.Type == GeometryMultiPoint
}

// IsLineString returns true with the geometry object is a LineString type.
func (g *Geometry) IsLineString() bool {
	return g.Type == GeometryLineString
}

// IsMultiLineString returns true with the geometry object is a LineString type.
func (g *Geometry) IsMultiLineString() bool {
	return g.Type == GeometryMultiLineString
}

// IsPolygon returns true with the geometry object is a Polygon type.
func (g *Geometry) IsPolygon() bool {
	return g.Type == GeometryPolygon
}

// IsMultiPolygon returns true with the geometry object is a MultiPolygon type.
func (g *Geometry) IsMultiPolygon() bool {
	return g.Type == GeometryMultiPolygon
}

// IsCollection returns true with the geometry object is a GeometryCollection type.
func (g *Geometry) IsCollection() bool {
	return g.Type == GeometryCollection
}

// mapIndexValue goes to the value behind the key in the map, and makes sure the kind of the value is not an interface
func mapIndexValue(mp reflect.Value, key reflect.Value) reflect.Value {
	return avoidInterface(mp.MapIndex(key))
}

// indexValue goes to the value behind the index in the iterable, and makes sure the kind of the value is not an interface
func indexValue(sl reflect.Value, i int) reflect.Value {
	return avoidInterface(sl.Index(i))
}

// avoidInterface makes sure the kind of the value is not an interface
func avoidInterface(value reflect.Value) reflect.Value {
	if value.Kind() != reflect.Interface {
		return value
	}
	return reflect.ValueOf(value.Interface())
}
