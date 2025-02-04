package utils

import (
	"math"
	"strconv"
	"strings"
	"unicode"
)

func alphanumericCount(name string) int {
	points := 0
	for _, c := range name {
		if unicode.IsDigit(c) || unicode.IsLetter(c) {
			points++
		}
	}
	return points
}

func calculateTotal(total string) int {
	points := 0
	price, _ := strconv.ParseFloat(total, 64)
	if math.Ceil(price) == price {
		points += 50
	}

	if math.Mod(price, 0.25) == 0 {
		points += 25
	}
	return points
}

func calculateItems(items []Item) int {
	points := 0
	points += (len(items) / 2) * 5
	for _, item := range items {
		trimDescription := strings.TrimSpace(*item.ShortDescription)
		if len(trimDescription)%3 == 0 {
			price, _ := strconv.ParseFloat(*item.Price, 64)
			price *= 0.2
			points += int(math.Ceil(price))
		}
	}
	return points
}

func calculateTime(date, time string) int {
	points := 0
	splitDate := strings.Split(date, "-")
	splitTime := strings.Split(time, ":")
	day, _ := strconv.Atoi(splitDate[2])
	if day%2 != 0 {
		points += 6
	}
	hour, _ := strconv.Atoi(splitTime[0])
	if hour == 14 {
		minute, _ := strconv.Atoi(splitTime[1])
		if minute > 0 {
			points += 10
		}
	} else if hour == 15 {
		points += 10
	}
	return points
}

func CalculatePoints(receipt Receipt) int {
	points := 0
	points += alphanumericCount(*receipt.Retailer)
	points += calculateTotal(*receipt.Total)
	points += calculateItems(receipt.Items)
	points += calculateTime(*receipt.PurchaseDate, *receipt.PurchaseTime)
	return points
}
