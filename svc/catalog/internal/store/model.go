package store

// Item represent Items entity in underlying db
type Item struct {
	ID          uint
	Name        *string
	Description *string
	Price       *float64
	PriceCode   *string
}
