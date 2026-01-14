package envloader

import (
	"fmt"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

func LoadEnvLazy[T any](target *T) error {
	return loadEnv[T](target, true)
}

func LoadEnv[T any](target *T) error {
	return loadEnv[T](target, false)
}

func loadEnv[T any](target *T, lazy bool) error {
	gc := reflect.ValueOf(&target).Elem().Elem()

	for i := 0; i < gc.NumField(); i++ {
		field := gc.Type().Field(i)
		envTag, ok := field.Tag.Lookup("env")
		if !ok && !lazy {
			continue
		}

		var envName string
		var envOptions []string
		var required bool
		envTagValues := strings.Split(envTag, ",")
		envName, envOptions = envTagValues[0], envTagValues[1:]
		if envName == "-" {
			continue
		}
		if envName == "" {
			if !lazy {
				continue
			}
			envName = ToSnakeCase(field.Name)
		}

		required = false
		for _, envOption := range envOptions {
			switch envOption {
			case "required":
				required = true
			case "-":
				continue
			case "":
				continue
			default:
				return fmt.Errorf("invalid envValue tag %q", envOption)
			}
		}

		envValue, ok := os.LookupEnv(envName)
		if !ok {
			if required {
				return fmt.Errorf("environment variable %s must be provided\n", envName)
			}
		}

		f := gc.Field(i)
		switch f.Kind() {
		case reflect.String:
			f.SetString(envValue)
		case reflect.Bool:
			boolVal, err := strconv.ParseBool(envValue)
			if err != nil {
				return fmt.Errorf("environment variable %s must be a valid boolean value\n", envName)
			}
			f.SetBool(boolVal)
		case reflect.Int:
			intVal, err := strconv.ParseInt(envValue, 10, 64)
			if err != nil {
				return fmt.Errorf("environment variable %s must be a valid integer value\n", envName)
			}
			f.SetInt(intVal)
		case reflect.Float64:
			float64Val, err := strconv.ParseFloat(envValue, 64)
			if err != nil {
				return fmt.Errorf("environment variable %s must be a valid integer value\n", envName)
			}
			f.SetFloat(float64Val)
		case reflect.Float32:
			float32Val, err := strconv.ParseFloat(envValue, 32)
			if err != nil {
				return fmt.Errorf("environment variable %s must be a valid integer value\n", envName)
			}
			f.SetFloat(float32Val)
		default:
			return fmt.Errorf("unsupported field type %s for field %s\n", f.Kind(), field.Name)
		}
	}
	return nil
}

func ToSnakeCase(str string) string {
	var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
	var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

	snake := strings.Replace(str, "-", "_", -1)
	snake = strings.Replace(snake, " ", "_", -1)

	snake = matchFirstCap.ReplaceAllString(snake, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")

	var multiUnderline = regexp.MustCompile("(__+)")
	snake = multiUnderline.ReplaceAllString(snake, "_")

	return strings.ToUpper(snake)
}
