package email

import (
	"cherry-web-api/models"
	"cherry-web-api/pkg/util"
	"encoding/json"
	"fmt"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	ses "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ses/v20201002"
	"strings"
)

func TencentSend(templateId uint64, secretId, secretKey, email, params string, tencentFrom, subject string) (int, string) {
	credential := common.NewCredential(
		secretId,
		secretKey,
	)
	// 实例化一个client选项，可选的，没有特殊需求可以跳过
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "ses.tencentcloudapi.com"
	// 实例化要请求产品的client对象,clientProfile是可选的
	client, _ := ses.NewClient(credential, "ap-hongkong", cpf)

	// 实例化一个请求对象,每个接口都会对应一个request对象
	request := ses.NewSendEmailRequest()
	request.FromEmailAddress = common.StringPtr(tencentFrom) // 发件人邮箱地址。不使用别名时请直接填写发件人邮箱地址
	request.Destination = common.StringPtrs([]string{email}) // 收件人邮箱地址，最多支持50人
	request.Subject = common.StringPtr(subject)              // 邮件主题
	request.Template = &ses.Template{                        // 模板相关信息
		TemplateID:   common.Uint64Ptr(templateId),
		TemplateData: common.StringPtr(params),
	}

	// 返回的resp是一个SendEmailResponse的实例，与请求对象对应
	response, err := client.SendEmail(request)

	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		msg := fmt.Sprintf("An API error has returned: %s", err)
		return -1, msg
	}
	if err != nil {
		//panic(err)
		return -1, "fail"
	}
	// 输出json格式的字符串回包
	msg := fmt.Sprintf("%s", response.ToJsonString())
	return 0, msg
}

func TencentSendEmail(email string, email_type int, params map[string]string, ip string) bool {
	tplSn := ""
	default_mail := models.GetConfigVal("default_email")
	confEmail := models.GetConfEmail(default_mail, email_type)
	tplSn = confEmail.TplId
	tplSnSubject := confEmail.Subject

	secretId := models.GetConfigVal("Tencent_SecretId")
	secretKey := models.GetConfigVal("Tencent_SecretKey")
	tencentFromConf := models.GetConfigVal("Tencent_From")
	tencentFrom := confEmail.FromTo

	if tencentFrom == "" {
		tencentFrom = tencentFromConf
	}
	if tplSn == "" || tplSnSubject == "" {
		return false
	}
	if secretId == "" || secretKey == "" || tencentFrom == "" {
		return false
	}

	templateID := uint64(util.StoI(tplSn)) // 125615 // 腾讯云模板ID，可前往腾讯云控制台查看
	if templateID == 0 {
		return false
	}

	paramsByte, _ := json.Marshal(params)
	paramsStr := string(paramsByte)

	sendRes, result := TencentSend(templateID, secretId, secretKey, email, paramsStr, tencentFrom, tplSnSubject)
	fmt.Println(sendRes)
	fmt.Println(result)
	models.AddLogEmail(email, paramsStr, tplSn, result, ip)

	if sendRes != 0 {
		return false
	}
	return true
}

func TencentSendEmailMarket(email string, code string, params map[string]string, ip string) bool {
	tplSn := ""
	default_mail := models.GetConfigVal("default_email")
	confEmail := models.GetConfEmailBy8(default_mail, code)
	tplSn = confEmail.TplId
	tplSnSubject := confEmail.Subject

	secretId := models.GetConfigVal("Tencent_SecretId")
	secretKey := models.GetConfigVal("Tencent_SecretKey")
	tencentFromConf := models.GetConfigVal("Tencent_From")
	tencentFrom := confEmail.FromTo

	if tencentFrom == "" {
		tencentFrom = tencentFromConf
	}
	if tplSn == "" || tplSnSubject == "" {
		fmt.Println("err ----1")
		return false
	}
	if secretId == "" || secretKey == "" || tencentFrom == "" {
		fmt.Println("err ----2")
		return false
	}

	templateID := uint64(util.StoI(tplSn)) // 125615 // 腾讯云模板ID，可前往腾讯云控制台查看
	if templateID == 0 {
		fmt.Println("err ----3")
		return false
	}

	if len(confEmail.PackageName) > 0 && len(confEmail.PackageUrl) > 0 {
		params["packageName"] = confEmail.PackageName
		params["packageUrl"] = strings.Replace(confEmail.PackageUrl, "https://www.cherryproxy.com/", "", -1)
	}
	paramsByte, _ := json.Marshal(params)
	paramsStr := string(paramsByte)

	sendRes, result := TencentSend(templateID, secretId, secretKey, email, paramsStr, tencentFrom, tplSnSubject)
	fmt.Println(sendRes)
	fmt.Println(result)
	models.AddLogEmail(email, paramsStr, tplSn, result, ip)

	if sendRes != 0 {
		return false
	}
	return true
}
