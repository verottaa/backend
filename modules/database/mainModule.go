package database

import (
	"fmt"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/mongo/options"
	"verottaa/config"
	"verottaa/constants"
)

func init() {
	databaseHost := config.GetConfiguration().GetDatabaseHost()
	err := mgm.SetDefaultConfig(nil, constants.DatabaseName, options.Client().ApplyURI(databaseHost))
	if err != nil {
		fmt.Println(err.Error())
	}
}
