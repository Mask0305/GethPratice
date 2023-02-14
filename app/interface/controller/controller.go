package controller

import (
	"scp/app/usecase"

	"github.com/gin-gonic/gin"
)

type controller struct {
	amtUsecase usecase.AMT
}

type Controller interface {
	Router(app *gin.Engine)
}

func New(amtUsecase usecase.AMT) Controller {
	return &controller{
		amtUsecase: amtUsecase,
	}
}

func (c *controller) Router(app *gin.Engine) {
	Group := app.Group("/amt")
	{
		Group.GET("/total_supply", c.TotalSupply)
		Group.GET("/owner", c.Owner)
		Group.GET("/name", c.Name)
		Group.GET("/symbol", c.Symbol)
		Group.GET("/balance_of/:addr", c.BalanceOf)
		Group.POST("/mint", c.Mint)
		Group.POST("/approve", c.Approve)
		Group.POST("/transfer_from", c.TransferFrom)
		Group.POST("/transfer", c.Transfer)

	}

}

func (c *controller) TotalSupply(ctx *gin.Context) {

	total := c.amtUsecase.TotalSupply()

	ctx.JSON(200, total)
}

func (c *controller) Owner(ctx *gin.Context) {
	owner := c.amtUsecase.Owner()
	ctx.JSON(200, owner)
}

func (c *controller) Name(ctx *gin.Context) {
	name := c.amtUsecase.Name()
	ctx.JSON(200, name)
}

func (c *controller) Symbol(ctx *gin.Context) {
	symbol := c.amtUsecase.Symbol()
	ctx.JSON(200, symbol)
}

func (c *controller) BalanceOf(ctx *gin.Context) {
	addr := ctx.Param("addr")

	balance := c.amtUsecase.BalanceOf(addr)

	ctx.JSON(200, balance)

}

func (c *controller) Mint(ctx *gin.Context) {

	toAddr := ctx.PostForm("toAddr")
	valueData := ctx.PostForm("valueData")
	nonce := ctx.PostForm("nonce")
	sigHex := ctx.PostForm("sigHex")

	if err := c.amtUsecase.Mint(toAddr, valueData, nonce, sigHex); err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(200, "success")
}

func (c *controller) Approve(ctx *gin.Context) {

	senderAddr := ctx.PostForm("senderAddr")
	spenderAddr := ctx.PostForm("spenderAddr")
	valueData := ctx.PostForm("valueData")

	if err := c.amtUsecase.Approve(senderAddr, spenderAddr, valueData); err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(200, "success")

}

func (c *controller) TransferFrom(ctx *gin.Context) {

	senderAddr := ctx.PostForm("senderAddr")
	fromAddr := ctx.PostForm("fromAddr")
	toAddr := ctx.PostForm("toAddr")
	valueData := ctx.PostForm("valueData")
	nonce := ctx.PostForm("nonce")
	sigHex := ctx.PostForm("sigHex")

	if err := c.amtUsecase.TransferFrom(senderAddr, fromAddr, toAddr, valueData, nonce, sigHex); err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(200, "success")
}

func (c *controller) Transfer(ctx *gin.Context) {

	senderAddr := ctx.PostForm("senderAddr")
	toAddr := ctx.PostForm("toAddr")
	valueData := ctx.PostForm("valueData")
	nonce := ctx.PostForm("nonce")
	sigHex := ctx.PostForm("sigHex")

	if err := c.amtUsecase.Transfer(senderAddr, toAddr, valueData, nonce, sigHex); err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(200, "success")
}
