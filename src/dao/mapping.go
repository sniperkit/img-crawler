package dao

import (
	"reflect"
	"strings"
)

func parseName(field reflect.StructField, tagName string) (tag, fieldName string) {
	// first, set the fieldName to the field's name
	fieldName = field.Name
	// if there's no tag to look for, return the field name
	if tagName == "" {
		return "", fieldName
	}
	// if this tag is not set using the normal convention in the tag,
	// then return the fieldname..  this check is done because according
	// to the reflect documentation:
	//    If the tag does not have the conventional format,
	//    the value returned by Get is unspecified.
	// which doesn't sound great.
	if !strings.Contains(string(field.Tag), tagName+":") {
		return "", fieldName
	}
	// at this point we're fairly sure that we have a tag, so lets pull it out
	tag = field.Tag.Get(tagName)
	// finally, split the options from the name
	parts := strings.Split(tag, ",")
	fieldName = parts[0]

	return tag, fieldName
}

func isZeroOfUnderlyingType(i interface{}) bool {
	return i == reflect.Zero(reflect.TypeOf(i)).Interface()
}

func GetMapping(r interface{}) map[string]interface{} {

	m := make(map[string]interface{})

	v := reflect.TypeOf(r)
	w := reflect.ValueOf(r)

	for i := 0; i < v.NumField(); i++ {
		tag, fieldName := parseName(v.Field(i), "db")
		// if tag has been filtered or field is AUTO_INCREMENT
		if tag == "" || strings.Contains(string(tag), "AUTO_INCREMENT") {
			continue
		}
		fv := w.FieldByName(v.Field(i).Name)
		// value is zero
		if isZeroOfUnderlyingType(fv.Interface()) {
			continue
		}
		m[fieldName] = fv.Interface()
	}
	return m
}
