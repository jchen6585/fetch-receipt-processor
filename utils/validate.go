package utils

func ValidateReceipt(receipt Receipt) bool {
	if receipt.Retailer == nil || receipt.PurchaseDate == nil || receipt.PurchaseTime == nil || receipt.Total == nil {
		return false
	}
	for _, item := range receipt.Items {
		if item.Price == nil || item.ShortDescription == nil {
			return false
		}
	}
	return true
}

func ConvertStringToPointer(str string) *string {
	temp := &str
	return temp
}
