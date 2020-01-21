# $name$ Frontend

**$desc$** by **$author$ - $organization$**, based on [govueadmin.g8](https://github.com/btnguyen2k/govueadmin.g8).

Copyright (C) by **$organization$**.

Latest release version: `$version$`. See [RELEASE-NOTES.md](../RELEASE-NOTES.md).

##Getting Started

**Project setup**
```
npm install
```

**Compiles and hot-reloads for development**
```
npm run serve
```

**Compiles and minifies for production**
```
npm run build
```

**Lints and fixes files**
```
npm run lint
```

> This frontend project is built on [CoreUI for Vue.js](https://coreui.io/vue/docs/introduction/).

##Application Configurations

`src/config.json`: application's main configuration file. Important config keys:
- `api_client.bo_api_base_url`: point to backend's base URL
- `api_client.app_id`: application id in order to authenticate with backend. _Must match between frontend and backend._
- `api_client.header_app_id` and `api_client.header_access_token`: name of HTTP headers passed along every API call for authentication. _Must match between frontend and backend._

## References

- [VueJS Configuration Reference](https://cli.vuejs.org/config/)
- [CoreUI for Vue.JS](https://coreui.io/vue/docs/introduction/)

## LICENSE & COPYRIGHT

See [LICENSE.md](../LICENSE.md).
