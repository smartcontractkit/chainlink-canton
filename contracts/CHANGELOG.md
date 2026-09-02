# Changelog

## [2.1.0](https://github.com/smartcontractkit/chainlink-canton/compare/contracts/v2.0.0...contracts/v2.1.0) (2026-09-02)


### Features

* **contracts:** Add PerPartyRouter InstanceId namespace ([#913](https://github.com/smartcontractkit/chainlink-canton/issues/913)) ([d3cff3b](https://github.com/smartcontractkit/chainlink-canton/commit/d3cff3bccaf5fe37b3fef0bfae91e7e754b008af))
* **contracts:** Deprecate BurnMintTokenPool transferTimeout ([#969](https://github.com/smartcontractkit/chainlink-canton/issues/969)) ([8662e97](https://github.com/smartcontractkit/chainlink-canton/commit/8662e97dae0987e5a85e3fc7a11b2b859fa64eee))
* **contracts:** Deprecate PerPartyRouter.feeTransferLifetime and use global lifetime of 14 days instead ([#947](https://github.com/smartcontractkit/chainlink-canton/issues/947)) ([9b53053](https://github.com/smartcontractkit/chainlink-canton/commit/9b53053c84ce6ddcbd3341ab2db86fc4d10f52ba))
* **contracts:** Deprecate TokenAdminRegistry entryCount field ([#941](https://github.com/smartcontractkit/chainlink-canton/issues/941)) ([46ab0cf](https://github.com/smartcontractkit/chainlink-canton/commit/46ab0cf037b1fe9b304800e1dd9a717771c77d12))
* **contracts:** Handle fractional seconds in RateLimiter ([#923](https://github.com/smartcontractkit/chainlink-canton/issues/923)) ([006f911](https://github.com/smartcontractkit/chainlink-canton/commit/006f911d9e7ce1b6b00417a87d27bdf0b7e43f3d))


### Bug Fixes

* **contracts:** Bind SetRateLimitConfig to the signed limiter address ([#915](https://github.com/smartcontractkit/chainlink-canton/issues/915)) ([19e20d1](https://github.com/smartcontractkit/chainlink-canton/commit/19e20d15c635b729af4b5a575e448bda3accc0da))
* **contracts:** LnR Pool - require transferTimeout to be non-zero ([#970](https://github.com/smartcontractkit/chainlink-canton/issues/970)) ([c4d1d35](https://github.com/smartcontractkit/chainlink-canton/commit/c4d1d3535f22ae98ebe803c1552bcc3d92eb7062))
* **contracts:** Remove PublicFetch calls from TAR SetTransferFactory & SetBurnMintFactory ([#926](https://github.com/smartcontractkit/chainlink-canton/issues/926)) ([87d9824](https://github.com/smartcontractkit/chainlink-canton/commit/87d982474ce84458c3ad91557b9be8ce2f8bb3c7))

## 2.0.0 (2026-08-24)


### ⚠ BREAKING CHANGES

* update /contracts to v2 ([#919](https://github.com/smartcontractkit/chainlink-canton/issues/919))

### Builds

* Update /contracts to v2 ([#919](https://github.com/smartcontractkit/chainlink-canton/issues/919)) ([48860c6](https://github.com/smartcontractkit/chainlink-canton/commit/48860c64ab3f55549cc220f3df8d05445c8445ec))


### Features

* Release 2.0.0 ([#917](https://github.com/smartcontractkit/chainlink-canton/issues/917)) ([72dc0ee](https://github.com/smartcontractkit/chainlink-canton/commit/72dc0ee1dee3bb7e0fa7e4eab4dc71cdc4d1fe43))
