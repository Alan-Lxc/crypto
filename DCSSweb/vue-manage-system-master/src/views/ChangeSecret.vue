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
            <el-form-item label="门限阈值" >
              {{secret.degree}}
            </el-form-item>
            <el-form-item label="原委员会成员数" >
              {{secret.oldcounter}}
            </el-form-item>
            <el-form-item label="新委员会成员数" prop="newcounter">
              <el-input-number v-model.number="secret.newcounter" :min="1" :max="100"></el-input-number>
            </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="onSubmit">确认提交</el-button>
            <el-button type="danger" @click="onReset">表单</el-button>
          </el-form-item>

        </el-form>

      </div>

    </div>
  </div>
</template>

<script>
import {ref, reactive} from "vue";
import { ElMessage } from "element-plus";
import {useRoute} from "vue-router";
import axios from "axios";

export default {
  name: "ChangeSecret",
  data() {

  },


  setup() {
    const route = useRoute();
    const secret = reactive({
      secretid: route.query.id,
      oldcounter : route.query.oldcounter,
      degree : route.query.degree,
      newcounter:1,
    })

    const changesecretRef = ref(null);
    // 提交
    const onSubmit = () => {
      // 表单校验
      changesecretRef.value.validate((valid) => {
        if (valid) {
          if (secret.degree*2+1>secret.newcounter){
            ElMessage.error("参数不符合规范");
            return false;
          }else {
            axios.get("http://localhost:8080/api/secret/updatesecretcounter",{
              params: {
                "newcounter": secret.newcounter,
                "secretid": secret.secretid,
              }
            }).then(function (res){
              console.log(res)
              if (res.status === 200){
                ElMessage.success("提交成功！");
              }else {

                ElMessage.error("failed");
              }
            })
          }
        } else {
          ElMessage.error("缺少必要的输入项");
        }
      });
    };
    // 重置
    const onReset = () => {
      changesecretRef.value.resetFields();
    };
    //检验t,n代数关系
    const checkNum = (rule, value, callback) => {
      if (secret.degree*2+1>secret.newcounter){
        callback(new Error("不对"));
      }else {
        console.log("对不对");
        callback();
      }
    }
    const rules = {
      newcounter: [
        {  validator: checkNum, message: "委员会成员数需要大于2×门限阈值+1", trigger: "change"}
      ],

    };
    return {
      rules,
      secretRef: changesecretRef,
      secret,
      onSubmit,
      onReset,
      checkNum,

    };
  },

};
</script>

<style scoped>

</style>
