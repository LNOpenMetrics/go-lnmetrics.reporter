# v0.0.6-rc2

## Fixed
- fix race condition caused by the gb (maybe) ([commit](https://github.com/LNOpenMetrics/go-lnmetrics.reporter/commit/f5999d9c8a56979dee5cc0b71f2c4cedc8c0cbe4)). @vincenzopalazzo 31-01-2023
- allow the rpc method to query only old score ([commit](https://github.com/LNOpenMetrics/go-lnmetrics.reporter/commit/0c98b7e99b9adec8caa8296de6448cf6724e73f6)). @vincenzopalazzo 31-01-2023
- code metrics code and rename metric one to rawLocalRep ([commit](https://github.com/LNOpenMetrics/go-lnmetrics.reporter/commit/4d6959e703e8cc96f0b92a3853e5ee8f9c4f731e)). @vincenzopalazzo 30-01-2023

## Added
- add ToMap function in the Metric interface ([commit](https://github.com/LNOpenMetrics/go-lnmetrics.reporter/commit/3d9d240e2cd3e9b612f0547ce4584e823367d4ce)). @vincenzopalazzo 30-01-2023
- refactoring metric one to raw local score ([commit](https://github.com/LNOpenMetrics/go-lnmetrics.reporter/commit/faff772f05bce6f6e2f3e05970ada271a8757b04)). @vincenzopalazzo 30-01-2023


# v0.0.6-rc1

## Fixed
- use the encoder inside the matric one ([commit](https://github.com/LNOpenMetrics/go-lnmetrics.reporter/commit/3e1674bf828c36ca4374ae2130075f8e198cf19d)). @vincenzopalazzo 13-12-2022
- fix the null node id that is inside the getNode ([commit](https://github.com/LNOpenMetrics/go-lnmetrics.reporter/commit/ae3481395e3d70437b4df3c2186059f89a9ddc94)). @vincenzopalazzo 13-12-2022
- fix rpc command to return the correct value ([commit](https://github.com/LNOpenMetrics/go-lnmetrics.reporter/commit/15fd6eff8eec462c5fbad81c5d1ba83fa5e7b33f)). @vincenzopalazzo 07-12-2022
- fix partial reading stream by reimplementing the scanner ([commit](https://github.com/LNOpenMetrics/go-lnmetrics.reporter/commit/303d666bedba7cf965b2d883f983565c84fba267)). @vincenzopalazzo 04-12-2022
- dependeces version ([commit](https://github.com/LNOpenMetrics/go-lnmetrics.reporter/commit/e59714dd79ecc1b60090d7898f5a1da14bd5024d)). @vincenzopalazzo 03-12-2022
- fix partial readin from socker un cln4go ([commit](https://github.com/LNOpenMetrics/go-lnmetrics.reporter/commit/22ae520094f04376ffc67748ffc104af81968784)). @vincenzopalazzo 03-12-2022
- fix underflow under uint64 ([commit](https://github.com/LNOpenMetrics/go-lnmetrics.reporter/commit/d648501120e4a348f58451e6131e078db8a4db43)). @vincenzopalazzo 20-11-2022
- unmarshall empty binary array ([commit](https://github.com/LNOpenMetrics/go-lnmetrics.reporter/commit/780ee6c4027f6895f8b772d966e90ee5a7ed148f)). @vincenzopalazzo 16-11-2022
- fix the metric one call and return the correct value ([commit](https://github.com/LNOpenMetrics/go-lnmetrics.reporter/commit/61ef1bf529efffb74647c7ecc02d6a9859d08209)). @vincenzopalazzo 12-11-2022
- introduce the model and the function call ([commit](https://github.com/LNOpenMetrics/go-lnmetrics.reporter/commit/82a5f6e5e74525acab6c57a090a97532091cf488)). @vincenzopalazzo 05-10-2022
- change the metric one plugin signature with cln4go client ([commit](https://github.com/LNOpenMetrics/go-lnmetrics.reporter/commit/b8800d364ea64afb25881fcf517ca131ad4364dc)). @vincenzopalazzo 05-10-2022
- clean the peer snapshot map each time ([commit](https://github.com/LNOpenMetrics/go-lnmetrics.reporter/commit/88b61941a4bfe189f1a30883749011cda66ebde9)). @vincenzopalazzo 05-10-2022
- clean the peer snapshot map each time ([commit](https://github.com/LNOpenMetrics/go-lnmetrics.reporter/commit/9f982aeeeb017ea2bd779c2a27209014cc13c4d0)). @vincenzopalazzo 05-10-2022

## Added
- add faster json encoder ([commit](https://github.com/LNOpenMetrics/go-lnmetrics.reporter/commit/cea9291ec80468f9cb54e90db5d04e19f631f060)). @vincenzopalazzo 14-11-2022
- log the paninc if any when the main will end ([commit](https://github.com/LNOpenMetrics/go-lnmetrics.reporter/commit/eac5c4832ab86c12d6e9d386df7d63c0c54d5fa6)). @vincenzopalazzo 14-11-2022
- migrate the last rpc command to cln4go ([commit](https://github.com/LNOpenMetrics/go-lnmetrics.reporter/commit/92834b7d0856d100b24369b00afdb8ea7c615342)). @vincenzopalazzo 29-10-2022


# v0.0.5-rc3

## Fixed
- forward payment direction ([commit](https://github.com/LNOpenMetrics/go-lnmetrics.reporter/commit/11ffa17f964ee7accb6fe7736973a9886c7b2e4c)). @vincenzopalazzo 23-08-2022


# v0.0.5-rc2

## Fixed
- be nicer with a node by ping it only one time ([commit](https://github.com/LNOpenMetrics/go-lnmetrics.reporter/commit/9231650b4575b8c39a3fef96931bf1def92680df)). @vincenzopalazzo 19-08-2022


# v0.0.5-rc1

## Added
- add rpc method `lnmetrics-force-update` ([commit](https://github.com/LNOpenMetrics/go-lnmetrics.reporter/commit/91499c4a7d8a4f12d5e6228e068724df42b41071)). @vincenzopalazzo 02-04-2022
- adding new rpc command `lnmetrics-cache` ([commit](https://github.com/LNOpenMetrics/go-lnmetrics.reporter/commit/27aac0f41f73b27fde4e2dceb2e04091302fa214)). @vincenzopalazzo 01-04-2022

## Fixed
- introduce the cache system to speed up the plugin with big node ([commit](https://github.com/LNOpenMetrics/go-lnmetrics.reporter/commit/d99abce1f3741207399b06922f9ab0f9fde5a767)). @vincenzopalazzo 01-04-2022
