import React from 'react';
import { Form, Input, Button, Select, Modal, message,InputNumber,DatePicker, Radio,Switch, Spin } from 'antd';
import modalAiFramework from "../modals/modalAiFramework";
const { Option } = Select;

//region Form相关
const data=[];
const data2 = [
    {
        "name": "learnrate",//formitem的名称
        "label": "学习率",//标题信息
        "defaultValue": 0.00261,//默认值
        "precision": 6, //小数点后保留几位
        "step": 0.00001, //步长
        "min": 0, //最小值
        "required": true, //是否必填
        "type": "InputNumber",
    }, {
        "name": "recalldatum",
        "label": "检出率基数",
        "defaultValue": 2,
        "precision": 2,
        "step": 0.1,
        "min": 0,
        "required": true,
        "type": "InputNumber",
    },{
        "name": "otherlabeltraintype",
        "label": "非当前标签图片训练方式",
        "defaultValue": 1,
        "step": 1,
        "min": 0,
        "required": true,
        "type": "InputNumber",
    },{
        "name": "batchSize",
        "label": "每次训练所选取的样本数",
        "defaultValue": 64,
        "step": 16,
        "min": 0,
        "required": true,
        "type": "InputNumber",
    },{
        "name": "subdivisionssize",
        "label": "GPU训练分批批次",
        "defaultValue": 2,
        "step": 2,
        "min": 0,
        "required": true,
        "type": "InputNumber",
    },{
        "name": "imageWidth",
        "label": "图像宽高",
        "defaultValue": -1,
        "step": 1,
        "min": -1,
        "required": true,
        "type": "InputNumber",
    }, {
        "name": "angle",
        "label": "图像随机旋转角度范围",
        "defaultValue": 360,
        "step": 1,
        "min": 0,
        "required": true,
        "type": "InputNumber",
    }, {
        "name": "split_ratio",
        "label": "训练样本占比",
        "defaultValue": 0.95,
        "precision": 2,
        "step": 0.01,
        "min": 0,
        "required": true,
        "type": "InputNumber",
    }, {
        "name": "maxIter",
        "label": "训练最大轮数",
        "defaultValue": -1,
        "step": 1,
        "min": -1,
        "required": true,
        "type": "InputNumber",
    }, {
        "name": "trainwithnolabelpic",
        "label": "最大负样本数",
        "defaultValue": 50000,
        "step": 1,
        "min": 0,
        "required": true,
        "type": "InputNumber",
    }, {
        "name": "cell_stride",
        "label": "正样本平移增强步长",
        "defaultValue": 1,
        "step": 1,
        "min": 0,
        "required": true,
        "type": "InputNumber",
    }, {
        "name": "otherlabelstride",
        "label": "负样本平移增强步长",
        "defaultValue": 1,
        "step": 1,
        "min": 0,
        "required": true,
        "type": "InputNumber",
    },  {
        "name": "cellsize",
        "label": "平移框大小",
        "defaultValue": 16,
        "step": 1,
        "min": 0,
        "required": true,
        "type": "InputNumber",
    },  {
        "names": ["learnratestepsratio0", "learnratestepsratio1"],
        "label": "学习率调整步数比",
        "defaultValue": [0.9,0.95],
        "precision": [2,2],
        "step": [0.01,0.01],
        "min": [0,0],
        "required": true,
        "type": "Range",
    },  {
        "names": ["expand_size0", "expand_size1"],
        "label": "扩展尺寸",
        "defaultValue": [8,8],
        "step": [1,1],
        "min": [0,0],
        "required": true,
        "type": "Range",
    }, {
        "names": ["ignore_size0", "ignore_size1"],
        "label": "忽略尺寸",
        "defaultValue": [6,6],
        "step": [1,1],
        "min": [0,0],
        "required": true,
        "type": "Range",
    }, {
        "names": ["resizearrange0", "resizearrange1"],
        "label": "Anchor调整变化幅度",
        "defaultValue": [0.3,1.6],
        "step": [0.01,0.01],
        "min": [0,0],
        "required": true,
        "type": "Range",
    }, {
        "name": "gpus",
        "label": "使用的GPU",
        "required": true,
        "type": "String",
        "defaultValue": "0,1"
    }, {
        "name": "trianType",
        "label": "训练类型",
        "required": true,
        "type": "RadioGroup",
        "defaultValue": 0,
        "radioValues": [0,1,2],
        "radioLabels": ["从头训练", "对应自训练", "漏检训练"]
    }, {
        "name": "rmgeneratedata",
        "label": "是否保留训练生成的临时数据",
        "checkedChildren": "保留",
        "unCheckedChildren": "删除",
        "required": true,
        "type": "Switch",
        "defaultChecked":  false,
    }, {
        "name": "isshuffle",
        "label": "是否打乱数据",
        "checkedChildren": "打乱",
        "unCheckedChildren": "不打乱",
        "required": true,
        "type": "Switch",
        "defaultChecked":  true,
    }, {
        "name": "pretrainweight",
        "label": "选择加载的预训练权重文件",
        "required": true,
        "type": "Select",
        "optionValues": [],
        "optionLabels":[],
    }
];
const layout = {
    labelCol: {
        xs: { span: 24 },
        sm: { span: 6 },
    },
    wrapperCol: {
        xs: { span: 24 },
        sm: { span: 14 },
    }
};
// formItem css 样式
const formItemLayout = {
    labelCol: {
        span: 24
    },
    wrapperCol: {
        span: 24
    }
};

// 保存按钮 css 样式
const tailFormItemLayout = {
    wrapperCol: {
        xs: {
            span: 24,
            offset: 0,
        },
        sm: {
            span: 14,
            offset: 6,
        },
    }
};

// form css 样式
const formLayout = {
    width: 400,
    marginTop: 100,
    marginLeft: 'auto',
    marginRight: 'auto'
}
//endregion


export default class AIForm extends React.Component {
    formRef = React.createRef();
    state = {
        visible: false,
        title: "",
        aiPars: [],
        loadingAiPars: false,
    };
    componentWillMount() {
    }

    componentWillUpdate(newProps,newState,newContext) {
        // console.error("更新了");
        // this.onFill(newState.data);
    }

    componentDidMount() {
    }
    loadAiPars = (data) => {
        this.setState({
            ...this.state,
            loadingAiPars: false,
            aiPars: data,
        },()=>{
        });
    };
    loadAiParsOld = (aiFrameworkId) => {
        this.setState({
            ...this.state,
            loadingAiPars: true,
        },()=>{
            const {dispatch} = this.props;
            dispatch({
                type: 'service/getAiFramework_v1',
                payload: {
                    query: encodeURI(`Id:${aiFrameworkId}`),
                },
                callback: (bb) => {
                    console.log(bb.Data[0].ParsJson);
                    this.setState({
                        ...this.state,
                        loadingAiPars: false,
                        aiPars: JSON.parse(bb.Data[0].ParsJson),
                    }, ()=>{
                        dispatch({
                            type: 'service/getAiFramework_v1',
                            payload: {
                                sortby: "CreateTime",
                                order: "desc",
                                limit: 200,
                            },
                        });
                    });
                },
            });
        });
    };
    closeAiPars = () => {
        this.setState({
            ...this.state,
            loadingAiPars: false,
            aiPars: [],
        });
    };

    doSubmit = () => {
        this.formRef.current.submit();
    };
    /**
     * 根据后台返回的 data 中 type 类型生成不同的组件
     * @param item  json
     * @param Component
     */
    switchItem(item) {
        const type = item.type;
        switch (type) {
            case 'Input':
                return <InputNumber style={{width: '100%'}}/>;
                break;
            case 'InputNumber':
                return <InputNumber style={{width: '100%'}}
                                    placeholder={item.placeholder} precision={item.precision} step={item.step}
                                    min={item.min}/>;
                break;
            case 'Range':
                return <Form.Item style={{marginBottom: 0}}>
                    <Form.Item initialValue={item.defaultValue} rules={[{ required: item.required, message: `请输入${item.label}`}]} name={item.names?item.names[0]:undefined}
                               style={{display: 'inline-block', width: '45%', marginBottom: 0}}>
                        <InputNumber style={{width: '100%', textAlign: 'center'}}
                                     precision={item.precision?item.precision[0]:undefined} step={item.step?item.step[0]:undefined} min={item.min?item.min[0]:undefined}/>
                    </Form.Item>
                    <Input className="site-input-split" style={{width: '10%', borderLeft: 0, borderRight: 0, pointerEvents: 'none', textAlign: 'center'}} placeholder="~" disabled/>
                    <Form.Item initialValue={item.defaultValue} rules={[{ required: item.required, message: `请输入${item.label}`}]} name={item.names?item.names[1]:undefined}
                               style={{display: 'inline-block', width: '45%', margin: '0 -1px', marginBottom: 0}}>
                        <InputNumber style={{width: '100%', textAlign: 'center'}}
                                     precision={item.precision?item.precision[1]:undefined} step={item.step?item.step[1]:undefined} min={item.min?item.min[1]:undefined}/>
                    </Form.Item>
                </Form.Item>;
                break;
            case 'String':
                return <Input style={{width: '100%'}}/>;
                break;
            case 'RadioGroup':
                return <Radio.Group>
                    {
                        item.radioLabels && item.radioLabels.map((lab, index) => {
                            return (<Radio key={index} value={item.radioValues[index]}>{lab}</Radio>)
                        })
                    }
                </Radio.Group>;
                break;
            case 'Switch':
                return <Switch defaultChecked={item.defaultChecked} checkedChildren={item.checkedChildren} unCheckedChildren={item.unCheckedChildren}/>;
                break;
            case 'Select':
                return (
                    <Select>
                        {
                            item.optionLabels &&  item.optionLabels.map((lab, index) => {
                                return (<Option key={index} value={item.optionValues[index]}>{lab}</Option>)
                            })
                        }
                    </Select>
                );
                break;
            default:
                return <Input/>;
                break;
        }
    }

    render(){ // 将页面通过return返回给调用者  必须有一个根节点包裹 同vue的使用
        return (
            <Spin spinning={this.state.loadingAiPars} tip="正在加载框架参数...">
            <Form
                ref={this.formRef}
                {...layout}
                layout='vertical'
                size={'small'}
                onFinish={(values) => {
                    this.props.onFinish(values);
                }}>
                {
                    this.state.aiPars.map((item, index) => {
                        // type 为 date 日期格式需要强制转化为 moment 格式
                        // item.value = item.type == 'date' ? moment(item.value).format('YYYY-MM-DD') : item.value;
                        return (
                            <Form.Item
                                name={item.name}
                                key={index}
                                {...formItemLayout}
                                style={{
                                    marginBottom: 1,
                                    padding: 0,
                                }}
                                label={item.label}
                                initialValue={item.defaultValue}
                                rules={[{ required: item.required, message: `请输入${item.label}` }]}
                            >
                                {this.switchItem(item)}
                            </Form.Item>
                        )
                    })
                }
                <Form.Item />
            </Form>
            </Spin>
        )
    }
}
