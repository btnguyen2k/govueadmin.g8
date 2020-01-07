function getUserToken() {
    return localStorage.getItem("utoken")
}

function setUserToken(token) {
    localStorage.setItem("utoken", token)
}

export default {
    getUserToken,
    setUserToken
}
