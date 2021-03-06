package vectors

import (
	"fmt"
	"reflect"
)

var Bool = reflect.TypeOf(true)
var Int = reflect.TypeOf(int(1))
var Int8 = reflect.TypeOf(int8(1))
var Int16 = reflect.TypeOf(int16(1))
var Int32 = reflect.TypeOf(int32(1))
var Int64 = reflect.TypeOf(int64(1))
var Uint = reflect.TypeOf(uint(1))
var Uint8 = reflect.TypeOf(uint8(1))
var Uint16 = reflect.TypeOf(uint16(1))
var Uint32 = reflect.TypeOf(uint32(1))
var Uint64 = reflect.TypeOf(uint64(1))
var Uintptr = reflect.TypeOf(uintptr(1))
var Float32 = reflect.TypeOf(float32(1.0))
var Float64 = reflect.TypeOf(float64(1))
var Complex64 = reflect.TypeOf(complex64(1))
var Complex128 = reflect.TypeOf(complex128(1))

//var Array
//var Chan
//var Func
//var Interface
//var Map
//var Ptr
//var Slice
var String = reflect.TypeOf(string(""))

//var Struct
//var UnsafePointer

type Vector struct {
	slice  reflect.Value
	typeof reflect.Type
}

//NewVector Creates a Vector of Type T and returns it
func NewVector(T reflect.Type) *Vector {

	return &Vector{
		slice:  reflect.MakeSlice(reflect.SliceOf(T), 0, 0),
		typeof: T,
	}
}

func newVector(t reflect.Type, len, cap int) *Vector {
	return &Vector{
		slice:  reflect.MakeSlice(reflect.SliceOf(t), len, cap),
		typeof: t,
	}
}

//Get Returns the value from the Index in the Vector
func (v *Vector) Get(index int) interface{} {
	return v.slice.Index(index)
}

//Put Sets multiple element in the vector
func (v *Vector) Put(elements ...interface{}) {

	for i := range elements {
		if reflect.ValueOf(elements[i]).Type() != v.slice.Type().Elem() {
			panic(fmt.Sprintf("Put: cannot put a %T into a vector of %s", elements[i], v.slice.Type().Elem()))
		}

		v.slice = reflect.Append(v.slice, reflect.ValueOf(elements[i]))
	}
}

//PutFront Sets multiple element in the front of the vector
func (v *Vector) PutFront(elements ...interface{}) {

	v2 := newVector(v.typeof, 0, 0)

	for i := range elements {
		if reflect.ValueOf(elements[i]).Type() != v.slice.Type().Elem() {
			panic(fmt.Sprintf("Put: cannot put a %T into a vector of %s", elements[i], v.slice.Type().Elem()))
		}

		v2.slice = reflect.Append(v2.slice, reflect.ValueOf(elements[i]))
	}

	v.slice = reflect.Append(v2.slice, v.slice)
}

//Copy Clones an entire Vector and returns it
func (v *Vector) Copy() *Vector {
	v2 := newVector(v.typeof, v.slice.Len(), v.slice.Cap())
	reflect.Copy(v2.slice, v.slice)
	return v2
}

//Cut Removes a section or slice from the Vector
func (v *Vector) Cut(i, j int) {
	lastItem := v.slice.Len()
	cutLen := j - i

	reflect.Copy(v.slice.Slice(i, lastItem), v.slice.Slice(j, lastItem))

	for n := v.slice.Len() - cutLen; n < v.slice.Len(); n++ {
		v.slice.Index(n).Set(reflect.Zero(v.typeof))
	}

	v.slice = v.slice.Slice(0, v.slice.Len()-cutLen)
}

//Delete Removes a single index from the vector
func (v *Vector) Delete(i int) {
	reflect.Copy(v.slice.Slice(i, v.slice.Len()), v.slice.Slice(i+1, v.slice.Len()))
	v.slice.Index(v.slice.Len() - 1).Set(reflect.Zero(v.typeof))
	v.slice = v.slice.Slice(0, v.slice.Len()-1)
}

//DeleteNoPreserveOrder Removes a single index from the vector without preserving order
func (v *Vector) DeleteNoPreserveOrder(i int) {
	v.slice.Index(i).Set(v.slice.Index(v.slice.Len() - 1))
	v.slice.Index(v.slice.Len() - 1).Set(reflect.Zero(v.typeof))
	v.slice = v.slice.Slice(0, v.slice.Len()-1)
}

//Expand Increases the size of the vector at the offset with the amount of indexes
func (v *Vector) Expand(offset, indexes int) {
	// Zeroed Out, Expander
	v2 := newVector(v.typeof, indexes, indexes)

	// Empty Vector
	v3 := newVector(v.typeof, 0, 0)

	//Before Offset
	bef := v.slice.Slice(0, offset)

	//After Offset
	aft := v.slice.Slice(offset, v.slice.Len())

	// Expand Operation
	v.slice = reflect.AppendSlice(v3.slice, bef)
	v.slice = reflect.AppendSlice(v.slice, v2.slice)
	v.slice = reflect.AppendSlice(v.slice, aft)
}

//Extend Increases the size of the vector by placing new indexes at the end
func (v *Vector) Extend(indexes int) {
	v.slice = reflect.Append(v.slice, newVector(v.typeof, indexes, indexes).slice)
}

//Insert Sets a element in the vector at the offset
func (v *Vector) Insert(offset int, element interface{}) {
	if reflect.ValueOf(element).Type() != v.slice.Type().Elem() {
		panic(fmt.Sprintf("Insert: cannot insert a %T into a vector of %s", element, v.slice.Type().Elem()))
	}

	v.slice = reflect.Append(v.slice, reflect.ValueOf(0))
	reflect.Copy(v.slice.Slice(offset+1, v.slice.Len()), v.slice.Slice(offset, v.slice.Len()))
	v.slice.Index(offset).Set(reflect.ValueOf(element))
}

//InsertVector Sets a vector in the vector at the offset
func (v *Vector) InsertVector(offset int, vec *Vector) {
	if vec.typeof != v.slice.Type().Elem() {
		panic(fmt.Sprintf("InsertVector: cannot insert a %T vector into a vector of %s", vec.slice.Interface(), v.slice.Type().Elem()))
	}

	v.slice = reflect.AppendSlice(v.slice.Slice(0, offset), reflect.AppendSlice(vec.slice, v.slice.Slice(offset, v.slice.Len())))
}

//Pop Removes the first element from a vector and returns it
func (v *Vector) Pop() interface{} {
	var x reflect.Value
	x, v.slice = v.slice.Index(0), v.slice.Slice(1, v.slice.Len())
	return x.Interface()
}

//PopBack Removes the last element from a vector and returns it
func (v *Vector) PopBack() interface{} {
	var x reflect.Value
	x, v.slice = v.slice.Index(v.slice.Len()-1), v.slice.Slice(0, v.slice.Len()-1)
	return x.Interface()
}

//PopOut Removes the specified element in the index from a vector and returns it
func (v *Vector) PopOut(i int) interface{} {
	x := v.slice.Index(i).Interface()
	reflect.Copy(v.slice.Slice(i, v.slice.Len()), v.slice.Slice(i+1, v.slice.Len()))
	v.slice.Index(v.slice.Len() - 1).Set(reflect.Zero(v.typeof))
	v.slice = v.slice.Slice(0, v.slice.Len()-1)
	return x
}

//Push Sets an element to the back of a vector
func (v *Vector) Push(element interface{}) {

	if reflect.ValueOf(element).Type() != v.slice.Type().Elem() {
		panic(fmt.Sprintf("Put: cannot put a %T into a vector of %s", element, v.slice.Type().Elem()))
	}

	v.slice = reflect.Append(v.slice, reflect.ValueOf(element))
}

//PushFront Sets an element to the front of a vector
func (v *Vector) PushFront(element interface{}) {
	v2 := newVector(v.typeof, 0, 0)
	v2.Push(element)
	v.slice = reflect.Append(v2.slice, v.slice)
}

//Modify Changes the index of a vector with a new element
func (v *Vector) Modify(index int, element interface{}) {

	if reflect.ValueOf(element).Type() != v.slice.Type().Elem() {
		panic(fmt.Sprintf("Modify: cannot change a %T into a type of %s", element, v.slice.Type().Elem()))
	}

	v.slice.Index(index).Set(reflect.ValueOf(element))
}
