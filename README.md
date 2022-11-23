# govueadmin.g8

[![Actions Status](https://github.com/btnguyen2k/govueadmin.g8/workflows/govueadmin/badge.svg)](https://github.com/btnguyen2k/govueadmin.g8/actions)
[![codecov](https://codecov.io/gh/btnguyen2k/govueadmin.g8/branch/master/graph/badge.svg?token=HVAP5A0R2Z)](https://codecov.io/gh/btnguyen2k/govueadmin.g8)
[![Release](https://img.shields.io/github/release/btnguyen2k/govueadmin.g8.svg?style=flat-square)](RELEASE-NOTES.md)

Giter8 template to build `Admin Control Panel` for Go with VueJS-based frontend.

Demo: https://demo-govueadmin.gpvcloud.com/.

## Features

- [Giter8](https://github.com/btnguyen2k/go-giter8) template.
- Single-page applicaton (SPA):
  - Frontend (GUI) built on [CoreUI for Vue.js](https://coreui.io/vue/).
  - Backend (API) built on [Echo framework](https://echo.labstack.com).
- I18n support:
  - FE i18n support with [Vue-i18n](https://kazupon.github.io/vue-i18n/).
  - BE i18n support with [goyai](https://github.com/btnguyen2k/goyai).
- Sample features:
  - Login page & Logout
  - Blog post management (list, create, update, delete)
  - Dashboard/Feed
    - Vote up/down on public blog posts
  - BO & DAO implementation in AWS DynamoDB, Azure Cosmos DB, SQLite3, MySQL, PostgreSQL and MongoDB.
- Sample `Dockerfile` to package application as Docker image.
- Sample [GitHub Actions](https://docs.github.com/actions) workflow.


## Getting Started

### Install `go-giter8`

This a Giter8 template, so it is meant to be used in conjunction with a giter8 tool.
Since this is a template for Go application, it makes sense to use [go-giter8](https://github.com/btnguyen2k/go-giter8).

See [go-giter8](https://github.com/btnguyen2k/go-giter8) website for installation guide.

### Create new project from template

```
g8 new btnguyen2k/govueadmin.g8
```

and follow the instructions.

> Note: This template requires `go-giter8 v0.4.2` or higher.

Upon successful project creation, 2 sub-projects are created:

- `fe-gui`: frontend, which is a VueJS project.
- `be-api`: backend, which is a Go project.

### Write application code

**Frontend**

The frontend is a VueJS-based project located under `fe-gui` directory and can be imported into any VueJS-supported IDE.
Feel free to use it as a reference or the starting point to build your awesome application's frontend.

> The frontend is built on [CoreUI for Vue.js](https://coreui.io/vue/).

**Backend**

The backend is a Go project built on [Echo framework](https://echo.labstack.com) located under `be-api` directory.
It purely provides APIs for frontend to call (e.g. it has no GUI) and can be opened by any Go-supported IDE.

> The backend is based on `goapi.g8` template. See [goapi.g8](https://github.com/btnguyen2k/goapi.g8) for API implementation guideline.

**Note: The frontend is GUI only, it needs to call backend APIs to retrieve and modify data.**


## LICENSE & COPYRIGHT

See [LICENSE.md](LICENSE.md) for details.


## Giter8 template

For information on giter8 templates, please see http://www.foundweekends.org/giter8/
