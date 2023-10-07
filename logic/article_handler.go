package logic

import (
	"context"
	"errors"
	"log"
	"news/common/utils"
	"news/dao/model"
	"news/dao/query"
	"news/global"

	"github.com/cheggaaa/pb/v3"
	"gorm.io/gorm"
)

// 处理文章
func ArticleHandler(ctx context.Context) error {
	// bar := pb.StartNew(count)
	// bar.Increment()
	// bar.Finish()

	db := global.DB
	a := db.Article
	articleDo := db.Article.WithContext(ctx)

	count, err := articleDo.Count()
	if err != nil {
		return err
	}
	bar := pb.StartNew(int(count))
	for i := 0; i < int(count); i++ {
		article, err := articleDo.Select(a.ID, a.Category, a.UUID, a.Tags).Offset(i).Limit(1).First()
		if errors.Is(err, ErrRecordNotFound) {
			break
		} else if err != nil {
			return err
		}
		do := articleDo.Where(a.ID.Eq(article.ID)).Limit(1)
		if article.UUID == "" {
			_, err = do.Update(a.UUID, utils.UUID())
			if err != nil {
				return err
			}
		}
		if article.Tags == "" {
			_, err = do.Update(a.Tags, a.Category)
			if err != nil {
				return err
			}
		}
		bar.Increment()
	}
	bar.Finish()
	return nil
}

func queryCategoryByName(ctx context.Context, db *query.Query, name string) (*model.Category, error) {
	categoryDo := db.Category.WithContext(ctx)
	res, err := categoryDo.Where(db.Category.Name.Eq(name)).First()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// fmt.Println("不存在")
		if err := categoryDo.Save(&model.Category{
			Name: name,
		}); err != nil {
			log.Fatalln(err)
		}
		//回归
		return queryCategoryByName(ctx, db, name)
	} else if err != nil {
		log.Fatalln(err)
	}
	return res, nil
}
