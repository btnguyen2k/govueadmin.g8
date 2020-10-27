# govueadmin.g8

Giter8 template to develop `Admin Control Panel` in Go with VueJS-based frontend.

Latest release: [template-v0.2.0](RELEASE-NOTES.md).

Demo: https://demo-govueadmin.herokuapp.com/.

## Features

- Create new project from template with [go-giter8](https://github.com/btnguyen2k/go-giter8).
- Frontend (GUI) built on [CoreUI for Vue.js v3](https://coreui.io/vue/).
- Backend (API) built on [Echo framework v4](https://echo.labstack.com).
- Sample features:
  - Login page & Logout
  - Blog post management (list, create, update, delete)
  - Dashboard/Feed
    - Vote up/down on public blog posts
  - BO & DAO:
    - [SQLite3](https://github.com/mattn/go-sqlite3) implementation.
    - PostgreSQL implementation.
- Sample `Dockerfile` to package application as Docker image.


## Getting Started

### Install `go-giter8`

This a Giter8 template, so it is meant to be used in conjunction with a giter8 tool.
Since this is a template for Go application, it make sense to use [go-giter8](https://github.com/btnguyen2k/go-giter8).

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

The frontend is a VueJS-based project located under `fe-gui` directory and can be imported into ay VueJS-supported IDE.
Feel free to use it to develop your application's frontend.

> The frontend is built on [CoreUI for Vue.js](https://coreui.io/vue/).

**Backend**

The backend is a Go project built on [Echo framework v4](https://echo.labstack.com) located under `be-api` directory.
It purely provides APIs for frontend to call (e.g. it has no GUI) and can be opened by any Go-supported IDE.

> The backend is based on `goapi.g8` template. See [goapi.g8](https://github.com/btnguyen2k/goapi.g8) for API implementation guideline.

**Note: The frontend is only GUI, it needs to call backend API to retrieve and modify data.**


## LICENSE & COPYRIGHT

See [LICENSE.md](LICENSE.md) for details.


## Giter8 template

For information on giter8 templates, please see http://www.foundweekends.org/giter8/
