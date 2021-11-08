<template>
  <div>
    <div class="crumbs">
      <el-breadcrumb separator="/">
        <el-breadcrumb-item>
          <i class="el-icon-lx-vipcard"></i>秘密详情
        </el-breadcrumb-item>
        <el-breadcrumb-item>变更秘密节点</el-breadcrumb-item>
      </el-breadcrumb>
    </div>
    <div class="container">
      <div class="form-box">
        <el-form ref="secretRef" :rules="rules" :model="secret" label-width="160px">
          <el-form-item label="秘密名称：" >
            <el-input :readonly="true" v-model="secret.name"></el-input>
          </el-form-item>
          <el-form-item label="秘密委员会t值：" prop="numberOfT">
            <el-input-number v-model.number="secret.numberOfT" :min="1" :max="100" prop="numberOfT"></el-input-number>
          </el-form-item>
          <el-form-item label="秘密委员会n值：" prop="numberOfN">
            <el-input-number v-model.number="secret.numberOfN" :min="1" :max="100" prop="numberOfN"></el-input-number>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="onSubmit">提交秘密</el-button>
            <el-button type="danger" @click="onReset">重置表单</el-button>
            <el-button type="danger" @click="cancel">取消</el-button>
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
  name: "ChangeSecret",
  data() {


  },
  setup() {

    const secretRef = ref(null);
    const secret = reactive({
      name: "abc",
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

    const validatorNumber = (rule, value, callback) => {
      if (!value) {
        return callback(new Error("请输入账户信息"));
      } else {
        if (1) {
          callback();
        } else {
          return callback(new Error('账号格式不正确'))
        }
      }
    };

    const rules = {
      name: [
        { required: true, message: "请输入秘密名称", trigger: "blur"},
      ],
      numberOfN: [
        {  validator: checkNum, message: "n需要大于2×t+1", trigger: "change"}
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
