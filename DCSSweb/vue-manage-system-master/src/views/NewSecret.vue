<template>
  <div>
    <div class="crumbs">
      <el-breadcrumb separator="/">
        <el-breadcrumb-item><i class="el-icon-lx-vipcard"></i> 新建秘密</el-breadcrumb-item>
      </el-breadcrumb>
    </div>
    <div class="container">
      <div class="form-box">
        <el-form ref="secretRef" :rules="rules" :model="secret" label-width="160px">
          <el-form-item label="秘密名称" prop="secretname">
            <el-input v-model="secret.secretname"></el-input>
          </el-form-item>
          <el-form-item label="门限阈值" >
            <el-input-number v-model.number="secret.degree" :min="1" :max="100"></el-input-number>
          </el-form-item>
          <el-form-item label="委员会成员数" prop="counter">
            <el-input-number v-model.number="secret.counter" :min="1" :max="100" ></el-input-number>
          </el-form-item>
          <el-form-item label="秘密值" >
            <el-input v-model.number="secret.secret"> </el-input>
            <!--            <div class="content-title">支持拖拽</div>-->
            <!--&lt;!&ndash;            <div class="plugins-tips">&ndash;&gt;-->
            <!--&lt;!&ndash;              Element UI自带上传组件。&ndash;&gt;-->
            <!--&lt;!&ndash;              访问地址：&ndash;&gt;-->
            <!--&lt;!&ndash;              <a href="http://element.eleme.io/#/zh-CN/component/upload" target="_blank">Element UI Upload</a>&ndash;&gt;-->
            <!--&lt;!&ndash;            </div>&ndash;&gt;-->
            <!--            <el-upload class="upload-demo" drag action="http://jsonplaceholder.typicode.com/api/posts/" multiple>-->
            <!--              <i class="el-icon-upload"></i>-->
            <!--              <div class="el-upload__text">-->
            <!--                将秘密文件拖到此处，或-->
            <!--                <em>点击上传</em>-->
            <!--              </div>-->
            <!--&lt;!&ndash;              <template #tip>&ndash;&gt;-->
            <!--&lt;!&ndash;                <div class="el-upload__tip">只能上传 jpg/png 文件，且不超过 500kb</div>&ndash;&gt;-->
            <!--&lt;!&ndash;              </template>&ndash;&gt;-->
            <!--            </el-upload>-->
          </el-form-item>
          <el-form-item  label="秘密描述">
            <el-input v-model="secret.description"></el-input>

          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="onSubmit">提交秘密</el-button>
            <el-button type="danger" @click="onReset">重置表单</el-button>
          </el-form-item>

        </el-form>

      </div>

    </div>
  </div>
</template>

<script>
import {ref, reactive} from "vue";
import { ElMessage } from "element-plus";
import {useRouter} from "vue-router";
import axios from "axios";
export default {
  name: "NewSecret",
  data() {
  },
  headers:{
    'Content-Type':'applicaion/x-www-form-urlencoded;charset=UTF-8'
  },
  methods:{

  },
  setup() {
    const router = useRouter();
    const secretRef = ref(null);
    const secret = reactive({
      secretname: "",
      degree: 0,//t值
      counter: 0,//n值
      userId: 1,
      secret: "",
      description: "",
    });

    // 提交:model="newcounter"
    const onSubmit = () => {
      // 表单校验
      console.log(secret.degree,secret.counter)
      const api = "http://localhost:8080/api/secret/newsecret"
      secretRef.value.validate((valid) => {
        //验证输入的数据符合格式规范
        if (valid) {
          //验证t，n是否符合规范
          if (secret.degree*2+1>secret.counter){
            // console.log(form);
            console.log("t*2+1>n");
            ElMessage.error("参数不符合规范");
            return false;
          }else {//如果t，n符合规范，就向后端传值
            console.log(secret.degree,secret.counter)
            axios({
              method:'post',
              url:api,
              data:secret,
              transformRequest:[function (data) {
                let param = '';
                for (var it in data){
                  param += it + '=' + data[it] + '&'
                }
                return param
              }]
            }).then(function (response) {
              console.log(response.data)
              //判断
            })
            ElMessage.success("提交成功！");
            router.push('/secretlist');
          }
        } else {
          ElMessage.error("缺少必要的输入项");
          // return false;
        }
      });
    };
    // 重置
    const onReset = () => {
      secretRef.value.resetFields();
    };
    //检验t,n代数关系
    const checkNum = (rule, value, callback) => {
      if (secret.counter<secret.degree*2+1){
        callback(new Error("不对"));
      }else {
        console.log("对不对");
        callback();
      }
    }
    const rules = {
      secretname: [
        { required: true, message: "请输入秘密名称", trigger: "blur"},
      ],
      counter: [
        {  validator: checkNum, message: "委员会成员数需要大于2倍的门限阈值", trigger: "change"}
      ],

    };
    return {
      rules,
      secret,
      secretRef,
      onSubmit,
      onReset,
      checkNum,
    };
  },

};
</script>

<style scoped>

</style>
