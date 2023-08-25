package main

// Hot take: the generic version of this is way too noisey
//
// Initially the API wasn't generic, and it was a bit cleaner.
//
import (
	"testing"
	"math"
)

const float64EqualityThreshold = 1e-9

type Number interface {
	int | int8 | int16 | int32 | int64 | float32 | float64
}

func nearly[T Number](l, r T) bool {
	if l == r {
		return true
	}
	
	return math.Abs(float64(l - r)) <= float64EqualityThreshold
}


func Eq[T Number](l, r []T) bool {
	if len(l) != len(r) {
		return false
	}

	for i := 0; i < len(l); i++ {
		if !nearly(l[i],  r[i]) {
			return false
		}
	}

	return true
}

func Eqs(l, r []string) bool {
	if len(l) != len(r) {
		return false
	}

	for i := 0; i < len(l); i++ {
		if l[i] !=  r[i] {
			return false
		}
	}

	return true
}

func Test_DAGPath(t *testing.T) {
	g := CreateDag(
		Connection[float32]{
			NewNode[float32]("m"),
			&Edge[float32]{3.28},
			NewNode[float32]("ft"),
		},
		Connection[float32]{
			NewNode[float32]("ft"),
			&Edge[float32]{12},
			NewNode[float32]("in"),
		},
	)

	tests := []struct{
		From     string
		To       string
		Expected []string
	}{
		{"m", "ft", []string{"m", "ft"}},
		{"m", "in", []string{"m", "ft", "in"}},
		{"m", "m",  []string{"m"}},
		{"m", "hr", []string{}},
	}

	for i, test := range tests {
		result := g.Path(test.From, test.To)
		if !Eqs(test.Expected, result) {
			t.Logf("index: %d\n", i)
			t.Logf("\texpected: %v\n", test.Expected)
			t.Logf("\tresult: %v", result)
			t.FailNow()
		}
	}
}

func Test_DAGWeights(t *testing.T) {
	g := CreateDag(
		Connection[float32]{
			NewNode[float32]("m"),
			&Edge[float32]{3.28},
			NewNode[float32]("ft"),
		},
		Connection[float32]{
			NewNode[float32]("ft"),
			&Edge[float32]{12},
			NewNode[float32]("in"),
		},
		Connection[float32]{
			NewNode[float32]("in"),
			&Edge[float32]{2.54},
			NewNode[float32]("cm"),
		},
	)

	tests := []struct{
		From     string
		To       string
		Expected []float32
	}{
		{"m", "ft", []float32{3.28}},
		{"m", "in", []float32{3.28, 12}},
		{"m", "m",  []float32{}},
		{"m", "hr", []float32{}},
	}

	for i, test := range tests {
		result := g.Weights(test.From, test.To)
		if !Eq(test.Expected, result) {
			t.Logf("index: %d\n", i)
			t.Logf("\texpected: %v\n", test.Expected)
			t.Logf("\tresult: %v", result)
			t.FailNow()
		}
	}
}

func Test_DGPath(t *testing.T) {
	g := CreateDg[float32](
		func(v float32) float32 {
			return 1.0 / v
		},
		Connection[float32]{NewNode[float32]("m"), &Edge[float32]{3.28}, NewNode[float32]("ft")},
		Connection[float32]{NewNode[float32]("ft"), &Edge[float32]{12}, NewNode[float32]("in")},
		Connection[float32]{NewNode[float32]("in"), &Edge[float32]{2.54}, NewNode[float32]("cm")},
	)

	tests := []struct{
		From     string
		To       string
		Expected []string
	}{
		{"m", "ft", []string{"m", "ft"}},
		{"m", "in", []string{"m", "ft", "in"}},
		{"m", "m",  []string{"m"}},
		{"in", "m", []string{"in", "ft", "m"}},
		{"m", "hr", []string{}},
	}


	for i, test := range tests {
		result := g.Path(test.From, test.To)
		if !Eqs(test.Expected, result) {
			t.Logf("index: %d\n", i)
			t.Logf("\texpected: %v\n", test.Expected)
			t.Logf("\tresult: %v", result)
			t.FailNow()
		}
	}
}

func Test_DGWeights(t *testing.T) {
	g := CreateDg(
		func(v float32) float32 {
			return 1.0 / v
		},
		Connection[float32]{NewNode[float32]("m"), &Edge[float32]{3.28}, NewNode[float32]("ft")},
		Connection[float32]{NewNode[float32]("ft"), &Edge[float32]{12}, NewNode[float32]("in")},
		Connection[float32]{NewNode[float32]("in"), &Edge[float32]{2.54}, NewNode[float32]("cm")},
	)

	tests := []struct{
		From     string
		To       string
		Expected []float32
	}{
		{"m", "ft", []float32{3.28}},
		{"m", "in", []float32{3.28, 12}},
		{"in", "m", []float32{(1.0 / 12), (1.0 / 3.28)}},
		{"m", "m",  []float32{}},
		{"m", "hr", []float32{}},
	}

	for i, test := range tests {
		result := g.Weights(test.From, test.To)
		if !Eq(test.Expected, result) {
			t.Logf("index: %d\n", i)
			t.Logf("\texpected: %v\n", test.Expected)
			t.Logf("\tresult: %v", result)
			t.FailNow()
		}
	}
}
