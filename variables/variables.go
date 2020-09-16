package variables

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

var (
	Ctx, _       = context.WithTimeout(context.Background(), 10*time.Second)
	DatabaseHost = "mongodb://localhost:27017"
	Client       *mongo.Client
)
