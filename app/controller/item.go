package controller

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
	"strings"

	"../model"
	//	"github.com/gin-gonic/contrib/static"
	csrf "github.com/utrack/gin-csrf"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/olahol/go-imageupload"
	//	"github.com/mattn/go-colorable"
)

// Article is controller for requests to articles.
type Item struct {
	DB *sql.DB
}

func (t *Item) Root(c *gin.Context) {
	itemJoinUserNew, err := model.NewItemsAll(t.DB)
	if err != nil {
		c.String(500, "%s", err)
		return
	}
	itemJoinUserRank, err := model.RankItemsAll(t.DB)
	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"title":     "ShoesBox",
		"items":     itemJoinUserNew,
		"rankItems": itemJoinUserRank,
		"context":   c,
	})
}

func (t *Item) ListGet(c *gin.Context) {
	itemJoinUser, err := model.ItemsAll(t.DB)
	if err != nil {
		c.String(500, "%s", err)
		return
	}
	c.HTML(http.StatusOK, "list.tmpl", gin.H{
		"title":   "ShoesBox",
		"items":   itemJoinUser,
		"csrf":    csrf.GetToken(c),
		"context": c,
	})
}

func (t *Item) SavePoint(c *gin.Context) {
	point := c.PostForm("point")
	point1, err := strconv.ParseInt(point, 10, 64)

	if err != nil {
		c.String(500, "%s", err)
		return
	}
	id := c.PostForm("id")
	aid, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.String(500, "%s", err)
		return
	}
	var item model.Item
	item.Point = point1 + 1
	item.ID = aid
	var m *model.Item
	m = &item
	TXHandler(c, t.DB, func(tx *sql.Tx) error {
		if _, err := m.PointUpdate(tx); err != nil {
			return err
		}
		return tx.Commit()
	})
	c.Redirect(301, fmt.Sprintf("/list"))
}

func (t *Item) PointUpdate(c *gin.Context, m *model.Item) {
	log.Printf("up:%v", m.ID)
	TXHandler(c, t.DB, func(tx *sql.Tx) error {
		if _, err := m.PointUpdate(tx); err != nil {
			return err
		}
		return tx.Commit()
	})
	c.Redirect(301, fmt.Sprintf("/list"))
}

func (t *Item) MyPageGet(c *gin.Context) {
	sess := sessions.Default(c)
	uid, ok := sess.Get("uid").(int64)
	if !ok {
		c.String(500, "")
		return
	}
	items, err := model.MyItemsAll(t.DB, uid)
	if err != nil {
		c.String(500, "%s", err)
		return
	}
	c.HTML(http.StatusOK, "mypage.tmpl", gin.H{
		"title":   "ShoesBox",
		"items":   items,
		"context": c,
	})
}

func (t *Item) Get(c *gin.Context) {
	id := c.Param("id")
	aid, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.String(500, "%s", err)
		return
	}
	item, err := model.ItemOne(t.DB, aid)
	if err != nil {
		c.String(500, "%s", err)
		return
	}
	c.HTML(http.StatusOK, "detail.tmpl", gin.H{
		"item":    item,
		"context": c,
		"csrf":    csrf.GetToken(c),
	})
}

func (t *Item) OpenGet(c *gin.Context) {
	id := c.Param("id")
	aid, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.String(500, "%s", err)
		return
	}
	item, err := model.ItemOne(t.DB, aid)
	if err != nil {
		c.String(500, "%s", err)
		return
	}
	c.HTML(http.StatusOK, "openDetail.tmpl", gin.H{
		"item":    item,
		"context": c,
		"csrf":    csrf.GetToken(c),
	})
}

func (t *Item) Make(c *gin.Context) {
	var item model.Item
	item.Body = c.PostForm("body")
	item.Title = c.PostForm("title")
	img, err := imageupload.Process(c.Request, "file")
	if err != nil {
		panic("sassasas")
	}
	nowString:=time.Now().Format("20060102_0000")
	pos := strings.Split(img.Filename, ".")
	item.Image = pos[0]+"_"+nowString+"."+pos[1]
	thumb, err := imageupload.ThumbnailJPEG(img, 300, 300, 100)
	if err != nil {
		thumb1, _ := imageupload.ThumbnailPNG(img, 300, 300)
		thumb1.Save("images/" + item.Image)
	} else {
		thumb.Save("images/" + item.Image)
	}
	t.New(c, &item)
}

func (t *Item) New(c *gin.Context, m *model.Item) {
	var id int64
	sess := sessions.Default(c)
	uid, ok := sess.Get("uid").(int64)
	if !ok {
		c.String(500, "")
		return
	}
	TXHandler(c, t.DB, func(tx *sql.Tx) error {
		result, err := m.Insert(tx, uid)
		if err != nil {
			return err
		}
		if err := tx.Commit(); err != nil {
			return err
		}
		id, err = result.LastInsertId()
		return err
	})
	c.Redirect(301, fmt.Sprintf("/item/%d", id))
}

// Edit indicates edit page for certain article.
func (t *Item) Edit(c *gin.Context) {
	id := c.Param("id")
	aid, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.String(500, "%s", err)
		return
	}
	item, err := model.ItemOne(t.DB, aid)
	if err != nil {
		c.String(500, "%s", err)
		return
	}
	c.HTML(http.StatusOK, "edit.tmpl", gin.H{
		"item":    item,
		"context": c,
		"csrf":    csrf.GetToken(c),
	})
}

func (t *Item) Upload(c *gin.Context) {
	id := c.PostForm("id")
	aid, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.String(500, "%s", err)
		return
	}
	img, err := imageupload.Process(c.Request, "file")

	if err != nil {
		panic(err)
	}
	var item model.Item
	item.Body = c.PostForm("body")
	item.Title = c.PostForm("title")
	item.ID = aid
	nowString :=time.Now().Format("20060102_00000")
	 pos := strings.Split(img.Filename, ".")
	 item.Image = pos[0]+"_"+nowString+"."+pos[1]
	thumb, err := imageupload.ThumbnailJPEG(img, 300, 300, 100)
	if err != nil {
		thumb1, _ := imageupload.ThumbnailPNG(img, 300, 300)
		thumb1.Save("images/" + item.Image)
	} else {
		thumb.Save("images/" + item.Image)
	}

	oneitem, err := model.ItemOne(t.DB, aid)
	if err != nil {
		c.String(500, "%s", err)
		return
	}
	if err := os.Remove("images/" + oneitem.Image); err != nil {
		fmt.Println(err)
	}
	t.Update(c, &item)
}

func (t *Item) Update(c *gin.Context, m *model.Item) {
	log.Printf("up:%v", m.ID)
	TXHandler(c, t.DB, func(tx *sql.Tx) error {
		if _, err := m.Update(tx); err != nil {
			return err
		}
		return tx.Commit()
	})
	c.Redirect(301, fmt.Sprintf("/mypage"))
}

func (t *Item) Delete(c *gin.Context) {
	var item model.Item
	id := c.PostForm("id")
	if id == "" {
		c.Abort()
		return
	}
	aid, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.String(500, "%s", err)
		return
	}
	item.ID = aid
	TXHandler(c, t.DB, func(tx *sql.Tx) error {
		if _, err := item.Delete(tx); err != nil {
			return err
		}
		return tx.Commit()
	})

	c.Redirect(301, "/mypage")
}
