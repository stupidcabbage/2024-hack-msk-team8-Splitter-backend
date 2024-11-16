package router

import (
	docs "example.com/m/docs"
	"example.com/m/internal/api/v1/adapters/controllers"
	"example.com/m/internal/api/v1/infrastructure/middlewares"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

const prefix string = "/api/v1"

func BindRoutes(e *gin.Engine, a *middlewares.AuthMiddleware, uc *controllers.UserController, ac *controllers.AuthController, m *controllers.MetricController, dc *controllers.DebtController, gc *controllers.GroupController) {
	docs.SwaggerInfo.BasePath = "/api/v1"
	docs.SwaggerInfo.Version = "v1"

	e.POST(prefix+"/users", uc.CreateUser)
	e.POST(prefix+"/auth", ac.AuthorizeUser)
	e.POST(prefix+"/auth/changePassword", a.Authenticate(), ac.ChangePassword)
	e.POST(prefix+"/groups", a.Authenticate(), gc.CreateGroup)
	e.POST(prefix+"/groups/:invite_code/join", a.Authenticate(), gc.AddToGroup)
	e.GET(prefix+"/groups/:group_id/users", a.Authenticate(), gc.GetUsersByGroup)
	e.GET(prefix+"/groups/:group_id/debts/incoming", a.Authenticate(), gc.GetIncomingDebtsByGroup)
	e.GET(prefix+"/groups/:group_id/debts/outcoming", a.Authenticate(), gc.GetOutcomingDebtsByGroup)
	e.GET(prefix+"/groups/created", a.Authenticate(), gc.GetCreatedGroups)
	e.GET(prefix+"/users/:username", a.Authenticate(), uc.GetUserByUsername)
	e.GET(prefix+"/users/me", a.Authenticate(), uc.GetUserProfile)
	e.GET(prefix+"/metrics", m.GetMetrics())
	e.POST(prefix+"/debts", a.Authenticate(), dc.CreateDebt)
	e.GET(prefix+"/debts/outcoming", a.Authenticate(), dc.GetOutcomingDebts)
	e.GET(prefix+"/debts/incoming", a.Authenticate(), dc.GetIncomingDebts)
	e.GET(prefix+"/debts/:id", a.Authenticate(), dc.GetDebtById)
	e.PATCH(prefix+"/debts/:id/close", a.Authenticate(), dc.CloseDebt)
	e.PUT(prefix+"/debts/:id/payedAmount/increase", a.Authenticate(), dc.IncreaseDebtPayedAmount)
	e.GET(prefix+"/debts/closed", a.Authenticate(), dc.GetClosedDebts)
	e.GET(prefix+"/debts/stats", a.Authenticate(), dc.GetDebtsMetrics)

	ginSwagger.WrapHandler(swaggerfiles.Handler,
		ginSwagger.URL("http://localhost:8000/swagger/doc.json"),
		ginSwagger.DefaultModelsExpandDepth(-1))
	e.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

}
