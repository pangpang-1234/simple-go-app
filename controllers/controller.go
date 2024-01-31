package controllers

import (
	"math"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type pagingResult struct {
	Page      int `json:"page"`
	Limit     int `json:"limit"`
	PrevPage  int `json:"prevPage"`
	NextPage  int `json:"nextPage"`
	Count     int `json:"count"`
	TotalPage int `json:"totalPage"`
}

type pagination struct {
	ctx     *gin.Context
	query   *gorm.DB
	records interface{}
}

func (p *pagination) paginate() *pagingResult {
	// 1. Get limit, page from ?limit=10&page=2
	limit, _ := strconv.Atoi(p.ctx.DefaultQuery("limit", "12"))
	page, _ := strconv.Atoi(p.ctx.DefaultQuery("page", "1"))

	// 2. count records using go routine
	ch := make(chan int64)
	go p.countRecords(ch)

	// 3. find records
	// limit, offset(first value to search)
	// limit => 10
	// page => 1, 1-10, offset => 0
	// page => 2, 11-20, offset => 10
	offset := (page - 1) * limit
	p.query.Limit(limit).Offset(offset).Find(p.records)

	// 4. count total pages
	count := <-ch
	totalPage := math.Ceil(float64(count) / float64(limit))

	// 5. find next page
	var nextPage int
	if page == int(totalPage) {
		nextPage = int(totalPage)
	} else {
		nextPage = page - 1
	}

	// 6. create pagingResult
	return &pagingResult{
		Page:      page,
		Limit:     limit,
		Count:     int(count),
		PrevPage:  page - 1,
		NextPage:  nextPage,
		TotalPage: int(totalPage),
	}
}

func (p *pagination) countRecords(ch chan int64) {
	var count int64
	p.query.Model(p.records).Count(&count)

	ch <- count
}
