package controllers

import (
	"dksv/models"
	"fmt"
	"os"
	"testing"
)

func TestNetworkController_Create(t *testing.T) {
	f, err := os.Open(models.DefaultNetworkPath + "bridxxx")
	fmt.Println(f, err)
}
