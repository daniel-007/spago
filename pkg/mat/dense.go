// Copyright 2019 spaGO Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mat

import (
	"brillion.io/spago/pkg/mat/internal/asm/f64"
	"fmt"
	"math"
)

type Dense struct {
	rows int
	cols int
	size int // rows*cols
	data []float64
}

// Following 'runtime: memmove sometimes faster than memclrNoHeapPointers (https:///golang/go/issues/23306)',
// I opted to preallocate an 'empty' array initialized to zero, preferring append operations rather than make.
// Other reference: https:///go101/go-benchmarks/tree/master/append-vs-make
var empty [10000000]float64

// NewDense returns a new rows x cols dense matrix populated with a copy of the elements.
// The elements cannot be nil, panic otherwise. Use NewEmptyDense to initialize an empty matrix.
func NewDense(rows, cols int, elements []float64) *Dense {
	if elements == nil {
		panic("mat: elements cannot be nil. Use NewEmptyDense() instead.")
	}
	return newDense(rows, cols, elements)
}

// NewVecDense returns a new column vector populated with the elements.
// The elements cannot be nil, panic otherwise. Use NewEmptyVecDense to initialize an empty matrix.
func NewVecDense(elements []float64) *Dense {
	if elements == nil {
		panic("mat: elements cannot be nil. Use NewEmptyVecDense() instead.")
	}
	return newDense(len(elements), 1, elements)
}

// NewScalar returns a new 1x1 matrix containing the input value.
func NewScalar(n float64) *Dense {
	data := []float64{n}
	return &Dense{
		rows: 1,
		cols: 1,
		size: 1,
		data: data,
	}
}

// NewEmptyVecDense returns a new vector of the given size, initialized to zeros.
func NewEmptyVecDense(size int) *Dense {
	return NewEmptyDense(size, 1)
}

// NewEmptyVecDense returns a new rows x cols matrix initialized to zeros.
func NewEmptyDense(rows, cols int) *Dense {
	return newDense(rows, cols, empty[:rows*cols])
}

// NewEmptyVecDense returns a new one-hot vector of the given size.
func OneHotVecDense(size int, oneAt int) *Dense {
	vec := NewEmptyVecDense(size)
	vec.Set(1.0, oneAt)
	return vec
}

// newDense returns a new rows x cols dense matrix populated with a copy of the elements.
func newDense(rows, cols int, elements []float64) *Dense {
	return &Dense{
		rows: rows,
		cols: cols,
		size: rows * cols,
		data: append([]float64(nil), elements...),
	}
}

// NewInitDense returns a new rows x cols dense matrix initialized with a constant value.
func NewInitDense(rows, cols int, val float64) *Dense {
	out := NewEmptyDense(rows, cols)
	f64.AddConst(val, out.data)
	return out
}

// NewInitDense returns a new size x 1 dense matrix initialized with a constant value.
func NewInitVecDense(size int, val float64) *Dense {
	return NewInitDense(size, 1, val)
}

// SetData sets the data
func (d *Dense) SetData(data []float64) {
	_ = append(d.data[:0], data...)
}

// ZerosLike returns a new Dense with of the same dimensions of the receiver, initialized with zeros.
func (d *Dense) ZerosLike() Matrix {
	return NewEmptyDense(d.rows, d.cols)
}

// ZerosLike returns a new Dense with of the same dimensions of the receiver, initialized to ones.
func (d *Dense) OnesLike() Matrix {
	buf := d.ZerosLike()
	f64.AddConst(1.0, buf.(*Dense).data)
	return buf
}

// Clone returns a new matrix copying the values of the receiver.
func (d *Dense) Clone() Matrix {
	return NewDense(d.rows, d.cols, d.data)
}

// Copy copies the data to the receiver.
func (d *Dense) Copy(other Matrix) {
	//if !EqualDims(d, other) {
	//	panic("mat: incompatible matrix dimensions.")
	//}
	if other, ok := other.(*Dense); !ok {
		//panic("mat: incompatible matrix types.")
	} else {
		_ = append(d.data[:0], other.data...)
	}
}

// Zeros set all the values to zeros.
func (d *Dense) Zeros() {
	_ = append(d.data[:0], empty[:d.size]...)
}

// Dims returns the number of rows and columns.
func (d *Dense) Dims() (r, c int) {
	return d.rows, d.cols
}

// Rows returns the number of rows.
func (d *Dense) Rows() int {
	return d.rows
}

// Columns returns the number of columns.
func (d *Dense) Columns() int {
	return d.cols
}

// Size returns the size of the matrix (rows * cols).
func (d *Dense) Size() int {
	return d.size
}

// LastIndex returns the last index.
func (d *Dense) LastIndex() int {
	return d.size - 1
}

// Data returns the underlying data.
func (d *Dense) Data() []float64 {
	return d.data
}

// IsVectors returns whether the matrix has one row or one column, or not.
func (d *Dense) IsVector() bool {
	return d.rows == 1 || d.cols == 1
}

// IsScalar returns whether the matrix contains a scalar, or not.
func (d *Dense) IsScalar() bool {
	return d.size == 1
}

// Scalar returns the scalar. It panics if the matrix contains more elements.
func (d *Dense) Scalar() float64 {
	//if !d.IsScalar() {
	//	panic("mat: expected scalar but the matrix contains more elements.")
	//}
	return d.data[0]
}

func (d *Dense) Set(v float64, i int, j ...int) {
	//if len(j) > 1 {
	//	panic("mat: invalid 'j' argument.")
	//}
	if len(j) > 0 {
		d.data[i*d.cols+j[0]] = v
	} else {
		d.data[i] = v
	}
}

func (d *Dense) At(i int, j ...int) float64 {
	//if len(j) > 1 {
	//	panic("mat: invalid 'j' argument.")
	//}
	if len(j) > 0 {
		return d.data[i*d.cols+j[0]]
	} else {
		return d.data[i]
	}
}

// T returns the transpose of the matrix.
func (d *Dense) T() Matrix {
	r, c := d.Dims()
	m := NewEmptyDense(c, r)
	for i := 0; i < r; i++ {
		for j := 0; j < c; j++ {
			m.data[j*d.rows+i] = d.data[i*d.cols+j]
		}
	}
	return m
}

// Reshape returns a copy of the matrix. It panics if the dimensions are not compatible.
func (d *Dense) Reshape(r, c int) Matrix {
	if d.Size() != r*c {
		panic("mat: incompatible sizes.")
	}
	return NewDense(r, c, d.data)
}

func (d *Dense) ApplyWithAlpha(fn func(i, j int, v float64, alpha ...float64) float64, a Matrix, alpha ...float64) {
	//if !EqualDims(d, a) {
	//	panic("mat: incompatible matrix dimensions.")
	//}
	for i := 0; i < d.rows; i++ {
		for j := 0; j < d.cols; j++ {
			d.data[i*d.cols+j] = fn(i, j, a.At(i, j), alpha...)
		}
	}
}

func (d *Dense) Apply(fn func(i, j int, v float64) float64, a Matrix) {
	//if !EqualDims(d, a) {
	//	panic("mat: incompatible matrix dimensions.")
	//}
	for i := 0; i < d.rows; i++ {
		for j := 0; j < d.cols; j++ {
			d.data[i*d.cols+j] = fn(i, j, a.At(i, j))
		}
	}
}

func (d *Dense) AddScalar(n float64) Matrix {
	out := d.Clone().(*Dense)
	f64.AddConst(n, out.data)
	return out
}

func (d *Dense) SubScalar(n float64) Matrix {
	out := d.Clone().(*Dense)
	f64.AddConst(-n, out.data)
	return out
}

// Add adds the scalar to the receiver.
func (d *Dense) AddScalarInPlace(n float64) {
	f64.AddConst(n, d.data)
}

// Sub subtracts the scalar to the receiver.
func (d *Dense) SubScalarInPlace(n float64) {
	for i := 0; i < len(d.data); i++ {
		d.data[i] -= n
	}
}

func (d *Dense) ProdScalarInPlace(n float64) {
	f64.ScalUnitary(n, d.data)
}

func (d *Dense) ProdMatrixScalarInPlace(m Matrix, n float64) {
	b := m.(*Dense)
	f64.ScalUnitaryTo(d.data, n, b.data)
}

func (d *Dense) ProdScalar(n float64) Matrix {
	out := d.ZerosLike().(*Dense)
	f64.ScalUnitaryTo(out.data, n, d.data)
	return out
}

func (d *Dense) Add(other Matrix) Matrix {
	//if !(EqualDims(d, other) ||
	//	(other.Columns() == 1 && other.Rows() == d.Rows()) ||
	//	(other.IsVector() && d.IsVector() && other.Size() == d.Size())) {
	//	panic("mat: matrices with not compatible size")
	//}
	b := other.(*Dense)
	out := d.ZerosLike().(*Dense)
	f64.AxpyUnitaryTo(out.data, 1.0, b.data, d.data)
	return out
}

func (d *Dense) AddInPlace(other Matrix) {
	//if !(EqualDims(d, other) ||
	//	(other.Columns() == 1 && other.Rows() == d.Rows()) ||
	//	(other.IsVector() && d.IsVector() && other.Size() == d.Size())) {
	//	panic("mat: matrices with not compatible size")
	//}
	b := other.(*Dense)
	f64.AxpyUnitary(1.0, b.data, d.data)
}

func (d *Dense) Sub(other Matrix) Matrix {
	//if !(EqualDims(d, other) ||
	//	(other.Columns() == 1 && other.Rows() == d.Rows()) ||
	//	(other.IsVector() && d.IsVector() && other.Size() == d.Size())) {
	//	panic("mat: matrices with not compatible size")
	//}
	out := d.ZerosLike().(*Dense)
	b := other.(*Dense)
	f64.AxpyUnitaryTo(out.data, -1.0, b.data, d.data)
	return out
}

func (d *Dense) SubInPlace(other Matrix) {
	//if !(EqualDims(d, other) ||
	//	(other.Columns() == 1 && other.Rows() == d.Rows()) ||
	//	(other.IsVector() && d.IsVector() && other.Size() == d.Size())) {
	//	panic("mat: matrices with not compatible size")
	//}
	switch other := other.(type) {
	case *Dense:
		f64.AxpyUnitary(-1.0, other.data, d.data)
	case *Sparse:
		other.DoNonZero(func(i, j int, k float64) {
			d.Set(d.At(i, j)-k, i, j)
		})
	}
}

func (d *Dense) Prod(other Matrix) Matrix {
	//if !(EqualDims(d, other) ||
	//	(other.Columns() == 1 && other.Rows() == d.Rows()) ||
	//	(other.IsVector() && d.IsVector() && other.Size() == d.Size())) {
	//	panic("mat: matrices with not compatible size")
	//}
	out := d.ZerosLike().(*Dense)
	b := other.(*Dense)
	for i, val := range d.data {
		out.data[i] = val * b.data[i]
	}
	return out
}

func (d *Dense) ProdInPlace(other Matrix) {
	//if !(EqualDims(d, other) ||
	//	(other.Columns() == 1 && other.Rows() == d.Rows()) ||
	//	(other.IsVector() && d.IsVector() && other.Size() == d.Size())) {
	//	panic("mat: matrices with not compatible size")
	//}
	b := other.(*Dense)
	for i, val := range b.data {
		d.data[i] *= val
	}
}

func (d *Dense) DivInPlace(other Matrix) {
	//if !(EqualDims(d, other) ||
	//	(other.Columns() == 1 && other.Rows() == d.Rows()) ||
	//	(other.IsVector() && d.IsVector() && other.Size() == d.Size())) {
	//	panic("mat: matrices with not compatible size")
	//}
	b := other.(*Dense)
	for i, val := range b.data {
		d.data[i] *= 1.0 / val
	}
}

func (d *Dense) Div(other Matrix) Matrix {
	//if !(EqualDims(d, other) ||
	//	(other.Columns() == 1 && other.Rows() == d.Rows()) ||
	//	(other.IsVector() && d.IsVector() && other.Size() == d.Size())) {
	//	panic("mat: matrices with not compatible size")
	//}
	out := d.ZerosLike().(*Dense)
	f64.DivTo(out.data, d.data, other.(*Dense).data)
	return out
}

func (d *Dense) Mul(other Matrix) Matrix {
	if d.Columns() != other.Rows() {
		panic("mat: matrices with not compatible size")
	}
	out := NewEmptyDense(d.Rows(), other.Columns())

	switch b := other.(type) {
	case *Dense:
		if out.cols == 1 {
			f64.GemvN(
				uintptr(d.rows), // m
				uintptr(d.cols), // n
				1.0,             // alpha
				d.data,          // a
				uintptr(d.cols), // lda
				b.data,          // x
				1.0,             // incX
				0.0,             // beta
				out.data,        // y
				1.0,             // incY
			)
		} else {
			f64.DgemmSerial(
				false,
				false,
				d.rows,   // m
				b.cols,   // n
				d.cols,   // k
				d.data,   // a
				d.cols,   // lda
				b.data,   // b
				b.cols,   // ldb
				out.data, // c
				out.cols, // ldc
				1.0,      // alpha
			)

			/*
				// parallel implementation
				f64.Dgemm(
					false,    // aTrans
					false,    // bTrans
					d.rows,   // m
					b.cols,   // n
					d.cols,   // k
					1.0,      // alpha
					d.data,   // a
					d.cols,   // lda
					b.data,   // b
					b.cols,   // ldb
					0.0,      // beta
					out.data, // c
					out.cols, // ldc
				)

			*/

		}

		return out

	case *Sparse:
		b.DoNonZero(func(k, j int, v float64) {
			for i := 0; i < d.Rows(); i++ {
				out.Set(out.At(i, j)+d.At(i, k)*v, i, j)
			}
		})
	}
	return out
}

func (d *Dense) MulT(other Matrix) Matrix {
	//if d.Columns() != other.Rows() {
	//	panic("mat: matrices with not compatible size")
	//}
	out := NewEmptyDense(d.Columns(), other.Columns())

	switch b := other.(type) {
	case *Dense:
		if out.cols == 1 {
			f64.GemvT(
				uintptr(d.rows), // m
				uintptr(d.cols), // n
				1.0,             // alpha
				d.data,          // a
				uintptr(d.cols), // lda
				b.data,          // x
				1.0,             // incX
				0.0,             // beta
				out.data,        // y
				1.0,             // incY
			)
		} else {
			panic("mat: matrices with not compatible size")
		}
		return out
	case *Sparse:
		panic("mat: matrices not compatible")
	}
	return out
}

func (d *Dense) DotUnitary(other Matrix) float64 {
	return f64.DotUnitary(d.data, other.Data())
}

// Clip performs the clip in place.
func (d *Dense) Clip(min, max float64) {
	d.Apply(func(i, j int, v float64) float64 {
		if v < min {
			return min
		} else if v > max {
			return max
		} else {
			return v
		}
	}, d)
}

// Abs returns a new matrix applying the abs function to all elements.
func (d *Dense) Abs() Matrix {
	out := d.ZerosLike()
	out.Apply(func(i, j int, v float64) float64 {
		return math.Abs(v)
	}, d)
	return out
}

// Pow returns a new matrix applying the pow function to all elements.
func (d *Dense) Pow(power float64) Matrix {
	out := d.Clone().(*Dense)
	out.Apply(func(i, j int, v float64) float64 {
		return math.Pow(v, power)
	}, d)
	return out
}

// Sqrt returns a new matrix applying the sqrt function to all elements.
func (d *Dense) Sqrt() Matrix {
	out := d.ZerosLike().(*Dense)
	for i, val := range d.data {
		out.data[i] = math.Sqrt(val)
	}
	return out
}

// Sum returns the sum of all values of the matrix.
func (d *Dense) Sum() float64 {
	sum := 0.0
	for _, elem := range d.data {
		sum += elem
	}
	return sum
}

// Max returns the max value of the matrix.
func (d *Dense) Max() float64 {
	r, c := d.Dims()
	max := math.Inf(-1)
	for i := 0; i < r; i++ {
		for j := 0; j < c; j++ {
			v := d.At(i, j)
			if v > max {
				max = v
			}
		}
	}
	return max
}

// Min returns the min value of the matrix.
func (d *Dense) Min() float64 {
	r, c := d.Dims()
	min := math.Inf(1)
	for i := 0; i < r; i++ {
		for j := 0; j < c; j++ {
			v := d.At(i, j)
			if v < min {
				min = v
			}
		}
	}
	return min
}

// start (inclusive) end (exclusive)
// TODO
func (d *Dense) Range(start, end int) Matrix {
	data := make([]float64, end-start)
	for k := 0; k < end-start; k++ {
		data[k] = d.At(start + k)
	}
	return NewVecDense(data)
}

func (d *Dense) SplitV(sizes ...int) []Matrix {
	//if d.Columns() != 1 {
	//	panic("mat: required vector, found matrix.")
	//}
	//if d.Size() != SumInt(sizes) {
	//	panic("mat: incompatible split sizes.")
	//}
	out := make([]Matrix, len(sizes))
	offset := 0
	for i := 0; i < len(sizes); i++ {
		startIndex := offset
		offset = startIndex + sizes[i]
		out[i] = d.Range(startIndex, offset)
	}
	return out
}

// Norm returns the vector norm.  Use pow = 2.0 for Euclidean.
func (d *Dense) Norm(pow float64) float64 {
	s := 0.0
	for _, x := range d.data {
		s += math.Pow(x, pow)
	}
	return math.Pow(s, 1/pow)
}

// String returns the string representation of the data.
func (d *Dense) String() string {
	return fmt.Sprintf("%v", d.data)
}