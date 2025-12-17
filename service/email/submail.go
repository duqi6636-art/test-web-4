package email

// 赛邮
import (
	"cherry-web-api/models"
	"cherry-web-api/pkg/util"
	"fmt"
)

type SubEmailSendResult struct {
	Status string `json:"status"`
	Code   int    `json:"code"`
	Msg    string `json:"msg"`
	Json   string `json:"json"`
}

func SendEmailCode(address, codeType string, params models.SignParam, ip, area, useEmail string) (bool, string) {
	verifyCode := util.RandStr("n", 4)
	if codeType == "find" {
		verifyCode = util.RandStr("y", 6)
	}

	var (
		result bool
		data   string
	)

	if useEmail == "aws_mail" { //亚马逊
		vars := make(map[string]string)
		vars["email"] = address
		vars["code"] = verifyCode
		result = AwsSendEmail(address, 1, vars, ip)
		fmt.Println(result)
	}
	if useEmail == "tencent_mail" { //腾讯
		vars := make(map[string]string)
		vars["email"] = address
		vars["code"] = verifyCode
		result = TencentSendEmail(address, 1, vars, ip)
		fmt.Println(result)
	}

	if !result {
		return false, "__T_EMAIL_SEND_ERROR"
	}

	var r = false
	r = createVerifyCode(address, codeType, verifyCode, ip, data, params, 0)
	return r, "__T_EMAIL_SENDED"
}

// / 发送谷歌登录新注册用户邮箱
func SendGoogleLoginNewUserEmail(address, password, ip, area, useEmail string) {
	var (
		result bool
	)

	if useEmail == "aws_mail" { //亚马逊
		vars := make(map[string]string)
		vars["email"] = address
		vars["password"] = password
		result = AwsSendEmail(address, 9, vars, ip)
		fmt.Println(result)
	}

	if useEmail == "tencent_mail" { //亚马逊
		vars := make(map[string]string)
		vars["email"] = address
		vars["password"] = password
		result = TencentSendEmail(address, 9, vars, ip)
		fmt.Println(result)
	}

}
func createVerifyCode(account, codeType, verifyCode, ip, data string, params models.SignParam, uid int) bool {
	nowTime := util.GetNowInt()
	expTime := nowTime + 86400
	verMap := models.EmailCode{
		ID:         0,
		Type:       codeType,
		Enable:     "true",
		EmailCode:  verifyCode,
		ExpireTime: expTime,
		CreateTime: nowTime,
		Email:      account,
		Ip:         ip,
	}
	return models.AddVerifyCode(verMap)
}

// 注册发送邮件
func SendEmail(address, password, ip string) (bool, string) {
	var res = false
	var msg = ""
	useEmail := models.GetConfigVal("default_email")

	if useEmail == "aws_mail" { //亚马逊
		vars := make(map[string]string)
		vars["email"] = address
		vars["password"] = password
		res = AwsSendEmail(address, 3, vars, ip)
		fmt.Println(res)
	}

	if useEmail == "tencent_mail" { //亚马逊
		vars := make(map[string]string)
		vars["email"] = address
		vars["password"] = password
		res = TencentSendEmail(address, 3, vars, ip)
		fmt.Println(res)
	}
	return res, msg
}
