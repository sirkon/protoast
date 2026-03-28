package core

// Double represents double builtin type.
type Double struct{ isComparableType }

// Float represents float builtin type.
type Float struct{ isComparableType }

// Int32 represents int32 builtin type.
type Int32 struct{ isComparableType }

// Int64 represents int64 builtin type.
type Int64 struct{ isComparableType }

// Uint32 represents uint32 builtin type.
type Uint32 struct{ isComparableType }

// Uint64 represents uint64 builtin type.
type Uint64 struct{ isComparableType }

// Sint32 represents sint32 builtin type.
type Sint32 struct{ isComparableType }

// Sint64 represents sint64 builtin type.
type Sint64 struct{ isComparableType }

// Fixed32 represents fixed32 builtin type.
type Fixed32 struct{ isComparableType }

// Fixed64 represents fixed64 builtin type.
type Fixed64 struct{ isComparableType }

// Sfixed32 represents sfixed32 builtin type.
type Sfixed32 struct{ isComparableType }

// Sfixed64 represents sfixed64 builtin type.
type Sfixed64 struct{ isComparableType }

// Bool represents bool builtin type.
type Bool struct{ isComparableType }

// String represents string builtin type.
type String struct{ isComparableType }

// Bytes represents bytes builtin type.
type Bytes struct{ isBuiltinType }

func (Double) String() string   { return "double" }
func (Float) String() string    { return "float" }
func (Int32) String() string    { return "int32" }
func (Int64) String() string    { return "int64" }
func (Uint32) String() string   { return "uint32" }
func (Uint64) String() string   { return "uint64" }
func (Sint32) String() string   { return "sint32" }
func (Sint64) String() string   { return "sint64" }
func (Fixed32) String() string  { return "fixed32" }
func (Fixed64) String() string  { return "fixed64" }
func (Sfixed32) String() string { return "sfixed32" }
func (Sfixed64) String() string { return "sfixed64" }
func (Bool) String() string     { return "bool" }
func (String) String() string   { return "string" }
func (Bytes) String() string    { return "bytes" }
