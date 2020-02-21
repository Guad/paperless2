//const host = 'https://qbuy.kolhos.chichasov.es'
const host = '';

export default {
    // called when the user attempts to log in
    login: async ({ username, password }) => {
        const resp = await fetch(host + '/api/login', {
            method: 'POST',
            body: JSON.stringify({
                'username': username,
                'password': password,
            }),
            headers: {
                'Content-Type': 'application/json',
            },
        });
        
        if (!resp.ok) {
            throw new Error('Wrong username or password');
        }

        const json = await resp.json();

        localStorage.setItem('session', json.token);
        
        return Promise.resolve();
    },
    // called when the user clicks on the logout button
    logout: async () => {
        // const key = localStorage.getItem('session');
        // await fetch(host + '/api/v1/session', {
        //     method: 'DELETE',            
        //     headers: {
        //         'Content-Type': 'application/json',
        //         'Authorization': 'Bearer ' + key
        //     },
        // });

        localStorage.removeItem('session');
        return Promise.resolve();
    },
    // called when the API returns an error
    checkError: ({ status }) => {
        if (status === 401 || status === 403) {
            localStorage.removeItem('session');
            return Promise.reject();
        }
        return Promise.resolve();
    },
    // called when the user navigates to a new location, to check for authentication
    checkAuth: () => {
        return localStorage.getItem('session')
            ? Promise.resolve()
            : Promise.reject();
    },
    // called when the user navigates to a new location, to check for permissions / roles
    getPermissions: () => Promise.resolve(),
};
