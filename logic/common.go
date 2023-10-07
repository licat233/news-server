package logic

import (
	"errors"
	"news/dao/query"

	"gorm.io/gen/field"
	"gorm.io/gorm"
)

var errNotFoundMsg = errors.New("not exsit")
var ErrRecordNotFound = gorm.ErrRecordNotFound

func articleListFields(DB *query.Query) []field.Expr {
	a := DB.Article
	return []field.Expr{a.ID, a.UUID, a.Title, a.Category, a.PublishAt, a.Author, a.Cover, a.ViewsCount, a.Tags}
}
