package handler

import (
	"fmt"
	"io/ioutil"
	dblayer "filestore-server/db"
	"filestore-server/util"
	//"io/ioutil"
	"net/http"
	"time"
)

const (
	pwd_salt="*#890"
)
func SignupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet{
		data, err := ioutil.ReadFile("./static/view/signup.html")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(data)
		return

	}
	r.ParseForm()
	username := r.Form.Get("username")
	passwd := r.Form.Get("password")
	if len(username)<3 || len(passwd)<5{
		w.Write([]byte("Invalid parameter"))
		return
	}
	enc_passwd := util.Sha1([]byte(passwd+pwd_salt))
	succ := dblayer.UserSignUp(username, enc_passwd)
	if succ{
		w.Write([]byte("Success"))

	}else{
		w.Write([]byte("failed"))
	}
}
//登录接口
func SignInHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet{
		http.Redirect(w, r, "/static/view/signin.html", http.StatusFound)
		return
	}
	r.ParseForm()
	username := r.Form.Get("username")
	password := r.Form.Get("password")
	encPasswd := util.Sha1([]byte(password+pwd_salt))
	//1、校验用户名以及密码
	pwdChecked := dblayer.UserLoginIn(username, encPasswd)
	if !pwdChecked {
		w.Write([]byte("FAILED"))
	}
	//2、生成访问凭证（token）
	token := GetToken(username)
	upRes := dblayer.UpdateToken(username, token)
	if !upRes{
		w.Write([]byte("FAILED"))
		return
	}
	//3、登录成功后重定向到首页
	resp := util.RespMsg{
		Code: 0,
		Msg:  "OK",
		Data: struct {
			Location string
			Username string
			Token    string
		}{
			Location: "http://" + r.Host + "/static/view/home.html",
			Username: username,
			Token:    token,
		},
	}
	w.Write(resp.JSONBytes())
}

func GetToken(username string) string{
	//md5(username+timestamp+token_salt)+timestamp[:8]
	ts := fmt.Sprintf("%x", time.Now().Unix())
	tokenPrefix := util.MD5([]byte(username+ts+"_tokensalt"))
	return tokenPrefix + ts[:8]

}

func UserInfoHandler(w http.ResponseWriter, r *http.Request) {
	//解析请求参数
	r.ParseForm()
	username := r.Form.Get("username")
	//token := r.Form.Get("token")
	//isValidToken := CheckTokenIsInvalid(username, token)
	//验证token是否过期
	//if isValidToken == false {
	//	http.Redirect(w, r, "/static/view/signin.html", http.StatusFound)
	//	return
	//}
	//查询用户信息
	user, err := dblayer.GetUserInfo(username)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	//组装并相应用户数据
	resp := util.RespMsg{
		Code:0,
		Msg:"OK",
		Data:user,
	}
	w.Write(resp.JSONBytes())

}

func CheckTokenIsInvalid(username, token string) bool {
	//校验token长度是否正确
	if len(token) != 40 {
		return false
	}
	//校验token时间戳是否过期
	tokenTs := token[32:40]
	if util.Hex2Dec(tokenTs) < time.Now().Unix() - 86400  {
		fmt.Println("token in expired:" + token)
		return false
	}
	//从db获取用户token
	//对比token是否一致
	if dblayer.GetUserToken(username) != token {
		return false
	}
	return true
}