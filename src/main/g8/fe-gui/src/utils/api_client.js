/*
Client to make call to API server using Axios.

@author Thanh Nguyen <btnguyen2k@gmail.com>
@since template-v0.4.r1
*/
import Axios from "axios"
import APP_CONFIG from "./app_config";

const apiClient = Axios.create({
    baseURL: APP_CONFIG.api_client.bo_api_base_url,
    timeout: 10000,
});

const headerAppId = APP_CONFIG.api_client.header_app_id
const headerAccessToken = APP_CONFIG.api_client.header_access_token
let appId = APP_CONFIG.api_client.app_id
let accessToken = ""
let apiLogin = "/api/login"
let apiCheckLoginToken = "/api/checkLoginToken"
let apiSystemInfo = "/api/systemInfo"

function apiDoGet(apiUri, callbackSuccessful, callbackError) {
    const headers = {}
    headers[headerAppId] = appId
    headers[headerAccessToken] = accessToken
    return apiClient.get(apiUri, {
        headers: headers
    }).then(res => {
        if (callbackSuccessful != null) {
            callbackSuccessful(res.data)
        }
    }).catch(err => {
        console.log(err)
        if (callbackError != null) {
            callbackError(err)
        }
    })
}

function apiDoPost(apiUri, data, callbackSuccessful, callbackError) {
    const headers = {}
    headers[headerAppId] = appId
    headers[headerAccessToken] = accessToken
    apiClient.post(apiUri, data, {
        headers: headers
    }).then(res => {
        if (callbackSuccessful != null) {
            callbackSuccessful(res.data)
        }
    }).catch(err => {
        console.log(err)
        if (callbackError != null) {
            callbackError(err)
        }
    })
}

export default {
    apiLogin,
    apiCheckLoginToken,
    apiSystemInfo,
    apiDoGet,
    apiDoPost
}
