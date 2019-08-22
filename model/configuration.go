package model

import (
	"github.com/brunoksato/golang-boilerplate/core"
)

type Configuration struct {
	Model
	MinValueBuy float64 `json:"min_value_buy"`
}

func (c Configuration) ValidateForCreate() core.DefaultError {
	err := ValidateStruct(c)
	if err != nil {
		return err
	}

	return nil
}

func (c Configuration) ValidateForUpdate() core.DefaultError {
	err := ValidateStruct(c)
	if err != nil {
		return err
	}
	return nil
}

func (c Configuration) ValidateForDelete(ctx *ModelCtx) core.DefaultError {
	return nil
}

func (c Configuration) ValidateField(f string) core.DefaultError {
	err := ValidateStructField(c, f)

	return err
}

// Restrictor

func (c Configuration) UserCanView(ctx *ModelCtx, viewer User) (bool, core.DefaultError) {
	return viewer.IsAdmin(), nil
}

func (c Configuration) UserCanCreate(ctx *ModelCtx, creator User) (bool, core.DefaultError) {
	return creator.IsAdmin(), nil
}

func (c Configuration) UserCanUpdate(ctx *ModelCtx, updater User, fields []string) (bool, core.DefaultError) {
	return updater.IsAdmin(), nil
}

func (c Configuration) UserCanDelete(ctx *ModelCtx, deleter User) (bool, core.DefaultError) {
	return deleter.IsAdmin(), nil
}

// Business methods

// Scopes
