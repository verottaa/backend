package utils

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
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
		log.WithFields(log.Fields{
			"package":  "utils",
			"function": "IdFromInterfaceToString",
			"error":    err,
			"cause":    "Compile",
		}).Error("Unexpected error")
		return ""
	}
	return r.FindString(str)
}

func ParseTime(s string) time.Time {
	result, err := time.Parse(time.RFC822, s)
	if err != nil {
		log.WithFields(log.Fields{
			"package":  "utils",
			"function": "ParseTime",
			"error":    err,
			"cause":    "time parsing",
		}).Error("Unexpected error")
	}
	return result
}
