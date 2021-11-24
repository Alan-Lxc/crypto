<template>
    <div>
        <el-row :gutter="20">
            <el-col :span="8">
                <el-card shadow="hover" class="mgb20" style="height:252px;">
                    <div class="user-info">
                        <img src="../assets/img/img.jpg" class="user-avator" alt />
                        <div class="user-info-cont">
                            <div class="user-info-name">{{ name }}</div>
                            <div>{{ role }}</div>
                        </div>
                    </div>
                    <div class="user-info-list">
                        上次登录时间：
                        <span>2021-10-19</span>
                    </div>
                    <div class="user-info-list">
                        上次登录地点：
                        <span>北京</span>
                    </div>
                </el-card>
                <el-card shadow="hover" style="height:252px;">
                    <template #header>
                        <div class="clearfix">
                            <span>秘密类型</span>
                        </div>
                    </template>

                    文件<el-progress :percentage="71.3" color="#42b983"></el-progress>
                    口令<el-progress :percentage="24.1" color="#f1e05a"></el-progress>
                    其他<el-progress :percentage="13.7"></el-progress>
                    JSON<el-progress :percentage="5.9" color="#f56c6c"></el-progress>
                </el-card>

            </el-col>
            <el-col :span="16">
                <el-row :gutter="20" class="mgb20">
                    <el-col :span="8">
                        <el-card shadow="hover" :body-style="{ padding: '0px' }">
                            <div class="grid-content grid-con-1">
                                <i class="el-icon-user-solid grid-con-icon"></i>
                                <div class="grid-cont-right">
                                    <div class="grid-num" >{{currentSecretNum}}</div>
                                    <div>存储秘密数</div>
                                </div>
                            </div>
                        </el-card>
                    </el-col>
                    <el-col :span="8">
                        <el-card shadow="hover" :body-style="{ padding: '0px' }">
                            <div class="grid-content grid-con-2">
                                <i class="el-icon-message-solid grid-con-icon"></i>
                                <div class="grid-cont-right">
                                    <div class="grid-num">2</div>
                                    <div>系统消息</div>
                                </div>
                            </div>
                        </el-card>
                    </el-col>
<!--                    <el-col :span="8">-->
<!--                        <el-card shadow="hover" :body-style="{ padding: '0px' }">-->
<!--                            <div class="grid-content grid-con-3">-->
<!--                                <i class="el-icon-s-goods grid-con-icon"></i>-->
<!--                                <div class="grid-cont-right">-->
<!--                                    <div class="grid-num">5000</div>-->
<!--                                    <div>数量</div>-->
<!--                                </div>-->
<!--                            </div>-->
<!--                        </el-card>-->
<!--                    </el-col>-->
                </el-row>
                <el-card shadow="hover" style="height:200px;">
                    <template #header>
                        <div class="clearfix">
                            <span>最新消息</span>
                            <el-button style="float: right; padding: 3px 0" type="text">添加</el-button>
                        </div>
                    </template>

                    <el-table :show-header="false" :data="todoList" style="width:100%;">
                        <el-table-column width="40">
                            <template #default="scope">
                                <el-checkbox v-model="scope.row.status"></el-checkbox>
                            </template>
                        </el-table-column>
                        <el-table-column>
                            <template #default="scope">
                                <div class="todo-item" :class="{
                                        'todo-item-del': scope.row.status,
                                    }">{{ scope.row.title }}</div>
                            </template>
                        </el-table-column>
                        <el-table-column width="60">
                            <template>
                                <i class="el-icon-edit"></i>
                                <i class="el-icon-delete"></i>
                            </template>
                        </el-table-column>
                    </el-table>
                </el-card>
              <el-card shadow="hover" style="height:500px;">
                <div class="rightullidiv">
                  <img
                      src="../assets/img/indexpage1.png"
                      alt=""
                      class="rightulliimg"
                  >
                </div>




              </el-card>
            </el-col>
        </el-row>
<!--        <el-row :gutter="20">-->
<!--            <el-col :span="12">-->
<!--                <el-card shadow="hover">-->
<!--                    <schart ref="bar" class="schart" canvasId="bar" :options="options"></schart>-->
<!--                </el-card>-->
<!--            </el-col>-->
<!--            <el-col :span="12">-->
<!--                <el-card shadow="hover">-->
<!--                    <schart ref="line" class="schart" canvasId="line" :options="options2"></schart>-->
<!--                </el-card>-->
<!--            </el-col>-->
<!--        </el-row>-->
    </div>
</template>

<script>
import Schart from "vue-schart";
import {reactive, ref} from "vue";
import axios from "axios";
export default {
    name: "dashboard",
    components: { Schart },
  data(){
      return{
        currentSecretNum : 0,
      }
  },

  created() {
      this.getCurrentSecretNum();
    },
    methods:{
      getCurrentSecretNum(){
        let that = this
        const url = "http://localhost:8080/api/secret/getsecretlist";
        axios({
          methods: 'get',
          url:url,
          params: {
            "userid": 1,
          },

        }).then(
            function (res) {
              console.log(res.data.data.total)
              that.currentSecretNum = res.data.data.total;
              console.log(that.currentSecretNum);
            }
        ).catch(err =>{
          console.log(err);
        })
      },
    },

  setup() {
        const name = localStorage.getItem("ms_username");
        const role = name === "admin" ? "超级管理员" : "普通用户";

        const data = reactive([
            {
                name: "2018/09/04",
                value: 1083,
            },
            {
                name: "2018/09/05",
                value: 941,
            },
            {
                name: "2018/09/06",
                value: 1139,
            },
            {
                name: "2018/09/07",
                value: 816,
            },
            {
                name: "2018/09/08",
                value: 327,
            },
            {
                name: "2018/09/09",
                value: 228,
            },
            {
                name: "2018/09/10",
                value: 1065,
            },
        ]);
        const options = {
            type: "bar",
            title: {
                text: "最近一周各品类销售图",
            },
            xRorate: 25,
            labels: ["周一", "周二", "周三", "周四", "周五"],
            datasets: [
                {
                    label: "家电",
                    data: [234, 278, 270, 190, 230],
                },
                {
                    label: "百货",
                    data: [164, 178, 190, 135, 160],
                },
                {
                    label: "食品",
                    data: [144, 198, 150, 235, 120],
                },
            ],
        };
        const options2 = {

        };
        const todoList = reactive([
            {
                title: "您已经成功登录",
                status: false,
            }, {
                title: "您已经成功上传一个秘密\n秘密名称:123\n详情请见秘密列表",
                status: false,
            },

        ]);




        return {
            name,
            data,
            options,
            options2,
            todoList,
            role,
        };
    },
};
</script>

<style scoped>
.rightullidiv {
  width: 100%;
  height:100%;
  background: #f2f2f2;
  display: flex;
  justify-content: center;
  align-items: center;

}

.rightulliimg {
  max-width: 100%;
  max-height: 100%;
}

.el-row {
    margin-bottom: 20px;
}

.grid-content {
    display: flex;
    align-items: center;
    height: 100px;
}

.grid-cont-right {
    flex: 1;
    text-align: center;
    font-size: 14px;
    color: #999;
}

.grid-num {
    font-size: 30px;
    font-weight: bold;
}

.grid-con-icon {
    font-size: 50px;
    width: 100px;
    height: 100px;
    text-align: center;
    line-height: 100px;
    color: #fff;
}

.grid-con-1 .grid-con-icon {
    background: rgb(45, 140, 240);
}

.grid-con-1 .grid-num {
    color: rgb(45, 140, 240);
}

.grid-con-2 .grid-con-icon {
    background: rgb(100, 213, 114);
}

.grid-con-2 .grid-num {
    color: rgb(45, 140, 240);
}

.grid-con-3 .grid-con-icon {
    background: rgb(242, 94, 67);
}

.grid-con-3 .grid-num {
    color: rgb(242, 94, 67);
}

.user-info {
    display: flex;
    align-items: center;
    padding-bottom: 20px;
    border-bottom: 2px solid #ccc;
    margin-bottom: 20px;
}

.user-avator {
    width: 120px;
    height: 120px;
    border-radius: 50%;
}

.user-info-cont {
    padding-left: 50px;
    flex: 1;
    font-size: 14px;
    color: #999;
}

.user-info-cont div:first-child {
    font-size: 30px;
    color: #222;
}

.user-info-list {
    font-size: 14px;
    color: #999;
    line-height: 25px;
}

.user-info-list span {
    margin-left: 70px;
}

.mgb20 {
    margin-bottom: 20px;
}

.todo-item {
    font-size: 14px;
}

.todo-item-del {
    text-decoration: line-through;
    color: #999;
}

.schart {
    width: 100%;
    height: 300px;
}
</style>
