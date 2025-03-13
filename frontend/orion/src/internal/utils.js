export function isLogged() {
    const token = localStorage.getItem("JWT")
    return !!token
}