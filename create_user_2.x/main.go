package main

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"unsafe"
)

var (

	help bool
	url string
	token string
	enabled bool
	me bool
	mustChangePassword bool
	password string
	name string
	username string
	clusterid string
	roletemplateld string
	clusterType string
	globalRoleId string

)


func init() {
	flag.BoolVar(&help, "h", false, "显示帮助")
	flag.StringVar(&token,"t","","Rancher Token")
	flag.StringVar(&url,"u","","Rancher Server地址，不用加https")
	flag.BoolVar(&enabled,"e",true,"是否启用用户")
	flag.BoolVar(&me,"m",false,"用户是否只有登陆访问权限")
	flag.BoolVar(&mustChangePassword,"c",true,"是否要求用户在第一次登陆时修改密码")
	flag.StringVar(&name,"n","rancher","显示名称")
	flag.StringVar(&password,"p","rancher123","初始密码")
	flag.StringVar(&globalRoleId,"gr", "user","创建用户默认角色")
	flag.StringVar(&username,"us","rancher-user","用户名")
	flag.StringVar(&clusterid,"ci","","集群id")
	flag.StringVar(&roletemplateld,"r","","集群角色模版id")
	flag.StringVar(&clusterType,"clt","clusterRoleTemplateBinding","类型")

	flag.Parse()
}

func main() {

	GloUrl := "https://"+url+"/v3/user"
	sent := make(map[string]interface{})
	sent["enabled"] = enabled
	sent["me"] = me
	sent["mustChangePassword"] = mustChangePassword
	sent["name"] = name
	sent["password"] = password
	sent["username"] = username

	encodeString := base64.StdEncoding.EncodeToString([]byte(token))

	bytesData, err := json.Marshal(sent)
	if err != nil {
		fmt.Println("用户创建失败:"+err.Error() )
		return
	}
	reader := bytes.NewReader(bytesData)

	req, _ := http.NewRequest("POST", GloUrl, reader)

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Basic "+encodeString)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify : true},
	}
	client := &http.Client{Transport: tr}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("用户创建失败:"+err.Error())
		return
	}

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("用户创建失败:"+err.Error())
		return
	}
	type body struct {
		Id      string `json:"id"`
	}
	var basket body
	err = json.Unmarshal(respBytes, &basket)
	if err != nil {
		fmt.Println("用户创建失败:"+err.Error())
	}
	str := (*string)(unsafe.Pointer(&respBytes))
	resCode := resp.StatusCode
	if resCode == 201 {
		fmt.Println("用户创建成功，用户ID为："+basket.Id+"\n开始绑定用户到集群角色\n")
	}else{
		fmt.Println("用户创建失败")
		fmt.Println(*str)
		return
	}


	sentGloRole := make(map[string]interface{})
	sentGloRole["globalRoleId"] = globalRoleId
	sentGloRole["userId"] = basket.Id

	GloUrl = "https://" + url + "/v3/globalrolebindings"

	GloBytesData, err := json.Marshal(sentGloRole)
	if err != nil {
		fmt.Println("用户全局绑定角色失败:\n"+err.Error() )
		return
	}
	GloReader := bytes.NewReader(GloBytesData)

	GloReq, _ := http.NewRequest("POST", GloUrl, GloReader)

	GloReq.Header.Add("Content-Type", "application/json")
	GloReq.Header.Add("Authorization", "Basic "+encodeString)

	GloClient := &http.Client{Transport: tr}

	GloResp, err := GloClient.Do(GloReq)
	if err != nil {
		fmt.Println("用户全局绑定角色失败:\n"+err.Error())
		return
	}
	_, err = ioutil.ReadAll(GloResp.Body)
	if err != nil {
		fmt.Println("用户全局绑定角色失败:\n:"+err.Error())
		return
	}
	GloStr := (*string)(unsafe.Pointer(&respBytes))
	GloResCode := resp.StatusCode
	if GloResCode == 201 {
		fmt.Println("用户全局成功绑定角色："+globalRoleId)
	}else{
		fmt.Println("用户全局绑定角色失败")
		fmt.Println(*GloStr)
		return
	}
	
	sentRole := make(map[string]interface{})
	sentRole["clusterId"] = clusterid
	sentRole["roleTemplateId"] = roletemplateld
	sentRole["type"] = clusterType
	sentRole["userPrincipalId"] = "local://"+basket.Id

	CluUrl := "https://"+url+"/v3/clusterroletemplatebinding"

	CluBytesData, err := json.Marshal(sentRole)
	if err != nil {
		fmt.Println("用户集群角色绑定失败:\n"+err.Error() )
		return
	}
	CluReader := bytes.NewReader(CluBytesData)

	CluReq, _ := http.NewRequest("POST", CluUrl, CluReader)

	CluReq.Header.Add("Content-Type", "application/json")
	CluReq.Header.Add("Authorization", "Basic "+encodeString)

	CluClient := &http.Client{Transport: tr}

	CluResp, err := CluClient.Do(CluReq)
	if err != nil {
		fmt.Println("用户集群角色绑定失败:\n"+err.Error())
		return
	}
	CluRespBytes, err := ioutil.ReadAll(CluResp.Body)
	if err != nil {
		fmt.Println("用户集群角色绑定失败:\n:"+err.Error())
		return
	}
	CluStr := (*string)(unsafe.Pointer(&CluRespBytes))
	CluResCode := resp.StatusCode
	if CluResCode == 201 {
		fmt.Println("用户集群角色成功绑定角色："+roletemplateld)
	}else{
		fmt.Println("用户集群角色绑定失败")
		fmt.Println(*CluStr)
		return
	}
}

