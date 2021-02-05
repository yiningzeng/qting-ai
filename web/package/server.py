# -*- coding=utf-8 -*-
# !/usr/bin/env python
import os
import glob
import pika
import json
import time
import socket
import psycopg2
import numpy as np
from wxpy import *
import logging
from logging import handlers
from retry import retry
from visdom import Visdom
from datetime import datetime
from flask_cors import CORS
from flask import Flask, request, Response
from urllib.parse import quote, unquote, urlencode
from apscheduler.schedulers.background import BackgroundScheduler


# 网络框架枚举
class net_framework:
    yolov4Tiny3l = {"name": "yolov4-tiny-3l", "modeldcfgname": "yolov4-tiny-3l.cfg", "train_docker_volume": "/Afinaltrain/SourDatas", "modelSavePath": "backup", "modelSuffix": ".weights", "configSuffix": ".cfg"}


def get_net_framework(net_name):
    if net_name == "yolov4-tiny-3l":
        return net_framework.yolov4Tiny3l
    else:
        return None


app = Flask(__name__)
CORS(app, resources=r'/*')
# rabbitmq 文档 https://pika.readthedocs.io/en/stable/modules/channel.html
# retry https://github.com/invl/retry
# pika https://pypi.org/project/pika/

#


'''
usage:  dockertrain  -p  映射到本地的端口 默认8097 如果被占用会自动分配，只检测端口占用情况，可能存在多个未开启的容器相同端口的情况
                      -n  项目名 默认 ""
                      -v  需要映射的素材目录(必填)
                      -r  docker镜像的地址 默认registry.cn-hangzhou.aliyuncs.com/baymin/ai-power:ai-power-wo-v3.6
                      -w  root密码 默认icubic-123
                      -g  复制脚本到/usr/local/bin/，后面执行可以全局dockertrain
                      -o  日志的输出目录默认/var/log/train
                      -t  docker的gup版本，默认是最新版本2，设置1：nvidia-docker，2：docker run --gpus all
                      -h  帮助
'''

'''
第一次运行一定要保证queue要存在，就是直接运行两次
'''
class Logger(object):
    level_relations = {
        'debug':logging.DEBUG,
        'info':logging.INFO,
        'warning':logging.WARNING,
        'error':logging.ERROR,
        'crit':logging.CRITICAL
    }  # 日志级别关系映射

    def __init__(self, filename, level='info', when='D', backCount=30, fmt='%(asctime)s - %(pathname)s[line:%(lineno)d] - %(levelname)s: %(message)s'):
        self.logger = logging.getLogger(filename)
        format_str = logging.Formatter(fmt)  # 设置日志格式
        self.logger.setLevel(self.level_relations.get(level))  # 设置日志级别
        sh = logging.StreamHandler()  # 往屏幕上输出
        sh.setFormatter(format_str)  # 设置屏幕上显示的格式
        th = handlers.TimedRotatingFileHandler(filename=filename, when=when, backupCount=backCount, encoding='utf-8')  # 往文件里写入#指定间隔时间自动生成文件的处理器
        # 实例化TimedRotatingFileHandler
        # interval是时间间隔，backupCount是备份文件的个数，如果超过这个个数，就会自动删除，when是间隔的时间单位，单位有以下几种：
        # S 秒
        # M 分
        # H 小时、
        # D 天、
        # W 每星期（interval==0时代表星期一）
        # midnight 每天凌晨
        th.setFormatter(format_str)  # 设置文件里写入的格式
        self.logger.addHandler(sh)  # 把对象加到logger里
        self.logger.addHandler(th)


class pikaqiu(object):

    def __init__(self, root_password='icubic-123', host='localhost', port=5672,
                 username='guest', password='guest', assets_base_path='/assets/Projects',
                 train_exchange='ai.train.topic', train_queue='ai.train.topic-queue', train_routing_key='train.start.#',
                 test_exchange='ai.test.topic', test_queue='ai.test.topic-queue', test_routing_key='test.start.#',
                 package_exchange='ai.package.topic', package_queue='ai.package.topic-queue',
                 package_routing_key='package.upload-done.#'
                 ):
        self.assets_base_path = assets_base_path
        self.root_password = root_password
        self.host = host
        self.port = port
        self.username = username
        self.password = password
        self.sql_host = 'localhost'
        # region 画图参数
        self.draw = True
        self.draw_windows = None
        self.draw_host = 'localhost'
        self.draw_port = 8097
        # endregion
        # region 训练队列参数
        self.train_exchange = train_exchange
        self.train_queue = train_queue
        self.train_routing_key = train_routing_key
        # endregion
        # region 测试队列参数
        self.test_exchange = test_exchange
        self.test_queue = test_queue
        self.test_routing_key = test_routing_key
        # endregion
        # region 训练素材包队列参数
        self.package_exchange = package_exchange
        self.package_queue = package_queue
        self.package_routing_key = package_routing_key
        # endregion
        self.parameters = pika.URLParameters("amqp://%s:%s@%s:%d" % (username, password, host, port))
        # region postgres
        self.postgres_conn = None
        # endregion
        # self.consume()

    def draw_chat(self, data=[{"x": 13.00, "y": 13.33, "win_id": "窗体名称->就是当前容器的id+该图标的含义", "title": "窗体显示的名称"}],
                  debug=False, err=False):
        # record 定义的格式{"x": 13.00, "y": 13.33, "win_id": "窗体名称->就是当前容器的id+该图标的含义", "title": "窗体显示的名称"}
        if debug:
            self.draw_windows = Visdom(env="test")
        if err:
            a, rows = ff.postgres_execute(
                "SELECT project_id, assets_directory_base, assets_directory_name, project_name"
                " FROM train_record WHERE status=2", True)
            if rows is not None and len(rows) > 0:
                assets_directory_name = rows[0][2]
                # os.system("echo 训练失败-梯度爆炸了 > '%s/%s/train_status.log'" % (self.package_base_path, assets_directory_name))  # 会自动退出，所以这里不需要了
                # os.system("echo '%s' | sudo -S docker stop `cat '%s/%s/train.dname'`" % (self.root_password, self.package_base_path, assets_directory_name))  # 会自动退出，所以这里不需要了
                self.postgres_execute("UPDATE train_record SET "
                                      "status=%d, project_name='%s'"
                                      " WHERE project_id='%s'" %
                                      (-1, str(rows[0][3]) + "-梯度爆炸了", rows[0][0]))
        else:
            for record in data:
                if self.draw_windows is None:
                    temp = record["win_id"]
                    pos = temp.rfind("-")
                    project_id = temp[:pos]
                    print(project_id)  # "C:/Python27/1"
                    a, rows = ff.postgres_execute(
                        "SELECT project_id, assets_directory_base, assets_directory_name, project_name"
                        " FROM train_record WHERE project_id='%s'" % project_id, True)
                    if rows is not None and len(rows) > 0:
                        assets_directory_name = rows[0][2]
                        self.postgres_execute("UPDATE train_record SET "
                                              "status=%d, project_name='%s'"
                                              " WHERE project_id='%s'" %
                                              (2, str(rows[0][3]).replace("-梯度爆炸了", ""), rows[0][0]))
                        if debug:
                            self.draw_windows = Visdom(env=project_id)
                        else:
                            draw_log = self.assets_base_path + "/" + assets_directory_name + "/draw.log"
                            self.draw_windows = Visdom(env=project_id, log_to_filename=draw_log)
                if self.draw_windows.win_exists(record["win_id"]):
                    self.draw_windows.line(
                        X=np.array([record["x"]]),
                        Y=np.array([record["y"]]),
                        win=record["win_id"],
                        opts=dict(title=record["title"], width=600, height=380),
                        update='append')
                else:
                    self.draw_windows.line(
                        win=record["win_id"],
                        X=np.array([0]),
                        Y=np.array([0]),
                        opts=dict(title=record["title"], width=600, height=380))

    '''
        - *dbname*: the database name
        - *database*: the database name (only as keyword argument)
        - *user*: user name used to authenticate
        - *password*: password used to authenticate
        - *host*: database host address (defaults to UNIX socket if not provided)
        - *port*: connection port number (defaults to 5432 if not provided)
    :return: 
    '''

    def postgres_connect(self, host='localhost', port=5432, user='postgres', password='baymin1024', dbname='power_ai'):
        self.postgres_conn = psycopg2.connect("host=%s port=%d user=%s password=%s dbname=%s" %
                                              (host, port, user, password, dbname))
        return True

    def postgres_execute(self, sql=None, select=False):
        if self.postgres_conn is None:
            self.postgres_connect(host=self.sql_host)
        result = None
        if sql is None:
            log.logger.info("sql null")
            return False, result
        log.logger.info(sql)
        try:
            cur = self.postgres_conn.cursor()
            cur.execute(sql)
            if select:
                result = cur.fetchall()
            self.postgres_conn.commit()
            cur.close()
        except Exception as e:
            self.postgres_conn.commit()  # 修复有一次插入或者查询出错就直接报错了
            log.logger.info(str(e))
            return False, result
        else:
            log.logger.info("sql执行成功")
            return True, result

    def postgres_disconnect(self):
        self.postgres_conn.close()
        return True

    @retry(pika.exceptions.AMQPConnectionError, delay=5, jitter=(1, 3))
    def init(self, sql=True, sql_host='localhost', draw=True, draw_host='localhost', draw_port=1121):
        self.draw = draw
        self.draw_host = draw_host
        self.draw_port = draw_port
        self.sql_host = sql_host
        if draw:
            os.system(
                "echo %s | sudo -S docker stop service-web-loss && sudo docker rm service-web-loss" % self.root_password)
            os.system("echo %s | sudo -S docker run \
            --name service-web-loss \
            -p %d:80 \
            -v /assets/Projects/:/usr/local/apache2/htdocs/ \
            --net ai --ip 10.10.0.99 \
            --restart=always \
            -d registry.cn-hangzhou.aliyuncs.com/baymin/remote-train:web-v3.5" % (self.root_password, draw_port))
        if sql:
            self.postgres_connect(host=sql_host)
        connection = pika.BlockingConnection(self.parameters)
        channel = connection.channel()
        try:
            # region创建训练队列
            channel.exchange_declare(self.train_exchange, "topic", passive=True, durable=True)
            channel.queue_declare(self.train_queue, passive=True, durable=True)
            channel.queue_bind(self.train_queue, self.train_exchange, self.train_routing_key)
            # endregion
            # region创建测试队列
            channel.exchange_declare(self.test_exchange, "topic", passive=True, durable=True)
            channel.queue_declare(self.test_queue, passive=True, durable=True)
            channel.queue_bind(self.test_queue, self.test_exchange, self.test_routing_key)
            # endregion
            # region创建训练素材解包队列
            channel.exchange_declare(self.package_exchange, "topic", passive=True, durable=True)
            channel.queue_declare(self.package_queue, passive=True, durable=True)
            channel.queue_bind(self.package_queue, self.package_exchange, self.package_routing_key)
            # endregion
        except Exception as e:
            # region创建训练队列
            channel = connection.channel()
            channel.exchange_declare(self.train_exchange, "topic", durable=True)
            channel.queue_declare(self.train_queue)
            channel.queue_bind(self.train_queue, self.train_exchange, self.train_routing_key)
            # endregion
            # region创建训练队列
            channel = connection.channel()
            channel.exchange_declare(self.test_exchange, "topic", durable=True)
            channel.queue_declare(self.test_queue)
            channel.queue_bind(self.test_queue, self.test_exchange, self.test_routing_key)
            # endregion
            # region创建训练素材解包队列
            channel = connection.channel()
            channel.exchange_declare(self.package_exchange, "topic", durable=True)
            channel.queue_declare(self.package_queue)
            channel.queue_bind(self.package_queue, self.package_exchange, self.package_routing_key)
            # endregion
        connection.close()
        # self.get_one(channel)


# 队列新增，每次重新连接，防止出现断开的情况，并且会有重试机制
@retry(pika.exceptions.ChannelWrongStateError, tries=3, delay=3, jitter=(1, 3))
def do_basic_publish(exchange, routing_key, body):
    connection = pika.BlockingConnection(ff.parameters)
    channel = connection.channel()
    channel.basic_publish(exchange, routing_key, body)
    connection.close()


# 队列单个获取，每次重新连接，防止出现断开的情况，并且会有重试机制
@retry(pika.exceptions.ChannelWrongStateError, tries=3, delay=3, jitter=(1, 3))
def do_basic_get(queue, auto_ack=False):
    connection = pika.BlockingConnection(ff.parameters)
    channel = connection.channel()
    method_frame, header_frame, body = channel.basic_get(queue=queue, auto_ack=auto_ack)
    return connection, channel, method_frame, header_frame, body


def net_is_used(port, ip='0.0.0.0'):
    s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    try:
        s.connect((ip, port))
        s.shutdown(2)
        print('%s:%d is used' % (ip, port))
        return True
    except:
        print('%s:%d is unused' % (ip, port))
        return False


# 获取单个训练队列数据
def get_train_one():
    log.logger.info('get_train_one:%s' % (datetime.now().strftime('%Y-%m-%d %H:%M:%S')))
    os.system("notify-send '%s' '%s' -t %d" % ('ceshi', '测试', 10000))
    connection, channel, method_frame, header_frame, body = do_basic_get(queue=ff.train_queue)
    # chan.basic_ack(msg.delivery_tag)
    # It can be empty if the queue is empty so don't do anything

    if method_frame is None:
        log.logger.info("训练：Empty Basic.Get Response (Basic.GetEmpty)")
        return None, None
        # We have data
    else:
        # 这里需要检查训练素材包是否已经解包，如果未解包，这里需要拒绝，让它重新排队ff.channel.basic_nack
        train_info = json.loads(body.decode('utf-8'))

        # 判断训练状态文件是否存在
        if not os.path.exists("%s/%s/train_status.log" % (ff.assets_base_path, train_info["projectName"])):
            log.logger.info('%s 等待训练' % (datetime.now().strftime('%Y-%m-%d %H:%M:%S')))
            os.system("echo '等待训练\c' > '%s/%s/train_status.log'" %
                      (ff.assets_base_path, train_info["projectName"]))
            channel.basic_nack(method_frame.delivery_tag)
            connection.close()
            get_train_one()
            log.logger.info("等待训练")
        else:
            status = os.popen("cat '%s/%s/train_status.log' | head -n 1" %
                              (ff.assets_base_path, train_info["projectName"])).read().replace('\n', '')
            if "等待训练" in status:
                channel.basic_nack(method_frame.delivery_tag)  # 告诉队列他要滚回队列去

                '''
                执行后会在 -v 目录下生成 容器的id container_id.log
                usage:  dockertrain  -p  映射到本地的端口 默认8097 如果被占用会自动分配，只检测端口占用情况，可能存在多个未开启的容器相同端口的情况
                                     -n  项目名 默认 ""
                                     -v  需要映射的素材目录(必填)
                                     -d  如果-f选择other那么-d必填，就是-v映射的容器内的目录
                                     -r  docker镜像的地址 默认registry.cn-hangzhou.aliyuncs.com/baymin/ai-power:ai-power-wo-auto-v3.6
                                     -f  训练使用的网络默认detectron,可选darknet和other
                                     -w  root密码 默认icubic-123
                                     -g  复制脚本到/usr/local/bin/，后面执行可以全局dockertrainD
                                     -o  日志的输出目录默认/var/log/train
                                     -t  docker的gup版本，默认是最新版本2，设置1：nvidia-docker，2：docker run --gpus all
                                     -h  帮助
                '''
                # 增加镜像地址进数据库，并且镜像地址外部(队列里的数据)传入
                image_url = None
                docker_volume = None
                if train_info['providerType'] == net_framework.yolov4Tiny3l["name"]:
                    image_url = train_info['image']
                    docker_volume = net_framework.yolov4Tiny3l["train_docker_volume"]
                os.system("echo %s | sudo -S docker pull %s" % (ff.root_password, image_url))
                train_cmd = "dockertrain -n '%s' -v '%s' -w '%s' -t 2 -r '%s' -f '%s' -d '%s'" % \
                            (train_info["taskId"],
                             ff.assets_base_path + "/" + train_info["projectName"],
                             ff.root_password,
                             image_url,
                             train_info['providerType'],
                             docker_volume)
                log.logger.info("\n\n**************************\n训练的命令: %s\n**************************\n" % train_cmd)

                if image_url is None or docker_volume is None:
                    log.logger.info(
                        "\n\n**************************\n镜像地址和映射的容器内目录不可为空\n**************************\n")
                    return

                res = os.popen(train_cmd).read().replace('\n', '')
                catStr = "cat %s/%s/container_id.log" % (ff.assets_base_path, train_info["projectName"])
                container_id = os.popen(catStr).read().replace('\n', '')
                if len(container_id) > 80:
                    container_id = "more than 80"
                # elif len(container_id) < 63:
                #     container_id = "less 63fasterRcnn2"
                if "train_done" not in res:
                    log.logger.info("训练有误: %s" % res)
                    # draw_url = 'http://%s:%d/env/%s' % (ff.draw_host, ff.draw_port, train_info['taskId'])
                    sql = "UPDATE train_record SET container_id='%s', status=%d where task_id='%s'" % \
                          (container_id, -1, train_info['taskId'])
                    log.logger.info("训练:" + sql)
                    ff.postgres_execute(sql)
                    os.system("echo '训练失败\c' > '%s/%s/train_status.log'" %
                              (ff.assets_base_path,
                               train_info["projectName"]))
                    channel.basic_ack(method_frame.delivery_tag)  # 告诉队列可以放行了
                    return
                # 如果res长度==64，那么就是container_id

                os.system("echo '正在训练\c' > '%s/%s/train_status.log'" %
                          (ff.assets_base_path,
                           train_info["projectName"]))

                # region 更新数据库
                sql = "UPDATE train_record SET container_id='%s', status=%d where task_id='%s'" % \
                      (container_id, 2, train_info['taskId'])
                log.logger.info("训练:" + sql)
                ff.postgres_execute(sql)
                # endregion

                # region 初始化画图visdom
                # if ff.draw:
                #     # 保留画图日志，下次打开可直接加载
                #     draw_log = ff.assets_base_path + "/" + train_info["projectName"] + "/draw.log"
                #     ff.draw_windows = Visdom(env=train_info['projectId'], log_to_filename=draw_log)
                #     if os.path.exists(draw_log):
                #         print("已经存在直接加载")
                #         ff.draw_windows.replay_log(draw_log)
                # endregion

            elif "正在训练" in status:
                channel.basic_nack(method_frame.delivery_tag)  # 告诉队列他要滚回队列去
                sql = "UPDATE train_record SET status=%d where task_id='%s'" % (2, train_info['taskId'])
                ff.postgres_execute(sql)
            elif "训练失败" in status:
                channel.basic_ack(method_frame.delivery_tag)  # 告诉队列可以放行了
                ff.postgres_execute(
                    "UPDATE train_record SET status=%d"
                    "where task_id='%s'" %
                    (-1, train_info['taskId']))
            elif "训练完成" in status:
                channel.basic_ack(method_frame.delivery_tag)  # 告诉队列可以放行了
                # region 更新数据库
                ff.postgres_execute(
                    "UPDATE train_record SET status=%d"
                    "where task_id='%s'" %
                    (4, train_info['taskId']))
                os.system("echo %s | sudo -S chmod -R 777 %s/%s" % (ff.root_password, ff.assets_base_path, train_info["projectName"]))
                # endregion
        log.logger.info("训练：%s Basic.GetOk %s delivery-tag %i: %s" % (datetime.now().strftime('%Y-%m-%d %H:%M:%S'),
                                                               header_frame.content_type,
                                                               method_frame.delivery_tag,
                                                               body.decode('utf-8')))
    connection.close()
    return method_frame.delivery_tag, body.decode('utf-8')


# 获取单个训练队列数据
def get_test_one():
    log.logger.info('get_test_one:%s' % (datetime.now().strftime('%Y-%m-%d %H:%M:%S')))
    connection, channel, method_frame, header_frame, body = do_basic_get(queue=ff.test_queue)
    if method_frame is None:
        return None, None
    else:
        test_info = json.loads(body.decode('utf-8'))
        testing = int(os.popen(
            "echo %s | sudo -S docker ps -a | grep %s |wc -l" % (ff.root_password, "power-ai-testing")).read().replace(
            '\n', ''))
        if testing == 0:
            if not os.path.exists("%s/%s/test_status" % (ff.assets_base_path, test_info["assetsDir"])):
                os.system("echo '等待检测\c' > '%s/%s/test_status'" %
                          (ff.assets_base_path, test_info["assetsDir"]))
                channel.basic_nack(method_frame.delivery_tag)
            else:
                status = os.popen("cat '%s/%s/test_status' | head -n 1" %
                                  (ff.assets_base_path, test_info["assetsDir"])).read().replace('\n', '')
                if "等待检测" in status:
                    channel.basic_nack(method_frame.delivery_tag)  # 告诉队列他要滚回队列去
                    image_url = None
                    docker_volume = None
                    if test_info['providerType'] == 'yolov3':
                        docker_volume = "/darknet/assets"
                    elif test_info['providerType'] == 'fasterRcnn':
                        docker_volume = "/Detectron/detectron/datasets/data"
                    elif test_info['providerType'] == 'maskRcnn':
                        docker_volume = "/Detectron/detectron/datasets/data"
                    elif test_info['providerType'] == 'other':
                        docker_volume = test_info['providerOptions']['docker_volume']
                    train_cmd = "dockertrain -n 'Power-Ai-%s' -v '%s' -w '%s' -t 2 -r '%s' -f '%s' -d '%s'" % \
                                (test_info["projectId"],
                                 ff.assets_base_path + "/" + test_info["assetsDir"],
                                 ff.root_password,
                                 image_url,
                                 test_info['providerType'],
                                 docker_volume)
                    log.logger.info("\n\n**************************\n测试的命令: %s\n**************************\n" % train_cmd)
                elif "正在测试" in status:
                    channel.basic_nack(method_frame.delivery_tag)  # 告诉队列他要滚回队列去
                elif "测试失败" in status:
                    channel.basic_ack(method_frame.delivery_tag)  # 告诉队列可以放行了
                elif "测试完成" in status:
                    os.system("echo %s | sudo -S chmod -R 777 %s/%s" % (
                        ff.root_password, ff.assets_base_path, test_info["assetsDir"]))
                    os.system("rm %s/%s/test_status" % (ff.assets_base_path, test_info["assetsDir"]))
                    os.system("tar -C %s -cf %s/%s.tar infer" % (ff.assets_base_path + "/" + test_info["assetsDir"],
                                                                 ff.assets_base_path + "/" + test_info["assetsDir"],
                                                                 test_info["assetsDir"] + "检测结果"))
                    os.system("rm -r %s/%s/infer" % (ff.assets_base_path, test_info["assetsDir"]))
                    channel.basic_ack(method_frame.delivery_tag)  # 告诉队列可以放行了
    connection.close()
    return method_frame.delivery_tag, body.decode('utf-8')


@app.route('/power-ai-train', methods=['POST'])
def do_power_ai_train_http():
    data = request.json  # 获取 JOSN 数据
    # data = data.get('obj')     #  以字典形式获取参数
    if data is not None:
        package_info = {"projectId": data["projectId"], "projectName": data["projectName"],
                        "packageDir": data["packageDir"], "packageName": data["packageName"]}
        # region 更新数据库
        # 这里插入前需要判断是否存在相同的项目
        suc, rows = ff.postgres_execute("SELECT * FROM train_record WHERE project_id='%s'" % package_info['projectId'],
                                        True)
        if rows is None or len(rows) <= 0:
            ff.postgres_execute("INSERT INTO train_record "
                                "(project_id, project_name,"
                                " status, assets_directory_base,"
                                " assets_directory_name, create_time) "
                                "VALUES ('%s', '%s', %d, '%s', '%s', '%s')" %
                                (package_info['projectId'],
                                 package_info['projectName'],
                                 0,
                                 ff.assets_base_path,
                                 package_info["packageDir"],
                                 datetime.now().strftime('%Y-%m-%d %H:%M:%S')))
        else:
            ff.postgres_execute("UPDATE train_record SET "
                                "project_name='%s', status=%d,"
                                " assets_directory_base='%s', assets_directory_name='%s',"
                                " create_time='%s' WHERE project_id='%s'" %
                                (package_info['projectName'],
                                 0,
                                 ff.assets_base_path,
                                 package_info["packageDir"],
                                 datetime.now().strftime('%Y-%m-%d %H:%M:%S'),
                                 package_info['projectId']))
        # endregion
        return Response(json.dumps({"res": "ok"}), mimetype='application/json')
    else:
        return Response(json.dumps({"res": "err"}), mimetype='application/json')


@app.route('/train', methods=['POST'])
def do_train_http():
    data = request.json  # 获取 JOSN 数据
    # data = data.get('obj')     #  以字典形式获取参数
    if data is not None:
        # region 第一步先生成配置文件到本地项目目录 pretrainweightpath: 说明如果为""的话就使用官方的预训练模型
        modeldcfgname = None
        if data['providerType'] == net_framework.yolov4Tiny3l["name"]:
            modeldcfgname = net_framework.yolov4Tiny3l["modeldcfgname"]
        pretraincfgname = "" if data["pretrainWeight"] == "" else data["pretrainWeight"].split("_")[0] + ".cfg"
        os.system(
            "tee %s <<-'EOF'\n"
            "bacthsize: %d\n"
            "maxiter: %d\n"
            "imagesie: [%d,%d]\n"
            "modelname: '%s'\n"
            "triantype: %d\n"
            "pretrainweight: '%s'\n"
            "pretraincfgname: '%s'\n"
            "modeldcfgname: '%s'\n"
            "gpus: '%s'\nEOF" % (
                '{}/{}/config.yaml'.format(ff.assets_base_path, data["projectName"]),
                data["bacthSize"], data["maxIter"], data["imageWidth"], data["imageHeight"], data["taskId"],
                data["trianType"], data["pretrainWeight"], pretraincfgname, modeldcfgname, data["gpus"]))
        os.system("echo 等待训练 > %s" % '{}/{}/train_status.log'.format(ff.assets_base_path, data["projectName"]))
        # endregion
        # region 插入训练队列
        do_basic_publish('ai.train.topic', "train.start.%s" % data['projectName'], json.dumps(data))
        # endregion
        # region 更新数据库
        # 这里插入前需要判断是否存在相同的项目
        suc, rows = ff.postgres_execute("SELECT * FROM train_record WHERE task_id='%s'" % data['taskId'],
                                        True)
        if rows is None or len(rows) <= 0:
            ff.postgres_execute("INSERT INTO train_record "
                                "(task_id, task_name, project_name,"
                                " status, assets_directory_base,"
                                " assets_directory_name, create_time,"
                                " net_framework, assets_type, draw_url, image_url) "
                                "VALUES ('%s', '%s', '%s', %d, '%s', '%s', '%s', '%s', '%s', '%s', '%s')" %
                                (data['taskId'],
                                 data['taskName'],
                                 data['projectName'],
                                 1,
                                 ff.assets_base_path,
                                 data["projectName"],
                                 datetime.now().strftime('%Y-%m-%d %H:%M:%S'),
                                 data["providerType"], data["assetsType"],
                                 'http://{}:{}/{}/train_{}/chart.png'.format(ff.draw_host, ff.draw_port, data["projectName"], data['taskId']),
                                 data["image"]
                                 ))
        else:
            ff.postgres_execute("UPDATE train_record SET "
                                "task_name='%s', project_name='%s', status=%d,"
                                " assets_directory_base='%s', assets_directory_name='%s',"
                                " create_time='%s' WHERE task_id='%s'" %
                                (data['taskName'],
                                 data['project_name'],
                                 1,
                                 ff.assets_base_path,
                                 data["projectName"],
                                 datetime.now().strftime('%Y-%m-%d %H:%M:%S'),
                                 data['task_id']))
        # endregion
        get_train_one()
        return Response(json.dumps({"res": "ok"}), mimetype='application/json')
    else:
        get_train_one()
        return Response(json.dumps({"res": "err"}), mimetype='application/json')


@app.route('/draw_chart', methods=['POST'])
def draw_chat_http():
    try:
        data = request.json
        if data is not None:
            if "爆炸" in json.dumps(data):
                ff.draw_chat(err=True)
            else:
                ff.draw_chat(data)
    except Exception as e:
        try:
            if wechat_monitor:
                my_friend = bot.friends().search('郭永龙')[0]
                my_friend.send('训练溃溃%s' % datetime.now().strftime('%Y-%m-%d %H:%M:%S'))
        except Exception as e:
            log.logger.error('send wechat err')
        # ff.draw_chat(err=True)
        log.logger.error(e)
    return Response(json.dumps({"res": "ok"}), mimetype='application/json')


@app.route('/train_list', methods=['GET'])
def get_train_list_http():
    num = request.args.get('num', type=int, default=20)
    page = request.args.get('page', type=int, default=0)
    offset = num * page
    ret_json = {"num": num, "page": page, "total": 0, "list": []}
    i, count = ff.postgres_execute("SELECT COUNT(*) FROM train_record", True)
    if count is None:
        ret_json["total"] = 0
    else:
        ret_json["total"] = count[0][0]
    a, rows = ff.postgres_execute(
        "SELECT id, task_id, container_id, project_name, status,"
        " net_framework, assets_type, assets_directory_base, assets_directory_name,"
        " is_jump,draw_url,image_url, to_char(create_time, 'YYYY-MM-DD HH24:MI:SS') as create_time, task_name"
        " FROM train_record order by create_time desc limit %d OFFSET %d" % (num, offset), True)
    if rows is None or len(rows) == 0:
        return json.dumps(ret_json)
    else:
        for row in rows:
            ret_json["list"].append(
                {'id': row[0], 'task_id': str(row[1]), 'task_name': str(row[13]), 'container_id': str(row[2]),
                 'project_name': str(row[3]), 'status': row[4], 'net_framework': str(row[5]),
                 'assets_type': str(row[6]), 'assets_directory_base': str(row[7]),
                 'assets_directory_name': str(row[8]), 'is_jump': row[9],
                 'draw_url': str(row[10]), 'image_url': str(row[11]),
                 'create_time': str(row[12])
                 })
        return Response(json.dumps(ret_json), mimetype='application/json')


# region 测试服务的所有接口
@app.route('/get_model_list/<framework_type>/<path>', methods=['GET'])
def get_model_list(framework_type, path):
    weights_list = []
    width = 0
    height = 0
    max_batches = 0
    # framework_type = "yolov3"
    search_path = "/assets/%s/backup/*.weights"
    if framework_type == 'yolov3':
        search_path = "/assets/%s/backup/*.weights"
        cmd = "sed -n '/width/p' %s/yolov3-voc.cfg | sed 's/width//g' |sed 's/=//g' |sed 's/ //g'" % (
                ff.assets_base_path + "/" + path)
        width = os.popen(cmd).read().replace('\n', '')
        height = os.popen("sed -n '/height/p' %s/yolov3-voc.cfg | sed 's/height//g' |sed 's/=//g' |sed 's/ //g'" % (
                ff.assets_base_path + "/" + path)).read().replace('\n', '')
        max_batches = os.popen(
            "sed -n '/max_batches/p' %s/yolov3-voc.cfg | sed 's/max_batches//g' |sed 's/=//g' |sed 's/ //g'" % (
                    ff.assets_base_path + "/" + path)).read().replace('\n', '')
    elif framework_type == 'fasterRcnn' or framework_type == 'maskRcnn':
        search_path = "/assets/%s/result/train/coco_2014_train/generalized_rcnn/*.pkl"
        max_batches = os.popen(
            "sed -n '/MAX_ITER/p' %s/train-config.yaml | sed 's/MAX_ITER//g' |sed 's/://g' |sed 's/ //g'" % (
                    ff.assets_base_path + "/" + path)).read().replace('\n', '')
    # region detectron2
    elif framework_type == 'fasterRcnn2' or framework_type == 'maskRcnn2' or framework_type == 'keypointRcnn2':
        search_path = "/assets/%s/output/*.pth"
        max_batches = os.popen(
            "sed -n '/MAX_ITER/p' %s/train-config.yaml | sed 's/MAX_ITER//g' |sed 's/://g' |sed 's/ //g'" % (
                    ff.assets_base_path + "/" + path)).read().replace('\n', '')
    # endregion
    elif framework_type == 'other':
        search_path = "不支持，后续开发"

    for item in sorted(glob.glob(search_path % path), key=os.path.getctime,
                       reverse=True):  # key 根据时间排序 reverse true表示倒叙
        filepath, tempfilename = os.path.split(item)
        if "server.pkl" in tempfilename or "test.weights" in tempfilename:
            continue
        weights_list.append({"path": item, "filename": tempfilename})

    return Response(json.dumps(
        {"res": "ok", "weights_list": weights_list, "width": width, "height": height, "max_batches": max_batches}),
                    mimetype='application/json')


@app.route('/get_val_path_list', methods=['GET'])
def get_val_path_list():
    path_list = []
    # framework_type = "yolov3"
    search_path = "/assets/StandardValidationData/*"
    for item in sorted(glob.glob(search_path), key=os.path.getctime,
                       reverse=True):  # key 根据时间排序 reverse true表示倒叙
        filepath, tempfilename = os.path.split(item)
        if os.path.isdir(item):
            path_list.append({"path": item, "dir_name": tempfilename})

    return Response(json.dumps({"res": "ok", "val_path_list": path_list}), mimetype='application/json')


@app.route('/get_voc_path_list', methods=['GET'])
def get_voc_path_list():
    path_list = []
    # framework_type = "yolov3"
    search_path = "/PowerAiData/*"
    for item in sorted(glob.glob(search_path), key=os.path.getctime,
                       reverse=True):  # key 根据时间排序 reverse true表示倒叙
        filepath, tempfilename = os.path.split(item)
        if os.path.isdir(item):
            path_list.append({"path": item, "dir_name": tempfilename})

    return Response(json.dumps({"res": "ok", "voc_path_list": path_list}), mimetype='application/json')


@app.route('/start_test', methods=['POST'])  # 这里新增  weights: undefined, valPath: undefined,
def start_test():
    try:
        data = request.json
        if data is not None:
            docker_volume = "/darknet/assets"
            docker_volume_model = "/darknet/assets/backup/yolov3-voc_last.weights"
            port = data['port']  # 检测服务的端口 容器内外都一致，java 中间件对外端口在这个基础上+1
            docker_name = "darknet-service-testing"  # 用于每次新建 删除
            if data['providerType'] == 'yolov3':
                docker_name = "darknet-service-testing"
                docker_volume = "/darknet/assets"
                docker_volume_model = "/darknet/assets/backup/yolov3-voc_last.weights"
            elif data['providerType'] == 'fasterRcnn':
                docker_name = "detectron-service-testing"
                docker_volume = "/Detectron/detectron/datasets/data"
                docker_volume_model = "/Detectron/detectron/datasets/data/result/train/coco_2014_train/generalized_rcnn/server.pkl"
            elif data['providerType'] == 'maskRcnn':
                docker_name = "detectron-service-testing"
                docker_volume = "/Detectron/detectron/datasets/data"
                docker_volume_model = "/Detectron/detectron/datasets/data/result/train/coco_2014_train/generalized_rcnn/server.pkl"
            elif data['providerType'] == 'fasterRcnn2' or data['providerType'] == 'maskRcnn2':
                docker_name = "detectron2-service-testing"
                docker_volume = "/detectron2/datasets"
                docker_volume_model = "/detectron2/datasets/model_test.pth"
            elif data['providerType'] == 'other':
                docker_name = "other-service-testing"
                docker_volume = data['docker_volume']
                docker_volume_model = data['docker_volume_model']

            os.system('cp %s/yolov3-voc.cfg %s/yolov3-voc-test.cfg' % (
                ff.assets_base_path + "/" + data["assetsDir"],
                ff.assets_base_path + "/" + data["assetsDir"]))

            os.system('sed -i "s/^batch.*/batch=1/g" %s/yolov3-voc-test.cfg' % (
                        ff.assets_base_path + "/" + data["assetsDir"]))
            os.system('sed -i "s/^subdivisions.*/subdivisions=1/g" %s/yolov3-voc-test.cfg' % (
                        ff.assets_base_path + "/" + data["assetsDir"]))

            os.system("echo '%s' | sudo -S docker stop '%s'" % (ff.root_password, docker_name))
            cmd = "echo %s | sudo -S docker run --gpus '\"device=5\"' \
                        --name %s \
                        -p %d:8070 \
                        -p %d:%d \
                        -v /opt/remote_train_web/aiimg/:/aiimg \
                        -v /opt/remote_train_web/excel/:/excel \
                        -v /etc/localtime:/etc/localtime:ro \
                        -v '%s':'%s' \
                        -v '%s':'%s' \
                        --add-host service-postgresql:10.10.0.4 \
                        --add-host service-rabbitmq:10.10.0.3 \
                        --add-host service-ftp:10.10.0.2 \
                        --add-host service-web:10.10.0.5 \
                        --rm -d %s" % (
                ff.root_password,
                docker_name,
                port + 1,
                port, port,
                ff.assets_base_path + "/" + data["assetsDir"], docker_volume,
                data["weights"], docker_volume_model,
                data["image"])
            log.logger.info("\n\n**************************\n测试的命令: %s\n**************************\n" % cmd)
            os.system(cmd)
            start_time = time.time()
            while 1:
                time.sleep(0.5)
                if net_is_used(int(port)):
                    break
                if (time.time() - start_time) > 60:
                    return Response(json.dumps({"res": "err", "msg": "开启超时，请手动查询状态"}), mimetype='application/json')
            # registry.cn-hangzhou.aliyuncs.com/baymin/ai-power:darknet_auto_test-service-ai-power-v4.5
    except Exception as e:
        log.logger.error(e)
        return Response(json.dumps({"res": "err", "msg": "开启失败"}), mimetype='application/json')
    return Response(json.dumps({"res": "ok", "msg": "开启成功"}), mimetype='application/json')


# endregion 测试服务的所有接口


@app.route('/stop_train', methods=['POST'])  # 此处提交的参数projectId改为taskId
def stop_train_http():
    try:
        data = request.json
        if data is not None:
            cmd = "echo '%s' | sudo -S docker stop `cat '%s/%s/train.dname'`" % (
                ff.root_password, ff.assets_base_path, data['assetsDir'])
            log.logger.info('停止训练:%s' % cmd)
            os.system(cmd)  # 会自动退出，所以这里不需要了
            ff.postgres_execute("UPDATE train_record SET "
                                "status=%d WHERE project_id='%s'" %
                                (3, data['projectId']))
    except Exception as e:
        log.logger.error(e)
        return Response(json.dumps({"res": "err"}), mimetype='application/json')
    return Response(json.dumps({"res": "ok"}), mimetype='application/json')


@app.route('/get_record_by_project_id', methods=['POST'])
def get_record_by_project_id():
    try:
        data = request.json
        if data is not None:
            suc, rows = ff.postgres_execute(
                "select net_framework as providerType, assets_directory_name as assets from train_record where project_id ='%s'" %
                data['projectId'],
                True)
            if suc:
                if rows is not None and len(rows) > 0:
                    return Response(json.dumps({"res": "ok", "providerType": rows[0][0], "assets": rows[0][1]}),
                                    mimetype='application/json')
                else:
                    return Response(json.dumps({"res": "ok", "providerType": "none", "assets": "none"}),
                                    mimetype='application/json')
    except Exception as e:
        log.logger.logger.error(e)
        return Response(json.dumps({"res": "err"}), mimetype='application/json')
    return Response(json.dumps({"res": "ok"}), mimetype='application/json')


@app.route('/restart_train', methods=['POST'])
def restart_train_http():
    try:
        data = request.json
        if data is not None:
            trainInfo = {"projectId": data["projectId"],
                         "projectName": data["projectName"],
                         "assetsDir": data["assetsDir"],
                         "assetsType": data["assetsType"],
                         "providerType": data["providerType"],
                         "providerOptions": {"yolov3Image": data["image"]}
                         }

            # config.write(open("test.cfg", "w"))
            docker_volume = "/darknet/assets"
            docker_volume_model = "/darknet/assets/yiningzeng.weights"
            if data['providerType'] == 'yolov3':
                docker_volume = "/darknet/assets"
                docker_volume_model = "/darknet/assets/yiningzeng.weights"
                if "width" in data:
                    os.system('sed -i "s/^width.*/width=%s/g" %s/yolov3-voc.cfg' % (
                        data["width"], ff.assets_base_path + "/" + data["assetsDir"]))
                if "height" in data:
                    os.system('sed -i "s/^height.*/height=%s/g" %s/yolov3-voc.cfg' % (
                        data["height"], ff.assets_base_path + "/" + data["assetsDir"]))
                if "max_batches" in data:
                    os.system('sed -i "s/^max_batches.*/max_batches=%d/g" %s/yolov3-voc.cfg' % (
                        data["max_batches"], ff.assets_base_path + "/" + data["assetsDir"]))
                    os.system('sed -i "s/^steps.*/steps=%d,%d/g" %s/yolov3-voc.cfg' % (
                        int(int(data["max_batches"]) * 0.8), int(int(data["max_batches"]) * 0.9),
                        ff.assets_base_path + "/" + data["assetsDir"]))
            # region detectron
            elif data['providerType'] == 'fasterRcnn' or data['providerType'] == 'maskRcnn':
                trainInfo["providerOptions"] = {"fasterRcnnImage": data["image"]}
                docker_volume = "/Detectron/detectron/datasets/data"
                docker_volume_model = "/Detectron/models/R-50.pkl"

                if "max_batches" in data:  # fasterRcnn 和 maskRcnn 暂时不先替换掉steps
                    os.system('sed -i "s/    MAX_ITER.*/    MAX_ITER: %d/g" %s/train-config.yaml' % (
                        data["max_batches"], ff.assets_base_path + "/" + data["assetsDir"]))
                    # os.system('sed -i "s/^STEPS.*/steps=%d,%d/g" %s/yolov3-voc.cfg' % (
                    #     int(int(data["max_batches"]) * 0.8), int(int(data["max_batches"]) * 0.9),
                    #     ff.package_base_path + "/" + data["assetsDir"]))
            # endregion
            # region detectron2
            elif data['providerType'] == 'fasterRcnn2' or data['providerType'] == 'maskRcnn2' or data[
                'providerType'] == 'keypointRcnn2':
                trainInfo["providerOptions"] = {"detectron2Image": data["image"]}
                docker_volume = "/detectron2/datasets"
                docker_volume_model = "/detectron2/models/R-50.pkl"

                if "max_batches" in data:  # fasterRcnn 和 maskRcnn 暂时不先替换掉steps
                    os.system('sed -i "s/  MAX_ITER.*/  MAX_ITER: %d/g" %s/train-config.yaml' % (
                        data["max_batches"], ff.assets_base_path + "/" + data["assetsDir"]))
                if "weights" in data:  # fasterRcnn 和 maskRcnn 暂时不先替换掉steps
                    os.system('sed -i "s/  WEIGHTS.*/  WEIGHTS: %s/g" %s/train-config.yaml' % (
                        data["weights"], ff.assets_base_path + "/" + data["assetsDir"]))
            # endregion
            elif data['providerType'] == 'other':
                trainInfo["providerOptions"] = {"otherImage": data["image"]}
                docker_volume = data['docker_volume']
                docker_volume_model = data['docker_volume_model']

            if "yolov3-voc_last.weights" not in data["assetsDir"]:
                os.system("echo %s | sudo -s rm %s/backup/yolov3-voc_last.weights" % (
                ff.root_password, ff.assets_base_path + "/" + data["assetsDir"]))
            # 删除数据转换的状态文件
            os.system("echo %s | sudo -s rm %s/train_log/convert_data.log" % (
            ff.root_password, ff.assets_base_path + "/" + data["assetsDir"]))

            # 写入正在训练，否则队列会重新执行
            os.system("echo '正在训练\c' > '%s/%s/train_status.log'" %
                      (ff.assets_base_path, data["assetsDir"]))

            # 加入到训练队列
            do_basic_publish('ai.train.topic', "train.start.%s" % data['projectId'], json.dumps(trainInfo))

            cmd = "echo %s | sudo -S docker run --shm-size 32G --memory-swap -1 --rm --gpus '\"device=0,1,2,3,4\"' \
                        --name %s \
                        -v /etc/localtime:/etc/localtime:ro \
                        -v '%s':'%s' \
                        -v '%s':'%s' \
                        --add-host service-postgresql:10.10.0.4 \
                        --add-host service-rabbitmq:10.10.0.3 \
                        --add-host service-ftp:10.10.0.2 \
                        --add-host service-web:10.10.0.5 \
                        --rm -d %s" % (
                ff.root_password,
                data['projectId'].replace("_", ""),
                ff.assets_base_path + "/" + data["assetsDir"], docker_volume,
                data["weights"], docker_volume_model,
                data["image"])
            os.system("echo '%s\c' > '%s/%s/train.dname'" %
                      (data['projectId'].replace("_", ""), ff.assets_base_path, data["assetsDir"]))

            log.logger.info("\n\n**************************\n重新训练: %s\n**************************\n" % cmd)
            os.system(cmd)
            ff.postgres_execute("UPDATE train_record SET "
                                "status=%d WHERE project_id='%s'" %
                                (2, data['projectId']))
    except Exception as e:
        log.logger.error(e)
        return Response(json.dumps({"res": "err"}), mimetype='application/json')
    return Response(json.dumps({"res": "ok"}), mimetype='application/json')


# region 新版本的新增接口
# 获取所属项目列表
@app.route('/get_local_projects', methods=['GET'])
def get_local_projects():
    path_list = []
    # framework_type = "yolov3"
    search_path = "/assets/Projects/*"
    for item in sorted(glob.glob(search_path), key=os.path.getctime,
                       reverse=True):  # key 根据时间排序 reverse true表示倒叙
        filepath, tempfilename = os.path.split(item)
        if os.path.isdir(item):
            path_list.append({"path": item, "dir_name": tempfilename})

    return Response(json.dumps({"res": "ok", "path_list": path_list}), mimetype='application/json')


# 获取所有模型
@app.route('/get_model_list_v2/<framework_type>/<project_name>', methods=['GET'])
def get_model_list_v4(framework_type, project_name):
    model_list = []
    # framework_type = "yolov3"
    search_path = "/assets/Projects/%s/%s/*%s"
    if framework_type == net_framework.yolov4Tiny3l["name"]:
        search_path = search_path % (project_name, net_framework.yolov4Tiny3l["modelSavePath"], net_framework.yolov4Tiny3l["modelSuffix"])

    for item in sorted(glob.glob(search_path), key=os.path.getctime,
                       reverse=True):  # key 根据时间排序 reverse true表示倒叙
        filepath, tempfilename = os.path.split(item)
        if "server.pkl" in tempfilename or "test.weights" in tempfilename:
            continue
        model_list.append({"path": item, "filename": tempfilename})

    return Response(json.dumps({"res": "ok", "model_list": model_list}), mimetype='application/json')


# 训练中心通过项目名称获取当下的模型列表
@app.route('/get_models/<project_name>', methods=['GET'])
def get_project_models(project_name):
    models = []
    # framework_type = "yolov3"
    search_path = "/assets/Projects"
    now_model = ""
    # 先获取当前发布的模型
    for item in sorted(glob.glob(search_path + "/%s/model_release/yolov4-tiny-3l/*.weights" % project_name), key=os.path.getctime, reverse=True):
        _, now_model = os.path.split(item)
    for item in sorted(glob.glob(search_path+"/%s/backup/*.weights" % project_name), key=os.path.getctime, reverse=True):  # key 根据时间排序 reverse true表示倒叙
        path, name = os.path.split(item)
        status = 0
        # if "final" in name or "last" in name or "best" in name:
        if "release" in name:
            status = 1
        if now_model == name:
            status = 2
        models.append({"name": name, "path": item, "status": status})
    return Response(json.dumps({"res": "ok", "message": "获取成功", "models": models}), mimetype='application/json')


# 训练中心通过项目名称获取当下的模型列表
@app.route('/get_release_models_history/<project_name>', methods=['GET'])
def get_project_relase_models_history(project_name):
    models = []
    # framework_type = "yolov3"
    search_path = "/assets/Projects"
    now_model = ""
    # 先获取当前发布的模型
    for item in sorted(glob.glob(search_path + "/%s/model_release/yolov4-tiny-3l/*.weights" % project_name), key=os.path.getctime, reverse=True):
        _, now_model = os.path.split(item)
    for item in sorted(glob.glob(search_path+"/%s/model_release_history/*.weights" % project_name), key=os.path.getmtime, reverse=True):  # key 根据时间排序 reverse true表示倒叙
        path, name = os.path.split(item)
        status = 0
        if now_model == name:  ## STATUS =0  表示 不是最新发布的版本， =1标识是当前最新的版本
            status = 1
        models.append({"name": name, "path": item, "status": status})
    return Response(json.dumps({"res": "ok", "message": "获取成功", "models": models}), mimetype='application/json')


# 模型删除
@app.route('/del_model', methods=['DELETE'])
def delete_model():
    p = unquote(request.args.get('p'))
    os.system("echo %s | sudo -S rm %s" % (ff.root_password, p.replace("backup", "model_release_history")))  # 重命名模型文件
    os.system("echo %s | sudo -S rm %s" % (ff.root_password, p))
    return Response(json.dumps({"res": "ok", "message": "成功"}), mimetype='application/json')


# 模型发布并且上线
@app.route('/online_model', methods=['PUT'])
def online_model():
    is_history = request.args.get('is_history', type=int, default=0)
    p = unquote(request.args.get('p'))
    release = "-" + datetime.now().strftime('%Y%m%d') + "-release"
    if is_history == 1:
        p = p.replace("model_release_history", "backup")
        release = ''
    basePath, name = os.path.split(p)
    fname, fename = os.path.splitext(name)
    finalFilename = "%s/%s%s%s" % (basePath, fname, release, fename)
    configFilename = "%s/%s.*" % (basePath, name.split('_')[0])
    os.system("echo %s | sudo -S rm %s/*" % (ff.root_password, basePath.replace("backup", "model_release/yolov4-tiny-3l")))  # 重命名模型文件
    os.system("echo %s | sudo -S mv %s %s" % (ff.root_password, p, finalFilename))  # 重命名模型文件
    os.system("echo %s | sudo -S ln %s %s" % (ff.root_password, finalFilename, basePath.replace("backup", "model_release/yolov4-tiny-3l")))  # 硬链接模型文件
    os.system("echo %s | sudo -S cp -r %s %s" % (ff.root_password, configFilename, basePath.replace("backup", "model_release/yolov4-tiny-3l")))  # 硬链接配置文件

    os.system("echo %s | sudo -S ln -s %s %s" % (ff.root_password, finalFilename, basePath.replace("backup", "model_release_history")))  # 硬链接模型文件
    os.system("echo %s | sudo -S ln -s %s %s" % (ff.root_password, configFilename, basePath.replace("backup", "model_release_history")))  # 硬链接配置文件
    return Response(json.dumps({"res": "ok", "message": "成功"}), mimetype='application/json')


# 这是给AOI的升级模型接口
@app.route('/get_models', methods=['GET'])
def get_models():
    projects = []
    # framework_type = "yolov3"
    search_path = "/assets/Projects"
    httpUrl = "http://192.168.31.102:1121"
    for item in sorted(glob.glob(search_path+"/*"), key=os.path.getctime, reverse=True):  # key 根据时间排序 reverse true表示倒叙
        _, project = os.path.split(item)
        if os.path.isdir(item):
            one_project = {"project_name": project, "list": []}
            f_path = search_path + "/" + project + "/model_release/"
            net_frameworks = []
            for one_framework in glob.glob(f_path + "*"):
                _, one_framework_dir_name = os.path.split(one_framework)
                fra = {"net_framework": one_framework_dir_name, "models": []}
                for one_model in sorted(glob.glob(f_path + one_framework_dir_name + "/*"), key=os.path.getctime, reverse=True):
                    baseUrl = one_model.replace(search_path, httpUrl)
                    fra["models"].append(baseUrl)
                one_project["list"].append(fra)
            projects.append(one_project)
    return Response(json.dumps({"res": 0, "message": "获取成功", "project_list": projects}), mimetype='application/json')
# endregion


if __name__ == '__main__':
    log = Logger('server.log', when="S")
    if not os.path.isfile("/usr/local/bin/dockertrain"):
        print("dockertrain 执行文件不存在")
        exit(99)
    debug = False
    wechat_monitor = False
    if debug:
        ff = pikaqiu(root_password='baymin1024', host='192.168.31.77', username='baymin', password='baymin1024',
                     assets_base_path='/assets/Projects')
    else:
        ff = pikaqiu(root_password='baymin1024', host='192.168.31.77', username='baymin', password='baymin1024',
                     assets_base_path='/assets/Projects')
    # init(self, sql=True, sql_host='localhost', draw=True, draw_host='localhost', draw_port=8097):
    # sql: 是否开启数据库，sql_host：数据库地址，draw：是否开启画图，draw_host：画图的服务地址，draw_port：画图的服务端口
    ff.init(sql=False, sql_host='192.168.31.77', draw_host='localhost', draw_port=1121)

    # 创建后台执行的 schedulers
    scheduler = BackgroundScheduler()
    # 添加调度任务

    # 提醒写日报
    # scheduler.add_job(remind, 'cron', second="0/2")

    '''
    weeks(int)  间隔几周
    days(int)   间隔几天
    hours(int)  间隔几小时
    minutes(int)        间隔几分钟
    seconds(int)        间隔多少秒
    start_date(datetime or str) 开始日期
    end_date(datetime or str)   结束日期
    timezone(datetime.tzinfo or   str)  时区
    '''
    get_train_one()
    if debug:
        scheduler.add_job(get_train_one, 'interval', minutes=10000)
        scheduler.add_job(get_test_one, 'interval', seconds=5)
    else:
        if wechat_monitor:
            bot = Bot(cache_path=True, console_qr=True)
            myself = bot.self
            my_friend = bot.friends().search('郭永龙')[0]
            # my_friend.send('微信监督开始')
            bot.file_helper.send('微信监督开始')
        scheduler.add_job(get_train_one, 'interval', minutes=10)
        scheduler.add_job(get_test_one, 'interval', seconds=5)
    # scheduler.add_job(get_package_one, 'interval', seconds=5)
    scheduler.start()

    log.logger.info("start")
    # ff.consume(ch, on_message_callback)
    app.run(host="0.0.0.0", port=8080)
    embed()
