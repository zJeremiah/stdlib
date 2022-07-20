package stats

import "errors"

// Labels is a special type to be used with StatsD and Prometheus.
//
// If using with Prometheus, there must be an even number of elements
// in this type.
type Labels []string

// ErrOddLabels is an error explaining that the labels must have
// an even number.
var ErrOddLabels = errors.New("labels must be an even number to convert to a map")

// AsMap converts Labels to a map. Meaning there must be an even number
// of elements in order to convert to key and values.
//
// This returns ErrOddLabels if there are an odd number of elements.
func (l Labels) AsMap() (map[string]string, error) {
	lNumElems := len(l)
	if lNumElems%2 != 0 {
		return nil, ErrOddLabels
	}
	if lNumElems == 0 {
		return nil, nil
	}

	m := make(map[string]string, lNumElems/2)

	for i := range l {
		if i%2 != 0 {
			continue
		}

		key := l[i]
		value := l[i+1]

		m[key] = value
	}

	return m, nil
}

// EmptyLabels returns an empty list of Labels.
func EmptyLabels() Labels {
	return Labels{}
}
