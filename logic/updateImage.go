package logic

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"news/common/resp"
	"news/common/utils"
	"news/dao/model"
	"news/global"
	"path/filepath"
	"strings"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/cheggaaa/pb/v3"
	"github.com/gin-gonic/gin"
	"github.com/gocolly/colly/v2"
)

func UpdateImage(ctx *gin.Context) {
	err := UpdateImageHandler(ctx)
	if err != nil {
		resp.IntErr(ctx, err)
		return
	}
	resp.Ok(ctx, nil)
}

var Err301 = ""

func UpdateImageHandler(ctx context.Context) error {
	a := global.DB.Article
	do := a.WithContext(ctx)
	count, err := do.Count()
	if err != nil {
		return err
	}
	if count == 0 {
		return nil
	}
	bucket, err := CreateOssBucket()
	if err != nil {
		return err
	}
	imageDo := global.DB.Image.WithContext(ctx)
	bar := pb.StartNew(int(count))

	var coreImageUrlHandler = func(i int) {
		defer func() {
			bar.Increment()
		}()
		article, err := do.Select().Limit(1).Offset(i).First()
		if errors.Is(err, ErrRecordNotFound) {
			return
		} else if err != nil {
			fmt.Println(err)
			return
		}
		cover := article.Cover
		if cover = strings.TrimSpace(cover); len(cover) == 0 {
			return
		}
		//判斷是否已經存儲與oss中
		endpoint := global.V.GetString("endpoint")
		if strings.Contains(cover, endpoint) {
			return
		}
		//查詢是否已經存在oss中
		image, err := imageDo.Where(global.DB.Image.Source.Eq(cover)).First()
		if errors.Is(err, ErrRecordNotFound) {
			newCover, err := uploadToOSS(bucket, cover)
			if err != nil {
				fmt.Println(err)
				return
			}
			if err = imageDo.Save(&model.Image{
				ID:     0,
				Source: cover,
				Src:    newCover,
			}); err != nil {
				fmt.Println(err)
				return
			}
			article.Cover = newCover
		} else if err != nil {
			fmt.Println(err)
			return
		} else {
			article.Cover = image.Src
		}
		_, err = do.Where(a.ID.Eq(article.ID)).UpdateColumn(a.Cover, article.Cover)
		if err != nil {
			fmt.Println(err)
			return
		}

	}
	var n = int(count)
	process(10, n, coreImageUrlHandler)

	bar.Finish()
	return nil
}

func uploadToOSS(bucket *oss.Bucket, imageURL string) (string, error) {
	body, status := getImageBody(imageURL)
	if status != 200 {
		return "", fmt.Errorf("\nrequst 301 error: %s\n", imageURL)
	}
	if len(body) == 0 {
		return "", fmt.Errorf("\nrequst image failed: %s\n", imageURL)
	}

	name := utils.UUID()
	ext := filepath.Ext(imageURL)
	pathname := fmt.Sprintf("images/%s.%s", name, ext)
	ossDomain := global.V.GetString("ossDomain")

	if err := bucket.PutObject(pathname, bytes.NewReader(body)); err != nil {
		return "", err
	}
	res := fmt.Sprintf("%s/%s", ossDomain, pathname)
	return res, nil
}

func getImageBody(imageURL string) (body []byte, status int) {
	// 发送HTTP GET请求获取图片内容
	c := colly.NewCollector()
	c.SetRequestTimeout(10 * time.Second)
	proxyURL := "http://127.0.0.1:7890"
	if err := c.SetProxy(proxyURL); err != nil {
		log.Fatalln(err)
	}
	c.OnResponse(func(r *colly.Response) {
		status = r.StatusCode
		if r.StatusCode != 200 {
			return
		}
		body = r.Body
	})
	c.Visit(imageURL)
	c.Wait()
	return
}

func CreateOssBucket() (*oss.Bucket, error) {
	/// 从环境变量中获取访问凭证。运行本代码示例之前，请确保已设置环境变量OSS_ACCESS_KEY_ID和OSS_ACCESS_KEY_SECRET。
	provider, err := oss.NewEnvironmentVariableCredentialsProvider()
	if err != nil {
		return nil, fmt.Errorf("Error read env: %s", err)
	}
	// 配置阿里云OSS的AccessKey ID和AccessKey Secret
	accessKeyID := global.V.GetString("aliOssId")
	accessKeySecret := global.V.GetString("aliOssSecret")

	// 配置OSS存储桶的名称、Endpoint和访问域名
	bucketName := global.V.GetString("bucketName")
	endpoint := global.V.GetString("endpoint") // 根据你的实际情况修改，例如：oss-cn-hangzhou.aliyuncs.com
	// ossDomain := global.V.GetString("ossDomain") // 根据你的实际情况修改，例如：https://your-bucket-name.oss-cn-hangzhou.aliyuncs.com

	// 创建OSS客户端对象
	client, err := oss.New(endpoint, accessKeyID, accessKeySecret, oss.SetCredentialsProvider(&provider))
	if err != nil {
		return nil, fmt.Errorf("Error creating OSS client: %s", err)
	}

	// 获取存储桶对象
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		return nil, fmt.Errorf("Error getting OSS bucket: %s", err)
	}
	return bucket, nil
}
