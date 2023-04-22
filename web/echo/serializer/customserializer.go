package serializer

import (
	"fmt"
	"net/http"

	"github.com/pedidopago/go-common/web/echo/json"

	"github.com/labstack/echo/v4"
)

// CustomJSONSerializer implements JSON encoding using encoding/json.
type CustomJSONSerializer struct{}

// Serialize converts an interface into a json and writes it to the response.
// You can optionally use the indent parameter to produce pretty JSONs.
func (d CustomJSONSerializer) Serialize(c echo.Context, i interface{}, indent string) error {
	return echo.DefaultJSONSerializer{}.Serialize(c, i, indent)
}

// Deserialize reads a JSON from a request body and converts it into an interface.
func (d CustomJSONSerializer) Deserialize(c echo.Context, i interface{}) error {
	err := json.NewDecoderContext(c, c.Request().Body).Decode(i)
	if ute, ok := err.(*json.UnmarshalTypeError); ok {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Unmarshal type error: expected=%v, got=%v, field=%v, offset=%v", ute.Type, ute.Value, ute.Field, ute.Offset)).SetInternal(err)
	} else if se, ok := err.(*json.SyntaxError); ok {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Syntax error: offset=%v, error=%v", se.Offset, se.Error())).SetInternal(err)
	}
	return err
}
