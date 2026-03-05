package model

import (
	"strings"
)

type ExternalPackages struct {
	// Maps packageId => package
	Packages map[string]ExternalPackage
}

type ExternalPackage struct {
	// The import to use to refer to this package, e.g. github.com/smartcontractkit/go-daml/etc
	Import string
	// The alias to import this package by
	Alias string
}

// Daml Types

type DamlType interface {
	GoType() string
	GoImport() *ExternalPackage
}

type noImport struct{}

func (noImport) GoImport() *ExternalPackage { return nil }

type List struct {
	Inner DamlType
}

func (t List) GoType() string {
	return "[]" + t.Inner.GoType()
}

func (t List) GoImport() *ExternalPackage {
	return t.Inner.GoImport()
}

type Party struct {
	noImport
}

func (t Party) GoType() string {
	return "types.PARTY"
}

type Text struct {
	noImport
}

func (t Text) GoType() string {
	return "types.TEXT"
}

type Int64 struct {
	noImport
}

func (t Int64) GoType() string {
	return "types.INT64"
}

type Bool struct {
	noImport
}

func (t Bool) GoType() string {
	return "types.BOOL"
}

type Decimal struct {
	noImport
}

func (t Decimal) GoType() string {
	return "types.DECIMAL"
}

type Numeric struct {
	noImport
}

func (t Numeric) GoType() string {
	return "types.NUMERIC"
}

type Date struct {
	noImport
}

func (t Date) GoType() string {
	return "types.DATE"
}

type Timestamp struct {
	noImport
}

func (t Timestamp) GoType() string {
	return "types.TIMESTAMP"
}

type Unit struct {
	noImport
}

func (t Unit) GoType() string {
	return "types.UNIT"
}

type Map struct {
	noImport
}

func (t Map) GoType() string {
	return "types.MAP"
}

type Optional struct {
	Inner DamlType
}

func (t Optional) GoType() string {
	return "*" + t.Inner.GoType()
}

func (t Optional) GoImport() *ExternalPackage {
	return t.Inner.GoImport()
}

type ContractId struct {
	noImport
}

func (t ContractId) GoType() string {
	return "types.CONTRACT_ID"
}

type GenMap struct {
	noImport
}

func (t GenMap) GoType() string {
	return "types.GENMAP"
}

type TextMap struct {
	noImport
}

func (t TextMap) GoType() string {
	return "types.TEXTMAP"
}

type BigNumeric struct {
	noImport
}

func (t BigNumeric) GoType() string {
	return "types.BIGNUMERIC"
}

type RoundingMode struct {
	noImport
}

func (t RoundingMode) GoType() string {
	return "types.ROUNDING_MODE"
}

type Any struct {
	noImport
}

func (t Any) GoType() string {
	return "any"
}

type RelTime struct {
	noImport
}

func (t RelTime) GoType() string {
	return "types.RELTIME"
}

type Set struct {
	noImport
}

func (t Set) GoType() string {
	return "SET"
}

type Enum struct {
	noImport
}

func (t Enum) GoType() string {
	return "string"
}

type Imported struct {
	Underlying      DamlType
	ExternalPackage ExternalPackage
}

func (t Imported) GoType() string {
	return t.ExternalPackage.Alias + "." + t.Underlying.GoType()
}

func (t Imported) GoImport() *ExternalPackage {
	return &t.ExternalPackage
}

type Unknown struct {
	String string
	noImport
}

func (t Unknown) GoType() string {
	// Retain previous behavior of stripping all underscores, matched what the `capitalise` function in the template does.
	return strings.ReplaceAll(t.String, "_", "")
}
