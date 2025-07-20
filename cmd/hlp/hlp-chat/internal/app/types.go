package app

import "context"

type IApp interface {
	Run(context.Context) error
}
