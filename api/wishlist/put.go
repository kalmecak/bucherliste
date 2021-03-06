package wishlist

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kalmecak/bucherliste/common"
	"github.com/kalmecak/bucherliste/sql"
	logger "github.com/kalmecak/go-error-logger"
	"gorm.io/gorm"
)

// Put Modifica el contenido de una wishlist
func Put(c *fiber.Ctx) error {

	wlid := c.Context().UserValue("paramID").(*sql.UID)
	uid := c.Context().UserValue("userID").(*sql.UID)
	books := c.Context().UserValue("books").(*[]sql.Book)
	remove := c.Context().UserValue("del").(*bool)

	// conectamos a la base de datos
	db, err := sql.GormDB()
	if err != nil {

		logger.Error(err, "api.wishlist.Put.db.GormDB")
		var res common.Response
		res.InternalError("could not connect to BD", "")
		return c.Status(fiber.StatusInternalServerError).JSON(&res)
	}

	// Buscamos la wishlist
	wl := sql.Wishlist{
		ID: *wlid,
	}
	tx := db.Where("user_id = ?", *uid).First(&wl)
	if err := tx.Error; err != nil {

		var res common.Response
		// manejamos el error en caso que no exista la wishlist
		if err == gorm.ErrRecordNotFound {

			res.NotFound("wishlist not found")
			c.Status(fiber.StatusNotFound)
		} else {

			// en caso que sea un error diferente
			logger.Error(err, "api.wishlist.Put.db.find.First")
			res.InternalError("could not find wishlist", "")
			c.Status(fiber.StatusInternalServerError)
		}

		return c.JSON(&res)
	}

	// Si la wishlist existe, modificamos el contenido

	// Según lo indicado, eliminamos/agregamos los libros de la wishlist
	if *remove {

		// Ejecutamos el médoto de la wishlist para eliminar la relación de los libros
		err = wl.Delete(books)
	} else {

		// Ejecutamos el médoto de la wishlist para agregar la relación de los libros
		err = wl.Append(books)
	}
	// verificamos el error de la operación
	if err != nil {

		logger.Error(err, "api.wishlist.Put.db.Association")
		var res common.Response
		res.InternalError("could not update wishlist", "")
		return c.Status(fiber.StatusInternalServerError).JSON(&res)
	}
	res := common.Response{}
	res.Ok("")
	return c.Status(fiber.StatusOK).JSON(&res)
}
