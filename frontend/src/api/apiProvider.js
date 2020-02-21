import { fetchUtils } from 'react-admin';
import { stringify } from 'query-string';


const host = 'https://paperless2.kolhos.chichasov.es'
// const host = '';

const apiUrl = host+'/api';
const sessionKey = () => localStorage.getItem('session');
const httpClient = (url, options) => fetch(url, Object.assign({}, {
    headers: {
        'Content-Type': 'application/json',
        'Authorization': 'Bearer ' + sessionKey(),
    },
    method: 'GET',
}, options));

function flattenObject(ob) {
    var toReturn = {};

    for (var i in ob) {
        if (!ob.hasOwnProperty(i)) continue;

        if ((typeof ob[i]) == 'object' && ob[i] !== null) {
            var flatObject = flattenObject(ob[i]);
            for (var x in flatObject) {
                if (!flatObject.hasOwnProperty(x)) continue;

                toReturn[(i + '.' + x).toLowerCase()] = flatObject[x];
            }
        } else {
            toReturn[i.toLowerCase()] = ob[i];
        }
    }
    return toReturn;
}

function applyTransforms(res) {
    return res;
}

export default {
    getList: async (resource, params) => {
        const { page, perPage } = params.pagination;
        const { field, order } = params.sort;
        const query = {            
            order: order,
            sort: field,
            offset: (page - 1) * perPage,
            limit: perPage,
            filter: JSON.stringify(flattenObject(params.filter)),
        };
        const url = `${apiUrl}/${resource}?${stringify(query)}`;
        const resp = await httpClient(url);
        const json = await resp.json();

        return {
            data: json.data.map(applyTransforms),
            total: json.total,
        };
    },

    getOne: async (resource, params) => {
        const resp = await httpClient(`${apiUrl}/${resource}/${params.id}`);
        const data = await resp.json();        
        return { data: applyTransforms(data) };
    },

    // TODO
    getMany: async (resource, params) => {
        const query = {
            filter: JSON.stringify({ id: params.ids.pop().toString() }),
            offset: 0,
            limit: 10,
            order: "ASC",
            sort: "id",
        };
        
        const url = `${apiUrl}/${resource}?${stringify(query)}`;
        const resp = await httpClient(url);
        const json = await resp.json();

        return {
            data: json.data.map(applyTransforms),
            total: json.total,
        };
    },

    // TODO
    getManyReference: (resource, params) => {
        const { page, perPage } = params.pagination;
        const { field, order } = params.sort;
        const query = {
            sort: JSON.stringify([field, order]),
            range: JSON.stringify([(page - 1) * perPage, page * perPage - 1]),
            filter: JSON.stringify({
                ...params.filter,
                [params.target]: params.id,
            }),
        };
        const url = `${apiUrl}/${resource}?${stringify(query)}`;

        return httpClient(url).then(({ headers, json }) => ({
            data: json,
            total: parseInt(headers.get('content-range').split('/').pop(), 10),
        }));
    },

    update: async (resource, params) => {
        const resp = await httpClient(`${apiUrl}/${resource}/${params.id}`, {
            method: 'PUT',
            body: JSON.stringify(params.data),
        });
        const data = await resp.json();

        return { data: applyTransforms(data) };
    },

    // TODO
    updateMany: (resource, params) => {
        const query = {
            filter: JSON.stringify({ id: params.ids}),
        };
        return httpClient(`${apiUrl}/${resource}?${stringify(query)}`, {
            method: 'PUT',
            body: JSON.stringify(params.data),
        }).then(({ json }) => ({ data: json }));
    },

    create: async (resource, params) => {
        const resp = await httpClient(`${apiUrl}/${resource}`, {
            method: 'POST',
            body: JSON.stringify(params.data),
        });
        const json = await resp.json();
        return { data: applyTransforms(json) };
    },

    delete: async (resource, params) => {
        const resp = await httpClient(`${apiUrl}/${resource}/${params.id}`, {
            method: 'DELETE',
        });
        const data = await resp.json();

        return { data: applyTransforms(data) };
    },

    deleteMany: async (resource, params) => {
        const resp = await httpClient(`${apiUrl}/${resource}/${params.ids.join(',')}`, {
            method: 'DELETE',
        });
        const data = await resp.json();
        return { data: applyTransforms(data) };
    }
};
