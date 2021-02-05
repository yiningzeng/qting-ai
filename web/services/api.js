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


export async function doTrain(params) {
    console.log("getList"+JSON.stringify(params));
    return request(`http://${ip}/train`, {
        method: 'POST',
        body: params
    });
}

export async function getValPathList() {
    return request(`http://${ip}/get_val_path_list`, {
        method: 'GET',
    });
}

export async function getVocPathList() {
    return request(`http://${ip}/get_voc_path_list`, {
        method: 'GET',
    });
}

export async function getModelList(params) {
    return request(`http://${ip}/get_model_list/${params.type}/${params.path}`, {
        method: 'GET',
    });
}

export async function startTest(params) {
    return request(`http://${ip}/start_test`, {
        method: 'POST',
        body: params
    });
}

export async function stopTrain(params) {
    return request(`http://${ip}/stop_train`, {
        method: 'POST',
        body: params
    });
}

export async function continueTrainTrain(params) {
    return request(`http://${ip}/restart_train`, {
        method: 'POST',
        body: params
    });
}

// region 新增接口
export async function getLocalPathList() {
    return request(`http://${ip}/get_local_projects`, {
        method: 'GET',
    });
}
export async function getLabelsByProject(params) {
    return request(`http://${ip}/get_labels/${params.project_name}`, {
        method: 'GET',
    });
}
export async function getLabelsWithScoreByProject(params) {
    return request(`http://${ip}/get_labels_with_info/${params.project_name}`, {
        method: 'GET',
    });
}
export async function get_release_models_history(params) {
    return request(`http://${ip}/get_release_models_history/${params.project_name}`, {
        method: 'GET',
    });
}
export async function getModelListV2(params) {
    return request(`http://${ip}/get_model_list_v2/${params.framework_type}/${params.project_name}`, {
        method: 'GET',
    });
}

export async function del_model(params) {
    console.log("getList"+JSON.stringify(params));
    return request(`http://${ip}/del_model?${stringify(params)}`, {
        method: 'DELETE',
    });
}

export async function get_model_size(params) {
    return request(`http://${ip}/get_model_size?${stringify(params)}`, {
        method: 'GET',
    });
}

export async function online_model(params) {
    return request(`http://${ip}/online_model?${stringify(params)}`, {
        method: 'PUT',
    });
}

export async function suggest_score_put(params) {
    return request(`http://${ip}/suggest_score?${stringify(params)}`, {
        method: 'PUT',
    });
}

export async function suggest_score_get(params) {
    return request(`http://${ip}/suggest_score?${stringify(params)}`, {
        method: 'GET',
    });
}

export async function offline_model(params) {
    console.log("getList"+JSON.stringify(params));
    return request(`http://${ip}/offline_model?${stringify(params)}`, {
        method: 'PUT',
    });
}

// region 多标签，单模型接口
export async function get_models_multilabel(params) {
    return request(`http://${ip}/get_models_multilabel/${params.project_name}`, {
        method: 'GET',
    });
}
export async function get_multilabel_by_model(params) {
    return request(`http://${ip}/get_multilabel_by_model?${stringify(params)}`, {
        method: 'GET',
    });
}
//http://localhost:8080/get_model_list_v2/yolov4-tiny-3l/后道

// endregion

