package repository

import (
	"verottaa/modules/common"
	interfaces "verottaa/modules/plans/common"
)

type Repository interface {
	common.Destroyable
	interfaces.Reader
	interfaces.Writer
}
