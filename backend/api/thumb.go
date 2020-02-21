package api

import (
	"net/http"

	"github.com/globalsign/mgo/bson"
	"github.com/guad/paperless2/backend/db"
	"github.com/guad/paperless2/backend/model"
	"github.com/guad/paperless2/backend/storage"
	"github.com/labstack/echo"
	"github.com/minio/minio-go"
)

func GetThumbnail(c echo.Context) error {
	doc := c.Param("doc")

	sesh := db.Ctx()
	defer sesh.Close()
	col := sesh.DB("paperless").C("documents")

	if !bson.IsObjectIdHex(doc) {
		return c.JSON(http.StatusBadRequest, struct{}{})
	}

	id := bson.ObjectIdHex(doc)

	var document model.Document

	err := col.FindId(id).One(&document)

	if err != nil {
		return err
	}

	if document.ThumbnailPath == "" {
		return c.String(http.StatusNotFound, "")
	}

	key := document.ThumbnailPath

	file, err := storage.S3.GetObject(
		storage.DocumentBucket,
		key,
		minio.GetObjectOptions{},
	)

	if err != nil {
		return err
	}

	defer file.Close()

	return c.Stream(http.StatusOK, "image/png", file)
}
