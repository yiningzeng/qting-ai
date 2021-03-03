import {
    delModel_v1,
    delRecord_v1,
    getAiFramework_v1,
    getLabels_v1,
    getList_v1,
    getModelsByLabelsAndMulti_v1,
    getProjects_v1, onlineModel_v1, postTrain_v1,
    AiFramework_v1,
    getVersion_v1,
    stopTrain_v1,
} from "./api";
const state = {
    // region v1.0 新版本的
    TrainRecords: {
        Code: undefined,
        Msg: undefined,
        Data: [],
    },
    Projects: {
        Code: undefined,
        Msg: undefined,
        Data: [],
    },
    Labels: {
        Code: undefined,
        Msg: undefined,
        Data: [],
    },
    AiFrameworks: {
        Code: undefined,
        Msg: undefined,
        Data: [],
    },
    Models: {
        Code: undefined,
        Msg: undefined,
        Data: [],
    },
    Normal: {
        Code: undefined,
        Msg: undefined,
        Data: undefined,
    },
    // endregion
};
const reducers = {
    // region v1 新版本

    TrainRecords(state, action) {
        return {
            ...state,
            TrainRecords: action.payload,
        };
    },
    Projects(state, action) {
        return {
            ...state,
            Projects: action.payload,
        };
    },
    Labels(state, action) {
        return {
            ...state,
            Labels: action.payload,
        };
    },
    AiFrameworks(state, action) {
        return {
            ...state,
            AiFrameworks: action.payload,
        };
    },
    Models(state, action) {
        return {
            ...state,
            Models: action.payload,
        };
    },
    Normal(state, action) {
        return {
            ...state,
            Normal: action.payload,
        };
    },
    // endregion
};
const effects = {
    // region v1 新版本接口
    * getVersion_v1({payload, callback}, {call, put}) {
        const response = yield call(getVersion_v1, payload);
        if (callback) callback(response);
    },
    * getList_v1({payload, callback}, {call, put}) {
        const response = yield call(getList_v1, payload);
        yield put({
            type: 'TrainRecords',
            payload: response,
        });
        if (callback) callback(response);
    },
    * getProjects_v1({payload, callback}, {call, put}) {
        const response = yield call(getProjects_v1, payload);
        yield put({
            type: 'Projects',
            payload: response,
        });
        if (callback) callback(response);
    },
    * getLabels_v1({payload, callback}, {call, put}) {
        const response = yield call(getLabels_v1, payload);
        yield put({
            type: 'Labels',
            payload: response,
        });
        if (callback) callback(response);
    },
    * getAiFramework_v1({payload, callback}, {call, put}) {
        const response = yield call(getAiFramework_v1, payload);
        yield put({
            type: 'AiFrameworks',
            payload: response,
        });
        if (callback) callback(response);
    },
    * AiFramework_v1({payload, callback}, {call, put}) {
        const response = yield call(AiFramework_v1, payload);
        yield put({
            type: 'Normal',
            payload: response,
        });
        if (callback) callback(response);
    },
    * stopTrain_v1({payload, callback}, {call, put}) {
        const response = yield call(stopTrain_v1, payload);
        yield put({
            type: 'Normal',
            payload: response,
        });
        if (callback) callback(response);
    },
    * getModelsByLabelsAndMulti_v1({payload, callback}, {call, put}) {
        const response = yield call(getModelsByLabelsAndMulti_v1, payload);
        yield put({
            type: 'Models',
            payload: response,
        });
        if (callback) callback(response);
    },
    * postTrain_v1({payload, callback}, {call, put}) {
        const response = yield call(postTrain_v1, payload);
        yield put({
            type: 'Normal',
            payload: response,
        });
        if (callback) callback(response);
    },
    * delRecord_v1({payload, callback}, {call, put}) {
        yield call(delRecord_v1, payload);
        if (callback) callback();
    },
    * delModel_v1({payload, callback}, {call, put}) {
        yield call(delModel_v1, payload);
        if (callback) callback();
    },
    * onlineModel_v1({payload, callback}, {call, put}) {
        const response = yield call(onlineModel_v1, payload);
        if (callback) callback(response);
    },
    // endregion
};
let Models = {
    namespace: 'service',
    state,
    reducers,
    effects,
};
export default Models
