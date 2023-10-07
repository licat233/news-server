package logic

import (
	"errors"
	"news/dao/query"

	"gorm.io/gen/field"
	"gorm.io/gorm"
)

var errNotFound = errors.New("not exsit")
var ErrRecordNotFound = gorm.ErrRecordNotFound

func articleListFields(DB *query.Query) []field.Expr {
	a := DB.Article
	return []field.Expr{a.ID, a.Title, a.AuthorName, a.Cover, a.PublishedAt, a.ColumnID, a.ColumnName, a.ViewsCount}
}
