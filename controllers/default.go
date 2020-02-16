package controllers

import (
	"class/models"
	"encoding/json"
	"github.com/astaxie/beego"
	"math"
	"path"
	"strconv"
	"time"

	//"github.com/astaxie/beego/orm"
	//"class/models"
	"github.com/astaxie/beego/orm"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) Get() {
	c.TplName = "register.html"
}
func (c *MainController) Post() {
	//1.拿到数据
	userName := c.GetString("userName")
	pwd := c.GetString("pwd")
	//2.对数据进行校验
	if userName == "" || pwd == "" {
		beego.Info("数据不能为空")
		c.Redirect("/register", 302)
		return
	}
	//3.插入数据库
	o := orm.NewOrm()

	user := models.User{}
	user.Name = userName
	user.Pwd = pwd
	_, err := o.Insert(&user)
	if err != nil {
		beego.Info("插入数据失败")
		c.Redirect("/register", 302)
		return
	}
	//4.返回登陆界面
	c.Redirect("/login", 302)
}

func (c *MainController) ShowLogin() {

	c.TplName = "login.html"
}

//登陆业务处理
func (c *MainController) HandleLogin() {
	//c.Ctx.WriteString("这是登陆的POST请求")
	//1.拿到数据
	userName := c.GetString("userName")
	pwd := c.GetString("password")
	//2.判断数据是否合法
	if userName == "" || pwd == "" {
		beego.Info("输入数据不合法")
		c.TplName = "login.html"
		return
	}
	//3.查询账号密码是否正确
	o := orm.NewOrm()
	user := models.User{}

	user.Name = userName
	err := o.Read(&user, "Name")
	if err != nil {
		beego.Info("查询失败")
		c.TplName = "login.html"
		return
	}
	//4.跳转
	c.Redirect("/index", 302)
}

//显示主页添加界面
func (c *MainController) ShowIndex() {
	//查询所有的文章
	o := orm.NewOrm()
	var articles []models.Article
	querySeter := o.QueryTable("Article")
	pageIndex := 1
	var errGetParam error
	pageIndex, errGetParam = c.GetInt("pageIndex")
	if errGetParam != nil {
		pageIndex = 1
	}
	pageSize := 2
	pageCount, _ := querySeter.Count()
	page := float64(pageCount) / float64(pageSize)
	//向上取整
	pageSum := math.Ceil(page)
	start := pageSize * (pageIndex - 1)
	querySeter.Limit(pageSize, start).All(&articles)
	//beego.Info(articles)
	var firstPage, endPage bool
	if pageIndex == 1 {
		firstPage = true
	}
	if float64(pageIndex) == pageSum {
		endPage = true
	}
	c.Data["firstPage"] = firstPage
	c.Data["endPage"] = endPage
	c.Data["pageCount"] = pageCount
	c.Data["pageSum"] = pageSum
	c.Data["pageIndex"] = pageIndex
	c.Data["articles"] = articles
	c.TplName = "index.html"
}

//显示文章添加界面
func (c *MainController) ShowAddArticle() {
	c.TplName = "add.html"
}

func (c *MainController) HandleAddArticle() {
	//1、获取数据信息
	var art models.Article;
	artName := c.GetString("articleName")
	if len(artName) > 20 {
		beego.Info("文章名称太长")
		c.Redirect("addArtcle", 302)
		return
	}
	art.ArticleName = artName
	art.ArticleType = c.GetString("articleType")
	art.ArticleContent = c.GetString("articleContent")
	_, head, _ := c.GetFile("uploadname")
	//获取文件的名称
	fileName := head.Filename
	//获取文件的后缀
	fileNameHouzhui := path.Ext(fileName)
	beego.Info("filename =%v,fileNameHouzhui =%v", fileName, fileNameHouzhui)
	//文件格式判断
	if fileNameHouzhui != ".jpg" && fileNameHouzhui != ".png" {
		beego.Info("文件格式错误")
		c.Redirect("addArtcle", 302)
		return
	}
	//文件大小限制，不能大于3M
	if head.Size > 1024*1024*3 {
		beego.Info("文件超出大小限制")
		c.Redirect("addArtcle", 302)
		return
	}
	nowTime := time.Now()
	uploadfileName := "./static/img/" + nowTime.Format("2006-01-02-15-04-05--") + strconv.Itoa(nowTime.Nanosecond()) + fileNameHouzhui;
	err := c.SaveToFile("uploadname", uploadfileName)
	if err != nil {
		beego.Info("上传失败")
		c.Redirect("addArtcle", 302)
		return
	}
	art.ArticleCreateTime = time.Now()
	art.ArticleImg = uploadfileName
	data, _ := json.Marshal(art)
	beego.Info("insert article success , art =%v", string(data))
	//2、数据库对象的
	o := orm.NewOrm()
	row, err := o.Insert(&art)
	if row == 0 && err != nil {
		beego.Info("insert article fail err %v", err)
		return
	} else {
		beego.Info("insert article success , art =%v", string(data))
	}
	c.Redirect("index", 302)
}

func (c *MainController) ShowContent() {
	id, err := c.GetInt("id")
	if err != nil {
		beego.Info("参数获取失败")
		return
	}
	o := orm.NewOrm()
	var art models.Article
	art.Id = id
	o.Read(&art, "Id")
	c.Data["art"] = art
	beego.Info(art)
	c.TplName = "content.html"
}

func (c *MainController) ShowUpdate() {
	id, err := c.GetInt("id")
	if err != nil {
		beego.Info("参数获取失败 =%v", err)
		return
	}
	art := models.Article{Id: id}
	o := orm.NewOrm()
	o.Read(&art)
	c.Data["art"] = art
	c.TplName = "update.html"
}
func (c *MainController) HandleUpdate() {
	id, errP := c.GetInt("id")
	if errP != nil {
		beego.Info("获取参数失败 %v", errP)
		c.Redirect("update", 302)
		return
	}
	//1、获取数据信息
	var art models.Article;
	artName := c.GetString("articleName")
	if len(artName) > 20 {
		beego.Info("文章名称太长")
		c.Redirect("update", 302)
		return
	}
	art.ArticleName = artName
	art.ArticleContent = c.GetString("articleContent")
	_, head, _ := c.GetFile("uploadname")
	//获取文件的名称
	fileName := head.Filename
	//获取文件的后缀
	fileNameHouzhui := path.Ext(fileName)
	beego.Info("filename =%v,fileNameHouzhui =%v", fileName, fileNameHouzhui)
	//文件格式判断
	if fileNameHouzhui != ".jpg" && fileNameHouzhui != ".png" {
		beego.Info("文件格式错误")
		c.Redirect("update", 302)
		return
	}
	//文件大小限制，不能大于3M
	if head.Size > 1024*1024*3 {
		beego.Info("文件超出大小限制")
		c.Redirect("update", 302)
		return
	}
	nowTime := time.Now()
	uploadfileName := "./static/img/" + nowTime.Format("2006-01-02-15-04-05--") + strconv.Itoa(nowTime.Nanosecond()) + fileNameHouzhui;
	err := c.SaveToFile("uploadname", uploadfileName)
	if err != nil {
		beego.Info("上传失败")
		c.Redirect("addArtcle", 302)
		return
	}
	art.ArticleImg = uploadfileName
	art.Id = id
	o := orm.NewOrm()
	o.Update(&art, "Id", "ArticleImg", "ArticleName", "ArticleContent")

	c.Redirect("index", 302)
}

func (c *MainController) HandleDelete() {
	id, err := c.GetInt("id")
	if err != nil {
		beego.Info("参数获取失败 %v", err)
		return
	}
	o := orm.NewOrm()
	var art models.Article
	art.Id = id
	row, err := o.Delete(&art, "Id")
	if err != nil {
		beego.Info("err", err)
		c.Redirect("index", 302)
		return
	}
	if row > 0 {
		beego.Info("数据删除成功")
	} else {
		beego.Info("数据删除失败")
	}
	c.Redirect("index", 302)
}
