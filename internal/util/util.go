package util

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"regexp"
	"strings"

	"github.com/fatih/color"
)

func SanitizeFileName(fileName string) string {
	// Replace any characters that aren't letters, numbers, or safe punctuation with a hyphen
	reg, err := regexp.Compile("[^a-zA-Z0-9-_]+")
	if err != nil {
		log.Fatal(err)
	}
	return reg.ReplaceAllString(fileName, "-")
}

var (
	greenColor  = color.New(color.FgGreen)
	yellowColor = color.New(color.FgYellow)
	redColor    = color.New(color.FgRed)
	blueColor   = color.New(color.FgBlue)
)

func Green(format string, args ...interface{}) {
	greenColor.Fprintf(os.Stdout, format, args...)
}

func Yellow(format string, args ...interface{}) {
	yellowColor.Fprintf(os.Stdout, format, args...)
}

func Red(format string, args ...interface{}) {
	redColor.Fprintf(os.Stdout, format, args...)
}

func Blue(format string, args ...interface{}) {
	blueColor.Fprintf(os.Stdout, format, args...)
}

func PrettyPrint(v interface{}) {
	printRecursive(v, 0)
}

func printRecursive(v interface{}, indent int) {
	value := reflect.ValueOf(v)
	typeOf := value.Type()

	prefix := strings.Repeat("  ", indent)

	switch value.Kind() {
	case reflect.Struct:
		for i := 0; i < value.NumField(); i++ {
			fieldName := typeOf.Field(i).Name
			fieldValue := value.Field(i).Interface()

			if reflect.TypeOf(fieldValue).Kind() == reflect.Slice && reflect.TypeOf(fieldValue).Elem().Kind() == reflect.Struct {
				// Special handling for slices of structs.
				Green("%s%s:\n", prefix, fieldName) // Directly using GreenPrint
				for j := 0; j < reflect.ValueOf(fieldValue).Len(); j++ {
					printRecursive(reflect.ValueOf(fieldValue).Index(j).Interface(), indent+1)
				}
			} else if reflect.TypeOf(fieldValue).Kind() == reflect.Struct {
				Green("%s%s:\n", prefix, fieldName) // Directly using GreenPrint
				printRecursive(fieldValue, indent+1)
			} else {
				combinedStr := fmt.Sprintf("%s%s: ", prefix, fieldName)
				Yellow("%s%v\n", combinedStr, fieldValue) // Use YellowPrint
			}
		}
	case reflect.Slice:
		for i := 0; i < value.Len(); i++ {
			elem := value.Index(i)
			Blue("%s- Element %d:\n", prefix, i) // Directly using BluePrint
			printRecursive(elem.Interface(), indent+1)
		}
	default:
		Yellow("%s%v\n", prefix, v) // Directly using YellowPrint
	}
}
