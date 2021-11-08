<template>
  <div>
    <div class="crumbs">
      <el-breadcrumb separator="/">
        <el-breadcrumb-item>
          <i class="el-icon-lx-cascades"></i> 秘密列表
        </el-breadcrumb-item>
      </el-breadcrumb>
    </div>

    <div class="container">
      <el-header height="40px" >
          <router-link to="/newsecret">
            <el-button @click="" type="mini" bgcolor="bule" >新建秘密</el-button>
          </router-link>
      </el-header>
      <el-table
          :data="tableData"
          style="width: 100%"
          :row-class-name="tableRowClassName">
        <el-table-column
            prop="name"
            label="秘密名称"
            width="180">
        </el-table-column>
        <el-table-column
            prop="description"
            label="描述"
            width="500">
        </el-table-column>
        <el-table-column
            prop="numberOfT"
            label="委员会t值">
        </el-table-column>
        <el-table-column
            prop="numberOfN"
            label="委员会n值">
        </el-table-column>
        <el-table-column
            fixed="right"
            label="操作"
            width="100">
<!--          <template slot-scope="scope">-->
            <router-link to="/secretinfo">
              <el-button @click="handleClick(scope.row)" type="text" size="small">查看</el-button>
            </router-link>

            <el-button type="text" size="small">编辑</el-button>
<!--          </template>-->
        </el-table-column>
      </el-table>
    </div>
  </div>

</template>

<style>
.el-table .warning-row {
  background: oldlace;
}

.el-table .success-row {
  background: #f0f9eb;
}
</style>

<script>
export default {
  methods: {
    tableRowClassName({row, rowIndex}) {
      if (rowIndex === 1) {
        return 'warning-row';
      } else if (rowIndex === 3) {
        return 'success-row';
      }
      return '';
    }
  },
  data() {
    return {
      tableData: [{
        name: '秘密1',
        description: '第一个秘密',
        numberOfN: 10,
        numberOfT: 3,
      },{
        name: '秘密2',
        description: '第二个秘密',
        numberOfN: 100,
        numberOfT: 20,


      },{
        name: '秘密3',
        description: '第三个秘密',
        numberOfN: 12,
        numberOfT: 4,

      },{
        name: '秘密4',
        description: '第四个秘密',
        numberOfN: 21,
        numberOfT: 9,
      },{
        name: '456',
        description: '',
        numberOfN: 3,
        numberOfT: 1,
      },],

    }
  }
}
</script>