package routers

import (
	"class/controllers"
	"github.com/astaxie/beego"
)

func init() {
	//访问通过访问localhost:8081 直接访问的注册页面
    beego.Router("/", &controllers.MainController{})
	beego.Router("/abc", &controllers.MainController{})

    //访问最注册页面的 action
	beego.Router("/register", &controllers.MainController{})
	//注意：当实现了自定义的请求方法，请求将不会访问默认方法
	beego.Router("/login", &controllers.MainController{},"get:ShowLogin;post:HandleLogin")

	//跳转主页
	beego.Router("/showIndex",&controllers.MainController{},"get:ShowIndex")

    //跳转 添加页面
	beego.Router("/addArtcle",&controllers.MainController{},"get:ShowAddArticle;post:HandleAddArticle")

    //跳转到 文章详情页面
    beego.Router("/content",&controllers.MainController{},"get:ShowContent")

    //跳转编辑文章详情页面
    beego.Router("/update",&controllers.MainController{},"get:ShowUpdate;post:HandleUpdate")

    //删除 文章
    beego.Router("/delete",&controllers.MainController{},"get:HandleDelete")
}
