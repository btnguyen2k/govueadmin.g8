/*
@author Thanh Nguyen <btnguyen2k@gmail.com>
@GoVueAdmin
*/
const APP_CONFIG = require('@/config.json')
if (!APP_CONFIG.api_client.bo_api_base_url) {
    APP_CONFIG.api_client.bo_api_base_url = process.env.VUE_APP_BO_API_BASE_URL
}
const APP_ID = APP_CONFIG.api_client.app_id
export default {
    APP_ID,
    APP_CONFIG
}
