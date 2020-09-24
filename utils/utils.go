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

func HandleError(err error) {
	if err != nil {
		fmt.Println("[ERROR]:", err)
	}
}

func GetContext() context.Context {
	ctx, _ := context.WithTimeout(globalContext, 5*time.Second)
	return ctx
}

func IdFromInterfaceToString(int interface{}) string {
	str := fmt.Sprintf("%v", int)
	r, err := regexp.Compile("\\w{24}")
	HandleError(err)
	return r.FindString(str)
}
