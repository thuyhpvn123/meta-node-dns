package handler

import (
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	"github.com/meta-node-blockchain/meta-node/cmd/meta-node-dns/internal/config"
	"github.com/meta-node-blockchain/meta-node/pkg/logger"
	"github.com/meta-node-blockchain/meta-node/pkg/storage"
)

type DnsHandler struct {
	dnsStorage storage.Storage
	config     *config.AppConfig
}

func NewDnsHandler(
	dnsStorage storage.Storage,
	config *config.AppConfig,
) *DnsHandler {
	return &DnsHandler{
		dnsStorage,
		config,
	}
}

func (h *DnsHandler) GetConnectionAddress(c *gin.Context) {
	address := c.Param("address")
	connectionAddress, _ := h.dnsStorage.Get(common.FromHex(address))
	rs := map[string]string{
		address: string(connectionAddress),
	}
	logger.DebugP(rs)
	c.JSON(http.StatusOK, rs)
}

func (h *DnsHandler) GetAllConnectionAddress(c *gin.Context) {
	rs := map[string]string{}
	iter := h.dnsStorage.GetIterator()
	for iter.Next() {
		rs[hex.EncodeToString(iter.Key())] = string(iter.Value())
	}
	logger.DebugP(rs)
	c.JSON(http.StatusOK, rs)
}

type CreateConnectionAddressReq struct {
	Address           string `json:"address"`
	ConnectionAddress string `json:"connection_address"`
	Secret            string `json:"secret"`
}

func (h *DnsHandler) CreateConnectionAddress(c *gin.Context) {
	req := &CreateConnectionAddressReq{}
	err := c.ShouldBindJSON(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Unable to bind body %v", err),
		})
		return
	}
	if req.Secret != h.config.Secret {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid secret",
		})
		return
	}

	h.dnsStorage.Put(common.FromHex(req.Address), []byte(req.ConnectionAddress))
	c.JSON(http.StatusOK, gin.H{
		"message":            "successful request",
		"address":            req.Address,
		"connection-address": req.ConnectionAddress,
	})
}
