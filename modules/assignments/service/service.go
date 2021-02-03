package service

import (
	"verottaa/modules/common"
	interfaces "verottaa/modules/plans/common"
)

type Service interface {
	common.Destroyable
	interfaces.Reader
	interfaces.Writer
}
