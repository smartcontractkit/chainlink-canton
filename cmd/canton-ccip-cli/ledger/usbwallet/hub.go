// Heavily adjusted from go-ethereum's usbwallet implementation:
// https://github.com/ethereum/go-ethereum/blob/b20778d77d50e2c9775cb4d2fa872579989c5421/accounts/usbwallet/hub.go

package usbwallet

import (
	"errors"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/log"
	"github.com/karalabe/hid"
)

// CantonLedgerScheme is the protocol scheme prefixing account and wallet URLs.
const CantonLedgerScheme = "cantonledger"

// refreshThrottling is the minimum time between wallet refreshes to avoid USB
// trashing.
const refreshThrottling = 500 * time.Millisecond

const (
	// deviceUsagePage identifies Ledger devices by HID usage page (0xffa0) on Windows and macOS.
	// See: https://github.com/LedgerHQ/ledger-live/blob/05a2980e838955a11a1418da638ef8ac3df4fb74/libs/ledgerjs/packages/hw-transport-node-hid-noevents/src/TransportNodeHid.ts
	deviceUsagePage = 0xffa0
	// deviceInterface identifies Ledger devices by USB interface number (0) on Linux.
	deviceInterface = 0
)

// Hub is a accounts.Backend that can find and handle generic USB hardware wallets.
type Hub struct {
	scheme     string   // Protocol scheme prefixing account and wallet URLs.
	vendorID   uint16   // USB vendor identifier used for device discovery
	productIDs []uint16 // USB product identifiers used for device discovery
	usageID    uint16   // USB usage page identifier used for macOS device discovery
	endpointID int      // USB endpoint identifier used for non-macOS device discovery

	refreshed time.Time      // Time instance when the list of wallets was last refreshed
	wallets   []CantonWallet // List of USB wallet devices currently tracking

	quit chan chan error

	stateLock sync.RWMutex // Protects the internals of the hub from racey access

	commsPend int           // Number of operations blocking enumeration
	commsLock sync.Mutex    // Lock protecting the pending counter and enumeration
	enumFails atomic.Uint32 // Number of times enumeration has failed
}

// NewCantonLedgerHub creates a new hardware wallet manager for Ledger devices.
func NewCantonLedgerHub() (*Hub, error) {
	if !hid.Supported() {
		return nil, errors.New("unsupported platform")
	}
	hub := &Hub{
		scheme:   CantonLedgerScheme,
		vendorID: 0x2c97,
		productIDs: []uint16{
			// Device definitions taken from
			// https://github.com/LedgerHQ/ledger-live/blob/595cb73b7e6622dbbcfc11867082ddc886f1bf01/libs/ledgerjs/packages/devices/src/index.ts

			// Original product IDs
			0x0000, /* Ledger Blue */
			0x0001, /* Ledger Nano S */
			0x0004, /* Ledger Nano X */
			0x0005, /* Ledger Nano S Plus */
			0x0006, /* Ledger Nano FTS */
			0x0007, /* Ledger Flex */
			0x0008, /* Ledger Nano Gen5 */

			0x0000, /* WebUSB Ledger Blue */
			0x1000, /* WebUSB Ledger Nano S */
			0x4000, /* WebUSB Ledger Nano X */
			0x5000, /* WebUSB Ledger Nano S Plus */
			0x6000, /* WebUSB Ledger Nano FTS */
			0x7000, /* WebUSB Ledger Flex */
			0x8000, /* WebUSB Ledger Nano Gen5 */
		},
		usageID:    deviceUsagePage,
		endpointID: deviceInterface,
		quit:       make(chan chan error),
	}
	hub.refreshWallets()

	return hub, nil
}

// Wallets returns all the currently tracked USB devices that appear to be hardware wallets.
func (hub *Hub) Wallets() []CantonWallet {
	// Make sure the list of wallets is up to date
	hub.refreshWallets()

	hub.stateLock.RLock()
	defer hub.stateLock.RUnlock()

	cpy := make([]CantonWallet, len(hub.wallets))
	copy(cpy, hub.wallets)

	return cpy
}

// refreshWallets scans the USB devices attached to the machine and updates the
// list of wallets based on the found devices.
func (hub *Hub) refreshWallets() {
	// Don't scan the USB like crazy it the user fetches wallets in a loop
	hub.stateLock.RLock()
	elapsed := time.Since(hub.refreshed)
	hub.stateLock.RUnlock()

	if elapsed < refreshThrottling {
		return
	}
	// If USB enumeration is continually failing, don't keep trying indefinitely
	if hub.enumFails.Load() > 2 {
		return
	}
	// Retrieve the current list of USB wallet devices
	var devices []hid.DeviceInfo

	if runtime.GOOS == "linux" {
		// hidapi on Linux opens the device during enumeration to retrieve some infos,
		// breaking the Ledger protocol if that is waiting for user confirmation. This
		// is a bug acknowledged at Ledger, but it won't be fixed on old devices so we
		// need to prevent concurrent comms ourselves. The more elegant solution would
		// be to ditch enumeration in favor of hotplug events, but that don't work yet
		// on Windows so if we need to hack it anyway, this is more elegant for now.
		hub.commsLock.Lock()
		if hub.commsPend > 0 { // A confirmation is pending, don't refresh
			hub.commsLock.Unlock()
			return
		}
	}
	infos, err := hid.Enumerate(hub.vendorID, 0)
	if err != nil {
		failcount := hub.enumFails.Add(1)
		if runtime.GOOS == "linux" {
			// See rationale before the enumeration why this is needed and only on Linux.
			hub.commsLock.Unlock()
		}
		log.Error("Failed to enumerate USB devices", "hub", hub.scheme,
			"vendor", hub.vendorID, "failcount", failcount, "err", err)

		return
	}
	hub.enumFails.Store(0)

	for _, info := range infos {
		for _, id := range hub.productIDs {
			// We check both the raw ProductID (legacy) and just the upper byte, as Ledger
			// uses `MMII`, encoding a model (MM) and an interface bitfield (II)
			mmOnly := info.ProductID & 0xff00
			// Windows and Macos use UsageID matching, Linux uses Interface matching
			if (info.ProductID == id || mmOnly == id) && (info.UsagePage == hub.usageID || info.Interface == hub.endpointID) {
				devices = append(devices, info)
				break
			}
		}
	}
	if runtime.GOOS == "linux" {
		// See rationale before the enumeration why this is needed and only on Linux.
		hub.commsLock.Unlock()
	}
	// Transform the current list of wallets into the new one
	hub.stateLock.Lock()
	defer hub.stateLock.Unlock()

	wallets := make([]CantonWallet, 0, len(devices))

	for _, device := range devices {
		url := accounts.URL{Scheme: hub.scheme, Path: device.Path}

		// Drop wallets in front of the next device or those that failed for some reason
		for len(hub.wallets) > 0 {
			// Abort if we're past the current device and found an operational one
			_, failure := hub.wallets[0].Status()
			if hub.wallets[0].URL().Cmp(url) >= 0 || failure == nil {
				break
			}
			// Drop the stale and failed devices
			hub.wallets = hub.wallets[1:]
		}
		// If there are no more wallets or the device is before the next, wrap new wallet
		if len(hub.wallets) == 0 || hub.wallets[0].URL().Cmp(url) > 0 {
			logger := log.New("url", url)
			wallet := &cantonWallet{
				log:    logger,
				hub:    hub,
				url:    &url,
				info:   device,
				device: nil,
				mux:    sync.Mutex{},
			}

			wallets = append(wallets, wallet)

			continue
		}
		// If the device is the same as the first wallet, keep it
		if hub.wallets[0].URL().Cmp(url) == 0 {
			wallets = append(wallets, hub.wallets[0])
			hub.wallets = hub.wallets[1:]
			continue
		}
	}
	// Drop any leftover wallets and set the new batch
	hub.refreshed = time.Now()
	hub.wallets = wallets
}
