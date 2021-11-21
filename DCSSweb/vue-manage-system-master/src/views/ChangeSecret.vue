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
        <el-form ref="secretRef" :rules="rules" :model="newcounter" label-width="160px">
            <el-form-item label="门限阈值" >
              {{degree}}
            </el-form-item>
            <el-form-item label="原委员会成员数" >
              {{oldcounter}}
            </el-form-item>
            <el-form-item label="新委员会成员数">
              <el-input-number v-model.number="newcounter" :min="1" :max="100"></el-input-number>
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
    return {
      newcounter : 1,
      secretid: this.$route.query.id,
      oldcounter:this.$route.query.oldcounter,
      degree:this.$route.query.degree,
    }
  },
  setup() {

    const newcounter = ref(1);
    const secretid = this.$route.query.id;
    const oldcounter=this.$route.query.oldcounter;
    const degree=this.$route.query.degree;
    const secretRef = ref(null);
    // 提交
    const onSubmit = () => {
      // 表单校验
      secretRef.value.validate((valid) => {
        if (valid) {
          if (degree*2+1>newcounter){
            ElMessage.error("参数不符合规范");
            return false;
          }else {
            ElMessage.success("提交成功！");
          }
        } else {
          ElMessage.error("缺少必要的输入项");
        }
      });
    };
    // 重置
    const onReset = () => {
      secretRef.value.resetFields();
    };
    //检验t,n代数关系
    const checkNum = (rule, value, callback) => {
      if (degree*2+1>newcounter){
        callback(new Error("不对"));
      }else {
        console.log("对不对");
        callback();
      }
    }


    const rules = {
      numberOfN: [
        {  validator: checkNum, message: "委员会成员数需要大于2×门限阈值+1", trigger: "change"}
      ],

    };
    return {
      rules,
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
