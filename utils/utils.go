package utils

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"regexp"
	"time"
)

var globalContext = context.Background()

func NewObjectId() primitive.ObjectID {
	return primitive.NewObjectID()
}

func GetContext() context.Context {
	ctx, _ := context.WithTimeout(globalContext, 5*time.Second)
	return ctx
}

func IdFromInterfaceToString(int interface{}) string {
	str := fmt.Sprintf("%v", int)
	r, err := regexp.Compile("\\w{24}")
	if err != nil {
		// TODO: логирование
		return ""
	}
	return r.FindString(str)
}

func ParseTime(s string) time.Time {
	result, err := time.Parse(time.RFC822, s)
	if err != nil {

	}
	return result
}
