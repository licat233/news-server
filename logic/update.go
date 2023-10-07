package logic

import (
	"context"
	"errors"
	"fmt"
	"log"
	"news/common/resp"
	"news/dao/model"
	"news/dao/query"
	"news/global"
	"time"

	"github.com/cheggaaa/pb/v3"
	"github.com/gin-gonic/gin"
	"github.com/tealeg/xlsx"
	"gorm.io/gorm"
)

func UpdateArticle(ctx *gin.Context) {
	err := UpdateArticleHandler(ctx)
	resp.New(ctx, nil, err)
}

func UpdateArticleHandler(ctx context.Context) error {
	startAt := time.Now()

	db := global.DB

	articleDo := db.Article.WithContext(ctx)
	excelFileName := global.V.GetString("file")
	xlFile, err := xlsx.OpenFile(excelFileName)
	if err != nil {
		return err
	}
	count := 0
	for _, sheet := range xlFile.Sheets {
		count += len(sheet.Rows)
	}
	fmt.Println("准备就绪，开始读取...")
	bar := pb.StartNew(count)
	for _, sheet := range xlFile.Sheets {
		for index, row := range sheet.Rows {
			if index == 0 {
				continue
			}
			// for _, cell := range row.Cells {
			// 	text := cell.String()
			// 	fmt.Printf("%s\n", text)
			// }
			if len(row.Cells) < 5 {
				continue
			}
			categoryName := row.Cells[0].String()
			//查询类id
			category, _ := queryCategoryByName(ctx, db, categoryName)

			title := row.Cells[1].String()
			publishAt := row.Cells[2].String()
			source := row.Cells[3].String()
			content := row.Cells[4].String()
			// url := row.Cells[5]
			cover := row.Cells[6].String()

			//存储前，先判断是否存在
			_, err := articleDo.Where(db.Article.Title.Eq(title)).First()
			if err == nil {
				//已存在，不存储
				bar.Increment()
				continue
			} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}
			article := model.Article{
				ID:            0,
				FromID:        "",
				Title:         title,
				Summary:       "",
				ColumnID:      fmt.Sprintf("%v", category.ID),
				ColumnName:    category.Name,
				AuthorName:    source,
				AuthorEmail:   "",
				AuthorAvatar:  "",
				PostID:        "",
				Cover:         cover,
				Content:       content,
				ViewsCount:    0,
				CommentsCount: 0,
				StoreAt:       publishAt,
				PublishedAt:   publishAt,
				Type:          "",
				CreatedAt:     time.Time{},
				UpdatedAt:     time.Time{},
			}
			if err := articleDo.Save(&article); err != nil {
				return err
			}
			bar.Increment()
		}
	}
	bar.Finish()
	subTime := time.Now().Sub(startAt)
	fmt.Println("耗时：", subTime)
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
