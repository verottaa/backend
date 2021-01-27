package service

import (
	"verottaa/modules/common"
	interfaces "verottaa/modules/users/common"
)

type Service interface {
	common.Destroyable
	interfaces.Reader
	interfaces.Writer
}
