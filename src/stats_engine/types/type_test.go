package types

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"testing"
)

var snakeCaseRegex = regexp.MustCompile(`^[a-z0-9]+(_[a-z0-9]+)*$`)

func isSnakeCase(s string) bool {
	return snakeCaseRegex.MatchString(s)
}

// assertSnakeCaseTags uses reflection to inspect a struct and verify that all
// `parquet:"name=..."` tags are in snake_case.
func assertSnakeCaseTags(t *testing.T, v interface{}) {
	typ := reflect.TypeOf(v)
	if typ.Kind() != reflect.Struct {
		t.Fatalf("assertSnakeCaseTags expected a struct, but got %v", typ.Kind())
	}

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)

		parquetTag := field.Tag.Get("parquet")
		if parquetTag == "" {
			continue
		}

		parts := strings.Split(parquetTag, ",")
		for _, part := range parts {
			if strings.HasPrefix(part, "name=") {
				columnName := strings.TrimPrefix(part, "name=")
				if !isSnakeCase(columnName) {
					t.Errorf("Field '%s' in struct '%s' has a non-snake_case parquet name tag: '%s'", field.Name, typ.Name(), columnName)
				}
				// Found the name tag, no need to check other parts for this field.
				break
			}
		}
	}
}

// TestParquetTagConventions is the main test function that runs the convention check
// on all relevant structs.
func TestParquetTagConventions(t *testing.T) {
	structsToTest := []struct {
		name   string
		target interface{}
	}{
		{"AnnualEarningRecordParquet", AnnualEarningRecordParquet{}},
		{"CombinedPriceRecordParquet", CombinedPriceRecordParquet{}},
		//! Add other Parquet structs here in the future
	}

	for _, tc := range structsToTest {
		t.Run(fmt.Sprintf("Given the %s struct, verify its parquet tags are snake_case", tc.name), func(t *testing.T) {
			assertSnakeCaseTags(t, tc.target)
		})
	}
}
