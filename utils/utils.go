package utils

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"regexp"
	"time"
	"verottaa/utils/logger"
)

const logTag = "UTILS"

var localLogger = logger.CreateLogger(logTag)

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
		localLogger.Error(err)
	}
	return r.FindString(str)
}
