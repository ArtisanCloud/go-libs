package helper

import (
	"context"
	"crypto/sha256"
	"fmt"
	. "github.com/ArtisanCloud/PowerLibs/v3/object"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/crypto/bcrypt"
	"log"
	"reflect"
)

/**
 * Transform Array Keys to Camel fmt
 *
 * @param  string $mx
 * @return strID
 */
func TransformArrayKeysToCamel(arrayData HashMap) HashMap {
	arrayTransformedKeys := HashMap{}
	for key, value := range arrayData {

		if reflect.TypeOf(value).Kind() == reflect.Map {
			value = TransformArrayKeysToCamel(value.(HashMap))
			arrayTransformedKeys[Camel(key)] = value
		}

	}
	return arrayTransformedKeys

}

/**
 * Transform Array Keys to Snake fmt
 *
 * @param  string $mx
 * @return strID
 */
func TransformArrayKeysToSnake(arrayData interface{}) HashMap {
	arrayTransformedKeys := HashMap{}

	for key, value := range arrayData.(HashMap) {

		if reflect.TypeOf(value).Kind() == reflect.Map {
			value = TransformArrayKeysToSnake(value.(HashMap))
			arrayTransformedKeys[Camel(key)] = value
		}

	}

	return arrayTransformedKeys

}

///*
// * This function will help you to convert your object from struct to map[string]interface{} based on your JSON tag in your structs.
// * Example how to use posted in sample_test.go file.
// */
//func StructToMap(item interface{}) map[string]interface{} {
//
//	res := map[string]interface{}{}
//	if item == nil {
//		return res
//	}
//	v := reflect.TypeOf(item)
//	reflectValue := reflect.ValueOf(item)
//	reflectValue = reflect.Indirect(reflectValue)
//
//	if v.Kind() == reflect.Ptr {
//		v = v.Elem()
//	}
//	for i := 0; i < v.NumField(); i++ {
//		tag := v.Field(i).Tag.Get("json")
//		field := reflectValue.Field(i).Interface()
//		if tag != "" && tag != "-" {
//			if v.Field(i).Type.Kind() == reflect.Struct {
//				res[tag] = StructToMap(field)
//			} else {
//				res[tag] = field
//			}
//		}
//	}
//	return res
//}

func HashPassword(plainPassword string) (hashedPassword string) {

	hashed := sha256.Sum256([]byte(plainPassword))
	strHashed := fmt.Sprintf("%x", hashed)

	//fmt.Println("Hashed password", strHashed)

	return strHashed

}

func EncodePassword(hashedPassword string) (encodedPassword string) {

	encoded, err := bcrypt.GenerateFromPassword([]byte(hashedPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println("Hashed password", string(encoded))

	return string(encoded)
}

func EncodePlainPassword(plainPassword string) (encodedPassword string) {
	hashedPassword := HashPassword(plainPassword)
	encodedPassword = EncodePassword(hashedPassword)

	return encodedPassword
}

func CheckPassword(hashedPassword string, password string) (isPasswordValid bool) {

	//fmt.Printf("hashedPassword %s\r\n", hashedPassword)
	//fmt.Printf("password %s\n", password)

	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
		fmt.Printf("%x", err)
		return false
	}

	return true
}

func SpanIDFromContext(ctx context.Context) string {
	spanCtx := trace.SpanContextFromContext(ctx)
	if spanCtx.HasSpanID() {
		return spanCtx.SpanID().String()
	}

	return ""
}

func TraceIDFromContext(ctx context.Context) string {
	spanCtx := trace.SpanContextFromContext(ctx)
	if spanCtx.HasTraceID() {
		return spanCtx.TraceID().String()
	}

	return ""
}
