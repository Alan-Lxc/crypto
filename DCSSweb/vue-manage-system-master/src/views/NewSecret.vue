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
          <el-form-item label="秘密名称" prop="name">
            <el-input v-model="secret.name"></el-input>
          </el-form-item>
          <el-form-item label="门限阈值" prop="numberOfT">
            <el-input-number v-model.number="secret.numberOfT" :min="1" :max="100" prop="numberOfT"></el-input-number>
          </el-form-item>
          <el-form-item label="委员会成员数" prop="numberOfN">
            <el-input-number v-model.number="secret.numberOfN" :min="1" :max="100" prop="numberOfN"></el-input-number>
          </el-form-item>
          <el-form-item label="秘密文件">
            <div class="content-title">支持拖拽</div>
<!--            <div class="plugins-tips">-->
<!--              Element UI自带上传组件。-->
<!--              访问地址：-->
<!--              <a href="http://element.eleme.io/#/zh-CN/component/upload" target="_blank">Element UI Upload</a>-->
<!--            </div>-->
            <el-upload class="upload-demo" drag action="http://jsonplaceholder.typicode.com/api/posts/" multiple>
              <i class="el-icon-upload"></i>
              <div class="el-upload__text">
                将秘密文件拖到此处，或
                <em>点击上传</em>
              </div>
<!--              <template #tip>-->
<!--                <div class="el-upload__tip">只能上传 jpg/png 文件，且不超过 500kb</div>-->
<!--              </template>-->
            </el-upload>
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
export default {
  name: "NewSecret",
  data() {


  },
  setup() {

    const secretRef = ref(null);
    const secret = reactive({
      name: "",
      numberOfT: 0,
      numberOfN: 0,
      test: "",

    });

    // 提交
    const onSubmit = () => {
      // 表单校验
      console.log(secret.numberOfN,secret.numberOfT)

      secretRef.value.validate((valid) => {
        if (valid) {
          if (secret.numberOfT*2+1>secret.numberOfN){
            // console.log(form);
            console.log("t*2+1>n");
            ElMessage.error("t*2+1>n");
            return false;
          }else {
            ElMessage.success("提交成功！");
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
      if (secret.numberOfN<secret.numberOfT*2+1){
        callback(new Error("不对"));
      }else {
        console.log("对不对");
        callback();
      }
    }

    // const validatorNumber = (rule, value, callback) => {
    //   if (!value) {
    //     return callback(new Error("请输入账户信息"));
    //   } else {
    //     if (1) {
    //       callback();
    //     } else {
    //       return callback(new Error('账号格式不正确'))
    //     }
    //   }
    // };

    const rules = {
      name: [
        { required: true, message: "请输入秘密名称", trigger: "blur"},
      ],
      numberOfN: [
        {  validator: checkNum, message: "委员会成员数需要大于2倍的门限阈值", trigger: "change"}
      ],
      // numberOfT:[
      //   {validator:(rule,value,callback)=>{
      //       if (value<=0){
      //         callback(new Error("门限阈值必须为正数"));
      //       }
      //       callback();
      //     },trigger:'blur'}
      // ],
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
