import React from 'react';
import { Form, Input, Button, Select, Modal, message,InputNumber,DatePicker, Radio,Switch, Spin } from 'antd';
import modalAiFramework from "../modals/modalAiFramework";
const { Option } = Select;

//region Form相关
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
                    <Form.Item initialValue={item.defaultValue[0]} rules={[{ required: item.required, message: `请输入${item.label}`}]} name={`Range#${item.name}#0`}
                               style={{display: 'inline-block', width: '45%', marginBottom: 0}}>
                        <InputNumber style={{width: '100%', textAlign: 'center'}}
                                     precision={item.precision?item.precision[0]:undefined} step={item.step?item.step[0]:undefined} min={item.min?item.min[0]:undefined}/>
                    </Form.Item>
                    <Input className="site-input-split" style={{width: '10%', borderLeft: 0, borderRight: 0, pointerEvents: 'none', textAlign: 'center'}} placeholder="~" disabled/>
                    <Form.Item initialValue={item.defaultValue[1]} rules={[{ required: item.required, message: `请输入${item.label}`}]} name={`Range#${item.name}#1`}
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
                style={{display: this.props.display}}
                ref={this.formRef}
                {...layout}
                layout='vertical'
                size={'small'}
                onFinish={(values) => {
                    // console.log(JSON.stringify(this.formRef.current.getFieldsValue(), null, 2));
                    // const aa = values.filter(v=>v.indexOf("Range#")>-1);
                    let finalValues = {};
                    for (const name in values) {
                        if (name.includes("Range#")) {
                            try{
                                const id = name.split('#')[2];
                                if (id === "0") {
                                    const nStr = name.split('#')[1];
                                    finalValues[nStr] = [values[name], values[name.replace("#0", "#1")]];
                                }
                            }
                            catch (e) {
                            }
                        }
                        else {
                            finalValues[name] =values[name];
                        }
                    }
                    // console.log(`Final: ${JSON.stringify(finalValues, null, 2)}`);
                    this.props.onFinish(finalValues);
                }}>
                {
                    this.state.aiPars.map((item, index) => {
                        // type 为 date 日期格式需要强制转化为 moment 格式
                        // item.value = item.type == 'date' ? moment(item.value).format('YYYY-MM-DD') : item.value;
                        return (
                            <Form.Item
                                name={item.type === "Range" ? undefined : item.name}
                                key={index}
                                {...formItemLayout}
                                style={{
                                    marginBottom: 1,
                                    padding: 0,
                                    display: item.visible ? "" : "none",
                                }}
                                label={item.label}
                                initialValue={item.type === "Range" ? undefined : item.defaultValue}
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
