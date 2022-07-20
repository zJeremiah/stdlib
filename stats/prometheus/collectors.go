package prometheus

import "github.com/prometheus/client_golang/prometheus"

type Collectors map[string]prometheus.Collector

func (k Collectors) Merge(o Collectors) Collectors {
	for key, value := range o {
		k[key] = value
	}

	return k
}
