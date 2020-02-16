package models

import (
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"time"
)


//表的设计
type User struct {
	Id int
	Name string `orm:"unique"`
	Pwd string
	Article []*Article `orm:"rel(m2m)"`
}
//文章结构体
type Article struct {
	Id int `orm:"pk;auto"`
	ArticleName string `orm:"size(20)"`
	ArticleCreateTime time.Time `orm:"auto_now"`
	ArticleCount int `orm:"default(0);null"`
	ArticleContent string
	ArticleImg string
	ArticleForType *ArticleType `orm:"rel(fk);index"`
	ArticleType string
}
type ArticleType struct{
	ArticleTypeId int `orm:"pk;auto;column(article_type_id)"`
	ArticleTypeName string
	Articles []*Article `orm:"reverse(many)"`
}

func init(){
	// 设置数据库基本信息
	orm.RegisterDataBase("default", "mysql", "root:2013.a@tcp(127.0.0.1:3306)/blog?charset=utf8")
	// 映射model数据
	orm.RegisterModel(new(User),new(Article),new(ArticleType))
	// 生成表
	orm.RunSyncdb("default", false, true)
}

