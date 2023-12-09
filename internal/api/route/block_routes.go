package route

import (
	"github.com/gin-gonic/gin"
	"github.com/meta-node-blockchain/meta-node/cmd/meta-node-dns/internal/api/handler"
)

func SetupBlockRoutes(
	router *gin.RouterGroup,
	dnsHandler *handler.DnsHandler,
) {
	routes := router.Group("/dns")
	{
		routes.GET("/connection-address/:address", dnsHandler.GetConnectionAddress)
		routes.GET("/all-connection-address", dnsHandler.GetAllConnectionAddress)
		routes.POST("/connection-address", dnsHandler.CreateConnectionAddress)
	}
}
