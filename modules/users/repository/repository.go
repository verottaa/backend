package repository

import (
	"verottaa/modules/common"
	interfaces "verottaa/modules/users/common"
)

type Repository interface {
	common.Destroyable
	interfaces.Reader
	interfaces.Writer
}
