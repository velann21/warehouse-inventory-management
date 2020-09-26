package database

type Article struct {
	ArtID         int    `json:"art_id, omitempty"`
	Name          string `json:"name, omitempty"`
	AvilableStock int    `json:"available_stock, omitempty"`
	SoldStock     int    `json:"sold_stock, omitempty"`
}

func NewArticle(artID, stock int, name string) *Article {
	return &Article{ArtID: artID, Name: name, AvilableStock: stock}
}

func (article *Article) SetArtID(artID int) {
	article.ArtID = artID
}

func (article *Article) SetName(name string) {
	article.Name = name
}

func (article *Article) SetAvailableStock(availStock int) {
	article.AvilableStock = availStock
}

func (article *Article) SetSoldStock(soldStock int) {
	article.SoldStock = soldStock
}

func (article *Article) GetArtID() int {
	return article.ArtID
}

func (article *Article) GetName() string {
	return article.Name
}

func (article *Article) GetAvailableStock() int {
	return article.AvilableStock
}

func (article *Article) GetSoldStock() int {
	return article.SoldStock
}
