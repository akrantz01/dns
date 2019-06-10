export default class Authentication {
    // Check if a user is authenticated
    static isAuthenticated = () => localStorage.getItem("token") !== null;

    // Getters and setters for the authentication token
    static setToken = (token) => localStorage.setItem("token", token);
    static getToken = () => localStorage.getItem("token") || "";

    // Getters and setters for user data
    static setUser = (user) => localStorage.setItem("user", JSON.stringify(user));
    static getUser = () => JSON.parse(localStorage.getItem("user")) || {name: ""};
}
