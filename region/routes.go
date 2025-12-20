package region

import (
	"github.com/gofiber/fiber/v2"
)

// SetReadOnlyRoutes sets the read-only routes for the region API.
func SetReadOnlyRoutes(app fiber.Router, container *RegionContainer) {
	regionGroup := app.Group("/region")
	regionGroup.Get("/subdistricts", container.SubdistrictHandler.Index)
	regionGroup.Get("/provinces", container.ProvinceHandler.Index)
	regionGroup.Get("/countryinfos", container.CountryinfoHandler.Index)
	regionGroup.Get("/districts", container.DistrictHandler.Index)
	regionGroup.Get("/citys", container.CityHandler.Index)
}

// SetCRUDRoutes sets the CRUD routes for the region API.
func SetCRUDRoutes(app fiber.Router, container *RegionContainer) {
	regionGroup := app.Group("/region")

	// subdistrict CRUD
	subdistrictGroup := regionGroup.Group("/subdistricts")
	subdistrictGroup.Post("/", container.SubdistrictHandler.Create)
	subdistrictGroup.Get("/", container.SubdistrictHandler.Index)
	subdistrictGroup.Get("/:id", container.SubdistrictHandler.View)
	subdistrictGroup.Put("/:id", container.SubdistrictHandler.Update)
	subdistrictGroup.Delete("/:id", container.SubdistrictHandler.Delete)

	// province CRUD
	provinceGroup := regionGroup.Group("/provinces")
	provinceGroup.Post("/", container.ProvinceHandler.Create)
	provinceGroup.Get("/", container.ProvinceHandler.Index)
	provinceGroup.Get("/:id", container.ProvinceHandler.View)
	provinceGroup.Put("/:id", container.ProvinceHandler.Update)
	provinceGroup.Delete("/:id", container.ProvinceHandler.Delete)

	// countryinfo CRUD
	countryinfoGroup := regionGroup.Group("/countryinfos")
	countryinfoGroup.Post("/", container.CountryinfoHandler.Create)
	countryinfoGroup.Get("/", container.CountryinfoHandler.Index)
	countryinfoGroup.Get("/:id", container.CountryinfoHandler.View)
	countryinfoGroup.Put("/:id", container.CountryinfoHandler.Update)
	countryinfoGroup.Delete("/:id", container.CountryinfoHandler.Delete)

	// district CRUD
	districtGroup := regionGroup.Group("/districts")
	districtGroup.Post("/", container.DistrictHandler.Create)
	districtGroup.Get("/", container.DistrictHandler.Index)
	districtGroup.Get("/:id", container.DistrictHandler.View)
	districtGroup.Put("/:id", container.DistrictHandler.Update)
	districtGroup.Delete("/:id", container.DistrictHandler.Delete)

	// city CRUD
	cityGroup := regionGroup.Group("/citys")
	cityGroup.Post("/", container.CityHandler.Create)
	cityGroup.Get("/", container.CityHandler.Index)
	cityGroup.Get("/:id", container.CityHandler.View)
	cityGroup.Put("/:id", container.CityHandler.Update)
	cityGroup.Delete("/:id", container.CityHandler.Delete)
}
