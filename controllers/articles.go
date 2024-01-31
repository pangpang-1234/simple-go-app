package controllers

import (
	"simplegoapp/models"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

type Articles struct {
	DB *gorm.DB
}

type createArticleForm struct {
	Title   string                `form:"title" binding:"required"`
	Body    string                `form:"body" binding:"required"`
	Excerpt string                `form:"excerpt" binding:"required"`
	Image   *multipart.FileHeader `form:"image" binding:"required"`
}

type updateArticleForm struct {
	Title   string                `form:"title"`
	Body    string                `form:"body"`
	Excerpt string                `form:"excerpt"`
	Image   *multipart.FileHeader `form:"image"`
}

type articleResponse struct {
	ID      uint   `json:"id"`
	Title   string `json:"title"`
	Excerpt string `json:"excerpt"`
	Body    string `json:"body"`
	Image   string `json:"image"`
}

// for pagination
type articlesPaging struct {
	Items  []articleResponse `json:"items"`
	Paging *pagingResult     `json:"paging"`
}

// Query all articles in database
func (a *Articles) FindAll(ctx *gin.Context) {
	var articles []models.Article

	// /articles => limit => 12, page => 1
	// /articles?limit=10 => limit => 10, page => 1
	// /articles?page=2&limit=4 => limit => 4, page => 2
	// Descending order
	pagination := pagination{ctx: ctx, query: a.DB.Order("id desc"), records: &articles}
	paging := pagination.paginate()

	var serializeArticles []articleResponse
	copier.Copy(&serializeArticles, &articles)
	ctx.JSON(http.StatusOK, gin.H{"articles": articlesPaging{Items: serializeArticles, Paging: paging}})
}

// Query article by ID
func (a *Articles) FindOne(ctx *gin.Context) {
	article, err := a.findArticleByID(ctx)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	serializeArticle := articleResponse{}
	copier.Copy(&serializeArticle, &article)
	ctx.JSON(http.StatusOK, gin.H{"article": article})
}

// Create new article
func (a *Articles) Create(ctx *gin.Context) {
	var form createArticleForm
	if err := ctx.ShouldBind(&form); err != nil {
		//422 unprocessable entity (bad validation)
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}
	// form => article
	var article models.Article
	copier.Copy(&article, &form)

	// model article => db
	if err := a.DB.Create(&article).Error; err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	a.setArticleImage(ctx, &article)
	serializeArticle := articleResponse{}
	copier.Copy(&serializeArticle, &article)

	ctx.JSON(http.StatusCreated, gin.H{"article": article})
}

func (a *Articles) Update(ctx *gin.Context) {
	var form updateArticleForm
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return 
	}

	article, err := a.findArticleByID(ctx)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if err := a.DB.Model(&article).Updates(&form).Error; err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error})
		return
	}

	a.setArticleImage(ctx, article)

	var serializeArticle  articleResponse
	copier.Copy(&serializeArticle, article)
	ctx.JSON(http.StatusOK, gin.H{"article": serializeArticle})
}

func (a *Articles) Delete(ctx *gin.Context) {
	article, err := a.findArticleByID(ctx)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	a.DB.Unscoped().Delete(&article)
	ctx.Status(http.StatusNoContent)
}

// Save article image
func (a *Articles) setArticleImage(ctx *gin.Context, article *models.Article) error {
	file, err := ctx.FormFile("image")
	if err != nil || file == nil {
		return err
	}
	// check and remove existing image
	// http://localhost:8080/upload/articles/<ID>/image.png
	if article.Image != "" {
		article.Image = strings.Replace(article.Image, os.Getenv("HOST"), "", 1)
		pwd, _ := os.Getwd()
		os.Remove(pwd + article.Image)
	}

	path := "uploads/articles/" + strconv.Itoa(int(article.ID))
	os.MkdirAll(path, 0755)
	filename := path + "/" + file.Filename
	if err := ctx.SaveUploadedFile(file, filename); err != nil {
		return err
	}

	article.Image = os.Getenv("HOST") + "/" + filename
	a.DB.Save(article)

	return nil
}

// function for finding article by ID
func (a *Articles) findArticleByID(ctx *gin.Context) (*models.Article, error) {
	var article models.Article
	id := ctx.Param("id")

	if err := a.DB.First(&article, id).Error; err != nil {
		return nil, err
	}
	return &article, nil
}


