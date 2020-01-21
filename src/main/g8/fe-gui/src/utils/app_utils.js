const lskeyLoginSession = "usession"
const lskeyLoginSessionLastCheck = "usession_lastcheck"

function localStorageGet(key) {
    return localStorage.getItem(key)
}

function localStorageGetAsInt(key) {
    let str = localStorage.getItem(key)
    if (str == null) {
        return 0
    }
    let v = Number.parseInt(str, 10)
    return isNaN(v) ? 0 : v
}

function localStorageSet(key, value) {
    if (value == null) {
        localStorage.removeItem(key)
    } else {
        localStorage.setItem(key, value)
    }
}

function getUnixTimestamp() {
    return Math.round((new Date()).getTime() / 1000)
}

function getLoginSession() {
    let str = localStorageGet(lskeyLoginSession)
    return str != null ? JSON.parse(str) : null
}

function saveLoginSession(session) {
    localStorageSet(lskeyLoginSession, JSON.stringify(session))
}

export default {
    lskeyLoginSession,
    lskeyLoginSessionLastCheck,

    localStorageGet,
    localStorageSet,
    localStorageGetAsInt,
    getUnixTimestamp,
    loadLoginSession: getLoginSession,
    saveLoginSession
}
