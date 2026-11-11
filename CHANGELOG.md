# Changelog

## [0.38.6](https://github.com/apeirora/showroom-msp-apeiro-hosting/compare/v0.38.5...v0.38.6) (2026-11-11)


### Bug Fixes

* **controller:** keep legacy SCPCs without account contract ([#69](https://github.com/apeirora/showroom-msp-apeiro-hosting/issues/69)) ([ef7ac59](https://github.com/apeirora/showroom-msp-apeiro-hosting/commit/ef7ac593fa0adde4c270d2c1ef4f35f715d995d8))
* **deps:** build with Go 1.26.11 ([#87](https://github.com/apeirora/showroom-msp-apeiro-hosting/issues/87)) ([099cbee](https://github.com/apeirora/showroom-msp-apeiro-hosting/commit/099cbeede92232a670dc74867f220ac155f731e7))
* **ironcore:** limit worker zones to the worker maximum ([c048258](https://github.com/apeirora/showroom-msp-apeiro-hosting/commit/c0482587f27e6f789770dc82e837ad36db44976d))
* **ui:** bump vite to 8.3.5 ([#72](https://github.com/apeirora/showroom-msp-apeiro-hosting/issues/72)) ([ad6629b](https://github.com/apeirora/showroom-msp-apeiro-hosting/commit/ad6629b99eeb9abe6bf1fa4bd3906810d6121131))

## [0.38.5](https://github.com/apeirora/showroom-msp-apeiro-hosting/compare/v0.38.4...v0.38.5) (2026-10-23)


### Bug Fixes

* **ui:** make openstack/ironcore shoot nodes navigable ([#67](https://github.com/apeirora/showroom-msp-apeiro-hosting/issues/67)) ([dfb6d69](https://github.com/apeirora/showroom-msp-apeiro-hosting/commit/dfb6d69b4696be92e645b943bfa69237078f3f9a))

## [0.38.4](https://github.com/apeirora/showroom-msp-apeiro-hosting/compare/v0.38.3...v0.38.4) (2026-10-22)


### Bug Fixes

* bump vulnerable dependencies (swag, vitest, nanoid, postcss, source-map-js) ([#62](https://github.com/apeirora/showroom-msp-apeiro-hosting/issues/62)) ([004c46e](https://github.com/apeirora/showroom-msp-apeiro-hosting/commit/004c46e883a0a477e4dbe3cca6e920b8e8f18d36))
* **controller:** restore account mode and adopt existing Shoots ([#65](https://github.com/apeirora/showroom-msp-apeiro-hosting/issues/65)) ([e1dc88c](https://github.com/apeirora/showroom-msp-apeiro-hosting/commit/e1dc88ccf76386dd4389b3e200836a85496d5901))
* **controller:** use account children for both hosting flavors ([#66](https://github.com/apeirora/showroom-msp-apeiro-hosting/issues/66)) ([0809964](https://github.com/apeirora/showroom-msp-apeiro-hosting/commit/080996472edf8206a8b488bda89fd3cdc875174c))

## [0.38.3](https://github.com/apeirora/showroom-msp-apeiro-hosting/compare/v0.38.2...v0.38.3) (2026-10-15)


### Bug Fixes

* use configured Shoot defaults and map related Secret namespaces ([#52](https://github.com/apeirora/showroom-msp-apeiro-hosting/issues/52)) ([9068462](https://github.com/apeirora/showroom-msp-apeiro-hosting/commit/90684626123e8dbab3a3170eb129ac557087215c))

## [0.38.2](https://github.com/apeirora/showroom-msp-apeiro-hosting/compare/v0.38.1...v0.38.2) (2026-10-14)


### Bug Fixes

* attach provider views to account and namespace dashboards ([#55](https://github.com/apeirora/showroom-msp-apeiro-hosting/issues/55)) ([7f54558](https://github.com/apeirora/showroom-msp-apeiro-hosting/commit/7f5455810cb5a3fa41b92b8999c7fc838dd5df19))

## [0.38.1](https://github.com/apeirora/showroom-msp-apeiro-hosting/compare/v0.38.0...v0.38.1) (2026-10-14)


### Bug Fixes

* use ApeiroRA logo for portal and marketplace icons ([#53](https://github.com/apeirora/showroom-msp-apeiro-hosting/issues/53)) ([04d7c62](https://github.com/apeirora/showroom-msp-apeiro-hosting/commit/04d7c623ab5784ed987cb1a35a647d927dc492c5))

## [0.38.0](https://github.com/apeirora/showroom-msp-apeiro-hosting/compare/v0.37.1...v0.38.0) (2026-10-14)


### ⚠ BREAKING CHANGES

* drop SAP naming, rename API group to apeiro.eu ([#42](https://github.com/apeirora/showroom-msp-apeiro-hosting/issues/42))

### Features

* drop SAP naming, rename API group to apeiro.eu ([#42](https://github.com/apeirora/showroom-msp-apeiro-hosting/issues/42)) ([95bc095](https://github.com/apeirora/showroom-msp-apeiro-hosting/commit/95bc0957a5e0e87733e1beccfa88fd741152357c))

## [0.37.1](https://github.com/apeirora/showroom-msp-apeiro-hosting/compare/v0.37.0...v0.37.1) (2026-10-14)


### Bug Fixes

* sanitize chart version labels ([#40](https://github.com/apeirora/showroom-msp-apeiro-hosting/issues/40)) ([0d26e0f](https://github.com/apeirora/showroom-msp-apeiro-hosting/commit/0d26e0f058f12679b4aa441149e466a21ba09a1e))

## [0.37.0](https://github.com/apeirora/showroom-msp-apeiro-hosting/compare/v0.36.0...v0.37.0) (2026-10-13)


### Features

* **rbac:** configure Secret reference validation ([#37](https://github.com/apeirora/showroom-msp-apeiro-hosting/issues/37)) ([c318f50](https://github.com/apeirora/showroom-msp-apeiro-hosting/commit/c318f50c58db796c8c5c297923f8572e4cddf782))


### Bug Fixes

* validate VM and Shoot credential references ([#38](https://github.com/apeirora/showroom-msp-apeiro-hosting/issues/38)) ([8dfe541](https://github.com/apeirora/showroom-msp-apeiro-hosting/commit/8dfe541cebfd86b62671269b4fe1e878c8149e26))

## [0.36.0](https://github.com/apeirora/showroom-msp-apeiro-hosting/compare/v0.35.4...v0.36.0) (2026-10-13)


### Features

* configure the IronCore Garden landscape through chart values ([#33](https://github.com/apeirora/showroom-msp-apeiro-hosting/issues/33)) ([2367fa9](https://github.com/apeirora/showroom-msp-apeiro-hosting/commit/2367fa9897417b0405f7fad50a27d23bad4daf4a))

## [0.35.4](https://github.com/apeirora/showroom-msp-apeiro-hosting/compare/v0.35.3...v0.35.4) (2026-09-24)


### Bug Fixes

* remove anonymous bind RBAC from pm-integration chart ([#29](https://github.com/apeirora/showroom-msp-apeiro-hosting/issues/29)) ([96fb6f4](https://github.com/apeirora/showroom-msp-apeiro-hosting/commit/96fb6f4bb02ab9501dac7bf815358234fbe40a92))

## [0.35.3](https://github.com/apeirora/showroom-msp-apeiro-hosting/compare/v0.35.2...v0.35.3) (2026-09-17)


### Bug Fixes

* **ui:** read namespace from Luigi path params ([#26](https://github.com/apeirora/showroom-msp-apeiro-hosting/issues/26)) ([eb24220](https://github.com/apeirora/showroom-msp-apeiro-hosting/commit/eb24220d25f70f60d52b104e9bcd3b3f3e2002e5))

## [0.35.2](https://github.com/apeirora/showroom-msp-apeiro-hosting/compare/v0.35.1...v0.35.2) (2026-09-15)


### Bug Fixes

* use SAP OpenStack credentials binding ([#24](https://github.com/apeirora/showroom-msp-apeiro-hosting/issues/24)) ([329145e](https://github.com/apeirora/showroom-msp-apeiro-hosting/commit/329145eb3815a9b8ae17b4ea0fd68226290e18c5))

## [0.35.1](https://github.com/apeirora/showroom-msp-apeiro-hosting/compare/v0.35.0...v0.35.1) (2026-09-15)


### Bug Fixes

* update SAP OpenStack Kubernetes version ([#22](https://github.com/apeirora/showroom-msp-apeiro-hosting/issues/22)) ([29a1bfb](https://github.com/apeirora/showroom-msp-apeiro-hosting/commit/29a1bfbf1a5cde0c6591fe0ddc7eff48cdbe2d29))

## [0.35.0](https://github.com/apeirora/showroom-msp-apeiro-hosting/compare/v0.34.1...v0.35.0) (2026-09-15)


### Features

* iron-core machine-type added `memory-intel1-4cpu` ([#20](https://github.com/apeirora/showroom-msp-apeiro-hosting/issues/20)) ([5b9d96f](https://github.com/apeirora/showroom-msp-apeiro-hosting/commit/5b9d96ff84cde858fe471094c88a58ffc23fa4b6))

## [0.34.1](https://github.com/apeirora/showroom-msp-apeiro-hosting/compare/v0.34.0...v0.34.1) (2026-08-06)


### Bug Fixes

* align ironcore shoot template with working shoots ([#18](https://github.com/apeirora/showroom-msp-apeiro-hosting/issues/18)) ([ba597dc](https://github.com/apeirora/showroom-msp-apeiro-hosting/commit/ba597dc8bdf1c48f98f27f6e07640a0231e24747))

## [0.34.0](https://github.com/apeirora/showroom-msp-apeiro-hosting/compare/v0.33.1...v0.34.0) (2026-08-04)


### Features

* disable delete for SimpleCluster-managed shoots and VMs ([#16](https://github.com/apeirora/showroom-msp-apeiro-hosting/issues/16)) ([9423987](https://github.com/apeirora/showroom-msp-apeiro-hosting/commit/94239875a990af61a90d9ad6c7007ef3449acb4a))

## [0.33.1](https://github.com/apeirora/showroom-msp-apeiro-hosting/compare/v0.33.0...v0.33.1) (2026-08-04)


### Bug Fixes

* default worker zone wdf4b in the embedded ironcore SCPC — wdf4a is decommissioned ([#14](https://github.com/apeirora/showroom-msp-apeiro-hosting/issues/14)) ([a430173](https://github.com/apeirora/showroom-msp-apeiro-hosting/commit/a430173101b5da4cb4a76fbd8a77ef2f11d81a22))

## [0.33.0](https://github.com/apeirora/showroom-msp-apeiro-hosting/compare/v0.32.0...v0.33.0) (2026-08-03)


### Features

* Label provider CR with owner-cluster for KCP up-sync ([#12](https://github.com/apeirora/showroom-msp-apeiro-hosting/issues/12)) ([32e2a3e](https://github.com/apeirora/showroom-msp-apeiro-hosting/commit/32e2a3e8e969c8ef76cbff3f15adc8ff8820c96f))

## [0.32.0](https://github.com/apeirora/showroom-msp-apeiro-hosting/compare/v0.31.1...v0.32.0) (2026-08-03)


### Features

* simplecluster local provider cr ([#10](https://github.com/apeirora/showroom-msp-apeiro-hosting/issues/10)) ([855750c](https://github.com/apeirora/showroom-msp-apeiro-hosting/commit/855750c9ed778eb03873536f47ba986962ff9fa7))

## [0.31.1](https://github.com/apeirora/showroom-msp-apeiro-hosting/compare/v0.31.0...v0.31.1) (2026-07-29)


### Bug Fixes

* **ui:** simplify sidebar resource labels ([#8](https://github.com/apeirora/showroom-msp-apeiro-hosting/issues/8)) ([d79a4b1](https://github.com/apeirora/showroom-msp-apeiro-hosting/commit/d79a4b15a45447149cb7a5a58a180a46940115b2))

## [0.31.0](https://github.com/apeirora/showroom-msp-apeiro-hosting/compare/v0.30.0...v0.31.0) (2026-07-29)


### Features

* add imagePullSecrets to ironcore-operator chart ([#4](https://github.com/apeirora/showroom-msp-apeiro-hosting/issues/4)) ([ae701f7](https://github.com/apeirora/showroom-msp-apeiro-hosting/commit/ae701f785f8b8274eeccfaac214cfb42fc44a1c9))


### Bug Fixes

* **ci:** release-please tags via apeirora-showroom-bot app token ([#5](https://github.com/apeirora/showroom-msp-apeiro-hosting/issues/5)) ([c7dfd11](https://github.com/apeirora/showroom-msp-apeiro-hosting/commit/c7dfd11a902656e84ff1224e408648ef2c6bf75f))
* rename API group to sap.eu (incl. IronCoreVM) ([#6](https://github.com/apeirora/showroom-msp-apeiro-hosting/issues/6)) ([abba79d](https://github.com/apeirora/showroom-msp-apeiro-hosting/commit/abba79d549716f94caf193cb9cbebcd0c18c1b25))

## [0.30.0](https://github.com/apeirora/showroom-msp-apeiro-hosting/compare/v0.29.0...v0.30.0) (2026-07-27)


### Miscellaneous Chores

* bootstrap standalone releases at 0.30.0 ([#2](https://github.com/apeirora/showroom-msp-apeiro-hosting/issues/2)) ([a069f1b](https://github.com/apeirora/showroom-msp-apeiro-hosting/commit/a069f1bc91768af99a8cbb7d764d7673a62bc013))
