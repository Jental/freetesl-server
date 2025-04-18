package common

import "fmt"

// MapErr is similar to lo.Map, but handles error in iteratee function
func MapErr[T any, R any](collection []T, iteratee func(T, int) (R, error)) ([]R, error) {
	result := make([]R, len(collection))

	for i, item := range collection {
		res, err := iteratee(item, i)
		if err != nil {
			newErr := fmt.Errorf("MapErr: error during parsing %d element: %s", i, err)
			return nil, newErr
		}
		result[i] = res
	}

	return result, nil
}
