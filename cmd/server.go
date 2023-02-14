package cmd

import (
	"context"
	"fmt"
	"log"
	"scp/app/infra/erc20"
	"scp/app/infra/ethereum"
	"scp/app/interface/controller"
	"scp/app/interface/router"
	"scp/app/usecase"

	"github.com/gin-gonic/contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"go-micro.dev/v4"
	"go.uber.org/fx"
)

// serverpCmd represents the serverp command
var serverpCmd = &cobra.Command{
	Use:   "server",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		server()
		fmt.Println("server called")
	},
}

func init() {
	rootCmd.AddCommand(serverpCmd)
}
func newGin() *gin.Engine {
	g := gin.New()
	//Cors 跨域设置
	g.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowedMethods:   []string{"*"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: false,
	}))
	g.Use(gin.Recovery())
	return g
}

func server() {

	app := fx.New(
		fx.NopLogger,
		fx.Provide(
			context.Background,
			newGin,
			erc20.NewAMTtoken,
			ethereum.New,
			usecase.New,
			controller.New,
			controller.NewMetaMaskController,
			router.New,
		),
		fx.Invoke(NewHTTPServer),
	)

	if err := app.Err(); err != nil {
		log.Print(err)
	}

	app.Run()

}

// NewHTTPServer -
func NewHTTPServer(lc fx.Lifecycle, g *gin.Engine, r router.Router) {

	srv := micro.NewService(
		micro.Name("crypto"),
		micro.Version("0.0.1"),
	)
	srv.Init()

	lc.Append(
		fx.Hook{
			OnStart: func(context.Context) error {
				r.Set()
				//go srv.Run()
				go g.Run(":8001")

				return nil
			},
			OnStop: func(ctx context.Context) error {
				return nil
			},
		})

}
