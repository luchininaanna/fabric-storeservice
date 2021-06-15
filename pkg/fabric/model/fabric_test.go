package model

import (
	"github.com/google/uuid"
	"testing"
)

func TestCreateFabricWithEmptyName(t *testing.T) {
	_, err := NewFabric(uuid.New(), "", 77, 5)
	if err != FabricWithoutNameError {
		t.Error("Create fabric without name")
	}
}

func TestCreateFabricWithBelowZeroQuantity(t *testing.T) {
	_, err := NewFabric(uuid.New(), "", 77, -5)
	if err != FabricWithoutNameError {
		t.Error("Create fabric with below zero quantity")
	}
}

func TestCreateOrderItemWithZeroQuantity(t *testing.T) {
	_, err := NewFabric(uuid.New(), "", 77, 0)
	if err != FabricWithoutNameError {
		t.Error("Create fabric with zero quantity")
	}
}

func TestCreateFabricWithBelowZeroCost(t *testing.T) {
	_, err := NewFabric(uuid.New(), "", -77, 5)
	if err != FabricWithoutNameError {
		t.Error("Create fabric with below zero cost")
	}
}

func TestCreateOrderItemWithZeroCost(t *testing.T) {
	_, err := NewFabric(uuid.New(), "cotton", 0, 5)
	if err != InvalidFabricCostError {
		t.Error("Create fabric with zero cost")
	}
}
