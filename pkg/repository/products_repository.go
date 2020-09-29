package repository

import (
	"context"
	"fmt"
	"github.com/velann21/warehouse-inventory-management/pkg/helpers"
	"github.com/velann21/warehouse-inventory-management/pkg/helpers/databases"
	"github.com/velann21/warehouse-inventory-management/pkg/models/database"
)

const (
	PRODUCTS_REPO_VERSION1 = "version1"
)

//1. Todo to reuse the code of transactions instead of doing it always in each functions
//2. TODO Make use of proper transaction isolation level based on type of transaction
type ProductsRepository interface {
	InsertProduct(ctx context.Context, product *database.Product) (int64, error)
	GetProduct(ctx context.Context, query string, args ...interface{}) (*database.Product, error)
	GetProducts(ctx context.Context) ([]*database.Product, error)
	UpdateProduct(ctx context.Context, product *database.Product) (int64, error)
	InsertProductsArticles(ctx context.Context, productArticle *database.ProductsArticles) (int64, error)
	GetArticleByID(ctx context.Context, query string, args ...interface{}) (*database.Article, error)
	UpdateproductArticles(ctx context.Context, productArticle *database.ProductsArticles) (int64, error)
	GetProductArticle(ctx context.Context, query string, args ...interface{}) (*database.ProductsArticles, error)
	GetProductArticleByProductId(ctx context.Context, productId int) ([]*database.ProductsArticles, error)
	PurchaseProductsTransaction(ctx context.Context, product *database.Product, productsAndArticles []*database.ProductsArticles) error
	DeleteByID(ctx context.Context, productID int) (int64, error)
	GetProductDetailsByID(ctx context.Context, productID int) ([]*database.ProductDetails, error)
	UpdateArticle(ctx context.Context, article *database.Article) (int64, error)
}

type ProductsRepositoryImpl struct {
	client databases.SqlClient
}

func NewProductsRepositoryFactory(version string, client databases.SqlClient) ProductsRepository {
	switch version {
	case PRODUCTS_REPO_VERSION1:
		return &ProductsRepositoryImpl{client: client}
	default:
		return &ProductsRepositoryImpl{client: client}
	}
}

func (productsRepo *ProductsRepositoryImpl) InsertProduct(ctx context.Context, product *database.Product) (int64, error) {
	options := productsRepo.client.BuildOptions(false, productsRepo.client.GetIsolationLevel(1))
	tx, err := productsRepo.client.BeginWithContext(ctx, &options)
	if err != nil {
		return -1, err
	}

	stmt, err := productsRepo.client.PrepareWithContext(ctx, tx, helpers.CreateProduct)
	if err != nil {
		_ = productsRepo.client.RollBack(tx)
		return -1, err
	}
	res, err := productsRepo.client.ExecWithContext(ctx, stmt, product.GetName(), product.GetDescription(), product.GetPrice(), product.GetQuantity())
	if err != nil {
		_ = productsRepo.client.RollBack(tx)
		return -1, err
	}
	id, err := productsRepo.client.LastInsertedID(res)
	if err != nil {
		_ = productsRepo.client.RollBack(tx)
		return -1, err
	}
	err = productsRepo.client.Commit(tx)
	if err != nil {
		err = productsRepo.client.RollBack(tx)
		if err != nil {
			return -1, err
		}
		return -1, err
	}
	return id, nil

}

func (productsRepo *ProductsRepositoryImpl) GetProduct(ctx context.Context, query string, args ...interface{}) (*database.Product, error) {
	result := productsRepo.client.QueryRowWithContext(ctx, query, args...)
	product := database.Product{}
	err := result.Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.Quantity)
	if err != nil {

		return nil, err
	}
	return &product, nil
}

func (productsRepo *ProductsRepositoryImpl) GetProducts(ctx context.Context) ([]*database.Product, error) {
	results, err := productsRepo.client.QueryWithContext(ctx, helpers.GetAllProducts)
	if err != nil {
		return nil, err
	}
	finalResult := []*database.Product{}
	for results.Next() {
		product := database.Product{}
		err := results.Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.Quantity)
		if err != nil {
			return nil, err
		}
		finalResult = append(finalResult, &product)
	}
	return finalResult, nil

}

func (productsRepo *ProductsRepositoryImpl) UpdateProduct(ctx context.Context, product *database.Product) (int64, error) {
	options := productsRepo.client.BuildOptions(false, productsRepo.client.GetIsolationLevel(1))
	tx, err := productsRepo.client.BeginWithContext(ctx, &options)
	if err != nil {
		return -1, err
	}
	stmt, err := productsRepo.client.PrepareWithContext(ctx, tx, helpers.UpdateProductByID)
	if err != nil {
		_ = productsRepo.client.RollBack(tx)
		return -1, err
	}
	fmt.Println("PID", product.ID)
	res, err := productsRepo.client.ExecWithContext(ctx, stmt, product.GetName(), product.GetDescription(), product.GetPrice(), product.GetQuantity(), product.ID)
	if err != nil {
		_ = productsRepo.client.RollBack(tx)
		return -1, err
	}
	id, err := productsRepo.client.LastInsertedID(res)
	if err != nil {
		_ = productsRepo.client.RollBack(tx)
		return -1, err
	}
	err = productsRepo.client.Commit(tx)
	if err != nil {
		err = productsRepo.client.RollBack(tx)
		if err != nil {
			return -1, err
		}
		return -1, err
	}
	return id, nil
}

func (productsRepo *ProductsRepositoryImpl) InsertProductsArticles(ctx context.Context, productArticle *database.ProductsArticles) (int64, error) {
	options := productsRepo.client.BuildOptions(false, productsRepo.client.GetIsolationLevel(1))
	tx, err := productsRepo.client.BeginWithContext(ctx, &options)
	if err != nil {
		return -1, err
	}

	stmt, err := productsRepo.client.PrepareWithContext(ctx, tx, helpers.CreateProductArticle)
	if err != nil {
		_ = productsRepo.client.RollBack(tx)
		return -1, err
	}
	res, err := productsRepo.client.ExecWithContext(ctx, stmt, productArticle.ProductID, productArticle.InventoryID, productArticle.EachQuantity, productArticle.TotalRequiredQuantity)
	if err != nil {
		_ = productsRepo.client.RollBack(tx)
		return -1, err
	}
	id, err := productsRepo.client.LastInsertedID(res)
	if err != nil {
		_ = productsRepo.client.RollBack(tx)
		return -1, err
	}
	err = productsRepo.client.Commit(tx)
	if err != nil {
		err = productsRepo.client.RollBack(tx)
		if err != nil {
			return -1, err
		}
		return -1, err
	}
	return id, nil
}

func (productsRepo *ProductsRepositoryImpl) UpdateproductArticles(ctx context.Context, productArticle *database.ProductsArticles) (int64, error) {
	options := productsRepo.client.BuildOptions(false, productsRepo.client.GetIsolationLevel(1))
	tx, err := productsRepo.client.BeginWithContext(ctx, &options)
	if err != nil {
		return -1, err
	}
	stmt, err := productsRepo.client.PrepareWithContext(ctx, tx, helpers.UpdateProductArticle)
	if err != nil {
		_ = productsRepo.client.RollBack(tx)
		return -1, err
	}
	res, err := productsRepo.client.ExecWithContext(ctx, stmt, productArticle.EachQuantity, productArticle.TotalRequiredQuantity, productArticle.ProductID, productArticle.InventoryID)
	if err != nil {
		_ = productsRepo.client.RollBack(tx)
		return -1, err
	}
	id, err := productsRepo.client.LastInsertedID(res)
	if err != nil {
		_ = productsRepo.client.RollBack(tx)
		return -1, err
	}
	err = productsRepo.client.Commit(tx)
	if err != nil {
		err = productsRepo.client.RollBack(tx)
		if err != nil {
			return -1, err
		}
		return -1, err
	}
	return id, nil
}

func (productsRepo *ProductsRepositoryImpl) GetProductArticle(ctx context.Context, query string, args ...interface{}) (*database.ProductsArticles, error) {
	result := productsRepo.client.QueryRowWithContext(ctx, query, args...)
	productArticle := database.ProductsArticles{}
	err := result.Scan(&productArticle.ProductID, &productArticle.InventoryID, &productArticle.EachQuantity, &productArticle.TotalRequiredQuantity)
	if err != nil {
		return nil, err
	}
	return &productArticle, nil
}

func (productsRepo *ProductsRepositoryImpl) GetProductArticleByProductId(ctx context.Context, productId int) ([]*database.ProductsArticles, error) {
	results, err := productsRepo.client.QueryWithContext(ctx, helpers.GetProductArticleByProductID, productId)
	if err != nil {
		return nil, err
	}
	finalResult := []*database.ProductsArticles{}
	for results.Next() {
		product := database.ProductsArticles{}
		err := results.Scan(&product.ProductID, &product.InventoryID, &product.EachQuantity, &product.TotalRequiredQuantity)
		if err != nil {
			return nil, err
		}
		finalResult = append(finalResult, &product)
	}
	return finalResult, nil
}

func (productsRepo *ProductsRepositoryImpl) PurchaseProductsTransaction(ctx context.Context, product *database.Product, productsAndArticles []*database.ProductsArticles) error {
	options := productsRepo.client.BuildOptions(false, productsRepo.client.GetIsolationLevel(1))
	tx, err := productsRepo.client.BeginWithContext(ctx, &options)
	if err != nil {
		return err
	}

	//Updating the product first
	stmt, err := productsRepo.client.PrepareWithContext(ctx, tx, helpers.UpdateProductByID)
	if err != nil {
		_ = productsRepo.client.RollBack(tx)
		return err
	}
	_, err = productsRepo.client.ExecWithContext(ctx, stmt, product.GetName(), product.GetDescription(), product.GetPrice(), product.GetQuantity(), product.ID)
	if err != nil {
		_ = productsRepo.client.RollBack(tx)
		return err
	}

	//Updating the ProductsAdnArticles Table
	for _, productAndArticle := range productsAndArticles {
		productAndArticle.TotalRequiredQuantity -= productAndArticle.EachQuantity
		stmt, err := productsRepo.client.PrepareWithContext(ctx, tx, helpers.UpdateProductArticle)
		if err != nil {
			_ = productsRepo.client.RollBack(tx)
			return err
		}
		_, err = productsRepo.client.ExecWithContext(ctx, stmt, productAndArticle.EachQuantity, productAndArticle.TotalRequiredQuantity, productAndArticle.ProductID, productAndArticle.InventoryID)
		if err != nil {
			_ = productsRepo.client.RollBack(tx)
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		_ = productsRepo.client.RollBack(tx)
		return err
	}
	return nil
}

func (productsRepo *ProductsRepositoryImpl) DeleteByID(ctx context.Context, productID int) (int64, error) {
	options := productsRepo.client.BuildOptions(false, productsRepo.client.GetIsolationLevel(1))
	tx, err := productsRepo.client.BeginWithContext(ctx, &options)
	if err != nil {
		return -1, err
	}
	stmt, err := productsRepo.client.PrepareWithContext(ctx, tx, helpers.DeleteByProductID)
	if err != nil {
		_ = productsRepo.client.RollBack(tx)
		return -1, err
	}
	res, err := productsRepo.client.ExecWithContext(ctx, stmt, productID)
	if err != nil {
		_ = productsRepo.client.RollBack(tx)
		return -1, err
	}
	id, err := productsRepo.client.LastInsertedID(res)
	if err != nil {
		_ = productsRepo.client.RollBack(tx)
		return -1, err
	}
	err = productsRepo.client.Commit(tx)
	if err != nil {
		err = productsRepo.client.RollBack(tx)
		if err != nil {
			return -1, err
		}
		return -1, err
	}
	return id, nil
}

func (productsRepo *ProductsRepositoryImpl) GetProductDetailsByID(ctx context.Context, productID int) ([]*database.ProductDetails, error) {
	results, err := productsRepo.client.QueryWithContext(ctx, helpers.GetProductDetailsByID, productID)
	if err != nil {
		return nil, err
	}
	finalResult := []*database.ProductDetails{}
	for results.Next() {
		product := database.ProductDetails{}
		err := results.Scan(&product.QuantityEach, &product.ArtID, &product.Name)
		if err != nil {
			return nil, err
		}
		finalResult = append(finalResult, &product)
	}
	return finalResult, nil
}

// Todo: Move this to article repo
func (productsRepo *ProductsRepositoryImpl) UpdateArticle(ctx context.Context, article *database.Article) (int64, error) {
	options := productsRepo.client.BuildOptions(false, productsRepo.client.GetIsolationLevel(1))
	tx, err := productsRepo.client.BeginWithContext(ctx, &options)
	if err != nil {
		return -1, err
	}
	stmt, err := productsRepo.client.PrepareWithContext(ctx, tx, helpers.UpdateArticleByID)
	if err != nil {
		_ = productsRepo.client.RollBack(tx)
		return -1, err
	}
	res, err := productsRepo.client.ExecWithContext(ctx, stmt, article.GetName(), article.GetAvailableStock(), article.GetSoldStock(), article.GetArtID())
	if err != nil {
		_ = productsRepo.client.RollBack(tx)
		return -1, err
	}
	id, err := productsRepo.client.LastInsertedID(res)
	if err != nil {
		_ = productsRepo.client.RollBack(tx)
		return -1, err
	}
	err = productsRepo.client.Commit(tx)
	if err != nil {
		err = productsRepo.client.RollBack(tx)
		if err != nil {
			return -1, err
		}
		return -1, err
	}
	return id, nil
}

// Todo: Move this to article repo
func (productsRepo *ProductsRepositoryImpl) GetArticleByID(ctx context.Context, query string, args ...interface{}) (*database.Article, error) {
	result := productsRepo.client.QueryRowWithContext(ctx, query, args...)
	article := database.Article{}
	err := result.Scan(&article.ArtID, &article.Name, &article.AvilableStock, &article.SoldStock)
	if err != nil {
		return nil, err
	}
	return &article, nil
}