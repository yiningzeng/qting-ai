import request from '../utils/request';
import { stringify } from 'qs';

// const ip="192.168.31.75:8080";//"10.50.102.166";
// const ip = `${localStorage.getItem("api.url") === null?"localhost":localStorage.getItem("api.url")}:${localStorage.getItem("api.port") === null?8080:localStorage.getItem("api.port")}`;
// const urlPrefix = "..";
const urlPrefix = "http://localhost:8080";
// const urlPrefix = "http://192.168.31.77:8080";
// region v1 新的接口
export async function getVersion_v1() {
    return request(`${urlPrefix}/v1/tools/version`, {
        method: 'GET',
    });
}
export async function getList_v1(params) {
    console.log("getList"+JSON.stringify(params));
    return request(`${urlPrefix}/v1/qt_train_record/?${stringify(params)}`, {
        method: 'GET',
    });
}
export async function getProjects_v1(params) {
    return request(`${urlPrefix}/v1/qt_projects/?${stringify(params)}`, {
        method: 'GET',
    });
}
export async function getLabels_v1(params) {
    return request(`${urlPrefix}/v1/qt_labels/?${stringify(params)}`, {
        method: 'GET',
    });
}
export async function addLabels_v1(params) {
    return request(`${urlPrefix}/v1/qt_labels/`, {
        method: 'POST',
        body: params
    });
}
export async function getAiFramework_v1(params) {
    return request(`${urlPrefix}/v1/qt_ai_framework/?${stringify(params)}`, {
        method: 'GET',
    });
}
export async function AiFramework_v1(params) {
    return request(`${urlPrefix}/v1/qt_ai_framework/${params.method === "PUT" || params.method === "DELETE" ? params.data.Id : ""}`, {
        method: params.method,
        body: params.data,
    });
}
export async function stopTrain_v1(params) {
    return request(`${urlPrefix}//v1/ai/stop/${params.TaskId}`, {
        method: "PUT",
    });
}
export async function postTrain_v1(params) {
    console.log("getList"+JSON.stringify(params));
    return request(`${urlPrefix}/v1/ai/`, {
        method: 'POST',
        body: params
    });
}
export async function delRecord_v1(params) {
    return request(`${urlPrefix}/v1/qt_train_record/${params.Id}`, {
        method: 'DELETE',
    });
}
export async function getModelsByLabelsAndMulti_v1(params) {
    return request(`${urlPrefix}/v1/qt_models/GetAllQtModelsByLabelsAndMulti/?${stringify(params)}`, {
        method: 'GET',
    });
}
export async function delModel_v1(params) {
    return request(`${urlPrefix}/v1/qt_models/${params.Id}`, {
        method: 'DELETE',
    });
}
export async function onlineModel_v1(params) {
    console.log("getList"+JSON.stringify(params));
    return request(`${urlPrefix}/v1/qt_models/OnlineModel/`, {
        method: 'POST',
        body: params
    });
}
// endregion

