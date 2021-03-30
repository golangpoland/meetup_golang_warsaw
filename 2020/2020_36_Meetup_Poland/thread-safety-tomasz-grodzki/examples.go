
// #1 OMIT

var s int

go func() {
	s += 1
}()

go func() {
	s += 1
}()

// ## OMIT

// #2 OMIT

var s int

go func() {
	i := s + 1
	// ...
}()

go func() {
	i := 2*s - 1
	// ...
}()

// ## OMIT

// #3 OMIT

var s = make([]int, 10)

go func() {
	s[0] += 1
}()

go func() {
	s[1] += 1
}()

// ## OMIT

// #4 OMIT

var weight = make(map[string]float64)

go func() {
	weight["gopher"] = 0.2
}()

go func() {
	weight["python"] = 12
}()

// ## OMIT

// #5 OMIT

var weight = map[string]float64{
	"gopher": 0.2,
	"python": 12,
}

go func() {
	w := weight["gopher"]
	// ...
}()

go func() {
	w := weight["python"]
	// ...
}()

// ## OMIT
