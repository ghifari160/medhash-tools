// MedHash Tools
// Copyright (c) 2023 GHIFARI160
// MIT License

package main

import "fmt"

const (
	B   Size = 1
	KiB Size = 1024 * B
	MiB Size = 1024 * KiB
	GiB Size = 1024 * MiB
	TiB Size = 1024 * GiB
	PiB Size = 1024 * TiB
	EiB Size = 1024 * PiB
)

type Size int64

func (s Size) String() string {
	var operand float64
	var unit string

	if s >= EiB {
		operand = s.EiB()
		unit = "EiB"
	} else if s >= PiB {
		operand = s.PiB()
		unit = "PiB"
	} else if s >= TiB {
		operand = s.TiB()
		unit = "TiB"
	} else if s >= GiB {
		operand = s.GiB()
		unit = "GiB"
	} else if s >= MiB {
		operand = s.MiB()
		unit = "MiB"
	} else if s >= KiB {
		operand = s.KiB()
		unit = "KiB"
	} else if s > B {
		operand = s.B()
		unit = "iB"
	}

	return fmt.Sprintf("%.4f %s", operand, unit)
}

func (s Size) B() float64 {
	return float64(s) / float64(B)
}

func (s Size) KiB() float64 {
	return float64(s) / float64(KiB)
}

func (s Size) MiB() float64 {
	return float64(s) / float64(MiB)
}

func (s Size) GiB() float64 {
	return float64(s) / float64(GiB)
}

func (s Size) TiB() float64 {
	return float64(s) / float64(TiB)
}

func (s Size) PiB() float64 {
	return float64(s) / float64(PiB)
}

func (s Size) EiB() float64 {
	return float64(s) / float64(EiB)
}
