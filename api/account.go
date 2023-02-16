package api

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/yagoernandes/simplebank/db/sqlc"
)

type CreateAccountParams struct {
	Owner    string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required,oneof=USD EUR"`
}

func (s *Server) createAccount(ctx *gin.Context) {
	var params CreateAccountParams
	if err := ctx.ShouldBindJSON(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorMessage(err))
		return
	}

	arg := db.CreateAccountParams{
		Owner:    params.Owner,
		Balance:  0,
		Currency: params.Currency,
	}

	account, err := s.store.CreateAccount(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorMessage(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

type GetAccountParams struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (s *Server) getAccount(ctx *gin.Context) {
	var params GetAccountParams
	if err := ctx.ShouldBindUri(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorMessage(err))
		return
	}

	account, err := s.store.GetAccount(ctx, params.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, ErrorMessage(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, ErrorMessage(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

type ListAccountParams struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=1,max=100"`
}

func (s *Server) listAccounts(ctx *gin.Context) {
	var params ListAccountParams
	if err := ctx.ShouldBindQuery(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorMessage(err))
		return
	}

	arg := db.ListAccountsParams{
		Limit:  params.PageSize,
		Offset: (params.PageID - 1) * params.PageSize,
	}

	accounts, err := s.store.ListAccounts(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorMessage(err))
		return
	}

	ctx.JSON(http.StatusOK, accounts)
}
