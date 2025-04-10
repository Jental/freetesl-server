package models

type Attribute struct {
	ID   int
	Name string
}

type AttributeSorter struct {
	Attributes []*Attribute
}

// func (sorter *AttributeSorter) Len() int {
// 	return len(sorter.Attributes)
// }

// // Swap is part of sort.Interface.
// func (sorter *AttributeSorter) Swap(i, j int) {
// 	sorter.Attributes[i], sorter.Attributes[j] = sorter.planets[j], s.planets[i]
// }

// // Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
// func (sorter *AttributeSorter) Less(i, j int) bool {
// 	return s.by(&s.planets[i], &s.planets[j])
// }
