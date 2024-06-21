package service

import (
	"context"
	"io"

	"gitub.com/wellingtonchida/gRPC/internal/database"
	"gitub.com/wellingtonchida/gRPC/internal/pb"
)

type CategoryService struct {
	pb.UnsafeCategoryServiceServer
	CategoryDb database.Category
}

func NewCategoryService(categoryDb database.Category) *CategoryService {
	return &CategoryService{CategoryDb: categoryDb}
}

func (c *CategoryService) CreateCategory(ctx context.Context, in *pb.CreateCategoryRequest) (*pb.Category, error) {
	category, err := c.CategoryDb.Create(in.Name, in.Description)
	if err != nil {
		return nil, err
	}

	categoryResponse := &pb.Category{
		Id:          category.ID,
		Name:        category.Name,
		Description: category.Description,
	}

	return categoryResponse, nil
}

func (c *CategoryService) GetCategories(ctx context.Context, in *pb.Blank) (*pb.CategoryListResponse, error) {
	categories, err := c.CategoryDb.FindAll()
	if err != nil {
		return nil, err
	}

	var categoryResponses []*pb.Category
	for _, category := range categories {
		categoryResponses = append(categoryResponses, &pb.Category{
			Id:          category.ID,
			Name:        category.Name,
			Description: category.Description,
		})
	}

	return &pb.CategoryListResponse{Categories: categoryResponses}, nil
}

func (c *CategoryService) GetCategory(ctx context.Context, in *pb.GetCatergoryRequest) (*pb.Category, error) {
	category, err := c.CategoryDb.Find(in.Id)
	if err != nil {
		return nil, err
	}

	categoryResponse := &pb.Category{
		Id:          category.ID,
		Name:        category.Name,
		Description: category.Description,
	}

	return categoryResponse, nil
}

func (c *CategoryService) CreateCategoryStream(stream pb.CategoryService_CreateCategoryStreamServer) error {
	categories := &pb.CategoryListResponse{}
	for {
		category, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(categories)
		}
		if err != nil {
			return err
		}

		createdCategory, err := c.CategoryDb.Create(category.Name, category.Description)
		if err != nil {
			return err
		}

		categories.Categories = append(categories.Categories, &pb.Category{
			Id:          createdCategory.ID,
			Name:        createdCategory.Name,
			Description: createdCategory.Description,
		})

	}
}

func (c *CategoryService) CreateCategoryStreamBidirectional(stream pb.CategoryService_CreateCategoryStreamBidirectionalServer) error {
	for {
		category, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		createdCategory, err := c.CategoryDb.Create(category.Name, category.Description)
		if err != nil {
			return err
		}

		err = stream.Send(&pb.Category{
			Id:          createdCategory.ID,
			Name:        createdCategory.Name,
			Description: createdCategory.Description,
		})

		if err != nil {
			return err
		}
	}
}
