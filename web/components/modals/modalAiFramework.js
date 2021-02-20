import React from 'react';
import { Form, Input, Button, Select, Modal, message } from 'antd';
const { Option } = Select;
const layout = {
    labelCol: {
        span: 6,
    },
    wrapperCol: {
        span: 16,
    },
};
const tailLayout = {
    wrapperCol: {
        offset: 6,
        span: 16,
    },
};
export default class ModalAiFramework extends React.Component {
    formRef = React.createRef();
    state = {
        visible: false,
        title: "",
    };
    componentWillMount() {
    }

    componentWillUpdate(newProps,newState,newContext) {
        // console.error("更新了");
        // this.onFill(newState.data);
    }

    componentDidMount() {

    }
    onVisible = (show: boolean, title: string, data) => {
        this.setState({
            ...this.state,
            visible: show,
            title,
            data,
        }, () => this.onFill());
    };
    onReset = () => {
        this.formRef.current.resetFields();
    };
    onFill = () => {
        if (this.state.data === null || this.state.data === undefined) {
            try{
                this.formRef.current.resetFields();
            }
            catch (e) {
                console.error(e);
            }
        } else {
            try{
                this.formRef.current.setFieldsValue({...this.state.data});
            }
            catch (e) {
                console.error(e);
            }
        }
    };
    render(){ // 将页面通过return返回给调用者  必须有一个根节点包裹 同vue的使用
        return(
            <Modal
                maskClosable={false}
                // destroyOnClose
                title={this.state.title}
                visible={this.state.visible}
                footer={null}
                onCancel={() => {
                    this.setState({
                        ...this.state,
                        visible: false,
                    });
                }}
            >
                <Form {...layout} ref={this.formRef} name="control-hooks" onFinish={(values) => {
                    // console.log(values);
                    const {dispatch} = this.props;
                    dispatch({
                        type: 'service/AiFramework_v1',
                        payload: {
                            method: this.state.title.includes("编辑") ? "PUT" : "POST",
                            data: values,
                        },
                        callback: (bb) => {
                            if (bb["Code"] === 200) {
                                message.success("操作成功");
                                this.setState({
                                    visible: false,
                                });
                                dispatch({
                                    type: 'service/getAiFramework_v1',
                                    payload: {
                                        sortby: "CreateTime",
                                        order: "desc",
                                        limit: 200,
                                    }
                                });
                            } else {
                                message.error(bb["Msg"]);
                            }
                        },
                    });
                }}>
                    <Form.Item name="Id" noStyle>
                        <Input type="hidden"/>
                    </Form.Item>
                    <Form.Item name="FrameworkName" label="框架名称" rules={[{ required: true }]}>
                        <Input />
                    </Form.Item>
                    <Form.Item name="Cfg" label="配置文件" rules={[{ required: true }]}>
                        <Input />
                    </Form.Item>
                    <Form.Item name="Volume" label="映射目录" rules={[{ required: true }]}>
                        <Input />
                    </Form.Item>
                    <Form.Item label="镜像地址" name="BaseImageUrl" rules={[{ required: true }]}>
                        <Input />
                    </Form.Item>
                    <Form.Item label="镜像版本" name="ImageVersion" rules={[{ required: true }]}>
                        <Input />
                    </Form.Item>
                    <Form.Item name="Remarks" label="备注">
                        <Input.TextArea />
                    </Form.Item>
                    <Form.Item {...tailLayout}>
                        <Button type="primary" htmlType="submit">
                            提交
                        </Button>
                        <Button style={{marginLeft: 5}} htmlType="button" onClick={this.onReset}>
                            清空
                        </Button>
                        {
                            this.state.title.includes("编辑") ? <Button type="link" htmlType="button" onClick={this.onFill}>
                                恢复默认
                            </Button> : undefined
                        }
                    </Form.Item>
                </Form>
            </Modal>
        )
    }
}