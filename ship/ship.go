package ship

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

var (
	ErrOrderQtyInvalidFormat  = errors.New("invalid format for ordered items input (should be an integer value)")
	ErrOrderQtyInvalidValue   = errors.New("invalid value for ordered items input (should be a strictly positive integer)")
	ErrPackSizesInvalidFormat = errors.New("invalid format for pack sizes input (should be a list of integer values)")
	ErrPackSizeInvalidValue   = errors.New("invalid value for pack size (every value should be strictly positive)")
)

// ExtractOrderQty extracts the number of items to be ordered from the request input
func ExtractOrderQty(orderQtyStr string) (int, error) {
	orderQty, err := strconv.Atoi(orderQtyStr)
	if err != nil {
		// This error is not handled anywhere but we know it's a format error
		return 0, fmt.Errorf("wrong input: %w: %v", ErrOrderQtyInvalidFormat, err)
	}

	if orderQty <= 0 {
		return 0, fmt.Errorf("wrong input: %w", ErrOrderQtyInvalidValue)
	}

	return orderQty, nil
}

// ExtractPackSizes extracts the list of available pack sizes from the request input
func ExtractPackSizes(packSizesStr string) ([]int, error) {
	// Split request input to get all available pack sizes (as strings)
	packSizesStrs := strings.Split(packSizesStr, ",")
	// Remove any extra spaces
	packSizesStrList := make([]string, len(packSizesStrs))
	for i, packSizeStr := range packSizesStrs {
		packSizesStrList[i] = strings.TrimSpace(packSizeStr)
	}
	// Try to convert the input strings into integers
	packSizes := make([]int, len(packSizesStrList))
	for i, packSizeStr := range packSizesStrList {
		packSize, err := strconv.Atoi(packSizeStr)
		if err != nil {
			// This error is not handled anywhere but we know it's a format error
			return nil, fmt.Errorf("wrong input: %w: %v", ErrPackSizesInvalidFormat, err)
		}
		if packSize <= 0 {
			return nil, fmt.Errorf("wrong input: %w", ErrPackSizeInvalidValue)
		}
		packSizes[i] = packSize
	}

	return packSizes, nil
}

// CalculatePacksToShip calculates how many packs to ship according to the following rules:
// 1. Only whole packs can be sent. Packs cannot be broken open.
// 2. Within the constraints of Rule 1 above, send out no more items than necessary to fulfil the order.
// 3. Within the constraints of Rules 1 and 2 above, send out as few packs as possible to fulfil each order.
func CalculatePacksToShip(orderQty int, packSizes []int) map[int]int {
	// Sort pack sizes in reverse order to fill up the largest ones first
	sort.Sort(sort.Reverse(sort.IntSlice(packSizes)))

	packSizesNeeded := make(map[int]int)
	remainingQty := orderQty

	for _, packSize := range packSizes {
		if remainingQty <= 0 {
			break
		}

		packs := remainingQty / packSize
		if packs > 0 {
			packSizesNeeded[packSize] = packs
			remainingQty -= packs * packSize
		}
	}

	if remainingQty > 0 {
		// If there's still some remaining quantity, use the smallest available pack size
		smallestPackSize := packSizes[len(packSizes)-1]
		packSizesNeeded[smallestPackSize]++
	}

	// Redistribute the order to use minimum number of packs
	sort.Sort(sort.IntSlice(packSizes))

	packSizesAll := make(map[int]struct{})
	for _, packSize := range packSizes {
		packSizesAll[packSize] = struct{}{}
	}

	for _, packSize := range packSizes {
		packSizeNeeded := packSizesNeeded[packSize]

		packsNeeded := packSizeNeeded * packSize
		if packsNeeded <= 0 {
			continue
		}

		if _, found := packSizesAll[packsNeeded]; found {
			delete(packSizesNeeded, packSize)
			packSizesNeeded[packsNeeded]++
		}
	}

	return packSizesNeeded
}
