package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/globalsign/mgo/bson"
	"github.com/streadway/amqp"

	"github.com/guad/paperless2/backend/api/user"
	"github.com/guad/paperless2/backend/broker"
	"github.com/guad/paperless2/backend/db"
	"github.com/guad/paperless2/backend/model"
	"github.com/labstack/echo"
)

func or(val, def string) string {
	if val == "" {
		return def
	}

	return val
}

func ListDocuments(c echo.Context) error {
	userid := user.GetUserID(c)

	sesh := db.Ctx()
	defer sesh.Close()

	col := sesh.DB("paperless").C("documents")

	offset, _ := strconv.Atoi(or(c.QueryParam("offset"), "0"))
	limit, _ := strconv.Atoi(or(c.QueryParam("limit"), "50"))
	order := or(c.QueryParam("order"), "ASC")
	field := or(c.QueryParam("sort"), "_id")

	filter := or(c.QueryParam("filter"), "{}")
	var filters map[string]interface{}
	_ = json.Unmarshal([]byte(filter), &filters)

	search := ""

	if _, ok := filters["q"]; ok {
		search = filters["q"].(string)
	}

	var m bson.M = bson.M{}

	if search != "" {
		if len(strings.Split(search, " ")) > 1 {
			m = bson.M{"$text": bson.M{
				"$search": search,
			}}
		} else {
			m = bson.M{
				"$or": []bson.M{
					bson.M{"content": bson.M{"$regex": search, "$options": "i"}},
					bson.M{"title": bson.M{"$regex": search, "$options": "i"}},
				},
			}
		}
	}

	if tags, ok := filters["tags"]; ok {
		split := strings.Split(tags.(string), " ")

		m = bson.M{
			"$and": []bson.M{
				m,
				bson.M{
					"tags": bson.M{
						"$all": split,
					},
				},
			},
		}
	}

	m["user_id"] = userid

	query := col.Find(m)

	count, _ := query.Count()

	if order == "ASC" {
		query = query.Sort(field)
	} else {
		query = query.Sort("-" + field)
	}

	query = query.
		Skip(offset).
		Limit(limit)

	var items []model.Document

	err := query.All(&items)

	if err != nil {
		return err
	}

	if items == nil {
		items = []model.Document{}
	}

	return c.JSON(http.StatusOK, struct {
		Data  []model.Document `json:"data"`
		Total int              `json:"total"`
	}{items, count})
}

func GetDocument(c echo.Context) error {
	id, err := getIdParam(c.Param("id"))
	userid := user.GetUserID(c)

	if err != nil {
		return err
	}

	sesh := db.Ctx()
	defer sesh.Close()

	col := sesh.DB("paperless").C("documents")

	var doc model.Document
	err = col.Find(bson.M{"_id": id, "user_id": userid}).One(&doc)

	if err != nil {
		return err
	}

	if doc.UserID.Hex() != userid {
		return c.JSON(http.StatusForbidden, struct{}{})
	}

	return c.JSON(http.StatusOK, doc)
}

func UpdateDocument(c echo.Context) error {
	sesh := db.Ctx()
	defer sesh.Close()

	col := sesh.DB("paperless").C("documents")
	userid := user.GetUserID(c)

	var newDoc model.Document
	err := c.Bind(&newDoc)

	if err != nil {
		return err
	}

	id, err := getIdParam(c.Param("id"))

	if err != nil {
		return err
	}

	var doc model.Document
	err = col.Find(bson.M{"_id": id, "user_id": userid}).One(&doc)

	if err != nil {
		return err
	}

	if doc.UserID.Hex() != userid {
		return c.JSON(http.StatusForbidden, struct{}{})
	}

	doc.Content = newDoc.Content
	doc.Title = newDoc.Title
	doc.Tags = newDoc.Tags

	err = col.UpdateId(id, doc)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, doc)
}

func DeleteDocument(c echo.Context) error {
	sesh := db.Ctx()
	defer sesh.Close()

	col := sesh.DB("paperless").C("documents")
	userid := user.GetUserID(c)
	id, err := getIdParam(c.Param("id"))

	if err != nil {
		return err
	}

	var doc model.Document

	err = col.Find(bson.M{"_id": id, "user_id": userid}).One(&doc)

	if err != nil {
		return err
	}

	if doc.UserID.Hex() != userid {
		return c.JSON(http.StatusForbidden, struct{}{})
	}

	err = col.RemoveId(id)

	if err != nil {
		return err
	}

	// Delete from S3 as well
	jsonbytes, _ := json.Marshal(doc)

	broker.RabbitMQ.Publish(
		broker.DocumentCleanupQueue, // Exchange
		"",                          // routing key
		false,                       // mandatory
		false,                       // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         jsonbytes,
			DeliveryMode: amqp.Persistent,
		},
	)

	return c.JSON(http.StatusOK, struct{}{})

}

func getIdParam(p string) (bson.ObjectId, error) {
	if !bson.IsObjectIdHex(p) {
		return "", fmt.Errorf("Not an id")
	}

	return bson.ObjectIdHex(p), nil
}
