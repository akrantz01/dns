import axios from 'axios';
import {API_URL} from './config';

export class ApiAuthorization {
    static Login = (username, password) => new Promise((resolve, reject) => {
        axios({
            method: "POST",
            url: `${API_URL}/users/login`,
            headers: {"Content-Type": "application/json"},
            data: {
                username: username,
                password: password
            }
        }).then(res => resolve(res.data)).catch(err => reject(err));
    });

    static Logout = (token) => new Promise((resolve, reject) => {
        axios({
            method: "GET",
            url: `${API_URL}/users/logout`,
            headers: {"Authorization": token}
        }).then(res => resolve(res.data)).catch(err => reject(err));
    });
}

export class ApiUsers {
    static Create = (name, username, password, role, token) => new Promise((resolve, reject) => {
        axios({
            method: "POST",
            url: `${API_URL}/users`,
            headers: {
                "Authorization": token,
                "Content-Type": "application/json"
            },
            data: {
                name: name,
                username: username,
                password: password,
                role: role
            }
        }).then(res => resolve(res.data)).catch(err => reject(err));
    });

    static Read = (token, username="") => new Promise((resolve, reject) => {
        axios({
            method: "GET",
            url: `${API_URL}/users${ (username !== "") ? "?user="+username : "" }`,
            headers: {"Authorization": token}
        }).then(res => resolve(res.data)).catch(err => reject(err));
    });

    static Update = (name, password, role, token, username="") => new Promise((resolve, reject) => {
        axios({
            method: "PUT",
            url: `${API_URL}/users${ (username !== "") ? "?user="+username : "" }`,
            headers: {
                "Authorization": token,
                "Content-Type": "application/json"
            },
            data: {
                name: name,
                password: password,
                role: role
            }
        }).then(res => resolve(res)).catch(err => reject(err));
    });

    static Delete = (token, username="") => new Promise((resolve, reject) => {
        axios({
            method: "DELETE",
            url: `${API_URL}/users${ (username !== "") ? "?user="+username : "" }`,
            headers: {"Authorization": token}
        }).then(res => resolve(res.data)).catch(err => reject(err));
    });
}

export class ApiRecords {
    static List = (type="", token) => new Promise((resolve, reject) => {
        axios({
            method: "GET",
            url: `${API_URL}/records${ (type !== "") ? "?type="+type : "" }`,
            headers: {"Authorization": token}
        }).then(res => resolve(res.data)).catch(err => reject(err));
    });

    static Create = (type, name, data, token) => new Promise((resolve, reject) => {
        axios({
            method: "POST",
            url: `${API_URL}/records`,
            headers: {
                "Content-Type": "application/json",
                "Authorization": token
            },
            data: {
                name: name,
                type: type,
                ...data
            }
        }).then(res => resolve(res.data)).catch(err => reject(err));
    });

    static Read = (name, type, token) => new Promise((resolve, reject) => {
        axios({
            method: "GET",
            url: `${API_URL}/records/${name}?type=${type}`,
            headers: {"Authorization": token}
        }).then(res => resolve(res.data)).catch(err => reject(err));
    });

    static Update = (name, type, data, token) => new Promise((resolve, reject) => {
        axios({
            method: "PUT",
            url: `${API_URL}/records/${name}`,
            headers: {
                "Authorization": token,
                "Content-Type": "application/json"
            },
            data: {
                type: type,
                ...data
            }
        }).then(res => resolve(res.data)).catch(err => reject(err));
    });

    static Delete = (name, type, token) => new Promise((resolve, reject) => {
        axios({
            method: "DELETE",
            url: `${API_URL}/records/${name}?type=${type}`,
            headers: {"Authorization": token}
        }).then(res => resolve(res.data)).catch(err => reject(err));
    });
}

export class ApiRoles {
    static List = (token) => new Promise((resolve, reject) => {
        axios({
            method: "GET",
            url: `${API_URL}/roles`,
            headers: {"Authorization": token}
        }).then(res => resolve(res.data)).catch(err => reject(err));
    });

    static Create = (name, filter, effect, token) => new Promise((resolve, reject) => {
        axios({
            method: "POST",
            url: `${API_URL}/roles`,
            headers: {
                "Authorization": token,
                "Content-Type": "application/json"
            },
            data: {
                name: name,
                filter: filter,
                effect: effect
            }
        }).then(res => resolve(res.data)).catch(err => reject(err));
    });

    static Read = (name, token) => new Promise((resolve, reject) => {
        axios({
            method: "GET",
            url: `${API_URL}/roles/${name}`,
            headers: {"Authorization": token}
        }).then(res => resolve(res.data)).catch(err => reject(err));
    });

    static Update = (name, filter, effect, token) => new Promise((resolve, reject) => {
        axios({
            method: "PUT",
            url: `${API_URL}/roles/${name}`,
            headers: {
                "Authorization": token,
                "Content-Type": "application/json"
            },
            data: {
                filter: filter,
                effect: effect
            }
        }).then(res => resolve(res.data)).catch(err => reject(err));
    });

    static Delete = (name, effect="", token) => new Promise((resolve, reject) => {
        axios({
            method: "DELETE",
            url: `${API_URL}/roles/${name}${ (effect !== "") ? "?effect="+effect : "" }`,
            headers: {"Authorization": token}
        }).then(res => resolve(res.data)).catch(err => reject(err));
    });
}
