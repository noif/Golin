package web

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"golin/global"
	"golin/run"
	"html/template"
	"net/http"
	"path/filepath"
)

var save bool

func Start(cmd *cobra.Command, args []string) {
	ip, _ := cmd.Flags().GetString("ip")
	port, _ := cmd.Flags().GetString("port")
	save, _ = cmd.Flags().GetBool("save")
	r := gin.Default()
	r.SetHTMLTemplate(IndexTemplate()) //加载首页
	golin := r.Group("/golin")
	{
		golin.GET("/index", GolinIndex)    //首页
		golin.POST("/submit", GolinSubmit) //提交任务
	}
	r.Run(ip + ":" + port)
}

// GolinIndex 首页
func GolinIndex(c *gin.Context) {
	c.HTML(http.StatusOK, "index", "")
}

// runtypepath 不同类型对应的采集完成目录的文件夹，拼接目录用
var runtypepath = map[string]string{"Linux": "Linux", "Mysql": "MySQL", "Redis": "Redis"}

// GolinSubmit 提交任务
func GolinSubmit(c *gin.Context) {
	name, ip, user, passwd, port, mode := c.PostForm("name"), c.PostForm("ip"), c.PostForm("user"), c.PostForm("password"), c.PostForm("port"), c.PostForm("run_mode")
	//fmt.Println(name, ip, user, passwd, port, mode)
	run.Onlyonerun(fmt.Sprintf("%s~~~%s~~~%s~~~%s~~~%s", name, ip, user, passwd, port), "~~~", mode)
	successfile := filepath.Join(global.Succpath, runtypepath[mode], name+"_"+ip+".log")
	if global.PathExists(successfile) {
		filename := fmt.Sprintf("%s_%s(%s).log)", name, ip, mode)
		c.Header("Content-Description", "File Transfer")
		c.Header("Content-Disposition", "attachment; filename="+filename)
		c.Header("Content-Type", "application/octet-stream")
		//返回文件
		c.File(successfile)
		//此处暂时不生效，待解决！
		//if !save {
		//	os.Remove(successfile)
		//}
	} else {
		c.String(200, "失败了哦客官～")
	}
}

// IndexTemplate 返回包含模板内容的模板结构体
func IndexTemplate() *template.Template {
	tmpl, err := template.New("index").Parse(IndexHtml())
	if err != nil {
		panic(err)
	}
	return tmpl
}
