package handler

import (
	"net/http"
	"strconv"

	"github.com/tencentad/martech/cmd/web/ginx"
	"github.com/tencentad/martech/pkg/orm"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func handleWithData(c *gin.Context, successMessage string, data interface{}, err error) {
	if err != nil {
		ginx.ResponseWithError(c, http.StatusInternalServerError, err)
		return
	}

	if data == nil {
		ginx.ResponseWithStatus200(c, successMessage)
	} else {
		ginx.ResponseWithStatus200(c, successMessage, data)
	}
}

func dbPageHandler(c *gin.Context, f func(*gorm.DB, *orm.StatementOption) (interface{}, error)) {
	db := orm.GetDB()
	paginationOpt, err := ginx.GetPaginationOption(c)
	if err != nil {
		ginx.ResponseWithError(c, http.StatusBadRequest, err)
		return
	}
	filterOpts, err := ginx.GetFilterOption(c)
	if err != nil {
		ginx.ResponseWithError(c, http.StatusBadRequest, err)
		return
	}

	data, err := f(db, &orm.StatementOption{Pagination: paginationOpt, Filter: filterOpts})
	handleWithData(c, "get success", data, err)
}

func dbListHandler(c *gin.Context, f func(*gorm.DB) (interface{}, error)) {
	data, err := f(orm.GetDB())
	handleWithData(c, "get success", data, err)
}

func dbGetByIdHandler(c *gin.Context, f func(*gorm.DB, uint64) (interface{}, error)) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		ginx.ResponseWithError(c, http.StatusBadRequest, err)
		return
	}

	db := orm.GetDB()
	data, err := f(db, id)
	handleWithData(c, "get success", data, err)
}

func dbGetHandler(c *gin.Context, item interface{}, f func(*gorm.DB) (interface{}, error)) {
	if err := c.ShouldBindJSON(item); err != nil {
		ginx.ResponseWithError(c, http.StatusBadRequest, err)
		return
	}

	db := orm.GetDB()
	data, err := f(db)
	handleWithData(c, "get success", data, err)
}

func dbDeleteMultiple(c *gin.Context, item interface{}, f func(*gorm.DB) error) {
	if err := c.ShouldBindJSON(item); err != nil {
		ginx.ResponseWithError(c, http.StatusBadRequest, err)
		return
	}

	db := orm.GetDB()
	err := f(db)
	handleWithData(c, "delete success", nil, err)
}

func dbEditHandler(c *gin.Context, item interface{}, f func(*gorm.DB) error) {
	if err := c.ShouldBindJSON(item); err != nil {
		ginx.ResponseWithError(c, http.StatusBadRequest, err)
		return
	}

	db := orm.GetDB()
	err := f(db)
	handleWithData(c, "add success", nil, err)
}

func dbDeleteHandler(c *gin.Context, f func(*gorm.DB, uint64) error) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		ginx.ResponseWithError(c, http.StatusBadRequest, err)
		return
	}

	db := orm.GetDB()
	err = f(db, id)
	handleWithData(c, "delete success", nil, err)
}
