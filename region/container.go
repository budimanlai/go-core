package region

import (
	"github.com/budimanlai/go-core/base"
	"github.com/budimanlai/go-core/region/handler"
	"github.com/budimanlai/go-core/region/repository"
	"github.com/budimanlai/go-core/region/service"
)

type RegionContainer struct {
	factory *base.Factory

	// Repositories
	SubdistrictRepository repository.SubdistrictRepository
	ProvinceRepository    repository.ProvinceRepository
	CountryinfoRepository repository.CountryinfoRepository
	DistrictRepository    repository.DistrictRepository
	CityRepository        repository.CityRepository

	// Services
	SubdistrictService service.SubdistrictService
	ProvinceService    service.ProvinceService
	CountryinfoService service.CountryinfoService
	DistrictService    service.DistrictService
	CityService        service.CityService

	// Handlers
	SubdistrictHandler *handler.SubdistrictHandler
	ProvinceHandler    *handler.ProvinceHandler
	CountryinfoHandler *handler.CountryinfoHandler
	DistrictHandler    *handler.DistrictHandler
	CityHandler        *handler.CityHandler
}

func NewRegionContainer(factory *base.Factory) *RegionContainer {
	region := &RegionContainer{
		factory: factory,
	}

	region.initRepositories()
	region.initServices()
	region.initHandlers()

	return region
}

func (c *RegionContainer) initRepositories() {
	c.SubdistrictRepository = repository.NewSubdistrictRepository(c.factory)
	c.ProvinceRepository = repository.NewProvinceRepository(c.factory)
	c.CountryinfoRepository = repository.NewCountryinfoRepository(c.factory)
	c.DistrictRepository = repository.NewDistrictRepository(c.factory)
	c.CityRepository = repository.NewCityRepository(c.factory)
}

func (c *RegionContainer) initServices() {
	c.SubdistrictService = service.NewSubdistrictService(c.SubdistrictRepository, c.factory.DB)
	c.ProvinceService = service.NewProvinceService(c.ProvinceRepository, c.factory.DB)
	c.CountryinfoService = service.NewCountryinfoService(c.CountryinfoRepository, c.factory.DB)
	c.DistrictService = service.NewDistrictService(c.DistrictRepository, c.factory.DB)
	c.CityService = service.NewCityService(c.CityRepository, c.factory.DB)
}

func (c *RegionContainer) initHandlers() {
	c.SubdistrictHandler = handler.NewSubdistrictHandler(c.SubdistrictService)
	c.ProvinceHandler = handler.NewProvinceHandler(c.ProvinceService)
	c.CountryinfoHandler = handler.NewCountryinfoHandler(c.CountryinfoService)
	c.DistrictHandler = handler.NewDistrictHandler(c.DistrictService)
	c.CityHandler = handler.NewCityHandler(c.CityService)
}
