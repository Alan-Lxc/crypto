<template>
  <div>
    <div class="crumbs">
      <el-breadcrumb separator="/">
        <el-breadcrumb-item>
          <i class="el-icon-lx-text"></i> 秘密详情
        </el-breadcrumb-item>
      </el-breadcrumb>
    </div>
    <div class="container">
      <el-row>
        <!--        <el-col :span="6"><div class="grid-content bg-purple"><el-button type="">修改门限阈值</el-button></div></el-col>-->
        <el-col :span="8">
          <div class="grid-content bg-purple-light">
            <el-button @click="updatecounter()" type="">修改委员会成员数</el-button>
          </div>
        </el-col>
        <el-col :span="8">
          <div class="grid-content bg-purple">
            <el-button @click="handoffsecret()" type="">交接秘密</el-button>
          </div>
        </el-col>
        <el-col :span="8">
          <div class="grid-content bg-purple-light">
            <el-button @click="reconstructsecret"  type="">重构秘密值</el-button>
          </div>
        </el-col>


      </el-row>
    </div>
    <div class="container">
      <div class="form-box">
        <el-form ref="secretRef" label-width="160px">

          <el-form-item label="门限阈值" >
            {{secretinfo.degree}}
          </el-form-item>
          <el-form-item label="委员会成员数" prop="numberOfN">
            {{secretinfo.counter}}
          </el-form-item>
          <el-form-item label="秘密描述" prop="secret">
            {{secretinfo.description}}
          </el-form-item>
        </el-form>
      </div>
    </div>

    <el-table
        :data="nodelist"
        style="width: 100%"
        :row-class-name="tableRowClassName"
        @row-click="handleClick">
      <el-table-column
          prop="unit_id"
          label="节点ID"
          width="180">
      </el-table-column>
      <el-table-column
          prop="unit_ip"
          label="节点IP"
          width="300">
      </el-table-column>
      <el-table-column>

      </el-table-column>


    </el-table>

  </div>

</template>

<script>
import { ref, reactive } from "vue";
import { ElMessage, ElMessageBox } from "element-plus";
import { fetchData } from "../api/index";
import axios from "axios";
export default {
  name: "SecretInfo",
  data() {
    return {
      secretinfo: {
        degree: 0,
        counter: 0,
        user_id: 1,
        description: "",
      },
      nodelist: [],
    }
  },
  created() {
    this.getsecretinfoAndUnitList()
  },
  methods :{
    getsecretinfoAndUnitList(){
      let arr= this;
      let secretid = arr.$route.params.id;
      console.log(secretid);
      const url = "http://localhost:8080/api/secret/getsecret";
      // axios({
      //   methods: 'get',
      //   url: url ,
      //   params: {
      //     "secretid" : secretid,
      //   }
      // }).then(function (res) {
      //   console.log(res.data.data.secret);
      //   arr.secretinfo=res.data.data.secret;
      // })
      axios.get("http://localhost:8080/api/secret/getsecret",{
        params: {
          "secretid": secretid,
        },
      }).then(

      ).catch(err =>{

      });
      axios.get("http://localhost:8080/api/unit/getunitlist",{
        params:{
          "secretid": secretid,
        }
      }).then(function (res){
        console.log(res.data.data.nodelist);
        this.nodelist = res.data.data.unitlist;
      })
    },
    handleClick(row){
      this.$router.push({path:"/unitinfo",query:{userid:row[unit_id],secretid:row[unit_ip]}})
    },
    updatesecret(){
      let secretid = this.$route.params.id;
      console.log(secretid);
      axios.get("http://localhost:8080/api/secret/updatesecret",{
        params: {
          "id": secretid,
          "counter":90,
        }
      }).then(
          function (res) {
            console.log(res.data.data.secret);
          }
      ).catch(err =>{

      });

    },
    handoffsecret(){
      let secretid = this.$route.params.id;
      console.log(secretid);
      axios.get("http://localhost:8080/api/secret/updatesecret",{
        params: {
          "secretid": secretid,
        }
      }).then(
          function (res) {
            console.log(res.data.data.secret);
            alert("秘密值"+res.data.data.secret)
          }
      ).catch(err =>{});

    },
    reconstructsecret(){
      let secretid = this.$route.params.id;
      console.log(secretid);
      axios.get("http://localhost:8080/api/secret/reconstructsecret",{
        params: {
          "secretid": secretid
        }
      }).then(
          function (res) {
            console.log(res.data.data.secret);
            alert("秘密值"+res.data.data.secret)
          }
      ).catch(err =>{});
    }
  },

}
</script>

<style scoped>

</style>