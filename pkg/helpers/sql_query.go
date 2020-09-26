package helpers

const (
	CreateArticle    = "Insert into inventory (name, available_stock, sold_stock) values (?, ?, ?);"
	GetArticleByName = "Select * from inventory where name=?;"
	UpdateArticle    = "update inventory set name =?, available_stock=?, sold_stock=? where art_id=?;"
)
