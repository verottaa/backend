package controller_helpers

import (
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"verottaa/models"
)

type variableReader struct {
}

type VariableReader interface {
	models.Destroyable
	GetVariableByName(r *http.Request, name string) string
	IdReader
}

type IdReader interface {
	GetObjectIDFromString(s string) primitive.ObjectID
	GetObjectIdFromQueryByName(r *http.Request, name string) primitive.ObjectID
}

var variableReaderDestroyCh = make(chan bool)
var variableReaderInstance *variableReader

func GetVariableReader() VariableReader {
	once.Do(func() {
		variableReaderInstance = createInstance()
		go func() {
			for
			{
				select {
				case <-variableReaderDestroyCh:
					return
				}
			}
		}()
	})

	return variableReaderInstance
}

func createInstance() *variableReader {
	return &variableReader{}
}

func (vReader variableReader) Destroy() {
	variableReaderDestroyCh <- true
	close(variableReaderDestroyCh)
	variableReaderInstance = nil
}

func (vReader variableReader) GetObjectIDFromString(s string) primitive.ObjectID {
	result, _ := primitive.ObjectIDFromHex(s)
	return result
}

func (vReader variableReader) GetVariableByName(r *http.Request, name string) string {
	vars := mux.Vars(r)
	variable := vars[name]
	return variable
}

func (vReader variableReader) GetObjectIdFromQueryByName(r *http.Request, name string) primitive.ObjectID {
	stringId := vReader.GetVariableByName(r, name)
	return vReader.GetObjectIDFromString(stringId)
}
