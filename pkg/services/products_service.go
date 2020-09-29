package services

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/velann21/warehouse-inventory-management/pkg/helpers"
	"github.com/velann21/warehouse-inventory-management/pkg/models/database"
	"github.com/velann21/warehouse-inventory-management/pkg/models/requests"
	"github.com/velann21/warehouse-inventory-management/pkg/repository"
	"strconv"
)

const (
	PRODUCTS_SERVICE_VERSION1 = "Version1"
)

type ProductsService interface {
	AddProducts(ctx context.Context, products *requests.AddProducts)
	AddProductsFromFile(ctx context.Context, decoder *json.Decoder, waitChannel chan bool) error
	GetAllProducts(ctx context.Context) ([]*database.Product, error)
	GetProductByID(ctx context.Context, productDetails *requests.GetProductDetails) ([]*database.ProductDetails, error)
	PurchaseProducts(ctx context.Context, products *requests.PurchaseProduct) error
}

type ProductsServiceImpl struct {
	repo repository.ProductsRepository
}

func NewProductsServiceFactory(version string, repo repository.ProductsRepository) ProductsService {
	switch version {
	case PRODUCTS_SERVICE_VERSION1:
		return &ProductsServiceImpl{repo: repo}
	default:
		return &ProductsServiceImpl{repo: repo}
	}
}

func (productsService *ProductsServiceImpl) AddProducts(ctx context.Context, products *requests.AddProducts) {
	for _, product := range products.Products {
		_, err := productsService.addProductJob(ctx, &product)
		if err != nil {
			//TODO Add into failed response
			logrus.WithError(err).Error("Something went wrong while addProductJob()")
		}
		//TODO Add into succed one response
	}
}

func (productsService *ProductsServiceImpl) AddProductsFromFile(ctx context.Context, decoder *json.Decoder, waitChannel chan bool) error {
	if decoder == nil {
		return nil
	}
	productsStreams := productsService.streamArticlesData(decoder)
	productsService.assignTask(ctx, productsStreams, waitChannel)
	return nil
}

func (productsService *ProductsServiceImpl) PurchaseProducts(ctx context.Context, product *requests.PurchaseProduct) error {
	dbProduct, err := productsService.getProduct(ctx, product.Name)
	if err != nil {
		if err.Error() == helpers.SQLRowNotFound {
			return err
		}
		return err
	}
	if dbProduct.Quantity <= 0 {
		return helpers.InvalidRequest
	}
	dbProduct.Quantity -= 1
	productIDInINT, err := strconv.Atoi(product.ID)
	productsAndArticles, err := productsService.getProductArticlesByProductID(ctx, productIDInINT)
	if err != nil {
		return err
	}
	err = productsService.repo.PurchaseProductsTransaction(ctx, dbProduct, productsAndArticles)
	if err != nil {
		return err
	}

	return nil
}

func (productsService *ProductsServiceImpl) DeleteByID(ctx context.Context, productID int) (int64, error) {
	return productsService.repo.DeleteByID(ctx, productID)
}

func (productsService *ProductsServiceImpl) GetProductByID(ctx context.Context, productDetails *requests.GetProductDetails) ([]*database.ProductDetails, error) {
	productID, _ := strconv.Atoi(productDetails.ID)
	productDetail, err := productsService.repo.GetProductDetailsByID(ctx, productID)
	if err != nil {
		return nil, err
	}
	return productDetail, nil
}

func (productsService *ProductsServiceImpl) GetAllProducts(ctx context.Context) ([]*database.Product, error) {
	productsList, err := productsService.repo.GetProducts(ctx)
	if err != nil {
		return nil, err
	}
	return productsList, nil
}

// All there are private functions
func (productsService *ProductsServiceImpl) streamArticlesData(decoder *json.Decoder) chan *requests.Products {
	streamChannel := make(chan *requests.Products, 15000)
	go func() {
		_, _ = decoder.Token()
		for decoder.More() {
			productsObj := requests.Products{}
			err := decoder.Decode(&productsObj)
			if err != nil {
				logrus.WithError(err).Error("Something wrong while decode inside streamProductsData()")
				continue
			}
			streamChannel <- &productsObj
		}
		close(streamChannel)
	}()
	return streamChannel
}

func (productsService *ProductsServiceImpl) getProduct(ctx context.Context, productName string) (*database.Product, error) {
	return productsService.repo.GetProduct(ctx, helpers.GetProductByName, productName)
}

func (productsService *ProductsServiceImpl) createNewProduct(ctx context.Context, product *requests.Products) (int64, error) {
	newProduct := database.NewProducts(0, product.Name, "Product "+product.Name, product.Price, 1)
	id, err := productsService.repo.InsertProduct(ctx, newProduct)
	if err != nil {
		return id, err
	}
	return id, nil
}

func (productsService *ProductsServiceImpl) updateProduct(ctx context.Context, existingproduct *database.Product) (int64, error) {
	existingproduct.Quantity += 1
	oldProduct := database.NewProducts(existingproduct.ID, existingproduct.Name, "Product "+existingproduct.Name, existingproduct.Price, existingproduct.Quantity)
	fmt.Println("old prod:", oldProduct)
	updatedID, err := productsService.repo.UpdateProduct(ctx, oldProduct)
	if err != nil {
		return updatedID, err
	}
	return updatedID, nil
}

func (productsService *ProductsServiceImpl) createProductArticleEntry(ctx context.Context, amount string, productID int, artID int) error {
	amountInInt, err := strconv.Atoi(amount)
	newProductArticles := database.NewProductsArticles(int(productID), artID, amountInInt, amountInInt)
	_, err = productsService.repo.InsertProductsArticles(ctx, newProductArticles)
	if err != nil {
		return err
	}
	return nil
}

//Todo Make this complete Mysql thing to one transaction so that we no need to hadle the rollback stuffs
func (productsService *ProductsServiceImpl) addProductJob(ctx context.Context, product *requests.Products) (int64, error) {
	//Checking whether the product is already available
	dbProduct, err := productsService.getProduct(ctx, product.Name)
	if err != nil {
		//If it is not available the it return SQLRowNotFound error
		if err.Error() == helpers.SQLRowNotFound {
			// If row is not found then I am creting the new product entry
			createdID, err := productsService.createNewProduct(ctx, product)
			if err != nil {
				return createdID, err
			}
			// addProductArticlesJob here
			err = productsService.addProductArticlesJob(ctx, createdID, product)
			if err != nil {
				// if any error happen to addProductArticlesJob then revert back the product as well.
				_, err = productsService.DeleteByID(ctx, int(createdID))
				if err != nil {
					return -1, err
				}
			}
			return createdID, nil
		}
		return -1, err
	}

	//This block execute when the product is already avialble
	if dbProduct != nil {
		// If product already exist trying to update the Quantity count
		updatedID, err := productsService.updateProduct(ctx, dbProduct)
		if err != nil {
			return updatedID, err
		}
		// Here I am updating the required Articles count as well
		for _, article := range product.Articles {
			artID, _ := strconv.Atoi(article.ArtID)
			// Checking whether the art id exist
			dbarticle, err := productsService.repo.GetArticleByID(ctx, helpers.GetArticleByID, artID)
			if err != nil {
				if err.Error() == helpers.SQLRowNotFound {
					// TODO What if the there are one Article not another one in Articles entry what should we do.
					return -1, err
				}
				return -1, err
			}
			// Checking whether the getProductArticles entry availabe for product ID and ART ID
			dbProductArticle, err := productsService.getProductArticles(ctx, dbProduct.ID, artID)
			if err != nil {
				continue
			}
			amount, err := strconv.Atoi(article.AmountOf)
			// If Entry is available the update the Articles counts
			dbProductArticle.TotalRequiredQuantity += amount
			_, err = productsService.updateProductArticles(ctx, dbProductArticle)
			if err != nil {
				//Todo roll back all Product count
				logrus.WithError(err).Error("Error occured while updateProductArticles() ArticleID: ", dbProductArticle.InventoryID, " ------  ProductID:", dbProductArticle.ProductID)
				continue
			}

			dbarticle.AvilableStock -= amount
			_, err = productsService.repo.UpdateArticle(ctx, dbarticle)
			if err != nil {
				//Todo roll back all Product count and Product_inventory entry
				return -1, err
			}

		}
		return updatedID, nil
	}
	return -1, nil
}

func (productsService *ProductsServiceImpl) addProductArticlesJob(ctx context.Context, productID int64, product *requests.Products) error {
	for _, productArticle := range product.Articles {
		article, err := productsService.repo.GetArticleByID(ctx, helpers.GetArticleByID, productArticle.ArtID)
		if err != nil {
			if err.Error() == helpers.SQLRowNotFound {
				// TODO What if the there are one Article not another one in Articles entry what should we do.
				return err
			}
			return err
		}

		err = productsService.createProductArticleEntry(ctx, productArticle.AmountOf, int(productID), article.ArtID)
		if err != nil {
			return err
		}
		amountInInt, err := strconv.Atoi(productArticle.AmountOf)
		article.AvilableStock -= amountInInt

		_, err = productsService.repo.UpdateArticle(ctx, article)
		if err != nil {
			//Todo roll back all Product count and Product_inventory entry
			return err
		}
	}
	return nil
}

func (productsService *ProductsServiceImpl) updateProductArticles(ctx context.Context, productsArticles *database.ProductsArticles) (int64, error) {
	ID, err := productsService.repo.UpdateproductArticles(ctx, productsArticles)
	if err != nil {
		return ID, err
	}
	return ID, nil
}

func (productsService *ProductsServiceImpl) getProductArticles(ctx context.Context, productID int, artID int) (*database.ProductsArticles, error) {
	return productsService.repo.GetProductArticle(ctx, helpers.GetProductArticle, productID, artID)
}

func (productsService *ProductsServiceImpl) getProductArticlesByProductID(ctx context.Context, productID int) ([]*database.ProductsArticles, error) {
	return productsService.repo.GetProductArticleByProductId(ctx, productID)
}

func (productsService *ProductsServiceImpl) assignTask(ctx context.Context, productStreams chan *requests.Products, waitChannel chan bool) {
	go func() {
		for productData := range productStreams {
			_, err := productsService.addProductJob(ctx, productData)
			if err != nil {
				logrus.WithError(err).Error("Something went wrong while addProductJob for" + productData.Name)
				//TODO  Add the error message into response struct
				return
			}
			//TODO Add the success message into response struct
		}
		waitChannel <- true
	}()
}
