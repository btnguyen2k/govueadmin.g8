/*
Client to make call to API server using Axios.

@author Thanh Nguyen <btnguyen2k@gmail.com>
@since template-v0.4.r1
*/
import Axios from "axios"
import APP_CONFIG from "./app_config"
import router from "@/router/index"
import utils from "@/utils/app_utils";

const apiClient = Axios.create({
    baseURL: APP_CONFIG.api_client.bo_api_base_url,
    timeout: 10000,
});

const headerAppId = APP_CONFIG.api_client.header_app_id
const headerAccessToken = APP_CONFIG.api_client.header_access_token
let appId = APP_CONFIG.api_client.app_id

let apiLogin = "/api/login"
let apiCheckLoginToken = "/api/checkLoginToken"
let apiSystemInfo = "/api/systemInfo"
let apiGroupList = "/api/groups"
let apiGroup = "/api/group"
let apiUsers = "/api/users"

function _apiOnSuccess(resp, apiUri, callbackSuccessful) {
    // if (apiUri != apiLogin && apiUri != apiCheckLoginToken && resp.hasOwnProperty("data") && resp.data.status == 403) {
    //     console.error("Error 403 from API [" + apiUri + "], redirecting to login page...")
    //     router.push({name: "Login", query: {returnUrl: router.currentRoute.fullPath}})
    //     return
    // }
    if (resp.hasOwnProperty("data") && resp.data.hasOwnProperty("extras") && resp.data.extras.hasOwnProperty("_access_token_")) {
        console.log("Update new access token from API [" + apiUri + "]")
        let tokens = resp.data.extras._access_token_.split(":")
        utils.saveLoginSession({
            uid: tokens[0],
            token: tokens[1],
            expiry: parseInt(tokens[2]),
        })
    }
    if (callbackSuccessful != null) {
        callbackSuccessful(resp.data)
    }
}

function _apiOnError(err, apiUri, callbackError) {
    console.error("Error calling api [" + apiUri + "]: " + err)
    if (callbackError != null) {
        callbackError(err)
    }
}

function apiDoGet(apiUri, callbackSuccessful, callbackError) {
    let session = utils.loadLoginSession()
    const headers = {}
    headers[headerAppId] = appId
    headers[headerAccessToken] = session != null ? session.uid + ":" + session.token : ""
    return apiClient.get(apiUri, {
        headers: headers
    }).then(res => _apiOnSuccess(res, apiUri, callbackSuccessful)).catch(err => _apiOnError(err, apiUri, callbackError))
}

function apiDoPost(apiUri, data, callbackSuccessful, callbackError) {
    let session = utils.loadLoginSession()
    const headers = {}
    headers[headerAppId] = appId
    headers[headerAccessToken] = session != null ? session.uid + ":" + session.token : ""
    apiClient.post(apiUri, data, {
        headers: headers
    }).then(res => _apiOnSuccess(res, apiUri, callbackSuccessful)).catch(err => _apiOnError(err, apiUri, callbackError))
}

function apiDoPut(apiUri, data, callbackSuccessful, callbackError) {
    let session = utils.loadLoginSession()
    const headers = {}
    headers[headerAppId] = appId
    headers[headerAccessToken] = session != null ? session.uid + ":" + session.token : ""
    apiClient.put(apiUri, data, {
        headers: headers
    }).then(res => _apiOnSuccess(res, apiUri, callbackSuccessful)).catch(err => _apiOnError(err, apiUri, callbackError))
}

function apiDoDelete(apiUri, callbackSuccessful, callbackError) {
    let session = utils.loadLoginSession()
    const headers = {}
    headers[headerAppId] = appId
    headers[headerAccessToken] = session != null ? session.uid + ":" + session.token : ""
    apiClient.delete(apiUri, {
        headers: headers
    }).then(res => _apiOnSuccess(res, apiUri, callbackSuccessful)).catch(err => _apiOnError(err, apiUri, callbackError))
}

export default {
    apiLogin,
    apiCheckLoginToken,
    apiSystemInfo,
    apiGroupList,
    apiGroup,
    apiUsers,

    apiDoGet,
    apiDoPost,
    apiDoPut,
    apiDoDelete,
}
