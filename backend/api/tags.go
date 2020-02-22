package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"

	"github.com/guad/paperless2/backend/db"
	"github.com/guad/paperless2/backend/model"
	"github.com/labstack/echo"
)

func ListTags(c echo.Context) error {
	sesh := db.Ctx()
	defer sesh.Close()

	col := sesh.DB("paperless").C("tags")

	offset, _ := strconv.Atoi(or(c.QueryParam("offset"), "0"))
	limit, _ := strconv.Atoi(or(c.QueryParam("limit"), "50"))
	order := or(c.QueryParam("order"), "ASC")
	field := or(c.QueryParam("sort"), "_id")

	filter := or(c.QueryParam("filter"), "{}")
	var filters map[string]interface{}
	_ = json.Unmarshal([]byte(filter), &filters)

	var query *mgo.Query

	if search, ok := filters["name"]; ok {
		query = col.Find(bson.M{"name": bson.M{"$regex": search, "$options": "i"}})
	} else {
		query = col.Find(bson.M{})
	}

	count, _ := query.Count()

	if order == "ASC" {
		query = query.Sort(field)
	} else {
		query = query.Sort("-" + field)
	}

	query = query.
		Skip(offset).
		Limit(limit)

	var items []model.Tag

	err := query.All(&items)

	if err != nil {
		return err
	}

	if items == nil {
		items = []model.Tag{}
	}

	return c.JSON(http.StatusOK, struct {
		Data  []model.Tag `json:"data"`
		Total int         `json:"total"`
	}{items, count})
}

func GetTag(c echo.Context) error {
	id, err := getIdParam(c.Param("id"))

	if err != nil {
		return err
	}

	sesh := db.Ctx()
	defer sesh.Close()

	col := sesh.DB("paperless").C("tags")

	var doc model.Tag
	err = col.FindId(id).One(&doc)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, doc)
}

func CreateTag(c echo.Context) error {
	sesh := db.Ctx()
	defer sesh.Close()

	col := sesh.DB("paperless").C("tags")

	var newDoc model.Tag
	err := c.Bind(&newDoc)

	if err != nil {
		return err
	}

	newDoc.ID = bson.NewObjectId()
	newDoc.Name = strings.ToLower(newDoc.Name)

	err = col.Insert(newDoc)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, newDoc)
}

func UpdateTag(c echo.Context) error {
	sesh := db.Ctx()
	defer sesh.Close()

	col := sesh.DB("paperless").C("tags")

	var newDoc model.Tag
	err := c.Bind(&newDoc)

	if err != nil {
		return err
	}

	id, err := getIdParam(c.Param("id"))

	if err != nil {
		return err
	}

	var doc model.Tag
	err = col.FindId(id).One(&doc)

	if err != nil {
		return err
	}

	doc.Name = strings.ToLower(newDoc.Name)
	doc.Regex = newDoc.Regex
	doc.Implies = newDoc.Implies

	err = col.UpdateId(id, doc)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, doc)
}

func DeleteTag(c echo.Context) error {
	sesh := db.Ctx()
	defer sesh.Close()

	col := sesh.DB("paperless").C("tags")

	id, err := getIdParam(c.Param("id"))

	if err != nil {
		return err
	}

	var doc model.Tag

	err = col.FindId(id).One(&doc)

	if err != nil {
		return err
	}

	err = col.RemoveId(id)

	if err != nil {
		return err
	}

	// TODO: Delete from all documents as well

	return c.JSON(http.StatusOK, struct{}{})

}
