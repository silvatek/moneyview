package main

import "fmt"

type Date struct {
	value    string
	original string
}

func DateFrom(s string) Date {
	date := Date{original: s}
	var day int
	var month int
	var year int
	fmt.Sscanf(s, "%2d/%2d/%4d", &day, &month, &year)
	date.value = fmt.Sprintf("%04d-%02d-%02d", year, month, day)
	return date
}

func (d Date) String() string {
	return d.value
}

func (d0 *Date) NotBefore(d1 Date) bool {
	return !(d0.value < d1.value)
}

func (d0 *Date) NotAfter(d1 Date) bool {
	return !(d0.value > d1.value)
}

func (d0 *Date) InRange(d1 Date, d2 Date) bool {
	return d0.NotBefore(d1) && d0.NotAfter(d2)
}
