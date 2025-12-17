package proxyManageController

import (
	"cherry-web-api/controller"
	"cherry-web-api/e"
	"cherry-web-api/models"
	"cherry-web-api/pkg/util"
	"fmt"
	"github.com/gin-gonic/gin"
)

func SetConfig(c *gin.Context) {
	resCode, msg, user := controller.DealUser(c) //处理用户信息
	if resCode != e.SUCCESS {
		controller.JsonReturn(c, resCode, msg, nil)
		return
	}
	defaultConfig := c.PostForm("default")
	columns := c.PostForm("columns")
	proxyLists := c.PostForm("proxy_lists")
	fmt.Println(user.Id)
	fmt.Println(defaultConfig, columns, proxyLists)
	browserType := c.PostForm("browser_type")
	uid := user.Id
	config, _ := models.GetProxyConfig(uid)
	if config.Id > 0 {
		models.UpdateProxyConfig(uid, models.ProxyConfig{Default: defaultConfig, Columns: columns, ProxyLists: proxyLists, BrowserType: browserType})
	} else {
		models.AddProxyConfig(models.ProxyConfig{Default: defaultConfig, Columns: columns, ProxyLists: proxyLists, BrowserType: browserType, Uid: uid, CreatedAt: util.GetNowInt()})
	}
	controller.JsonReturn(c, e.SUCCESS, "success", nil)
	return
}

func GetConfig(c *gin.Context) {
	resCode, msg, user := controller.DealUser(c) //处理用户信息
	if resCode != e.SUCCESS {
		controller.JsonReturn(c, resCode, msg, nil)
		return
	}
	uid := user.Id
	config, _ := models.GetProxyConfig(uid)
	if config.Id == 0 {
		controller.JsonReturn(c, -1, "not found config", nil)
		return
	}
	controller.JsonReturn(c, e.SUCCESS, "success", config)
	return
}
