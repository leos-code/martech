package router

import (
	"github.com/tencentad/martech/cmd/web/handler"
	"github.com/gin-gonic/gin"
)

func setupDashboard(router gin.IRouter) error {
	router.Use(handler.AuthLoginMiddleware())
	user := router.Group("/user")
	{
		user.GET("/info", handler.UserInfoGetHandler)
		user.POST("/info", handler.UserInfoPostHandler)
		user.POST("/tenant", handler.UserTenantHandler)
		user.GET("/authority", handler.UserAuthorityHandler)
		user.POST("/search", handler.UserSearchHandler)
	}

	router.Use(handler.AuthAPIMiddleware())
	management := router.Group("/management")
	{
		management.GET("/superadmin", handler.SuperAdminGetHandler)
		management.POST("/superadmin", handler.SuperAdminPostHandler)
		management.DELETE("/superadmin/:id", handler.SuperAdminDeleteHandler)
		management.GET("/tenant", handler.TenantGetHandler)
		management.POST("/tenant", handler.TenantPostHandler)
		management.DELETE("/tenant/:id", handler.TenantDeleteHandler)
		management.GET("/frontend", handler.FrontendGetHandler)
		management.POST("/frontend", handler.FrontendPostHandler)
		management.DELETE("/frontend/:id", handler.FrontendDeleteHandler)
		management.GET("/backend", handler.BackendGetHandler)
		management.POST("/backend", handler.BackendPostHandler)
		management.DELETE("/backend/:id", handler.BackendDeleteHandler)
		management.GET("/feature", handler.FeatureGetHandler)
		management.POST("/feature", handler.FeaturePostHandler)
		management.DELETE("/feature/:id", handler.FeatureDeleteHandler)
		management.GET("/feature/sync", handler.FeatureSyncHandler)
		management.POST("/feature/relate", handler.FeatureRelateHandler)
	}

	organization := router.Group("/organization")
	{
		organization.GET("/info", handler.OrganizationInfoHandler)
		organization.GET("/user", handler.OrganizationUserGetHandler)
		organization.POST("/user", handler.OrganizationUserPostHandler)
		organization.GET("/role", handler.OrganizationRoleGetHandler)
		organization.POST("/role", handler.OrganizationRolePostHandler)
		organization.POST("/role/user", handler.OrganizationRoleUserHandler)
		organization.DELETE("/role/:id", handler.OrganizationRoleDeleteHandler)
		organization.GET("/object", handler.OrganizationObjectGetHandler)
		organization.POST("/object", handler.OrganizationObjectPostHandler)
		organization.DELETE("/object/:id", handler.OrganizationObjectDeleteHandler)
		organization.GET("/policy", handler.OrganizationPolicyGetHandler)
		organization.POST("/policy", handler.OrganizationPolicyPostHandler)
		organization.GET("/develop", handler.DeveloperAppGetHandler)
		organization.POST("develop", handler.DeveloperAppPostHandler)
		organization.DELETE("/develop/:id", handler.DeveloperAppDeleteHandler)
	}

	router.GET("/advertiser", handler.AdvertiserGetHandler)
	router.PATCH("/advertiser", handler.AdvertiserPatchHandler)
	router.DELETE("/advertiser/:id", handler.AdvertiserDeleteHandler)
	router.GET("/advertiser/authorize", handler.AdvertiserAuthorizeHandler)
	router.GET("/advertiser/callback", handler.AdvertiserCallbackHandler)

	schema := router.Group("/schema")
	{
		schema.GET("/", handler.SchemaGetHandler)
	}

	bindStrategy := router.Group("/bind_strategy")
	{
		bindStrategy.GET("/list", handler.StrategyListHandler)
		bindStrategy.POST("/edit", handler.StrategyEditHandler)
		bindStrategy.GET("/delete/:id", handler.StrategyDeleteHandler)
	}

	targeting := router.Group("/targeting")
	{
		targeting.GET("/list", handler.RTATargetingListHandler)
		targeting.POST("/edit", handler.RTATargetingEditHandler)
		targeting.GET("/delete/:id", handler.RTATargetingDeleteHandler)
	}

	experiment := router.Group("/experiment")
	experimentParameter := experiment.Group("/parameter")
	{
		experimentParameter.GET("/list", handler.ExperimentParameterListHandler)
		experimentParameter.POST("/edit", handler.ExperimentParameterEditHandler)
		experimentParameter.GET("/delete/:id", handler.ExperimentParameterDeleteHandler)
	}

	rtaAccount := experiment.Group("/account")
	{
		rtaAccount.GET("/list", handler.RtaAccountListHandler)
		rtaAccount.POST("/edit", handler.RtaAccountEditHandler)
		rtaAccount.GET("/delete/:id", handler.RtaAccountDeleteHandler)
		rtaAccount.POST("/sync", handler.RtaAccountSyncHandler)
	}

	rtaExp := experiment.Group("/rta_exp")
	{
		rtaExp.POST("/list", handler.RtaExpListHandler)
	}

	experimentGroup := experiment.Group("/group")
	{
		experimentGroup.GET("/list", handler.ExperimentGroupListHandler)
		experimentGroup.GET("/get/:id", handler.ExperimentGroupGetHandler)
		experimentGroup.POST("/edit", handler.ExperimentGroupEditHandler)
		experimentGroup.GET("/prompt/:id", handler.ExperimentGroupPromptHandler)
		experimentGroup.GET("/stop/:id", handler.ExperimentGroupStopHandler)
	}

	experimentReport := experiment.Group("/report")
	{
		experimentReport.POST("/get", handler.ExperimentReportGetHandler)
		experimentReport.GET("/get_attribution", handler.ExperimentReportGetAttributionHandler)
	}

	material := router.Group("/material")
	{
		material.POST("/edit", handler.MaterialPostHandler)
		material.GET("/list", handler.MaterialGetHandler)
		material.POST("/delete_many", handler.MaterialDeleteManyHandler)

		audit := material.Group("/audit")
		{
			audit.GET("/list", handler.MaterialAuditGetHandler)
			audit.POST("/submit", handler.MaterialAuditPostHandler)
		}

		file := material.Group("/file")
		{
			file.POST("/upload", handler.MaterialFileUploadHandler)
		}
	}

	return nil
}
