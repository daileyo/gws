# Changelog

## [2.6.0](https://github.com/daileyo/gws/compare/v2.5.0...v2.6.0) (2026-03-02)


### Features

* add deprecation layer for old flag forms ([47af8ea](https://github.com/daileyo/gws/commit/47af8ea36f2c3169e2923883ffb84feed9f5beb3))
* create core subcommands (list, init, add, refresh, print-workspace) ([6ef1318](https://github.com/daileyo/gws/commit/6ef1318d8cf7b9e7b86273c444a1f2fe019d5763))
* scope filter flags to list subcommand ([cccb7a5](https://github.com/daileyo/gws/commit/cccb7a53505d73fb708b6b3a3854c4044da54e5e))
* update shell templates to route subcommand names to binary ([727c62b](https://github.com/daileyo/gws/commit/727c62b1dfb93a8c82aabc06da5c70bb95d2d47a))


### Bug Fixes

* resolve lint issues in main.go and deprecated_test.go ([cb06381](https://github.com/daileyo/gws/commit/cb06381755b83e878f5a5024233cc4dab79f958c))

## [2.5.0](https://github.com/daileyo/gws/compare/v2.4.1...v2.5.0) (2026-03-02)


### Features

* add repo user update/delete commands and profile sync ([a067e0e](https://github.com/daileyo/gws/commit/a067e0e6782a35a59493d1ea12114eab15d2ea6c))

## [2.4.1](https://github.com/daileyo/gws/compare/v2.4.0...v2.4.1) (2026-02-27)


### Bug Fixes

* correct license badge in docds ([697ccbe](https://github.com/daileyo/gws/commit/697ccbecd41e0f8fdd0488526af625c388a01dab))

## [2.4.0](https://github.com/daileyo/gws/compare/v2.3.0...v2.4.0) (2026-02-27)


### Features

* add logo and tagline to README ([5ec212d](https://github.com/daileyo/gws/commit/5ec212d93e261c375f7ddba377e707e8636dbe8c))
* add logo to mkdocs nav bar, favicon, and homepage hero ([84f8b84](https://github.com/daileyo/gws/commit/84f8b84a4c78e9637287f06476d20ee3dce4c02f))
* add optimized logo variants for branding ([09a19e5](https://github.com/daileyo/gws/commit/09a19e54acf6f4acc8e7eba1ef42810aa5bdb93f))
* configure mkdocs theme with branded colors and fonts ([c60b0aa](https://github.com/daileyo/gws/commit/c60b0aa4b9eea74e91dc9384ff48eae82a3435b5))

## [2.3.0](https://github.com/daileyo/gws/compare/v2.2.0...v2.3.0) (2026-02-27)


### Features

* add commit-msg hook for conventional commit alignment ([f0ef8bd](https://github.com/daileyo/gws/commit/f0ef8bda2f9311821e62f7e1a8019ce7d87eee77))
* update setup-hooks messaging and document in README ([d7dd3d7](https://github.com/daileyo/gws/commit/d7dd3d7ca598d93ffe7016bb184b6d915ed539c8))

## [2.2.0](https://github.com/daileyo/gws/compare/v2.1.1...v2.2.0) (2026-02-27)


### Features

* **config:** extend model with user profiles and repository user fields ([5fd4bf7](https://github.com/daileyo/gws/commit/5fd4bf7f9b8eaca8c68375f4d51d982c82f6e5dd))
* **git:** implement user detection from repositories ([11f7310](https://github.com/daileyo/gws/commit/11f73103a7295f6d91c28c48a45a48be93f72768))
* **list:** add --user flag to display git user information ([00476f9](https://github.com/daileyo/gws/commit/00476f9adf6cf402a4eaf55e3b0a405bb5358cb9))
* **user:** add assign and sync commands for repository user management ([3dda5f1](https://github.com/daileyo/gws/commit/3dda5f14c774c463b9bf19d2a5e5b6e1125d7547))
* **user:** add user profile management commands ([418e3ed](https://github.com/daileyo/gws/commit/418e3ed2edeb64de0e1a2235ab450c0720137876))
* **user:** implement profile auto-detection from gitconfig ([c560843](https://github.com/daileyo/gws/commit/c5608432b5b1118d07cd0c436f62364bcd8c7fe2))


### Bug Fixes

* resolve linting issues ([c846fe0](https://github.com/daileyo/gws/commit/c846fe0b26585215f4a73ea01b5dd33e111e6027))

## [2.1.1](https://github.com/daileyo/gws/compare/v2.1.0...v2.1.1) (2026-02-26)


### Bug Fixes

* correct project structure paths and remove premature brew install references ([ec1208e](https://github.com/daileyo/gws/commit/ec1208efb1ef1f52193e5735fa923cc7fd9187d1))

## [2.1.0](https://github.com/daileyo/gws/compare/v2.0.0...v2.1.0) (2026-02-26)


### Features

* refactor init/add, shell completions, and shell-init integration ([#6](https://github.com/daileyo/gws/issues/6)) ([e9d1564](https://github.com/daileyo/gws/commit/e9d156448ac012460050bca78dd2f06ff4d63bb5))

## [2.0.0](https://github.com/daileyo/gws/compare/v1.1.0...v2.0.0) (2026-02-25)


### ⚠ BREAKING CHANGES

* binary is now named git-workspace

### Features

* add repository navigation and rename binary to git-workspace ([#4](https://github.com/daileyo/gws/issues/4)) ([fb4a6ff](https://github.com/daileyo/gws/commit/fb4a6ff133241fb1f640278fe83a6bfe351f6713))

## [1.1.0](https://github.com/daileyo/gws/compare/v1.0.0...v1.1.0) (2026-02-13)


### Features

* rework CLI to flag-based command pattern ([#2](https://github.com/daileyo/gws/issues/2)) ([7256c97](https://github.com/daileyo/gws/commit/7256c97df4a19d18317dad84f2b40ff59b1d2a32))

## 1.0.0 (2026-01-06)


### Features

* complete questionare ([9d83e51](https://github.com/daileyo/gws/commit/9d83e51154943a387e89279b1ccde362d82471ca))
* generate initial spec ([49666c7](https://github.com/daileyo/gws/commit/49666c7ff4874050bcdd9d7807f832f38338fbc7))
* implement step 6 ([8787dae](https://github.com/daileyo/gws/commit/8787dae2d3bb6153297907cd0316b82e7f38d520))
* implement task 4 ([1b95999](https://github.com/daileyo/gws/commit/1b95999c99984a2496517342d35dfe389ee58abb))
* project setup and basic CLI framework ([8eb3ed1](https://github.com/daileyo/gws/commit/8eb3ed15636660a9d5dd6056dbc5e2653a3983d5))
* repository discovery and workspace initialization ([15cf8b4](https://github.com/daileyo/gws/commit/15cf8b474d9ca1c3beb0e5f5298cd6f8a4723343))


### Bug Fixes

* addres linting issues ([05c57bc](https://github.com/daileyo/gws/commit/05c57bcbe9e30ac80295e67ff300fdacb550e652))
* update target go version ([91636ae](https://github.com/daileyo/gws/commit/91636ae24f3436e2f6ec643947784cb8e6f9efe9))
