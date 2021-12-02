package models

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/beego/beego/v2/client/httplib"
	"github.com/beego/beego/v2/core/logs"
	"github.com/buger/jsonparser"
)

type UserInfoResult struct {
	Data struct {
		JdVvipCocoonInfo struct {
			JdVvipCocoon struct {
				DisplayType   int    `json:"displayType"`
				HitTypeList   []int  `json:"hitTypeList"`
				Link          string `json:"link"`
				Price         string `json:"price"`
				Qualification int    `json:"qualification"`
				SellingPoints string `json:"sellingPoints"`
			} `json:"JdVvipCocoon"`
			JdVvipCocoonStatus string `json:"JdVvipCocoonStatus"`
		} `json:"JdVvipCocoonInfo"`
		JdVvipInfo struct {
			JdVvipStatus string `json:"jdVvipStatus"`
		} `json:"JdVvipInfo"`
		AssetInfo struct {
			AccountBalance string `json:"accountBalance"`
			BaitiaoInfo    struct {
				AvailableLimit     string `json:"availableLimit"`
				BaiTiaoStatus      string `json:"baiTiaoStatus"`
				Bill               string `json:"bill"`
				BillOverStatus     string `json:"billOverStatus"`
				Outstanding7Amount string `json:"outstanding7Amount"`
				OverDueAmount      string `json:"overDueAmount"`
				OverDueCount       string `json:"overDueCount"`
				UnpaidForAll       string `json:"unpaidForAll"`
				UnpaidForMonth     string `json:"unpaidForMonth"`
			} `json:"baitiaoInfo"`
			BeanNum    string `json:"beanNum"`
			CouponNum  string `json:"couponNum"`
			CouponRed  string `json:"couponRed"`
			RedBalance string `json:"redBalance"`
		} `json:"assetInfo"`
		FavInfo struct {
			FavDpNum    string `json:"favDpNum"`
			FavGoodsNum string `json:"favGoodsNum"`
			FavShopNum  string `json:"favShopNum"`
			FootNum     string `json:"footNum"`
			IsGoodsRed  string `json:"isGoodsRed"`
			IsShopRed   string `json:"isShopRed"`
		} `json:"favInfo"`
		GrowHelperCoupon struct {
			AddDays     int     `json:"addDays"`
			BatchID     int     `json:"batchId"`
			CouponKind  int     `json:"couponKind"`
			CouponModel int     `json:"couponModel"`
			CouponStyle int     `json:"couponStyle"`
			CouponType  int     `json:"couponType"`
			Discount    float64 `json:"discount"`
			LimitType   int     `json:"limitType"`
			MsgType     int     `json:"msgType"`
			Quota       float64 `json:"quota"`
			RoleID      int     `json:"roleId"`
			State       int     `json:"state"`
			Status      int     `json:"status"`
		} `json:"growHelperCoupon"`
		KplInfo struct {
			KplInfoStatus string `json:"kplInfoStatus"`
			Mopenbp17     string `json:"mopenbp17"`
			Mopenbp22     string `json:"mopenbp22"`
		} `json:"kplInfo"`
		OrderInfo struct {
			CommentCount     string        `json:"commentCount"`
			Logistics        []interface{} `json:"logistics"`
			OrderCountStatus string        `json:"orderCountStatus"`
			ReceiveCount     string        `json:"receiveCount"`
			WaitPayCount     string        `json:"waitPayCount"`
		} `json:"orderInfo"`
		PlusPromotion struct {
			Status int `json:"status"`
		} `json:"plusPromotion"`
		UserInfo struct {
			BaseInfo struct {
				AccountType    string `json:"accountType"`
				BaseInfoStatus string `json:"baseInfoStatus"`
				CurPin         string `json:"curPin"`
				DefinePin      string `json:"definePin"`
				HeadImageURL   string `json:"headImageUrl"`
				LevelName      string `json:"levelName"`
				Nickname       string `json:"nickname"`
				Pinlist        string `json:"pinlist"`
				UserLevel      string `json:"userLevel"`
			} `json:"baseInfo"`
			IsHideNavi     string `json:"isHideNavi"`
			IsHomeWhite    string `json:"isHomeWhite"`
			IsJTH          string `json:"isJTH"`
			IsKaiPu        string `json:"isKaiPu"`
			IsPlusVip      string `json:"isPlusVip"`
			IsQQFans       string `json:"isQQFans"`
			IsRealNameAuth string `json:"isRealNameAuth"`
			IsWxFans       string `json:"isWxFans"`
			Jvalue         string `json:"jvalue"`
			OrderFlag      string `json:"orderFlag"`
			PlusInfo       struct {
			} `json:"plusInfo"`
			XbScore string `json:"xbScore"`
		} `json:"userInfo"`
		UserLifeCycle struct {
			IdentityID      string `json:"identityId"`
			LifeCycleStatus string `json:"lifeCycleStatus"`
			TrackID         string `json:"trackId"`
		} `json:"userLifeCycle"`
	} `json:"data"`
	Msg       string `json:"msg"`
	Retcode   string `json:"retcode"`
	Timestamp int64  `json:"timestamp"`
}

func initCookie() {
	cks := GetJdCookies()
	for i := range cks {
		time.Sleep(time.Second * time.Duration(Config.Later))
		if cks[i].Available == True && !CookieOK(&cks[i]) {
			logs.Info("开始禁用")
			cks[i].OutPool()
		}
	}
	//for i := 0; i < l-1; i++ {
	//	if cks[i].Available == True && !CookieOK(&cks[i]) {
	//		if pt_key, err := cks[i].OutPool(); err == nil && pt_key != "" {
	//			i = i - 1
	//			logs.Info("正常操作")
	//			logs.Info(cks[i].PtPin)
	//			logs.Info(i)
	//		}
	//	}
	//}
	go func() {
		Save <- &JdCookie{}
	}()
}

func cleanCookie() {
	cks := GetJdCookies()
	(&JdCookie{}).Push("开始清理过期账号")
	xx := 0
	for i := range cks {
		if cks[i].Available == False {
			xx++
			cks[i].Removes(cks[i])
		}
	}
	(&JdCookie{}).Push(fmt.Sprintf("所有CK清理，共%d个", xx))
}

func cleanWck() {
	cks := GetJdCookies()
	xx := 0
	(&JdCookie{}).Push("开始清空Wskey")
	for i := range cks {
		if len(cks[i].WsKey) > 0 {
			ck := cks[i]
			ck.Update(WsKey, "")
			xx++
		}
	}
	(&JdCookie{}).Push(fmt.Sprintf("已清理WCK，一共%d", xx))
}

func updateCookie() {
	cks := GetJdCookies()
	l := len(cks)
	logs.Info(l)
	xx := 0
	yy := 0
	(&JdCookie{}).Push("开始定时更新转换Wskey")
	for i := range cks {
		if len(cks[i].WsKey) > 0 {
			time.Sleep(10 * time.Second)
			ck := cks[i]
			//JdCookie{}.Push(fmt.Sprintf("更新账号账号，%s", ck.Nickname))
			var pinky = fmt.Sprintf("pin=%s;wskey=%s;", ck.PtPin, ck.WsKey)
			rsp, err := getKey(pinky)
			if err != nil {
				logs.Error(err)
			}
			if strings.Contains(rsp, "fake") {
				ck.Push(fmt.Sprintf("Wskey失效账号，%s", ck.PtPin))
				(&JdCookie{}).Push(fmt.Sprintf("Wskey失效，%s", ck.PtPin))
			} else {
				ptKey := FetchJdCookieValue("pt_key", rsp)
				ptPin := FetchJdCookieValue("pt_pin", rsp)
				ck := JdCookie{
					PtKey: ptKey,
					PtPin: ptPin,
				}
				if ptPin != "" || ptKey != "" {
					if nck, err := GetJdCookie(ck.PtPin); err == nil {
						xx++
						nck.InPool(ck.PtKey)
						nck.Update(Available, True)
						//msg := fmt.Sprintf("定时更新账号，%s", ck.PtPin)
						////不再发送成功提醒
						//(&JdCookie{}).Push(msg)
						//logs.Info(msg)
					} else {
						yy++
						ck.Update(Available, False)
						(&JdCookie{}).Push(fmt.Sprintf("查无匹配得ptpin，%s", ck.PtPin))
					}
				}
				go func() {
					Save <- &JdCookie{}
				}()
			}
		}
	}
	(&JdCookie{}).Push(fmt.Sprintf("所有CK转换完成，共%d个,转换失败个数共%d个", xx, yy))
}

func CookieOK(ck *JdCookie) bool {
	cookie := "pt_key=" + ck.PtKey + ";pt_pin=" + ck.PtPin + ";"
	// fmt.Println(cookie)
	// jdzz(cookie, make(chan int64))
	if ck == nil {
		return true
	}
	uri, err := url.Parse("http://118.190.244.234:3128/")

	if err != nil {
		log.Fatal("parse url error: ", err)
	}
	log.Println(uri.User)

	client := http.Client{
		Transport: &http.Transport{
			// 设置代理
			Proxy: http.ProxyURL(uri),
		},
	}
	//模拟ios客户端发送get请求
	// client := new(http.Client)
	reg, err := http.NewRequest("GET", `https://me-api.jd.com/user_new/info/GetJDUserInfoUnion`, nil)
	// reg, err := http.NewRequest("GET", `https://www.cip.cc/`, nil)
	if err != nil {
		fmt.Println("Error1:", err)
		return true
	}
	reg.Header.Add(`HTTP`, `2.0`)
	reg.Header.Add(`Accept`, `*/*`)
	reg.Header.Add(`Accept-Language`, `zh-cn`)
	reg.Header.Add(`User-Agent`, `AppStore/2.0 iOS/7.1.2 model/iPod5,1 build/11D257 (4; dt:81)`)
	reg.Header.Add(`Host`, `itunes.apple.com`)
	reg.Header.Add(`Connection`, `keep-alive`)
	reg.Header.Add(`X-Apple-Store-Front`, `143465-19,21 t:native`)
	reg.Header.Add(`X-Dsid`, `932530590`)
	reg.Header.Add("Cookie", cookie)
	// reg.Header.Set("Accept", "*/*")
	// reg.Header.Set("Accept-Language", "zh-cn,")
	// reg.Header.Set("Connection", "keep-alive,")
	// reg.Header.Set("Referer", "https://home.m.jd.com/myJd/newhome.action?sceneval=2&ufc=&")
	// reg.Header.Set("Host", "me-api.jd.com")
	// reg.Header.Set("User-Agent", "jdapp;iPhone;9.4.4;14.3;network/4g;Mozilla/5.0 (iPhone; CPU iPhone OS 14_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148;supportJDSHWK/1")
	//reg.Header.Set(`Cookie`, `xp_ci=3z1E7umazD0Dz5SBzCwNzB7weVKgD; s_vi=[CS]v1|29F61F868501299A-60000114E000452B[CE]; Pod=20; itspod=20; xt-src=b; xt-b-ts-932530590=1408324780262; mz_at_ssl-932530590=AwUAAAFRAAER1gAAAABT8vlrAo2EAZQvwAJjChIlGtIxIKYErLQ=; mz_at0-932530590=AwQAAAFRAAER1gAAAABT8VSrdHM0dXgdzosavj4+sT0AJfhYBx4=; wosid-lite=qQmZVeBH9vj91TakAeKEZg; ns-mzf-inst=35-163-80-118-68-8171-202429-20-nk11; X-Dsid=932530590`)

	//执行get请求
	resp, err := client.Do(reg)
	if err != nil {
		return true
	}
	defer resp.Body.Close()
	if err != nil {
		fmt.Println("Error2:", err.Error())
		os.Exit(1)
		return true
	}
	//分隔符
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return true
	}
	// req := httplib.Get("https://me-api.jd.com/user_new/info/GetJDUserInfoUnion")
	// req.Header("Cookie", cookie)
	// req.Header("Accept", "*/*")
	// req.Header("Accept-Language", "zh-cn,")
	// req.Header("Connection", "keep-alive,")
	// req.Header("Referer", "https://home.m.jd.com/myJd/newhome.action?sceneval=2&ufc=&")
	// req.Header("Host", "me-api.jd.com")
	// req.Header("User-Agent", "jdapp;iPhone;9.4.4;14.3;network/4g;Mozilla/5.0 (iPhone; CPU iPhone OS 14_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148;supportJDSHWK/1")
	// data, err := req.Bytes()
	if err != nil {
		return true
	}
	ui := &UserInfoResult{}
	if nil != json.Unmarshal(data, ui) {
		//if !Config.IFC {
		//	(&JdCookie{}).Push("第一个接口失效，切换到第二个接口，可能黑IP，会导致NickName获取失败，可能会自行恢复。")
		//	Config.IFC = true
		//}

		return av2(ck)
	}
	//if Config.IFC {
	//	(&JdCookie{}).Push("第一个接口恢复，切换回第一接口，恭喜你IP洗白白了")
	//	Config.IFC = false
	//}
	switch ui.Retcode {
	//case "1001": //ck.BeanNum
	//	if ui.Msg == "not login" {
	//		return false
	//	}
	case "0":
		if url.QueryEscape(ui.Data.UserInfo.BaseInfo.CurPin) != ck.PtPin {
			return av2(ck)
		}
		if ui.Data.UserInfo.BaseInfo.Nickname != ck.Nickname || ui.Data.AssetInfo.BeanNum != ck.BeanNum || ui.Data.UserInfo.BaseInfo.UserLevel != ck.UserLevel || ui.Data.UserInfo.BaseInfo.LevelName != ck.LevelName {
			ck.Updates(JdCookie{
				Nickname:  ui.Data.UserInfo.BaseInfo.Nickname,
				BeanNum:   ui.Data.AssetInfo.BeanNum,
				Available: True,
				UserLevel: ui.Data.UserInfo.BaseInfo.UserLevel,
				LevelName: ui.Data.UserInfo.BaseInfo.LevelName,
			})
			ck.UserLevel = ui.Data.UserInfo.BaseInfo.UserLevel
			ck.LevelName = ui.Data.UserInfo.BaseInfo.LevelName
			ck.Nickname = ui.Data.UserInfo.BaseInfo.Nickname
			ck.BeanNum = ui.Data.AssetInfo.BeanNum
		}
		return true
	}
	//(&JdCookie{}).Push("第一个接口失效，切换到第二个接口，可能黑IP")
	return av2(ck)
}

func av2(ck *JdCookie) bool {
	cookie := "pt_key=" + ck.PtKey + ";pt_pin=" + ck.PtPin + ";"
	req := httplib.Get(`https://m.jingxi.com/user/info/GetJDUserBaseInfo?_=1629334995401&sceneval=2&g_login_type=1&g_ty=ls`)
	req.Header("User-Agent", ua)
	req.Header("Host", "m.jingxi.com")
	req.Header("Accept", "*/*")
	req.Header("Connection", "keep-alive")
	req.Header("Accept-Language", "zh-cn")
	req.Header("Accept-Encoding", "gzip, deflate, br")
	req.Header("Referer", "https://st.jingxi.com/my/userinfo.html?&ptag=7205.12.4")
	req.Header("Cookie", cookie)
	data, err := req.Bytes()
	if err != nil {
		return true
	}
	if ck.Nickname == "" {
		ck.Nickname, _ = jsonparser.GetString(data, "nickname")
		ck.Update("Nickname", ck.Nickname)
		logs.Info("开始补齐NickName")
	}
	return !strings.Contains(string(data), "login")
}
