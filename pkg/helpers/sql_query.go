package helpers

const (
	CreateArticle     = "Insert into inventory (name, available_stock, sold_stock) values (?, ?, ?);"
	GetAllArticles    = "Select * from inventory;"
	GetArticleByName  = "Select * from inventory where name=?;"
	GetArticleByID    = "Select * from inventory where art_id=?;"
	UpdateArticleByID = "update inventory set name =?, available_stock=?, sold_stock=? where art_id=?;"

	CreateProduct     = "Insert into products (name, description, price, available_quantity) Values (?,?,?,?);"
	GetProductByName  = "Select * from products where name=?;"
	UpdateProductByID = "Update products set name =?, description = ?, price = ?, available_quantity = ? where id = ?;"
	GetAllProducts    = "Select * from products;"
	DeleteByProductID = "Delete from products where id = ?;"
	GetProductDetailsByID        = "Select  pa.quantity_each,  inventory.art_id, inventory.name from products as prod Inner Join product_inventory pa on prod.id = pa.product_id Inner Join (select * from inventory) inventory on pa.inventory_id=inventory.art_id where prod.id=?  ORDER BY id;"

	CreateProductArticle         = "Insert into product_inventory (product_id, inventory_id, quantity_each, total_required_quantity) Values (?, ?, ?, ?);"
	UpdateProductArticle         = "update product_inventory set quantity_each = ?, total_required_quantity = ? where product_id = ? and inventory_id = ?;"
	GetProductArticle            = "Select * from product_inventory where product_id=? and inventory_id=?;"
	GetProductArticleByProductID = "Select * from product_inventory where product_id = ?;"
)
