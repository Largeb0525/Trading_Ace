package eth

import (
	"sync"

	"github.com/ethereum/go-ethereum/common"
)

var (
	addresses     = make(map[common.Address]struct{})
	addressesLock sync.Mutex
	notifyChannel = make(chan struct{}, 1)
)

func AddAddresses(addrs []string) {
	addressesLock.Lock()
	defer addressesLock.Unlock()

	for _, addr := range addrs {
		addresses[common.HexToAddress(addr)] = struct{}{}
	}
	notifyChange()
}

func RemoveAddresses(addrs []string) {
	addressesLock.Lock()
	defer addressesLock.Unlock()

	for _, addr := range addrs {
		delete(addresses, common.HexToAddress(addr))
	}
	notifyChange()
}

func GetAddresses() []common.Address {
	addressesLock.Lock()
	defer addressesLock.Unlock()
	var addrList []common.Address
	for addr := range addresses {
		addrList = append(addrList, addr)
	}
	return addrList
}

func notifyChange() {
	select {
	case notifyChannel <- struct{}{}:
	default:
	}
}

func GetNotifyChannel() <-chan struct{} {
	return notifyChannel
}
