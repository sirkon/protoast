package core

import (
	"iter"

	"github.com/emicklei/proto"
)

// Reserved represents reserved construct in named types.
// TODO implement
type Reserved struct {
	isNode

	proto *proto.Reserved
}

func (r *Reserved) Values() iter.Seq[ReservedValue] {
	return func(yield func(ReservedValue) bool) {
		for _, rng := range r.proto.Ranges {
			var val ReservedValue
			switch {
			case rng.Max:
				val = &ReservedValueFrom{
					From: rng.From,
				}
			case rng.From == rng.To:
				val = &ReservedValueSingleNumber{
					Value: rng.From,
				}
			default:
				val = &ReservedValueRange{
					From: rng.From,
					To:   rng.To,
				}
			}
			if !yield(val) {
				return
			}
		}

		for _, name := range r.proto.FieldNames {
			val := &ReservedValueString{
				Value: name,
			}
			if !yield(val) {
				return
			}
		}
	}
}

type ReservedValue interface {
	isReserved()
}

type ReservedValueSingleNumber struct {
	Value int
}

type ReservedValueRange struct {
	From int
	To   int
}

type ReservedValueFrom struct {
	From int
}

type ReservedValueString struct {
	Value string
}

func (*ReservedValueSingleNumber) isReserved() {}
func (*ReservedValueRange) isReserved()        {}
func (*ReservedValueFrom) isReserved()         {}
func (*ReservedValueString) isReserved()       {}
