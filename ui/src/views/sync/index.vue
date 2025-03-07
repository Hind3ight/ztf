<template>
  <div v-if="currProject.type === 'unit'" class="panel">
    此为单元测试项目，不需要同步。
  </div>

  <div v-if="currProject.type === 'func'">
    <div class="main">

  <a-card title="从禅道同步用例信息">
    <a-form :label-col="labelCol" :wrapper-col="wrapperCol">
      <a-form-item label="产品" v-bind="validateInfos.productId">
        <a-select v-model:value="model.productId" @change="selectProduct">
          <a-select-option key="" value="">&nbsp;</a-select-option>
          <a-select-option v-for="item in products" :key="item.id" :value="item.id">{{item.name}}</a-select-option>
        </a-select>
      </a-form-item>
      <a-form-item label="模块" v-bind="validateInfos.moduleId">
        <a-select v-model:value="model.moduleId">
          <a-select-option key="" value="">&nbsp;</a-select-option>
          <a-select-option v-for="item in modules" :key="item.id" :value="item.id"><span v-html="item.name"></span></a-select-option>
        </a-select>
      </a-form-item>
      <a-form-item label="套件" v-bind="validateInfos.suiteId">
        <a-select v-model:value="model.suiteId">
          <a-select-option key="" value="">&nbsp;</a-select-option>
          <a-select-option v-for="item in suites" :key="item.id" :value="item.id">{{ item.name }}</a-select-option>
        </a-select>
      </a-form-item>
      <a-form-item label="任务" v-bind="validateInfos.taskId">
        <a-select v-model:value="model.taskId">
          <a-select-option key="" value="">&nbsp;</a-select-option>
          <a-select-option v-for="item in tasks" :key="item.id" :value="item.id">{{ item.name }}</a-select-option>
        </a-select>
      </a-form-item>
      <a-form-item label="语言" v-bind="validateInfos.lang">
        <a-select v-model:value="model.lang">
          <a-select-option v-for="item in langs" :key="item.code" :value="item.code">{{ item.name }}</a-select-option>
        </a-select>
      </a-form-item>
      <a-form-item label="期待结果为独立文件">
        <a-switch v-model:checked="model.independentFile" />
      </a-form-item>

      <a-form-item :wrapper-col="{ span: 14, offset: 4 }">
        <a-button type="primary" @click.prevent="syncFromZentaoSubmit">提交</a-button>
        <a-button style="margin-left: 10px" @click="resetFields">重置</a-button>
      </a-form-item>
    </a-form>
  </a-card>

  <a-card title="同步用例信息到禅道">
    <a-form :label-col="labelCol" :wrapper-col="wrapperCol">
      <a-form-item label="产品" v-bind="validateInfosCommit.productId">
        <a-select v-model:value="modelCommit.productId">
          <a-select-option key="" value="">&nbsp;</a-select-option>
          <a-select-option v-for="item in products" :key="item.id" :value="item.id">{{item.name}}</a-select-option>
        </a-select>
      </a-form-item>

      <a-form-item :wrapper-col="{ span: 14, offset: 4 }">
        <a-button type="primary" @click.prevent="syncToZentaoSubmit">提交</a-button>
      </a-form-item>
    </a-form>
  </a-card>

  </div>
  </div>
</template>
<script lang="ts">
import {defineComponent, ref, reactive, computed, watch, ComputedRef} from "vue";
import { useI18n } from "vue-i18n";

import { Props, validateInfos } from 'ant-design-vue/lib/form/useForm';
import {message, Form, notification} from 'ant-design-vue';
const useForm = Form.useForm;

import { SyncSettings } from './data.d';
import {useStore} from "vuex";
import {ProjectData} from "@/store/project";
import {ZentaoData} from "@/store/zentao";
import {syncFromZentao, syncToZentao} from "@/views/sync/service";

interface ConfigFormSetupData {
  currProject: ComputedRef;

  formRef: any
  model: SyncSettings
  rules: any

  labelCol: any
  wrapperCol: any
  validate: any
  validateInfos: validateInfos
  resetFields:  () => void;
  syncFromZentaoSubmit:  () => void;

  langs: ComputedRef<any[]>;
  products: ComputedRef<any[]>;
  modules: ComputedRef<any[]>;
  suites: ComputedRef<any[]>;
  tasks: ComputedRef<any[]>;
  selectProduct:  (item) => void;

  modelCommit: SyncSettings
  rulesCommit: any
  syncToZentaoSubmit:  () => void;
  validateCommit: any
  validateInfosCommit: validateInfos
  resetFieldsCommit:  () => void;
}

export default defineComponent({
  name: 'ConfigFormForm',
  components: {
  },
  setup(props): ConfigFormSetupData {
    const { t } = useI18n();

    const storeProject = useStore<{ project: ProjectData }>();
    const currConfig = computed<any>(() => storeProject.state.project.currConfig);
    const currProject = computed<any>(() => storeProject.state.project.currProject);

    const store = useStore<{zentao: ZentaoData}>();
    const langs = computed<any[]>(() => store.state.zentao.langs);
    const products = computed<any[]>(() => store.state.zentao.products);
    const modules = computed<any[]>(() => store.state.zentao.modules);
    const suites = computed<any[]>(() => store.state.zentao.suites);
    const tasks = computed<any[]>(() => store.state.zentao.tasks);

    store.dispatch('zentao/fetchLangs')
    store.dispatch('zentao/fetchProducts')
    watch(currConfig, (currConfig)=> {
      store.dispatch('zentao/fetchLangs')
      store.dispatch('zentao/fetchProducts')
    })

    const formRef = ref();

    const model = reactive<SyncSettings>({
      productId: '',
      lang: 'python',
      independentFile: false
    } as SyncSettings);

    const modelCommit = reactive<SyncSettings>({
      productId: '',
    } as SyncSettings);

    const rules = reactive({
      productId: [
        { required: true, message: '请选择产品'},
      ],
      lang: [
        { required: true, message: '请选择语言', trigger: 'change'}
      ],
    });
    const rulesCommit = reactive({
      productId: [
        { required: true, message: '请选择产品' },
      ]
    })

    const { resetFields, validate, validateInfos } = useForm(model, rules);

    const commitForm = useForm(modelCommit, rulesCommit);
    const resetFieldsCommit = commitForm.resetFields
    const validateCommit = commitForm.validate
    const validateInfosCommit = commitForm.validateInfos

    const selectProduct = (item) => {
      console.log('selectProduct', item)
      if (!item) return

      store.dispatch('zentao/fetchModules', item)
      store.dispatch('zentao/fetchSuites', item)
      store.dispatch('zentao/fetchTasks', item)
    };

    const syncFromZentaoSubmit = () => {
      console.log('syncFromZentaoSubmit')

      validate()
        .then(() => {
          syncFromZentao(model).then((json) => {
            console.log('json', json)
            if (json.code === 0) {
              notification.success({
                message: `同步成功`,
              });
            } else {
              notification.error({
                message: `同步失败`,
                description: json.msg,
              });
            }
          })
        })
        .catch(err => {console.log('validate fail', err)});
    };

    const syncToZentaoSubmit = () => {
      console.log('syncToZentaoSubmit')

      validateCommit()
        .then(() => {
          console.log('then', modelCommit);
          syncToZentao(modelCommit.productId).then((json) => {
            console.log('json', json)
            if (json.code === 0) {
              notification.success({
                message: `同步成功`,
              });
            } else {
              notification.error({
                message: `同步失败`,
                description: json.msg,
              });
            }
          })
        })
        .catch(err => {
          console.log('error', err);
        });
    };

    return {
      currProject,

      formRef,
      labelCol: { span: 6 },
      wrapperCol: { span: 12 },
      rules,
      validate,
      validateInfos,
      resetFields,
      syncFromZentaoSubmit,

      model,
      langs,
      products,
      modules,
      suites,
      tasks,
      selectProduct,

      modelCommit,
      rulesCommit,
      validateCommit,
      validateInfosCommit,
      resetFieldsCommit,
      syncToZentaoSubmit,
    }

  }
})
</script>

<style lang="less" scoped>
.panel {
  padding: 20px;
}

.main {
  padding: 0 20%;
}
</style>