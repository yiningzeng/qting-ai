import React from 'react';
import ReactDOM from 'react-dom';
import Iframe from 'react-iframe';
import { SettingOutlined, SmileTwoTone, CloudUploadOutlined } from '@ant-design/icons';
import '@ant-design/compatible/assets/index.css';
import dva, { connect } from 'dva';
import {
    InputNumber,
    Tabs,
    Tag,
    Row,
    Modal,
    Spin,
    Col,
    Collapse,
    Table,
    message,
    PageHeader,
    Button,
    Typography,
    Drawer,
    Divider,
    Select,
    Switch,
    Input,
    notification,
    Radio,
    Badge,
    Popconfirm,
    Image,
    Form,
    Empty
} from 'antd';
// 由于 antd 组件的默认文案是英文，所以需要修改为中文
import zhCN from 'antd/lib/locale-provider/zh_CN';
import moment from 'moment';
import 'moment/locale/zh-cn';
import { getList_v1,
    getProjects_v1,
    getLabels_v1,
    getAiFramework_v1,
    postTrain_v1,
    getModelsByLabelsAndMulti_v1,

    delModel_v1,
    delRecord_v1,
    onlineModel_v1,

    getModelList, getValPathList, getVocPathList,
    startTest, stopTrain, continueTrainTrain, getLocalPathList, getModelListV2,
    get_release_models_history, del_model, online_model, offline_model,
    getLabelsWithScoreByProject, suggest_score_get, suggest_score_put, get_model_size, get_models_multilabel, get_multilabel_by_model } from './services/api';
const { Title, Paragraph, Text } = Typography;
const { Option } = Select;
const InputGroup = Input.Group;
moment.locale('zh-cn');
const { confirm } = Modal;
const { TabPane } = Tabs;
/**
 *
 */
class FreeFish extends React.Component {
    state = {
        fuck: false,
        lossImgPreviewVisible: false,
        showSettingsModal: false,
        imageLossTimer: "zengyining",
        test: {
            frontImage: "registry.cn-hangzhou.aliyuncs.com/baymin/darknet-test:",
            baseImage: "latest",
            showTestDrawer: false,
            showTestDrawerUrl: "",
            showTestModal: false,
            showStandardValidationData: false,
            loading: false,
            tips: "载入可使用的权重文件... ",
            doTest: {
                providerType: "QTing-tiny-3l-single",
                assetsDir: "", // nowAssetsDir
                weights: undefined,
                valPath: undefined,
                port: 8100,
                javaUrl: "ai.8101.api.qtingvision.com",
                javaPort: 888,
                image: "registry.cn-hangzhou.aliyuncs.com/baymin/darknet-test:latest",
                projectId: undefined,
            },
        },
        timer: null,
        refreshInterval: localStorage.getItem("refreshInterval") === null?5000:localStorage.getItem("refreshInterval"),
        refreshTime: moment().format("YYYY-MM-DD HH:mm:ss"),
        selectedRowKeys: [], // Check here to configure the default column
        project_labels: [],
        loadingChart: false,
        pagination: {defaultPageSize:100, current:1},
        loading: false,
        leftVisible: false,
        rightVisible: false,
        doChangeAssetsDir: true,

        api: {
            url: localStorage.getItem("api.url") === null?"localhost":localStorage.getItem("api.url"),
            port: localStorage.getItem("api.port") === null?8080:localStorage.getItem("api.port"),
        },
        train : {
            frontImage: "registry.cn-hangzhou.aliyuncs.com/qtingvision/auto-train-tiny:",
            baseImage: "latest",
            showAiPar: false,
            loading: false,
            doTrain: {
                aiFrameworkId: undefined, // 框架id
                taskId: undefined, // 项目id
                taskName: undefined, // 训练任务名称
                projectName: undefined, // 项目名称
                projectId: undefined, // 项目id
                assetsDir: undefined, // 素材文件夹，和packageDir相同
                assetsType: "powerAi", // 素材的类型，pascalVoc和coco和other
                providerType: "QTing-tiny-3l-single", // 框架的类型yolov3
                batchSize: 64,
                imageWidth: -1,
                imageHeight: -1,
                maxIter: -1, // 训练最大轮数
                pretrainWeight: "", // 预训练权重文件
                gpus: "0,1", // 使用的gpu id
                trianType: 0,  // 0对应从头训练 1对应自训练 2 漏检训练
                singleTrain: [], // 单类训练名称 ‘’则全类训练，不为空则训练输入的单类，确保单类名在标记标签中

                angle: 360,
                cell_stride: 1, //正平移步长
                cellsize: 16, //平移框大小
                expand_size: [8, 8], //扩展尺寸
                ignore_size: [6, 6],//忽略尺寸
                resizearrange: [0.4, 1.2], // anchor  调整变化幅度
                trainwithnolabelpic: 1000, //最大负样本数

                subdivisionssize: 2,
                rmgeneratedata: 0, // 是否保留训练生成的临时数据
                split_ratio: 0.95, // 训练样本占比
                recalldatum: 2, // 检出率基准
                otherlabeltraintype: 1, // 非当前标签图片训练方式
                mergeTrainSymbol: 0, // 是否合并多标签训练

                learnrate: 0.00261, // 学习率
                otherlabelstride: 1, // 负样本平移增强步长
                isshuffle: true, // 是否打乱数据
            },
        },
        suggestScore: {
            loading: false,
            maxDetPerdm: undefined,
            pixel2realLength: undefined,
        },
        continueTrain: {
            showModal: false,
            loading: false,
            frontImage: "registry.cn-hangzhou.aliyuncs.com/qtingvision/auto-train:",
            baseImage: "latest",
            width: undefined,
            height: undefined,
            max_batches: undefined,
            projectId: undefined, // 项目id
            assetsType: "powerAi", // 素材的类型，pascalVoc和coco和other
            projectName: undefined, // 项目名称
            providerType: "QTing-tiny-3l-single", // 框架的类型yolov3 fasterRcnn maskRcnn
            image: "registry.cn-hangzhou.aliyuncs.com/qtingvision/auto-train:latest", // 镜像路径
            assetsDir: "", //nowAssetsDir
            weights: undefined,
        },
        modelManagerSingle: {
            loadingProjects: false,
            expandedRowKeys: undefined,
            nowEditProjectName: undefined,
            loadingModels: true,
            firstVisible: false,
            secondVisible: false,
            secondReleaseManagerVisible: false,
        },
        modelManagerMultilabel: {
            loadingModels: true,
            loadingLabels: false,
            expandedRowKeys: undefined,
            nowEditProjectName: undefined,
            firstVisible: false,
            secondVisible: false,
            secondReleaseManagerVisible: false,
        },
        publishModal: {
            visible: false,
            ModelWidth: undefined,
            ModelHeight: undefined,
            LabelName: undefined,
            ModelId: undefined,
            refresh: undefined,

            // 下面这两个主要用于发布之后的更新页面操作
            label: undefined,
            projectId: undefined,
        },
    };

    componentDidMount() {

    }

    componentWillMount() {
        message.success(`正在加载`);
        const {dispatch} = this.props;
        dispatch({
            type: 'service/getList_v1',
            payload: {
                sortby: "CreateTime",
                order: "desc",
                offset: 0,
                limit: 200,
            },
            callback: (v) => {
                // console.log(`加载：${JSON.stringify(v)}`);
                //<Pagination
                //       total={85}
                //       showTotal={total => `Total ${total} items`}
                //       pageSize={20}
                //       defaultCurrent={1}
                //     />
                // noinspection JSAnnotator
                this.setState({
                    ...this.state,
                    pagination:{
                        ...this.state.pagination,
                        total: v["total"],
                        pageSize: v["num"],
                        current: 1,
                    }
                });
            },
        });
        this.state.timer=setInterval(()=>{
            dispatch({
                type: 'service/getList_v1',
                payload: {
                    sortby: "CreateTime",
                    order: "desc",
                    offset: 0,
                    limit: 200,
                    // offset: this.state.pagination.current,
                    // limit: this.state.pagination.defaultPageSize,
                },
                callback: (v) => {
                    console.log(`加载：${JSON.stringify(v)}`);
                    this.setState({
                        ...this.state,
                        refreshTime: moment().format("YYYY-MM-DD HH:mm:ss"),
                        pagination:{
                            ...this.state.pagination,
                            total: v["total"],
                            pageSize: v["num"],
                        }
                    });
                },
            });
            console.log("我是定时任务");
        }, this.state.refreshInterval);
    }

    componentWillUnmount() {
        if(this.state.timer!= null) {
            clearInterval(this.state.timer);
        }
    }

    handleTableChange = (pagination, filters, sorter) => {
        const pager = { ...this.state.pagination };
        pager.current = pagination.current;
        this.setState({
            pagination: pager,
        });
        console.log(`页码：${JSON.stringify(pager)}`);

        const { dispatch } = this.props;
        dispatch({
            type: 'service/getList_v1',
            payload: {
                sortby: "CreateTime",
                order: "desc",
                offset: 0,
                limit: 200,
                // page: pager.current,
                // num: pager.pageSize,
            },
            callback: (v) => {
                console.log(`${JSON.stringify(v)}`);
            },
        });
    };

    showLeftDrawer = () => {
        this.setState({
            leftVisible: true,
        });
    };

    hideLeftDrawer = () => {
        this.setState({
            leftVisible: false,
        });
    };



    render() {
        const {
            service: {modelList, modelListV2, valPathList, vocPathList, modelByProject, labelsByProject,labelsWithScoreByProject, release_models_history_res, models_multilabel, multilabel_by_model}
        } = this.props;

        const expandedRowRender = (record) => {
            return <div>
                    <Row>
                        <Button type="primary" size="small" style={{marginLeft: 10}}
                                disabled={record.status !== 2} onClick={() => {
                            confirm({
                                title: '提示',
                                content: '确定要停止训练么?',
                                onOk: () => {
                                    const {dispatch} = this.props;
                                    dispatch({
                                        type: 'service/stopTrain',
                                        payload: {
                                            task_id: record.task_id,
                                            project_name: record.project_name
                                        },
                                        callback: (v) => {
                                            if (v.res === "ok") {
                                                message.success("已经停止训练");
                                                dispatch({
                                                    type: 'service/getList_v1',
                                                    payload: {
                                                        sortby: "CreateTime",
                                                        order: "desc",
                                                        offset: 0,
                                                        limit: 200,
                                                        // page: this.state.pagination.current,
                                                        // num: this.state.pagination.defaultPageSize,
                                                    },
                                                    callback: (v) => {
                                                        console.log(`加载：${JSON.stringify(v)}`);
                                                        this.setState({
                                                            ...this.state,
                                                            refreshTime: moment().format("YYYY-MM-DD HH:mm:ss"),
                                                            pagination: {
                                                                ...this.state.pagination,
                                                                total: v["total"],
                                                                pageSize: v["num"],
                                                            }
                                                        });
                                                    },
                                                });
                                            } else {
                                                message.error("停止训练失败");
                                            }
                                        }
                                    });
                                },
                                onCancel() {
                                },
                            });
                        }}>停止训练</Button>
                        <Button type="primary" size="small" style={{marginLeft: 10}} onClick={() => {
                            this.setState({
                                ...this.state,
                                imageLossTimer: moment().valueOf(),
                            })
                        }}>刷新Loss</Button>
                        <Button type="primary" size="small" style={{marginLeft: 10}} onClick={() => {
                            this.setState({
                                ...this.state,
                                lossImgPreviewVisible: true,
                            })
                        }}>预览图表</Button>
                        <Popconfirm
                            disabled={record.status === 2 || record.status === 0 || record.status === 1}
                            title="只是删除记录，该记录训练的模型不会被删除，确定要删除么？"
                            onConfirm={() => {
                                const {dispatch} = this.props;
                                dispatch({
                                    type: 'service/delRecord_v1',
                                    payload: {
                                        Id: encodeURI(record.Id)
                                    },
                                    callback: (aa) => {
                                        message.success("删除成功");
                                        location.reload();
                                    }
                                });
                            }}
                            okText="确定"
                            cancelText="取消"
                        >
                            <Button type="primary" disabled={record.status === 2 || record.status === 0 || record.status === 1} danger size="small" style={{marginLeft: 10}} >删除该记录</Button>
                        </Popconfirm>
                        {/*<Button type="primary" size="small" style={{marginLeft: 10}}>日志</Button>*/}
                    </Row>
                    <Row>
                        {
                            record.AiFrameworkId.FrameworkName === "QTing-tiny-3l-multilabel" &&
                            <Tabs defaultActiveKey="1" tabPosition="left" style={{ height: 720, marginTop: 20 }}>
                                <TabPane tab={`所有标签`}>
                                    <Col span={12} offset={4}>
                                        <Image
                                            width={"150%"}
                                            height={700}
                                            src={record.DrawUrl + "?id=" + moment().valueOf() + this.state.imageLossTimer}
                                            preview={{
                                                visible: this.state.lossImgPreviewVisible,
                                                onVisibleChange: (visible, prevVisible) => {
                                                    this.setState({
                                                        ...this.state,
                                                        lossImgPreviewVisible: visible,
                                                    })
                                                },
                                            }}
                                            // fallback=""
                                        />
                                    </Col>
                                </TabPane>
                            </Tabs>
                        }
                        {
                            record.AiFrameworkId.FrameworkName !== "QTing-tiny-3l-multilabel" && <Spin spinning={this.state.loadingChart}>
                                <Tabs defaultActiveKey="1" tabPosition="left" style={{ height: 720, marginTop: 20 }}>
                                    {this.state.project_labels.map(i => (
                                        <TabPane tab={`${i}`} key={i}>
                                            <Col span={12} offset={4}>
                                                <Image
                                                    width={"150%"}
                                                    height={700}
                                                    src={record.DrawUrl.replace(/\/chart.png/, `-${i}/chart.png`) + "?id=" + moment().valueOf() + this.state.imageLossTimer}
                                                    preview={{
                                                        visible: this.state.lossImgPreviewVisible,
                                                        onVisibleChange: (visible, prevVisible) => {
                                                            this.setState({
                                                                ...this.state,
                                                                lossImgPreviewVisible: visible,
                                                            })
                                                        },
                                                    }}
                                                    // fallback=""
                                                />
                                            </Col>
                                        </TabPane>
                                    ))}
                                </Tabs>
                            </Spin>
                        }
                    </Row>
                </div>
        };

        const expandedModelsRowRender = (mainRecord, index, indent, expanded) => {
            return  <div>
                <Spin tip="正在加载..." spinning={this.state.modelManagerSingle.loadingModels}>
                    {/*<Badge status="processing" text="Running" />*/}
                    <Table columns={[
                        {
                            title: '模型名称',
                            key: "ModelName",
                            dataIndex: 'ModelName',
                        },  {
                            title: '推荐置信度',
                            render: (text, record) => {
                                try {
                                    return JSON.parse(record.SuggestScore)[mainRecord.LabelName];
                                } catch (e) {
                                    return "";
                                }
                            },
                        }, {
                            title: '发布状态',
                            key: "status",
                            dataIndex: 'Status',
                            render: v => {
                                if (v === 2) return <Badge status="processing" text="已发布"/>;
                                else if (v === 1) return <Badge status="warning" text="已下线"/>;
                            },
                        }, {
                            title: '发布日期',
                            render: (text, record) => {
                                if (record.Status === 2) return <Tag icon={<CloudUploadOutlined/>} color="success">{record.PublishTime}</Tag>;
                                else return record.PublishTime === "0001-01-01T00:00:00Z" ? "" : record.PublishTime;
                            }
                        }, {
                            title: '操作',
                            render: (text, record) => (
                                <span>
                                    {
                                        record.Status !== 2 &&  <Popconfirm
                                            title="发布会把最新的模型替换为当前发布的模型，确定发布么？"
                                            onConfirm={() => {
                                                this.setState({
                                                    ...this.state,
                                                    publishModal: {
                                                        ...this.state.publishModal,
                                                        visible: true,
                                                        ModelWidth: record.ProjectId.ImageWidth,
                                                        ModelHeight: record.ProjectId.ImageHeight,
                                                        LabelName: mainRecord.LabelName,
                                                        ModelId: record.Id,
                                                        refresh: "单标签",
                                                        projectId: record.ProjectId.Id,
                                                        label: mainRecord.LabelName,
                                                    }
                                                });
                                            }}
                                            okText="确定并打开设置"
                                            cancelText="取消"
                                        >
                                            <a>{record.Status === 1 ? "上线" : "发布"}</a>
                                        </Popconfirm>
                                    }
                                    {
                                        record.Status !== 2 && <Divider type="vertical"/>
                                    }
                                    {
                                        record.Status !== 2 && <Popconfirm
                                            title="确定要删除么？"
                                            onConfirm={() => {
                                                const {dispatch} = this.props;
                                                dispatch({
                                                    type: 'service/delModel_v1',
                                                    payload: {
                                                        Id: encodeURI(record.Id)
                                                    },
                                                    callback: (aa) => {
                                                        const {dispatch} = this.props;
                                                        dispatch({
                                                            type: 'service/getModelsByLabelsAndMulti_v1',
                                                            payload: {
                                                                projectId: record.ProjectId.Id,
                                                                label: mainRecord.LabelName,
                                                                isMultilabel: 0,
                                                            },
                                                        });
                                                        notification.success({
                                                            message: "恭喜",
                                                            description: "删除成功",
                                                        });
                                                    }
                                                });
                                            }}
                                            okText="确定"
                                            cancelText="取消"
                                        >
                                            <a>删除</a>
                                        </Popconfirm>
                                    }


                                </span>),
                        }]} dataSource={this.props.service.Models === undefined ? [] : this.props.service.Models.Data}/>
                </Spin>
            </div>;
        };
        const expandedLabelsRowRender = (mainRecord, index, indent, expanded) => {
            const data = [];
            try {
                const strTemp = mainRecord.LabelStr.toString();
                const strSp = strTemp.substring(1, strTemp.length - 1).split(",");
                strSp.map(value => {
                    try{
                        data.push({"LabelName": value, "SuggestScore": JSON.parse(mainRecord.SuggestScore)[value]})
                    } catch (e) {
                        console.error(e)
                    }
                });
            } catch (e) {
                console.error(e)
            }

            return  <div>
                <Spin tip="正在加载..." spinning={this.state.modelManagerMultilabel.loadingLabels}>
                    {/*<Badge status="processing" text="Running" />*/}
                    <Table columns={[
                        {
                            title: '标签名称',
                            key: "LabelName",
                            dataIndex: 'LabelName',
                        },  {
                            title: '推荐置信度',
                            key: "SuggestScore",
                            dataIndex: 'SuggestScore',
                        }]} dataSource={data}/>
                </Spin>
            </div>;
        };

        return (
            <PageHeader backIcon={false}
                title="训练中心"
                subTitle="管理后台"
                tags={<Tag color="green">在线</Tag>}
                extra={[
                    <span>{`页面刷新时间:`}</span>,
                    <Text mark>{`${this.state.refreshTime}`}</Text>,
                    <span>刷新间隔(秒):</span>,
                    <Select defaultValue={`${this.state.refreshInterval / 1000}s`} style={{width: 120}}
                            onChange={(v) => {
                                localStorage.setItem("refreshInterval", v);
                                location.reload();
                            }}>
                        <Option value="5000">5s</Option>
                        <Option value="10000">10s</Option>
                        <Option value="30000">30s</Option>
                        <Option value="60000">60s</Option>
                        <Option value="6000000">6000s</Option>
                    </Select>,
                    <Button key="1" type="primary" onClick={() => {
                        this.setState({
                            ...this.state,
                            modelManagerSingle: {
                                ...this.state.modelManagerSingle,
                                firstVisible: true
                            },
                        }, () => {
                            // 首先加载相应的数据
                            const {dispatch} = this.props;
                            dispatch({
                                type: 'service/getProjects_v1',
                                callback: (aa) => {
                                    this.setState({
                                        ...this.state,
                                        suggestScore: {
                                            ...this.state.suggestScore,
                                            maxDetPerdm: undefined,
                                            pixel2realLength: undefined,
                                        }});
                                }
                            });
                        });
                    }}>
                        模型管理
                    </Button>,
                    <Button.Group style={{marginLeft: 7}}>
                        <Button key="1" type="primary" onClick={() => {
                            // 首先加载相应的数据
                            const {dispatch} = this.props;
                            dispatch({
                                type: 'service/getProjects_v1',
                                callback: (aa) => {
                                    this.setState({
                                        train: {
                                            ...this.state.train,
                                            doTrain: {
                                                ...this.state.train.doTrain,
                                                projectId: aa.Data.length > 0 ?  aa.Data[0].Id : undefined,
                                                projectName: aa.Data.length > 0 ? aa.Data[0].ProjectName : undefined,
                                                taskId: `${moment().format('YYYYMMDDHHmmss')}`
                                            }
                                        }
                                    }, () => {
                                        // 获取标签信息
                                        dispatch({
                                            type: 'service/getLabels_v1',
                                            payload: {
                                                query: encodeURI(`ProjectId:${this.state.train.doTrain.projectId}`),
                                                limit: 1000,
                                            },
                                            callback: (bb) => {
                                                // 获取框架信息
                                                dispatch({
                                                    type: 'service/getAiFramework_v1',
                                                    callback: (cc) => {
                                                        const singleTrain = [];
                                                        if (bb.Data !== null) {
                                                            bb.Data.map(v => {
                                                                singleTrain.push(v.LabelName);
                                                            });
                                                        }
                                                        this.setState({
                                                            ...this.state,
                                                            rightVisible: true,
                                                            leftVisible: true,
                                                            train: {
                                                                ...this.state.train,
                                                                doTrain: {
                                                                    ...this.state.train.doTrain,
                                                                    singleTrain: singleTrain,
                                                                    aiFrameworkId: cc.Data !== null ? cc.Data[0].Id : undefined,
                                                                    mergeTrainSymbol: 0,
                                                                    providerType: "QTing-tiny-3l-single",
                                                                }
                                                            }
                                                        });
                                                    }
                                                });
                                            }
                                        });
                                    });
                                }
                            });
                        }}>
                            新增训练任务
                        </Button>
                        <Button type="primary" onClick={() => this.setState({
                            ...this.state,
                            showSettingsModal: true,
                        })}>
                            <SettingOutlined />
                        </Button>
                    </Button.Group>,
                ]}
                footer={
                    <Table
                        rowKey="TaskId"
                        columns={[{
                            title: '任务标识',
                            dataIndex: 'TaskId',
                        }, {
                            title: '训练任务名称',
                            dataIndex: 'TaskName',
                        }, {
                            title: '对应项目名称',
                            dataIndex: 'ProjectId',
                            render: v => {
                                return v.ProjectName;
                            }
                        }, {
                            title: '网络框架',
                            dataIndex: 'AiFrameworkId',
                            render: v => {
                                return v.FrameworkName;
                            }
                        }, {
                            title: '数据类型',
                            dataIndex: 'AssetsType',
                        }, {
                            title: '创建时间',
                            dataIndex: 'CreateTime',
                        }, {
                            title: '状态',
                            dataIndex: 'Status',
                            render: v => {
                                if (v === 0) return <Tag color="#FFA500">准备完成</Tag>;
                                else if (v === 1) return <Tag color="#8A2BE2">等待训练</Tag>;
                                else if (v === 2) return <Button type="primary" loading>正在训练</Button>;
                                else if (v === 3) return <Tag color="#D3D3D3">停止训练</Tag>;
                                else if (v === 4) return <div><Tag color="#008000">训练完成</Tag><SmileTwoTone /></div>;
                                else if (v === -1) return <Tag color="#FF0000">训练出错</Tag>;
                                else return <Tag>未知</Tag>;
                            }
                        }]}
                        dataSource={this.props.service.TrainRecords.Data}
                        onChange={this.handleTableChange}
                        expandable={{
                            expandedRowRender,
                            expandRowByClick: true,
                            expandedRowKeys: this.state.selectedRowKeys,
                            onExpandedRowsChange: (expandedRows) => {
                            },
                            onExpand: (expanded, record) => {
                                if (!expanded) {
                                    this.setState({
                                        ...this.state,
                                        selectedRowKeys: [],
                                        loadingChart: false,
                                    });
                                    return;
                                }
                                this.setState({
                                    ...this.state,
                                    selectedRowKeys: [record.TaskId],
                                    loadingChart: true,
                                }, () => {
                                    const {dispatch} = this.props;
                                    dispatch({
                                        type: 'service/getLabels_v1',
                                        payload: {
                                            query: encodeURI(`ProjectId:${record.ProjectId.Id}`),
                                            limit: 1000,
                                        },
                                        callback: (bb) => {
                                            const singleTrain = [];
                                            if (bb.Data !== null) {
                                                bb.Data.map(v => {
                                                    singleTrain.push(v.LabelName);
                                                });
                                            }
                                            this.setState({
                                                ...this.state,
                                                project_labels: singleTrain,
                                                loadingChart: false,
                                            });
                                        }
                                    });
                                });
                            }
                        }}
                        pagination={this.state.pagination}
                    />
                }
            >
                <div className="wrap">
                    <Modal
                        title="设置"
                        okText="保存"
                        cancelText="取消"
                        destroyOnClose
                        visible={this.state.showSettingsModal}
                        onOk={() => {
                            localStorage.setItem("api.url", this.state.api.url);
                            localStorage.setItem("api.port", this.state.api.port);
                            message.success("保存成功");
                            location.reload();
                        }}
                        onCancel={() => {
                            this.setState({
                                ...this.state,
                                showSettingsModal: false,
                            });
                        }}
                    >
                        接口地址:
                        <InputGroup style={{marginTop: "10px", marginBottom: "20px"}} compact>
                            <Input style={{width: '70%'}} addonBefore="http://" value={this.state.api.url}
                                   onChange={e => {
                                       this.setState({
                                           ...this.state,
                                           api: {
                                               ...this.state.api,
                                               url: e.target.value,
                                           }
                                       });
                                   }}
                                   placeholder="网址不带http://" allowClear/>
                            <Input
                                style={{
                                    width: 30,
                                    borderLeft: 0,
                                    pointerEvents: 'none',
                                    backgroundColor: '#fff',
                                }}
                                placeholder=":"
                                disabled
                            />
                            <Input style={{width: '15%', textAlign: 'center', borderLeft: 0}}
                                   value={this.state.api.port} onChange={(e) => {
                                this.setState({
                                    ...this.state,
                                    api: {
                                        ...this.state.api,
                                        port: e.target.value,
                                    }
                                });
                            }} defaultValue={this.state.apiPort} placeholder="port"/>
                        </InputGroup>
                    </Modal>
                    <Modal
                        title="参数设置"
                        okText="打开测试服务"
                        cancelText="取消"
                        destroyOnClose
                        width={1000}
                        visible={this.state.test.showTestModal}
                        onOk={() => {
                            const {dispatch} = this.props;
                            this.setState({
                                    ...this.state,
                                    test: {
                                        ...this.state.test,
                                        tips: "正在打开服务",
                                        // showTestModal: false,
                                        loading: true,
                                    }
                                },
                                () => {
                                    dispatch({
                                        type: 'service/startTest',
                                        payload: {
                                            ...this.state.test.doTest
                                        },
                                        callback: (v) => {
                                            if (v.res !== "ok") {
                                                message.error(v.msg);
                                            }

                                            this.setState({
                                                ...this.state,
                                                test: {
                                                    ...this.state.test,
                                                    showTestDrawer: true,
                                                    showTestModal: false,
                                                    loading: false,
                                                    doTest: {
                                                        ...this.state.test.doTest,
                                                        weights: undefined,
                                                    },
                                                    // showTestDrawerUrl: `/test?javaUrl=${}&javaPort=${}&providerType=${this.state.test.doTest.providerType}&port=${this.state.test.port}&assets=${this.state.test.doTest.assetsDir}`,
                                                    showTestDrawerUrl: `/test?projectId=${this.state.test.doTest.projectId}&javaUrl=${this.state.test.doTest.javaUrl}&javaPort=${this.state.test.doTest.javaPort}&providerType=${this.state.test.doTest.providerType}&port=${this.state.test.doTest.port}&assets=${this.state.test.doTest.assetsDir}`,
                                                }
                                            });

                                            // const tempwindow=window.open();
                                            // tempwindow.location=`/test?port=8100&assets=${this.state.test.nowAssetsDir}`;
                                            // window.open(`/test?port=8100&assets=${this.state.test.nowAssetsDir}`, "_blank");
                                            // window.open(`/test?port=8100&assets=${this.state.test.nowAssetsDir}`, "_blank", "scrollbars=yes,resizable=1,modal=false,alwaysRaised=yes");
                                        }
                                    });
                                });
                        }}
                        onCancel={() => {
                            this.setState({
                                ...this.state,
                                test: {
                                    ...this.state.test,
                                    showTestModal: false,
                                    doTest: {
                                        ...this.state.test.doTest,
                                        weights: undefined,
                                    }
                                }
                            });
                        }}
                        okButtonProps={{disabled: this.state.test.doTest.weights === undefined}}
                        cancelButtonProps={{disabled: false}}
                    >
                        <Spin spinning={this.state.test.loading} tip={this.state.test.tips} delay={500}>
                            网络框架:
                            <Input style={{width: '100%', marginTop: "10px", marginBottom: "10px"}}
                                   placeholder="Basic usage" disabled value={this.state.test.doTest.providerType}/>
                            服务端口:
                            <Input style={{width: '100%', marginTop: "10px", marginBottom: "10px"}}
                                   placeholder="Basic usage" disabled value={this.state.test.doTest.port}/>
                            选择加载的权重文件:
                            <Select
                                style={{width: '100%', marginTop: "10px", marginBottom: "10px"}}
                                onChange={(v) => {
                                    this.setState({
                                        ...this.state,
                                        test: {
                                            ...this.state.test,
                                            doTest: {
                                                ...this.state.test.doTest,
                                                weights: v
                                            }
                                        }
                                    });
                                }}>
                                {modelList.weights_list.map(d => (
                                    <Option key={d.path}>{d.filename}</Option>
                                ))}
                            </Select>
                            镜像地址:
                            <Input style={{marginTop: "10px", marginBottom: "20px"}} placeholder="tar压缩包名"
                                   addonBefore={this.state.test.frontImage}
                                   value={this.state.test.baseImage} allowClear
                                   onChange={(e) => this.setState({
                                       ...this.state,
                                       test: {
                                           ...this.state.test,
                                           doTest: {
                                               ...this.state.test.doTest,
                                               image: `${this.state.test.frontImage}${e.target.value}`
                                           }
                                       }
                                   })}/>
                            标准验证集:&nbsp;&nbsp;
                            <Switch checkedChildren="使用" unCheckedChildren="不使用"
                                    onChange={(c) => {
                                        this.setState({
                                            ...this.state,
                                            test: {
                                                ...this.state.test,
                                                showStandardValidationData: c,
                                            }
                                        })
                                    }}/>
                            <br/>
                            {
                                this.state.test.showStandardValidationData && <div>选择加载的标准验证集目录:<Select
                                    style={{width: '100%', marginTop: "10px", marginBottom: "10px"}}
                                    onChange={(v) => {
                                        this.setState({
                                            ...this.state,
                                            test: {
                                                ...this.state.test,
                                                doTest: {
                                                    ...this.state.test.doTest,
                                                    valPath: v
                                                }
                                            }
                                        });
                                    }}>
                                    {valPathList.val_path_list.map(d => (
                                        <Option key={d.path}>{d.dir_name}</Option>
                                    ))}
                                </Select>
                                </div>
                            }
                        </Spin>
                    </Modal>
                    <Modal
                        maskClosable={false}
                        destroyOnClose
                        title="模型发布设置"
                        visible={this.state.publishModal.visible}
                        onOk={() => {
                            const {dispatch} = this.props;
                            dispatch({
                                type: 'service/onlineModel_v1',
                                payload: {
                                    ModelHeight: this.state.publishModal.ModelHeight,
                                    ModelId: this.state.publishModal.ModelId,
                                    ModelWidth: this.state.publishModal.ModelWidth,
                                    Label: this.state.publishModal.label,
                                },
                                callback: (aa) => {
                                    this.setState({
                                        ...this.state,
                                        publishModal: {
                                            ...this.state.publishModal,
                                            visible: false,
                                        },
                                    });
                                    if (aa["Code"] == 200) {
                                        notification.success({
                                            message: "提醒",
                                            description: "已经发布成功，线上的AOI软件可以通过更新直接获取当前模型",
                                        });
                                    } else {
                                        notification.error({
                                            message: "发布失败",
                                            description: aa["Data"],
                                        });
                                        return;
                                    }
                                    if (this.state.publishModal.refresh === "单标签") {
                                        dispatch({
                                            type: 'service/getModelsByLabelsAndMulti_v1',
                                            payload: {
                                                projectId: this.state.publishModal.projectId,
                                                label: this.state.publishModal.label,
                                                isMultilabel: 0,
                                            }});
                                    } else if (this.state.publishModal.refresh === "多标签") {
                                        dispatch({
                                            type: 'service/getModelsByLabelsAndMulti_v1',
                                            payload: {
                                                projectId: this.state.publishModal.projectId,
                                                isMultilabel: 1,
                                            }});
                                    }
                                }
                            });
                        }}
                        onCancel={() => {
                            this.setState({
                                ...this.state,
                                publishModal: {
                                    ...this.state.publishModal,
                                    visible: false,
                                },
                            });
                        }}
                        okText="发布"
                        cancelText="取消"
                    >
                        网络图像宽度:
                        <InputNumber style={{marginTop: "5px", marginBottom: "10px", width: "100%"}} value={this.state.publishModal.ModelWidth} placeholder="留空表示新建项目时保留的值"
                                     min={1}
                                     allowClear onChange={(value) => this.setState({
                            ...this.state,
                            publishModal: {
                                ...this.state.publishModal,
                                ModelWidth: value
                            }
                        })}/>
                        网络图像高度:
                        <InputNumber style={{marginTop: "5px", marginBottom: "10px", width: "100%"}} value={this.state.publishModal.ModelHeight} placeholder="留空表示新建项目时保留的值"
                                     min={1}
                                     allowClear onChange={(value) => this.setState({
                            ...this.state,
                            publishModal: {
                                ...this.state.publishModal,
                                ModelHeight: value
                            }
                        })}/>
                    </Modal>
                    <Drawer
                        title="在线测试"
                        placement="left"
                        width="100%"
                        height="1500px"
                        closable={true}
                        onClose={() => {
                            this.setState({
                                ...this.state,
                                test: {
                                    ...this.state.test,
                                    showTestDrawer: false,
                                }
                            })
                        }}
                        visible={this.state.test.showTestDrawer}
                    >
                        <Iframe url={this.state.test.showTestDrawerUrl}
                                width="100%"
                                height="850px"
                                id="myId"
                                frameBorder={0}
                                className="myClassname"
                                display="initial"
                                position="relative"/>
                        {/*<Iframe url={this.state.test.showTestDrawerUrl}*/}
                        {/*width="100%"*/}
                        {/*height="500px"*/}
                        {/*id="myId"*/}
                        {/*frameBorder={0}*/}
                        {/*className="myClassname"*/}
                        {/*display="initial"*/}
                        {/*position="relative"/>*/}
                    </Drawer>

                    <Drawer
                        destroyOnClose={true}
                        title="新增训练任务"
                        placement="right"
                        width="40%"
                        closable={true}
                        maskClosable={true}
                        onClose={() => {
                            this.setState({
                                rightVisible: false,
                                leftVisible: false,
                            });
                        }}
                        visible={this.state.rightVisible}
                    >
                        训练任务名称:
                        <Input style={{marginTop: "5px", marginBottom: "10px"}} placeholder="训练任务名称"
                               allowClear onChange={(e) => this.setState({
                            train: {
                                ...this.state.train,
                                doTrain: {
                                    ...this.state.train.doTrain,
                                    taskName: e.target.value
                                }
                            }
                        })}/>
                        所属项目:
                        <Select style={{marginTop: "5px", marginBottom: "10px", width: "100%"}}
                                defaultValue={this.props.service.Projects.Data === undefined ? "" : this.props.service.Projects.Data.length > 0 ? this.props.service.Projects.Data[0].ProjectName : ""}
                                onChange={(value) => this.setState({
                                    train: {
                                        ...this.state.train,
                                        doTrain: {
                                            ...this.state.train.doTrain,
                                            projectId: Number(value),
                                        }
                                    }
                                }, () => {

                                    const {dispatch} = this.props;
                                    dispatch({
                                        type: 'service/getLabels_v1',
                                        payload: {
                                            query: encodeURI(`ProjectId:${this.state.train.doTrain.projectId}`),
                                            limit: 1000,
                                        },
                                        callback: (bb) => {
                                            const singleTrain = [];
                                            if (bb.Data !== null) {
                                                bb.Data.map(v => {
                                                    singleTrain.push(v.LabelName);
                                                });
                                            }
                                            this.setState({
                                                ...this.state,
                                                train: {
                                                    ...this.state.train,
                                                    doTrain: {
                                                        ...this.state.train.doTrain,
                                                        singleTrain: singleTrain,
                                                    }
                                                }
                                            });
                                        }
                                    });
                                })}>
                            {
                                // 这里循环
                                this.props.service.Projects.Data.map(d => (
                                    <Option key={d.Id}>{d.ProjectName}</Option>
                                ))
                            }
                        </Select>
                        数据格式:
                        <Select style={{marginTop: "5px", marginBottom: "10px", width: "100%"}}
                                defaultValue="powerAi"
                                onChange={(value) => this.setState({
                                    train: {
                                        ...this.state.train,
                                        doTrain: {
                                            ...this.state.train.doTrain,
                                            assetsType: value
                                        }
                                    }
                                })}>
                            <Option value="powerAi">powerAi</Option>
                            <Option value="pascalVOC">pascalVOC</Option>
                            <Option value="coco">coco</Option>
                            <Option value="other">other</Option>
                        </Select>
                        是否开启多标签单模型训练:&nbsp;&nbsp;
                        <Switch checkedChildren="开启" unCheckedChildren="关闭"
                                checked={this.state.train.doTrain.mergeTrainSymbol === 1}
                                onChange={(c) => {
                                    this.setState({
                                        ...this.state,
                                        train: {
                                            ...this.state.train,
                                            doTrain: {
                                                ...this.state.train.doTrain,
                                                mergeTrainSymbol: c ? 1 : 0,
                                                singleTrain: c ? [] : this.state.train.doTrain.singleTrain,
                                                providerType: c ? "QTing-tiny-3l-multilabel" : "QTing-tiny-3l-single",
                                            }
                                        }
                                    })
                                }}/><br/>
                        训练的缺陷标签(留空也表示全部标签训练):
                        <Select  disabled={this.state.train.doTrain.mergeTrainSymbol === 1}
                                 style={{marginTop: "5px", marginBottom: "10px", width: "100%"}}
                                 allowClear
                                 value={this.state.train.doTrain.singleTrain}
                                 mode="multiple" placeholder="留空也表示全部标签训练" onChange={(value) => this.setState({
                            train: {
                                ...this.state.train,
                                doTrain: {
                                    ...this.state.train.doTrain,
                                    singleTrain: value,
                                    mergeTrainSymbol: 0,
                                    providerType: "QTing-tiny-3l-single",
                                }
                            }
                        })}>
                            {
                                this.props.service.Labels.Data !== null && this.props.service.Labels.Data.map(d => (
                                    <Option key={d.LabelName}>{d.LabelName}</Option>
                                ))
                            }
                        </Select>
                        使用的框架:
                        <Select style={{marginTop: "5px", marginBottom: "10px", width: "100%"}}
                                value={this.state.train.doTrain.providerType}
                                onChange={(value) => {
                                    this.setState({
                                        ...this.state,
                                        train: {
                                            ...this.state.train,
                                            doTrain: {
                                                ...this.state.train.doTrain,
                                                aiFrameworkId: Number(value.toString().split('|')[0]),
                                                providerType: value.toString().split('|')[1],
                                            }
                                        }
                                    });
                                }}>
                            {
                                this.props.service.AiFrameworks.Data.map(d => (
                                    <Option key={`${d.Id}|${d.FrameworkName}`}>{d.FrameworkName}</Option>
                                ))
                            }
                            {/*<Option value="QTing-tiny-3l-single">QTing-tiny-3l-single</Option>*/}
                            {/*<Option value="QTing-tiny-3l-multilabel">QTing-tiny-3l-multilabel</Option>*/}
                        </Select>
                        {// region 专业人员操作
                        }
                        AI参数(选填)-专业人员操作:&nbsp;&nbsp;
                        <Switch checkedChildren="已打开调参" unCheckedChildren="已关闭调参"
                                onChange={(c) => {
                                    // const {dispatch} = this.props;
                                    // dispatch({
                                    //     type: 'service/getModelListV2',
                                    //     payload: {
                                    //         framework_type: this.state.train.doTrain.providerType,
                                    //         project_name: encodeURI(this.state.train.doTrain.projectName)
                                    //     },
                                    // });
                                    this.setState({
                                        ...this.state,
                                        train: {
                                            ...this.state.train,
                                            showAiPar: c,
                                        }
                                    })
                                }}/>
                        {
                            this.state.train.showAiPar && <div>
                                学习率:
                                <InputNumber style={{width: '100%'}}
                                             placeholder={this.state.train.doTrain.learnrate}
                                             precision={6}
                                             step={0.00001}
                                             min={0}
                                             onChange={(value) => {
                                                 if (value === "" || value === null || value === undefined) value = 0.00261;
                                                 this.setState({
                                                     ...this.state,
                                                     train: {
                                                         ...this.state.train,
                                                         doTrain: {
                                                             ...this.state.train.doTrain,
                                                             learnrate: value,
                                                         },
                                                     }
                                                 })
                                             }}/>
                                检出率基数:
                                <InputNumber style={{width: '100%'}}
                                             placeholder={this.state.train.doTrain.recalldatum}
                                             precision={2}
                                             step={0.1}
                                             onChange={(value) => {
                                                 if (value === "" || value === null || value === undefined) value = 2;
                                                 this.setState({
                                                     ...this.state,
                                                     train: {
                                                         ...this.state.train,
                                                         doTrain: {
                                                             ...this.state.train.doTrain,
                                                             recalldatum: value,
                                                         },
                                                     }
                                                 }, () => {
                                                     console.log(`ducker do train: callback ${this.state.train.doTrain.recalldatum}`);
                                                 })
                                             }}/>
                                非当前标签图片训练方式:
                                <InputNumber style={{width: '100%'}}
                                             placeholder={this.state.train.doTrain.otherlabeltraintype}
                                             precision={0}
                                             step={1}
                                             min={0}
                                             onChange={(value) => {
                                                 if (value === "" || value === null || value === undefined) value = 1;
                                                 this.setState({
                                                     ...this.state,
                                                     train: {
                                                         ...this.state.train,
                                                         doTrain: {
                                                             ...this.state.train.doTrain,
                                                             otherlabeltraintype: value,
                                                         },
                                                     }
                                                 })
                                             }}/>
                                每次训练所选取的样本数:
                                <InputNumber style={{width: '100%'}}
                                             placeholder={this.state.train.doTrain.batchSize}
                                             min={0}
                                             onChange={(value) => {
                                                 if (value === "" || value === null || value === undefined) value = 64;
                                                 this.setState({
                                                     ...this.state,
                                                     train: {
                                                         ...this.state.train,
                                                         doTrain: {
                                                             ...this.state.train.doTrain,
                                                             batchSize: value,
                                                         },
                                                     }
                                                 }, () => {
                                                     console.log(`ducker do train: callback ${this.state.train.doTrain.batchSize}`);
                                                 })
                                             }}/>
                                GPU训练分批批次:
                                <InputNumber style={{width: '100%'}}
                                             placeholder={this.state.train.doTrain.subdivisionssize}
                                             min={0}
                                             onChange={(value) => {
                                                 if (value === "" || value === null || value === undefined) value = 16;
                                                 this.setState({
                                                     ...this.state,
                                                     train: {
                                                         ...this.state.train,
                                                         doTrain: {
                                                             ...this.state.train.doTrain,
                                                             subdivisionssize: value,
                                                         },
                                                     }
                                                 }, () => {
                                                     console.log(`ducker do train: callback ${this.state.train.doTrain.subdivisionssize}`);
                                                 })
                                             }}/>
                                图像宽高:
                                <InputNumber style={{width: '100%'}}
                                             placeholder={"默认自适应图像宽高"}
                                             precision={0}
                                             step={1}
                                             min={0}
                                             onChange={(value) => {
                                                 if (value === "" || value === null || value === undefined || value === -1)
                                                 {
                                                     this.setState({
                                                         ...this.state,
                                                         train: {
                                                             ...this.state.train,
                                                             doTrain: {
                                                                 ...this.state.train.doTrain,
                                                                 imageWidth: -1,
                                                                 imageHeight: -1,
                                                             },
                                                         }
                                                     });
                                                     return;
                                                 } else {
                                                     this.setState({
                                                         ...this.state,
                                                         train: {
                                                             ...this.state.train,
                                                             doTrain: {
                                                                 ...this.state.train.doTrain,
                                                                 imageWidth: value,
                                                                 imageHeight: value,
                                                             },
                                                         }
                                                     });
                                                 }
                                             }}/>
                                图像随机旋转角度范围:
                                <InputNumber style={{width: '100%'}}
                                             placeholder={this.state.train.doTrain.angle}
                                             precision={0}
                                             step={1}
                                             min={0}
                                             onChange={(value) => {
                                                 if (value === "" || value === null || value === undefined) value = 360;
                                                 this.setState({
                                                     ...this.state,
                                                     train: {
                                                         ...this.state.train,
                                                         doTrain: {
                                                             ...this.state.train.doTrain,
                                                             angle: value,
                                                         },
                                                     }
                                                 })
                                             }}/>
                                训练样本占比:
                                <InputNumber style={{width: '100%'}}
                                             placeholder={this.state.train.doTrain.split_ratio}
                                             precision={2}
                                             step={0.01}
                                             min={0}
                                             onChange={(value) => {
                                                 if (value === "" || value === null || value === undefined) value = 0.95;
                                                 this.setState({
                                                     ...this.state,
                                                     train: {
                                                         ...this.state.train,
                                                         doTrain: {
                                                             ...this.state.train.doTrain,
                                                             split_ratio: value,
                                                         },
                                                     }
                                                 })
                                             }}/>
                                训练最大轮数:
                                <InputNumber style={{width: '100%'}}
                                             placeholder={"默认自适应训练轮数"}
                                             precision={0}
                                             step={1}
                                             min={0}
                                             onChange={(value) => {
                                                 if (value === "" || value === null || value === undefined || value === -1) value = -1;
                                                 this.setState({
                                                     ...this.state,
                                                     train: {
                                                         ...this.state.train,
                                                         doTrain: {
                                                             ...this.state.train.doTrain,
                                                             maxIter: value,
                                                         },
                                                     }
                                                 })
                                             }}/>
                                最大负样本数:
                                <InputNumber style={{width: '100%'}}
                                             placeholder={this.state.train.doTrain.trainwithnolabelpic}
                                             precision={0}
                                             step={1}
                                             min={0}
                                             onChange={(value) => {
                                                 if (value === "" || value === null || value === undefined) value = 1000;
                                                 this.setState({
                                                     ...this.state,
                                                     train: {
                                                         ...this.state.train,
                                                         doTrain: {
                                                             ...this.state.train.doTrain,
                                                             trainwithnolabelpic: value,
                                                         },
                                                     }
                                                 })
                                             }}/>
                                正样本平移增强步长:
                                <InputNumber style={{width: '100%'}}
                                             placeholder={this.state.train.doTrain.cell_stride}
                                             precision={0}
                                             step={1}
                                             min={0}
                                             onChange={(value) => {
                                                 if (value === "" || value === null || value === undefined) value = 1;
                                                 this.setState({
                                                     ...this.state,
                                                     train: {
                                                         ...this.state.train,
                                                         doTrain: {
                                                             ...this.state.train.doTrain,
                                                             cell_stride: value,
                                                         },
                                                     }
                                                 })
                                             }}/>
                                负样本平移增强步长:
                                <InputNumber style={{width: '100%'}}
                                             placeholder={this.state.train.doTrain.otherlabelstride}
                                             precision={0}
                                             step={1}
                                             min={0}
                                             onChange={(value) => {
                                                 if (value === "" || value === null || value === undefined) value = 1;
                                                 this.setState({
                                                     ...this.state,
                                                     train: {
                                                         ...this.state.train,
                                                         doTrain: {
                                                             ...this.state.train.doTrain,
                                                             otherlabelstride: value,
                                                         },
                                                     }
                                                 })
                                             }}/>
                                平移框大小:
                                <InputNumber style={{width: '100%'}}
                                             placeholder={this.state.train.doTrain.cellsize}
                                             precision={0}
                                             step={1}
                                             min={0}
                                             onChange={(value) => {
                                                 if (value === "" || value === null || value === undefined) value = 16;
                                                 this.setState({
                                                 ...this.state,
                                                 train: {
                                                     ...this.state.train,
                                                     doTrain: {
                                                         ...this.state.train.doTrain,
                                                         cellsize: value,
                                                     },
                                                 }
                                             })}}/>


                                扩展尺寸:
                                <Input.Group compact>
                                    <InputNumber style={{width: '46%', textAlign: 'center'}}
                                                 placeholder={this.state.train.doTrain.expand_size[0]}
                                                 precision={0}
                                                 step={1}
                                                 min={0}
                                                 onChange={(value) => {
                                                     if (value === "" || value === null || value === undefined) value = 8;
                                                     this.setState({
                                                         ...this.state,
                                                         train: {
                                                             ...this.state.train,
                                                             doTrain: {
                                                                 ...this.state.train.doTrain,
                                                                 expand_size: [value, this.state.train.doTrain.expand_size[1]],
                                                             },
                                                         }
                                                     })
                                                 }}/>
                                    <Input
                                        className="site-input-split"
                                        style={{
                                            width: '8%',
                                            borderLeft: 0,
                                            borderRight: 0,
                                            pointerEvents: 'none',
                                            textAlign: 'center'
                                        }}
                                        placeholder="~"
                                        disabled
                                    />
                                    <InputNumber style={{width: '46%', textAlign: 'center'}}
                                                 placeholder={this.state.train.doTrain.expand_size[1]}
                                                 precision={0}
                                                 step={1}
                                                 min={0}
                                                 onChange={(value) => {
                                                     if (value === "" || value === null || value === undefined) value = 8;
                                                     this.setState({
                                                         ...this.state,
                                                         train: {
                                                             ...this.state.train,
                                                             doTrain: {
                                                                 ...this.state.train.doTrain,
                                                                 expand_size: [this.state.train.doTrain.expand_size[0], value],
                                                             },
                                                         }
                                                     })
                                                 }}/>
                                </Input.Group>


                                忽略尺寸:
                                <Input.Group compact>
                                    <InputNumber style={{width: '46%', textAlign: 'center'}}
                                                 placeholder={this.state.train.doTrain.ignore_size[0]}
                                                 precision={0}
                                                 step={1}
                                                 min={0}
                                                 onChange={(value) => {
                                                     if (value === "" || value === null || value === undefined) value = 6;
                                                     this.setState({
                                                         ...this.state,
                                                         train: {
                                                             ...this.state.train,
                                                             doTrain: {
                                                                 ...this.state.train.doTrain,
                                                                 ignore_size: [value, this.state.train.doTrain.ignore_size[1]],
                                                             },
                                                         }
                                                     })
                                                 }}/>
                                    <Input
                                        className="site-input-split"
                                        style={{
                                            width: '8%',
                                            borderLeft: 0,
                                            borderRight: 0,
                                            pointerEvents: 'none',
                                            textAlign: 'center'
                                        }}
                                        placeholder="~"
                                        disabled
                                    />
                                    <InputNumber style={{width: '46%', textAlign: 'center'}}
                                                 placeholder={this.state.train.doTrain.ignore_size[1]}
                                                 precision={0}
                                                 step={1}
                                                 min={0}
                                                 onChange={(value) => {
                                                     if (value === "" || value === null || value === undefined) value = 6;
                                                     this.setState({
                                                         ...this.state,
                                                         train: {
                                                             ...this.state.train,
                                                             doTrain: {
                                                                 ...this.state.train.doTrain,
                                                                 ignore_size: [this.state.train.doTrain.ignore_size[0], value],
                                                             },
                                                         }
                                                     })
                                                 }}/>
                                </Input.Group>

                                Anchor调整变化幅度:
                                <Input.Group compact>
                                    <InputNumber style={{width: '46%', textAlign: 'center'}}
                                                 placeholder={this.state.train.doTrain.resizearrange[0]}
                                                 precision={2}
                                                 step={0.01}
                                                 min={0}
                                                 onChange={(value) => {
                                                     if (value === "" || value === null || value === undefined) value = 0.3;
                                                     this.setState({
                                                         ...this.state,
                                                         train: {
                                                             ...this.state.train,
                                                             doTrain: {
                                                                 ...this.state.train.doTrain,
                                                                 resizearrange: [value, this.state.train.doTrain.resizearrange[1]],
                                                             },
                                                         }
                                                     })
                                                 }}/>
                                    <Input
                                        className="site-input-split"
                                        style={{
                                            width: '8%',
                                            borderLeft: 0,
                                            borderRight: 0,
                                            pointerEvents: 'none',
                                            textAlign: 'center'
                                        }}
                                        placeholder="~"
                                        disabled
                                    />
                                    <InputNumber style={{width: '46%', textAlign: 'center'}}
                                                 placeholder={this.state.train.doTrain.resizearrange[1]}
                                                 precision={2}
                                                 step={0.01}
                                                 min={0}
                                                 onChange={(value) => {
                                                     if (value === "" || value === null || value === undefined) value = 1.6;
                                                     this.setState({
                                                         ...this.state,
                                                         train: {
                                                             ...this.state.train,
                                                             doTrain: {
                                                                 ...this.state.train.doTrain,
                                                                 resizearrange: [this.state.train.doTrain.resizearrange[0], value],
                                                             },
                                                         }
                                                     })
                                                 }}/>
                                </Input.Group>


                                使用的GPU:
                                <Input style={{width: '100%'}}
                                       placeholder={this.state.train.doTrain.gpus}
                                       onChange={(e) => {
                                           let value = e.target.value;
                                           if (value === "" || value === null || value === undefined) value = "0,1";
                                           this.setState({
                                               ...this.state,
                                               train: {
                                                   ...this.state.train,
                                                   doTrain: {
                                                       ...this.state.train.doTrain,
                                                       gpus: value,
                                                   },
                                               }
                                           })
                                       }}/>
                                训练类型:&nbsp;&nbsp;
                                <Radio.Group defaultValue={this.state.train.doTrain.trianType}
                                             onChange={(e) => {
                                                 let value = e.target.value;
                                                 if (value === "" || value === null || value === undefined) value = 0;
                                                 this.setState({
                                                     ...this.state,
                                                     train: {
                                                         ...this.state.train,
                                                         doTrain: {
                                                             ...this.state.train.doTrain,
                                                             trianType: value,
                                                         },
                                                     }
                                                 })
                                             }}>
                                    <Radio value={0}>从头训练</Radio>
                                    <Radio value={1}>对应自训练</Radio>
                                    <Radio value={2}>漏检训练</Radio>
                                </Radio.Group>
                                <br/>
                                是否保留训练生成的临时数据:&nbsp;&nbsp;
                                <Switch checkedChildren="保留" unCheckedChildren="删除" defaultChecked={false} onChange={(v) => {
                                    this.setState({
                                        ...this.state,
                                        train: {
                                            ...this.state.train,
                                            doTrain: {
                                                ...this.state.train.doTrain,
                                                rmgeneratedata: v ? 1 : 0,
                                            },
                                        }
                                    });
                                }} />
                                <br/>
                                是否打乱数据:&nbsp;&nbsp;
                                <Switch checkedChildren="打乱" unCheckedChildren="不打乱" defaultChecked={false} onChange={(v) => {
                                    this.setState({
                                        ...this.state,
                                        train: {
                                            ...this.state.train,
                                            doTrain: {
                                                ...this.state.train.doTrain,
                                                isshuffle: v,
                                            },
                                        }
                                    });
                                }} />
                                <br/>
                                选择加载的预训练权重文件:
                                <Select
                                    style={{width: '100%'}}
                                    placeholder={"不选择的话默认使用初始的预训练文件"}
                                    onChange={(value) => {
                                        this.setState({
                                            ...this.state,
                                            train: {
                                                ...this.state.train,
                                                doTrain: {
                                                    ...this.state.train.doTrain,
                                                    pretrainweight: value,
                                                },
                                            }
                                        });
                                    }}>
                                    {modelListV2.model_list.map(d => (
                                        <Option key={d.filename}>{d.filename}</Option>
                                    ))}
                                </Select>
                                <br/>
                                <br/>
                                <br/>
                            </div>
                        }
                        {// endregion 专业人员操作
                        }
                        <div
                            style={{
                                position: 'absolute',
                                left: 0,
                                bottom: 0,
                                width: '100%',
                                borderTop: '1px solid #e9e9e9',
                                padding: '10px 16px',
                                background: '#fff',
                                textAlign: 'right',
                            }}
                        >
                            <Button onClick={() => {
                                this.setState({
                                    rightVisible: false,
                                    leftVisible: false,
                                });
                            }} style={{marginRight: 8}}>
                                关闭
                            </Button>

                            <Button type="primary"
                                    loading={this.state.train.loading}
                                    onClick={() => {
                                        if (!this.state.train.doTrain.projectName) {
                                            notification.error({
                                                message: "不存在项目",
                                                description: "当前未找到相关的项目，请返回首页新建项目再进行训练任务",
                                            });
                                            return;
                                        }
                                        if (this.state.train.doTrain.taskName) {
                                            this.setState({
                                                ...this.state,
                                                train: {
                                                    ...this.state.train,
                                                    loading: true,
                                                }
                                            }, () => {

                                                console.log(`ducker do train: ${JSON.stringify(this.state.train.doTrain)}`);
                                                const {dispatch} = this.props;
                                                dispatch({
                                                    type: 'service/postTrain_v1',
                                                    payload: this.state.train.doTrain,
                                                    callback: (v) => {
                                                        if (v["Code"] === 200) {
                                                            message.success("成功加入训练队列");
                                                            this.setState({
                                                                ...this.state,
                                                                rightVisible: false,
                                                                leftVisible: false,
                                                                train: {
                                                                    ...this.state.train,
                                                                    loading: false,
                                                                }
                                                            });
                                                        } else {
                                                            message.error("加入训练队列失败");
                                                            this.setState({
                                                                ...this.state,
                                                                train: {
                                                                    ...this.state.train,
                                                                    loading: false,
                                                                }
                                                            });
                                                        }
                                                    },
                                                });
                                            });

                                        } else {
                                            notification.error({
                                                message: "参数有误",
                                                description: "有部分参数未输入！请输完后重试",
                                            });
                                        }
                                    }}>
                                新增任务
                            </Button>
                        </div>
                    </Drawer>

                    <Drawer
                        destroyOnClose={true}
                        title="项目列表"
                        width="50%"
                        maskClosable={true}
                        onClose={() => {
                            this.setState({
                                ...this.state,
                                modelManagerSingle: {
                                    ...this.state.modelManagerSingle,
                                    firstVisible: false
                                },
                            });
                        }}
                        visible={this.state.modelManagerSingle.firstVisible}
                    >
                        <Table columns={[
                            {
                                title: '项目名',
                                key: "ProjectName",
                                dataIndex: 'ProjectName',
                                render: text => <Badge status="processing" text={text}/>,
                            },
                            {
                                title: '操作',
                                // dataIndex: 'ProjectName',
                                render: record => (
                                    <span>
                                    <a onClick={() => {
                                        this.setState({
                                            ...this.state,
                                            modelManagerSingle: {
                                                ...this.state.modelManagerSingle,
                                                secondVisible: true,
                                                nowEditProjectName: record.ProjectName,
                                                expandedRowKeys: [],
                                            },
                                        }, () => {
                                            this.setState({
                                               ...this.state,
                                               modelManagerSingle: {
                                                   ...this.state.modelManagerSingle,
                                                   loadingProjects: true,
                                               },
                                            }, () => {
                                                const {dispatch} = this.props;
                                                dispatch({
                                                    type: 'service/getLabels_v1',
                                                    payload: {
                                                        query: encodeURI(`ProjectId:${record.Id}`),
                                                        limit: 1000,
                                                    },
                                                    callback: (bb) => {
                                                        this.setState({
                                                            ...this.state,
                                                            modelManagerSingle: {
                                                                ...this.state.modelManagerSingle,
                                                                loadingProjects: false,
                                                            },
                                                        });
                                                    }
                                                });
                                            });
                                        });
                                    }}>单类别模型</a>
                                        <Divider type="vertical"/>
                                    <a  onClick={() => {
                                        this.setState({
                                            ...this.state,
                                            modelManagerMultilabel: {
                                                ...this.state.modelManagerMultilabel,
                                                secondVisible: true,
                                                nowEditProjectName: record.ProjectName,
                                                expandedRowKeys: [],
                                            },
                                        }, () => {
                                            this.setState({
                                                ...this.state,
                                                modelManagerMultilabel: {
                                                    ...this.state.modelManagerMultilabel,
                                                    loadingModels: true,
                                                },
                                            }, () => {
                                                const {dispatch} = this.props;
                                                dispatch({
                                                    type: 'service/getModelsByLabelsAndMulti_v1',
                                                    payload: {
                                                        projectId: record.Id,
                                                        isMultilabel: 1,
                                                    },
                                                    callback: (v) => {
                                                        this.setState({
                                                            ...this.state,
                                                            modelManagerMultilabel: {
                                                                ...this.state.modelManagerMultilabel,
                                                                loadingModels: false,
                                                            },
                                                        });
                                                    }
                                                });
                                            });
                                        });
                                    }}>多类别模型</a>
                                    </span>
                                ),
                            }]} dataSource={this.props.service.Projects.Data}/>
                        <Drawer
                            destroyOnClose={true}
                            title={`${this.state.modelManagerSingle.nowEditProjectName}-单类别模型列表`}
                            width="50%"
                            maskClosable={true}
                            onClose={() => {
                                this.setState({
                                    ...this.state,
                                    modelManagerSingle: {
                                        ...this.state.modelManagerSingle,
                                        secondVisible: false
                                    },
                                });
                            }}
                            visible={this.state.modelManagerSingle.secondVisible}
                        >
                            <Spin spinning={this.state.modelManagerSingle.loadingProjects}>
                            <Table
                                rowKey={"LabelName"}
                                columns={[{
                                    title: '标签名',
                                    key: "LabelName",
                                    dataIndex: 'LabelName',
                                }, {
                                    title: '备注',
                                    key: "Remarks",
                                    dataIndex: 'Remarks',
                                }]}
                            dataSource={this.props.service.Labels.Data}
                                pagination={false}
                                expandable={{
                                    expandedRowKeys: this.state.modelManagerSingle.expandedRowKeys,
                                    expandedRowRender: expandedModelsRowRender,
                                    expandRowByClick: true,
                                    onExpandedRowsChange: (expandedRows) => {
                                    },
                                    onExpand: (expanded, record) => {
                                        if (!expanded) {
                                            this.setState({
                                                ...this.state,
                                                modelManagerSingle: {
                                                    ...this.state.modelManagerSingle,
                                                    expandedRowKeys: [],
                                                },
                                            });
                                            return;
                                        }
                                        this.setState({
                                            ...this.state,
                                            modelManagerSingle: {
                                                ...this.state.modelManagerSingle,
                                                loadingModels: true,
                                                expandedRowKeys: [record.LabelName],
                                            },
                                        }, () => {
                                            const {dispatch} = this.props;
                                            dispatch({
                                                type: 'service/getModelsByLabelsAndMulti_v1',
                                                payload: {
                                                    projectId: record.ProjectId.Id,
                                                    label: record.LabelName,
                                                    isMultilabel: 0,
                                                },
                                                callback: (v) => {
                                                    this.setState({
                                                        ...this.state,
                                                        modelManagerSingle: {
                                                            ...this.state.modelManagerSingle,
                                                            loadingModels: false,
                                                        },
                                                    });
                                                }
                                            });
                                        });
                                    }
                                }}
                            />
                            </Spin>
                        </Drawer>
                        <Drawer
                            destroyOnClose={true}
                            title={`${this.state.modelManagerMultilabel.nowEditProjectName}-多类别模型列表`}
                            width="50%"
                            maskClosable={true}
                            onClose={() => {
                                this.setState({
                                    ...this.state,
                                    modelManagerMultilabel: {
                                        ...this.state.modelManagerMultilabel,
                                        secondVisible: false
                                    },
                                });
                            }}
                            visible={this.state.modelManagerMultilabel.secondVisible}
                        >
                            <Spin tip="正在加载..." spinning={this.state.modelManagerMultilabel.loadingModels}>
                                {/*<Badge status="processing" text="Running" />*/}
                                <Table
                                    rowKey={"Id"}
                                    columns={[
                                    {
                                        title: '模型名称',
                                        key: "ModelName",
                                        dataIndex: 'ModelName',
                                    }, {
                                        title: '发布状态',
                                        key: "Status",
                                        dataIndex: 'Status',
                                        render: v => {
                                            if (v === 1) return <Badge status="warning" text="已下线"/>;
                                            else if (v === 2) return <Badge status="processing" text="已发布"/>;
                                        },
                                    }, {
                                        title: '发布日期',
                                        render: (text, record) => {
                                            if (record.Status === 2) return <Tag icon={<CloudUploadOutlined/>} color="success">{record.PublishTime}</Tag>;
                                            else return record.PublishTime === "0001-01-01T00:00:00Z" ? "" : record.PublishTime;
                                        }
                                    }, {
                                        title: '操作',
                                        render: (text, record) => (
                                            <span>
                                    {
                                        record.Status !== 2 &&  <Popconfirm
                                            title="发布会把最新的模型替换为当前发布的模型，确定发布么？"
                                            onConfirm={() => {
                                                this.setState({
                                                    ...this.state,
                                                    publishModal: {
                                                        ...this.state.publishModal,
                                                        visible: true,
                                                        ModelWidth: record.ProjectId.ImageWidth,
                                                        ModelHeight: record.ProjectId.ImageHeight,
                                                        ModelId: record.Id,
                                                        refresh: "多标签",
                                                        label: "我是空的",
                                                        projectId: record.ProjectId.Id,
                                                    }
                                                });
                                            }}
                                            okText="确定并打开设置"
                                            cancelText="取消"
                                        >
                                            <a>{record.Status === 1 ? "上线" : "发布"}</a>
                                        </Popconfirm>
                                    }
                                                {
                                                    record.Status !== 2 && <Divider type="vertical"/>
                                                }
                                                {
                                                    record.Status !== 2 && <Popconfirm
                                                        title="确定要删除么？"
                                                        onConfirm={() => {
                                                            const {dispatch} = this.props;
                                                            dispatch({
                                                                type: 'service/delModel_v1',
                                                                payload: {
                                                                    Id: encodeURI(record.Id)
                                                                },
                                                                callback: (aa) => {
                                                                    const {dispatch} = this.props;
                                                                    dispatch({
                                                                        type: 'service/getModelsByLabelsAndMulti_v1',
                                                                        payload: {
                                                                            projectId: record.ProjectId.Id,
                                                                            isMultilabel: 1,
                                                                        },
                                                                    });
                                                                    notification.success({
                                                                        message: "恭喜",
                                                                        description: "删除成功",
                                                                    });
                                                                }
                                                            });
                                                        }}
                                                        okText="确定"
                                                        cancelText="取消"
                                                    >
                                                        <a>删除</a>
                                                    </Popconfirm>
                                                }


                                </span>),
                                    }]}
                                       dataSource={this.props.service.Models.Data}
                                       pagination={false}
                                       expandable={{
                                           // expandedRowKeys: this.state.modelManagerMultilabel.expandedRowKeys,
                                           expandedRowRender: expandedLabelsRowRender,
                                           expandRowByClick: true,
                                           onExpandedRowsChange: (expandedRows) => {
                                           },
                                           onExpand: (expanded, record) => {
                                               // if (!expanded) {
                                               //     this.setState({
                                               //         modelManagerMultilabel: {
                                               //             ...this.state.modelManagerMultilabel,
                                               //             expandedRowKeys: [],
                                               //         },
                                               //     });
                                               //     return;
                                               // }
                                               // this.setState({
                                               //     modelManagerMultilabel: {
                                               //         ...this.state.modelManagerMultilabel,
                                               //         expandedRowKeys: [record.Id],
                                               //     },
                                               // });
                                           }
                                       }}
                                />
                            </Spin>
                        </Drawer>
                    </Drawer>
                    {/*<div className="content padding">{content}</div>*/}
                    {/*<div className="content padding">{extraContent}</div>*/}
                </div>
            </PageHeader>

        );
    }
}
// 1. Initialize
const app = dva();
console.log(2);
// 2. Model
// app.model(require('./src/models/service').default);
app.model({
    namespace: 'service',
    state: {
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
        res: {
            code: undefined,
            status: undefined,
            msg: '',
            data: [],
        },
        modelList: {
            res: '',
            weights_list: [],
            width: undefined,
            height: undefined,
            max_batches: undefined,
        },
        modelListV2: {
            res: '',
            model_list: [],
        },
        valPathList: {
            res: '',
            val_path_list: [],
        },
        vocPathList: {
            res: '',
            voc_path_list: [],
        },
        modelByProject: {
            res: '',
            message: "",
            models: [],
        },
        labelsByProject: {
            res: '',
            message: "",
            labels: [],
        },
        labelsWithScoreByProject:  {
            res: '',
            message: "",
            labels: [],
        },
        dotrain:{},
        testRes: {
            res: '',
        },
        allres: {
            res: '',
        },
        get_release_models_history_res: {
            res: '',
            message: "",
            models: [],
        },
        del_model_res: {
            res: '',
            message: "",
        },
        online_model_res: {
            res: '',
            message: "",
        },
        offline_model_res: {
            res: '',
            message: "",
        },
        //region 多标签单模型
        models_multilabel: {
            res: '',
            message: "",
            models: [],
        },
        multilabel_by_model: {
            res: '',
            message: "",
            laebles: [],
        },
        //endregion
    },
    reducers: {
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
        res(state, action) {
            return {
                ...state,
                res: action.payload,
            };
        },
        modelList(state, action) {
            return {
                ...state,
                modelList: action.payload,
            };
        },
        valPathList(state, action) {
            return {
                ...state,
                valPathList: action.payload,
            };
        },
        vocPathList(state, action) {
            return {
                ...state,
                vocPathList: action.payload,
            };
        },
        dotrain(state, action) {
            return {
                ...state,
                dotrain: action.payload,
            };
        },
        testRes(state, action) {
            return {
                ...state,
                testRes: action.payload,
            };
        },
        allres(state, action) {
            return {
                ...state,
                allres: action.payload,
            };
        },
        // region 新增接口 自训练
        localPathList(state, action) {
            return {
                ...state,
                localPathList: action.payload,
            };
        },
        modelByProject(state, action) {
            return {
                ...state,
                modelByProject: action.payload,
            };
        },
        labelsByProject(state, action) {
            return {
                ...state,
                labelsByProject: action.payload,
            };
        },
        labelsWithScoreByProject(state, action) {
            return {
                ...state,
                labelsWithScoreByProject: action.payload,
            };
        },
        modelListV2(state, action) {
            return {
                ...state,
                modelListV2: action.payload,
            };
        },
        get_release_models_history_res(state, action) {
            return {
                ...state,
                get_release_models_history_res: action.payload,
            };
        },
        del_model_res(state, action) {
            return {
                ...state,
                del_model_res: action.payload,
            };
        },
        online_model_res(state, action) {
            return {
                ...state,
                online_model_res: action.payload,
            };
        },

        //region 多标签单模型
        models_multilabel(state, action) {
            return {
                ...state,
                models_multilabel: action.payload,
            };
        },
        multilabel_by_model(state, action) {
            return {
                ...state,
                multilabel_by_model: action.payload,
            };
        },
        //endregion
        // endregion
    },
    effects: {
        // region v1 新版本接口
        *getList_v1({ payload,callback}, { call, put }) {
            const response = yield call(getList_v1,payload);
            yield put({
                type: 'TrainRecords',
                payload: response,
            });
            if (callback)callback(response);
        },
        *getProjects_v1({ payload,callback}, { call, put }) {
            const response = yield call(getProjects_v1,payload);
            yield put({
                type: 'Projects',
                payload: response,
            });
            if (callback)callback(response);
        },
        *getLabels_v1({ payload,callback}, { call, put }) {
            const response = yield call(getLabels_v1,payload);
            yield put({
                type: 'Labels',
                payload: response,
            });
            if (callback)callback(response);
        },
        *getAiFramework_v1({ payload,callback}, { call, put }) {
            const response = yield call(getAiFramework_v1,payload);
            yield put({
                type: 'AiFrameworks',
                payload: response,
            });
            if (callback)callback(response);
        },
        *getModelsByLabelsAndMulti_v1({ payload,callback}, { call, put }) {
            const response = yield call(getModelsByLabelsAndMulti_v1,payload);
            yield put({
                type: 'Models',
                payload: response,
            });
            if (callback)callback(response);
        },
        *postTrain_v1({ payload,callback}, { call, put }) {
            const response = yield call(postTrain_v1,payload);
            yield put({
                type: 'Normal',
                payload: response,
            });
            if (callback)callback(response);
        },
        *delRecord_v1({ payload,callback}, { call, put }) {
            yield call(delRecord_v1,payload);
            if (callback)callback();
        },
        *delModel_v1({ payload,callback}, { call, put }) {
            yield call(delModel_v1,payload);
            if (callback)callback();
        },
        *onlineModel_v1({ payload,callback}, { call, put }) {
            const response = yield call(onlineModel_v1,payload);
            if (callback)callback(response);
        },
        // endregion

        *getModelList({ payload,callback}, { call, put }) {
            const response = yield call(getModelList,payload);
            yield put({
                type: 'modelList',
                payload: response,
            });
            if (callback)callback(response);
        },
        *getVocPathList({ payload,callback}, { call, put }) {
            const response = yield call(getVocPathList,payload);
            yield put({
                type: 'vocPathList',
                payload: response,
            });
            if (callback)callback(response);
        },
        *getValPathList({ payload,callback}, { call, put }) {
            const response = yield call(getValPathList,payload);
            yield put({
                type: 'valPathList',
                payload: response,
            });
            if (callback)callback(response);
        },
        *doTrain({ payload,callback}, { call, put }) {
            const response = yield call(doTrain,payload);
            yield put({
                type: 'dotrain',
                payload: response,
            });
            if (callback)callback(response);
        },
        *startTest({ payload,callback}, { call, put }) {
            const response = yield call(startTest,payload);
            yield put({
                type: 'testRes',
                payload: response,
            });
            if (callback)callback(response);
        },
        *stopTrain({ payload,callback}, { call, put }) {
            const response = yield call(stopTrain,payload);
            yield put({
                type: 'allres',
                payload: response,
            });
            if (callback)callback(response);
        },
        *continueTrainTrain({ payload,callback}, { call, put }) {
            const response = yield call(continueTrainTrain,payload);
            yield put({
                type: 'allres',
                payload: response,
            });
            if (callback)callback(response);
        },
        // region 新增接口 自训练
        *getLocalPathList({ payload,callback}, { call, put }) {
            const response = yield call(getLocalPathList,payload);
            yield put({
                type: 'localPathList',
                payload: response,
            });
            if (callback)callback(response);
        },
        *getLabelsWithScoreByProject({ payload,callback}, { call, put }) {
            const response = yield call(getLabelsWithScoreByProject,payload);
            yield put({
                type: 'labelsWithScoreByProject',
                payload: response,
            });
            if (callback)callback(response);
        },
        *suggest_score_get({ payload,callback}, { call, put }) {
            const response = yield call(suggest_score_get,payload);
            if (callback)callback(response);
        },
        *suggest_score_put({ payload,callback}, { call, put }) {
            const response = yield call(suggest_score_put,payload);
            if (callback)callback(response);
        },
        *getModelListV2({ payload,callback}, { call, put }) {
            const response = yield call(getModelListV2,payload);
            yield put({
                type: 'modelListV2',
                payload: response,
            });
            if (callback)callback(response);
        },
        *getReleaseModelsHistory({ payload,callback}, { call, put }) {
            const response = yield call(get_release_models_history,payload);
            yield put({
                type: 'get_release_models_history_res',
                payload: response,
            });
            if (callback)callback(response);
        },
        *delModel({ payload,callback}, { call, put }) {
            const response = yield call(del_model,payload);
            yield put({
                type: 'del_model_res',
                payload: response,
            });
            if (callback)callback(response);
        },
        *getModelSize({ payload,callback}, { call, put }) {
            const response = yield call(get_model_size,payload);
            if (callback)callback(response);
        },
        *offlineModel({ payload,callback}, { call, put }) {
            const response = yield call(offline_model,payload);
            yield put({
                type: 'offline_model_res',
                payload: response,
            });
            if (callback)callback(response);
        },

        //region 多标签单模型
        *getModelsMultilabel({ payload,callback}, { call, put }) {
            const response = yield call(get_models_multilabel,payload);
            yield put({
                type: 'models_multilabel',
                payload: response,
            });
            if (callback)callback(response);
        },
        *getMultilabelByModel({ payload,callback}, { call, put }) {
            const response = yield call(get_multilabel_by_model,payload);
            yield put({
                type: 'multilabel_by_model',
                payload: response,
            });
            if (callback)callback(response);
        },
        //endregion
        // endregion
    },
});
// 3. View

const App = connect(({ service }) => ({
    service
}))(function(props) {
    const { dispatch } = props;
    return (
        <div>
            <FreeFish {...props}/>
        </div>
    );
});

// 4. Router
app.router(() => <App />);

// 5. Start
app.start('#root');



// ReactDOM.render(<App />, document.getElementById('root'));