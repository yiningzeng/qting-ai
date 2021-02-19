import request from '../utils/request';
import { stringify } from 'qs';

// const ip="192.168.31.75:8080";//"10.50.102.166";
const ip = `${localStorage.getItem("api.url") === null?"localhost":localStorage.getItem("api.url")}:${localStorage.getItem("api.port") === null?8080:localStorage.getItem("api.port")}`;
// region v1 新的接口
export async function getList_v1(params) {
    console.log("getList"+JSON.stringify(params));
    return request(`http://${ip}/v1/qt_train_record/?${stringify(params)}`, {
        method: 'GET',
    });
}
export async function getProjects_v1() {
    return request(`http://${ip}/v1/qt_projects/`, {
        method: 'GET',
    });
}
export async function getLabels_v1(params) {
    return request(`http://${ip}/v1/qt_labels/?${stringify(params)}`, {
        method: 'GET',
    });
}
export async function getAiFramework_v1(params) {
    return request(`http://${ip}/v1/qt_ai_framework/?${stringify(params)}`, {
        method: 'GET',
    });
}
export async function postTrain_v1(params) {
    console.log("getList"+JSON.stringify(params));
    return request(`http://${ip}/v1/ai/`, {
        method: 'POST',
        body: params
    });
}
export async function delRecord_v1(params) {
    return request(`http://${ip}/v1/qt_train_record/${params.Id}`, {
        method: 'DELETE',
    });
}
export async function getModelsByLabelsAndMulti_v1(params) {
    return request(`http://${ip}/v1/qt_models/GetAllQtModelsByLabelsAndMulti/?${stringify(params)}`, {
        method: 'GET',
    });
}
export async function delModel_v1(params) {
    return request(`http://${ip}/v1/qt_models/${params.Id}`, {
        method: 'DELETE',
    });
}
export async function onlineModel_v1(params) {
    console.log("getList"+JSON.stringify(params));
    return request(`http://${ip}/v1/qt_models/OnlineModel/`, {
        method: 'POST',
        body: params
    });
}
// endregion

