/*
@author Thanh Nguyen <btnguyen2k@gmail.com>
@GoVueAdmin
*/
// cannot use require in vite https://github.com/vitejs/vite/issues/1241#issuecomment-762367047
// const APP_CONFIG = require('@/config.json')
import APP_CONFIG from '@/config.json'
if (!APP_CONFIG.api_client.bo_api_base_url) {
    APP_CONFIG.api_client.bo_api_base_url = process.env.VUE_APP_BO_API_BASE_URL
}
const APP_ID = APP_CONFIG.api_client.app_id
export default {
    APP_ID,
    APP_CONFIG
}
