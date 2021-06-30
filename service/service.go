package service

import (

	//"github.com/confluentinc/confluent-kafka-go/kafka"

	"fmt"
	"github/rossi1/go-api-microservice-example/domain/entity"
	"github/rossi1/go-api-microservice-example/domain/repository"
	"github/rossi1/go-api-microservice-example/internal/search"
	"net/http"

	"github.com/gofrs/uuid"
	"github.com/olivere/elastic"
	"gorm.io/gorm"
)

type Service struct {
	catRepo  repository.CategoryRepository
	prodRepo repository.ProductRepository
}

func New(db *gorm.DB, search interface{}, pub interface{}, topic string) *Service {
	var (
		service Service
		//producer publisher.Pub

	)
	//switch pub := pub.(type) {
	//case *kafka.Producer:
	//	producer = publisher.NewKaftaPublisher(pub, topic)
	//}
	//client := getSearcher(search)
	//fmt.Printf("type = %T\n", client)
	service.catRepo = repository.NewCategoryRepo(db)
	service.prodRepo = repository.NewProductRepo(db)

	return &service
}

func getSearcher(s interface{}) search.Searcher {
	switch s := s.(type) {
	case *elastic.Client:
		fmt.Println("testing sake here", s)
		return nil //search.NewElasticSearcher(s)
	default:
		fmt.Println("no+")
		return nil
	}
}

func (s *Service) CreateCategory(name string, product []entity.Product) (*entity.Category, error) {
	result, err := s.catRepo.Create(name, product)
	if err != nil {
		return nil, err
	}
	return result, nil

}

func (s *Service) UpdateCategory(name, categoryID string) (*entity.Category, error) {
	id, err := uuid.FromString(categoryID)
	if err != nil {
		return nil, err
	}
	result, err := s.catRepo.Update(id, name)
	if err != nil {
		return nil, err
	}

	return result, nil

}

func (s *Service) GetAllCategories(r *http.Request) ([]entity.Category, error) {
	result, err := s.catRepo.FindAll(r)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *Service) GetCategory(categoryID string) (*entity.Category, error) {
	id, err := uuid.FromString(categoryID)
	if err != nil {
		return nil, err
	}
	result, err := s.catRepo.Find(id)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *Service) DeleteCategory(categoryID string) error {
	id, err := uuid.FromString(categoryID)
	if err != nil {
		return err
	}
	err = s.catRepo.Delete(id)
	if err != nil {
		return nil
	}
	return nil

}

func (s *Service) CreateProduct() error {
	return nil
}

func (s *Service) UpdateProduct(productID string, params map[string]interface{}) (*entity.Product, error) {
	id, err := uuid.FromString(productID)
	if err != nil {
		return nil, err
	}
	result, err := s.prodRepo.Update(id, params)
	if err != nil {
		return nil, err
	}

	return result, nil

}

func (s *Service) GetProduct(productID string) (*entity.Product, error) {
	id, err := uuid.FromString(productID)
	if err != nil {
		return nil, err
	}
	result, err := s.prodRepo.Find(id)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *Service) DeleteProduct(productID string) error {
	id, err := uuid.FromString(productID)
	if err != nil {
		return err
	}
	err = s.prodRepo.Delete(id)
	if err != nil {
		return nil
	}
	return nil
}

func (s *Service) GetAllProducts(categoryID string, r *http.Request) ([]entity.Product, error) {
	id, err := uuid.FromString(categoryID)
	if err != nil {
		return nil, err
	}
	result, err := s.prodRepo.FindAll(id, r)
	if err != nil {
		return nil, err
	}
	return result, nil
}
