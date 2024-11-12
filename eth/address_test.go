package eth

import (
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

func TestAddAddresses(t *testing.T) {
	resetAddresses()

	addressesToAdd := []string{"0x123", "0x456", "0x789"}
	expected := []common.Address{
		common.HexToAddress("0x123"),
		common.HexToAddress("0x456"),
		common.HexToAddress("0x789"),
	}

	AddAddresses(addressesToAdd)

	result := GetAddresses()

	assert.ElementsMatch(t, result, expected, "Added addresses do not match expected addresses")
}

func TestRemoveAddresses(t *testing.T) {
	resetAddresses()

	addressesToAdd := []string{"0x123", "0x456", "0x789"}
	addressesToRemove := []string{"0x456"}

	AddAddresses(addressesToAdd)
	RemoveAddresses(addressesToRemove)

	expected := []common.Address{
		common.HexToAddress("0x123"),
		common.HexToAddress("0x789"),
	}

	result := GetAddresses()

	assert.ElementsMatch(t, result, expected, "Addresses after removal do not match expected addresses")
}

func TestNotifyChannel(t *testing.T) {
	resetAddresses()

	addressesToAdd := []string{"0x123"}
	done := make(chan struct{})

	go func() {
		<-GetNotifyChannel()
		done <- struct{}{}
	}()

	AddAddresses(addressesToAdd)

	select {
	case <-done:
	case <-time.After(1 * time.Second):
		t.Error("Expected notification, but channel was not notified")
	}
}

func resetAddresses() {
	addressesLock.Lock()
	defer addressesLock.Unlock()
	addresses = make(map[common.Address]struct{})
	select {
	case <-notifyChannel:
	default:
	}
}
