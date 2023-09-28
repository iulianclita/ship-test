package ship_test

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/iulianclita/ship-test/ship"
)

func TestExtractOrderQty(t *testing.T) {
	testCases := map[string]struct {
		inputOrderQty    string
		expectedOrderQty int
		expectedError    error
	}{
		"input is not a valid integer": {
			inputOrderQty:    "XXX",
			expectedOrderQty: 0,
			expectedError:    ship.ErrOrderQtyInvalidFormat,
		},
		"input is a negative integer": {
			inputOrderQty:    "-999",
			expectedOrderQty: 0,
			expectedError:    ship.ErrOrderQtyInvalidValue,
		},
		"input is a valid integer": {
			inputOrderQty:    "999",
			expectedOrderQty: 999,
			expectedError:    nil,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			gotOrderQty, gotErr := ship.ExtractOrderQty(tc.inputOrderQty)
			if gotErr != nil {
				if !errors.Is(gotErr, tc.expectedError) {
					t.Errorf("got error = %v; want %v", gotErr, tc.expectedError)
				}
			} else {
				if gotOrderQty != tc.expectedOrderQty {
					t.Errorf("got ExtractOrderQty(%v) = %v; want %v", tc.inputOrderQty, gotOrderQty, tc.expectedOrderQty)
				}
			}
		})
	}
}

func TestExtractPackSizes(t *testing.T) {
	testCases := map[string]struct {
		inputPackSizes    string
		expectedPackSizes []int
		expectedError     error
	}{
		"input is not a valid list of integers": {
			inputPackSizes:    "XXX",
			expectedPackSizes: nil,
			expectedError:     ship.ErrPackSizesInvalidFormat,
		},
		"input has an invalid (negative) integer value": {
			inputPackSizes:    "1,2,3,-4",
			expectedPackSizes: nil,
			expectedError:     ship.ErrPackSizeInvalidValue,
		},
		"input is a valid list of integers": {
			inputPackSizes:    "1,2,3,4",
			expectedPackSizes: []int{1, 2, 3, 4},
			expectedError:     nil,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			gotPackSizes, gotErr := ship.ExtractPackSizes(tc.inputPackSizes)
			if gotErr != nil {
				if !errors.Is(gotErr, tc.expectedError) {
					t.Errorf("got error = %v; want %v", gotErr, tc.expectedError)
				}
			} else {
				if !cmp.Equal(gotPackSizes, tc.expectedPackSizes, cmpopts.SortSlices(
					func(a, b int) bool {
						return a < b
					},
				)) {
					t.Errorf("got ExtractPackSizes(%v) = %v; want %v\n(diff: %#v)",
						tc.inputPackSizes, gotPackSizes, tc.expectedPackSizes,
						cmp.Diff(tc.inputPackSizes, gotPackSizes),
					)
				}
			}
		})
	}
}

func TestCalculatePacksToShip(t *testing.T) {
	testCases := map[string]struct {
		inputOrderQty       int
		inputPackSizes      []int
		expectedPacksToShip map[int]int
	}{
		"1 item ordered": {
			inputOrderQty:  1,
			inputPackSizes: []int{250, 500, 1000, 2000, 5000},
			expectedPacksToShip: map[int]int{
				250: 1,
			},
		},
		"250 items ordered": {
			inputOrderQty:  250,
			inputPackSizes: []int{250, 500, 1000, 2000, 5000},
			expectedPacksToShip: map[int]int{
				250: 1,
			},
		},
		"251 items ordered": {
			inputOrderQty:  251,
			inputPackSizes: []int{250, 500, 1000, 2000, 5000},
			expectedPacksToShip: map[int]int{
				500: 1,
			},
		},
		"501 items ordered": {
			inputOrderQty:  501,
			inputPackSizes: []int{250, 500, 1000, 2000, 5000},
			expectedPacksToShip: map[int]int{
				250: 1,
				500: 1,
			},
		},
		"12001 items ordered": {
			inputOrderQty:  12001,
			inputPackSizes: []int{250, 500, 1000, 2000, 5000},
			expectedPacksToShip: map[int]int{
				250:  1,
				2000: 1,
				5000: 2,
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			gotPacksToShip := ship.CalculatePacksToShip(tc.inputOrderQty, tc.inputPackSizes)

			if !cmp.Equal(gotPacksToShip, tc.expectedPacksToShip, cmpopts.SortMaps(
				func(a, b int) bool {
					return a < b
				},
			)) {
				t.Errorf("got CalculatePacksToShip(orderQry = %v, packSizes = %v) = %v; want %v\n(diff: %#v)",
					tc.inputOrderQty, tc.inputPackSizes, tc.inputPackSizes, gotPacksToShip,
					cmp.Diff(tc.expectedPacksToShip, gotPacksToShip),
				)
			}
		})
	}
}
