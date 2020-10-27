# govueadmin.g8 Release Notes

## 2020-10-28: template-v0.2.0

- Use [Gravatar](https://gravatar.com/) for user profile picture.
- AB#21: Migrate to CoreUI-vue v3.1.1.
- AB#22: Integrate with [Exter](https://github.com/btnguyen2k/exter).
- AB#25: Use JWT as login-token.
- EP#25: Rebuild sample app.
- Others:
  - Vuejs: replace `data` with `computed` (https://techformist.com/data-vs-computed-vs-watcher-in-vue-components/) if possible.
  - Other fixes and enhancements.


## 2020-01-30: template-v0.1.2

- PostgreSQL implementation of BO & DAO.


## 2020-01-21: template-v0.1.0

- Frontend using [CoreUI for Vue.js](https://coreui.io/vue/) `v3.0.0-beta.3`.
- Backend using [Echo](https://echo.labstack.com) `v4.1.x`.
- Sample features:
  - Login page & Logout
  - Dashboard
  - User group management (list, create, update, delete)
  - User management (list, create, update, delete, change password)
  - BO & DAO implementation using [SQLite3](https://github.com/mattn/go-sqlite3)
- Sample `.gitlab-ci.yaml` & `Dockerfile`.
