package handlers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/emarifer/calorie-tracker-app/pkg/api/config"
	"github.com/emarifer/calorie-tracker-app/pkg/api/models"
	"github.com/gin-gonic/gin"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func AddEntry(c *gin.Context) {
	entryCollection := config.OpenCollection(config.Client, "calories")
	entry := models.Entry{}
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	if err := c.BindJSON(&entry); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		fmt.Println(err)
		return
	}

	// We validate that the fields are not empty
	if err := models.ValidateStruct(entry); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		fmt.Println(err)
		return
	}

	// entry.ID = primitive.NewObjectID()
	entry.CreatedAt = time.Now().UTC()
	result, err := entryCollection.InsertOne(ctx, entry)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "order item was not created",
		})
		fmt.Println(err)
		return
	}

	c.JSON(http.StatusCreated, result)
}

func GetEntries(c *gin.Context) {
	entryCollection := config.OpenCollection(config.Client, "calories")
	entries := []models.Entry{}
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	// VER nota abajo para las options
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})
	cursor, err := entryCollection.Find(ctx, bson.M{}, opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		fmt.Println(err)
		return
	}

	timeZone := c.GetHeader("X-TimeZone") // VER nota abajo
	loc, _ := time.LoadLocation(timeZone)
	for cursor.Next(ctx) {
		currentEntry := new(models.Entry)
		cursor.Decode(currentEntry)
		currentEntry.CreatedAt = currentEntry.CreatedAt.In(loc)
		entries = append(entries, *currentEntry)
	}

	c.JSON(http.StatusOK, entries)
}

func GetEntryById(c *gin.Context) {
	entryCollection := config.OpenCollection(config.Client, "calories")
	entryID := c.Params.ByName("id")
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	docID, err := primitive.ObjectIDFromHex(entryID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		fmt.Println(err)
		return
	}

	entry := new(models.Entry)
	if err := entryCollection.FindOne(ctx, bson.M{"_id": docID}).Decode(entry); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		fmt.Println(err)
		return
	}

	timeZone := c.GetHeader("X-TimeZone") // VER nota abajo
	loc, _ := time.LoadLocation(timeZone)
	entry.CreatedAt = entry.CreatedAt.In(loc)

	c.JSON(http.StatusOK, entry)
}

func GetEntriesByIngredient(c *gin.Context) {
	entryCollection := config.OpenCollection(config.Client, "calories")
	entries := []models.Entry{}
	ingredient := c.Params.ByName("ingredient")
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})
	cursor, err := entryCollection.Find(ctx, bson.M{"ingredients": ingredient}, opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		fmt.Println(err)
		return
	}

	timeZone := c.GetHeader("X-TimeZone") // VER nota abajo
	loc, _ := time.LoadLocation(timeZone)
	for cursor.Next(ctx) {
		currentEntry := new(models.Entry)
		cursor.Decode(currentEntry)
		currentEntry.CreatedAt = currentEntry.CreatedAt.In(loc)
		entries = append(entries, *currentEntry)
	}

	c.JSON(http.StatusOK, entries)
}

func UpdateEntry(c *gin.Context) {
	entryCollection := config.OpenCollection(config.Client, "calories")
	entryID := c.Params.ByName("id")
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	docID, err := primitive.ObjectIDFromHex(entryID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		fmt.Println(err)
		return
	}

	entry := models.Entry{}
	if err := c.BindJSON(&entry); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"errores": err.Error(),
		})
		fmt.Println(err)
		return
	}

	// VER nota abajo sobre UpdateOne/FindOneAndUpdate
	result, err := entryCollection.UpdateOne(
		ctx,
		bson.M{"_id": docID},
		bson.D{{Key: "$set", Value: entry}},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		fmt.Println(err)
		return
	}

	if result.ModifiedCount != 0 {
		c.JSON(http.StatusOK, gin.H{
			"updated_entry": result.ModifiedCount,
		})
	} else {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "entry not found in database",
		})
		fmt.Println(err)
	}
}

func UpadateIngredient(c *gin.Context) {
	entryCollection := config.OpenCollection(config.Client, "calories")
	entryID := c.Params.ByName("id")
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	docID, err := primitive.ObjectIDFromHex(entryID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		fmt.Println(err)
		return
	}

	ingredients := models.IngredientsUpdate{}
	if err := c.BindJSON(&ingredients); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		fmt.Println(err)
		return
	}

	// We validate that the "ingredients" field is not empty
	if err := models.ValidateStruct(ingredients); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		fmt.Println(err)
		return
	}

	// VER nota abajo sobre UpdateOne/FindOneAndUpdate
	result, err := entryCollection.UpdateOne(
		ctx,
		bson.M{"_id": docID},
		bson.D{{Key: "$set", Value: ingredients}},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		fmt.Println(err)
		return
	}

	if result.ModifiedCount != 0 {
		c.JSON(http.StatusOK, gin.H{
			"updated_ingredients": result.ModifiedCount,
		})
	} else {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "entry not found in database",
		})
		fmt.Println(err)
	}
}

func DeleteEntry(c *gin.Context) {
	entryCollection := config.OpenCollection(config.Client, "calories")
	entryID := c.Params.ByName("id")
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	docID, err := primitive.ObjectIDFromHex(entryID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		fmt.Println(err)
		return
	}

	result, err := entryCollection.DeleteOne(ctx, bson.M{
		"_id": docID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		fmt.Println(err)
		return
	}

	if result.DeletedCount != 0 {
		c.JSON(http.StatusOK, gin.H{
			"deleted_entry": result.DeletedCount,
		})
	} else {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "entry not found in database",
		})
		fmt.Println(err)
	}
}

/* Fecha y Hora local. VER:
https://stackoverflow.com/questions/1091372/getting-the-clients-time-zone-and-offset-in-javascript#34602679
https://rapidapi.com/guides/request-headers-fetch
https://stackoverflow.com/questions/49913266/how-to-get-header-data-of-postman-using-gin-package-in-golang#74414677
https://stackoverflow.com/questions/25318154/convert-utc-to-local-time-in-go#25319349

Ordenación descendente según campo (en este caso por fecha y hora). VER:
https://www.mongodb.com/docs/drivers/go/current/fundamentals/crud/read-operations/sort/#descending

Actualización únicamente de los campos enviado. VER:
https://stackoverflow.com/questions/54087093/mongodb-update-only-changed-fields-in-an-object-instead-of-replacing-the-whole-o
https://www.mongodb.com/docs/drivers/go/current/fundamentals/crud/write-operations/upsert/#insert-or-update-in-a-single-operation
https://pkg.go.dev/go.mongodb.org/mongo-driver@v1.12.1/mongo#Collection.UpdateOne
https://pkg.go.dev/go.mongodb.org/mongo-driver@v1.12.1/mongo#Collection.FindOneAndUpdate
https://stackoverflow.com/questions/36362457/bson-m-to-struct
*/
