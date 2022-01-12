import { Mutation, Action } from 'vuex';
import { StoreModuleType } from "@/utils/store";
import { ResponseData } from '@/utils/request';
import { queryProduct, queryModule, querySuite, queryTask } from "./service";

export interface ZentaoData {
    products: any[]
    modules: any[]
    suites: any[]
    tasks: any[]
}

export interface ModuleType extends StoreModuleType<ZentaoData> {
    state: ZentaoData;
    mutations: {
        saveProducts: Mutation<any>;
        saveModules: Mutation<any>;
        saveSuites: Mutation<any>;
        saveTasks: Mutation<any>;
    };
    actions: {
        fetchProducts: Action<ZentaoData, ZentaoData>;
        fetchModules: Action<ZentaoData, ZentaoData>;
        fetchSuites: Action<ZentaoData, ZentaoData>;
        fetchTasks: Action<ZentaoData, ZentaoData>;
    };
}

const initState: ZentaoData = {
    products: [],
    modules: [],
    suites: [],
    tasks: [],
}

const StoreModel: ModuleType = {
    namespaced: true,
    name: 'zentao',
    state: {
        ...initState
    },
    mutations: {
        saveProducts(state, payload) {
            console.log('payload', payload)
            state.products = payload
            state.modules = []
            state.suites = []
            state.tasks = []
        },
        saveModules(state, payload) {
            console.log('payload', payload)
            state.modules = payload
        },
        saveSuites(state, payload) {
            console.log('payload', payload)
            state.suites = payload
        },
        saveTasks(state, payload) {
            console.log('payload', payload)
            state.tasks = payload
        },
    },
    actions: {
        async fetchProducts({ commit }) {
            try {
                const response: ResponseData = await queryProduct();
                const { data } = response;
                commit('saveProducts', data)

                return true;
            } catch (error) {
                return false;
            }
        },
        async fetchModules({ commit }, productId) {
            try {
                const response: ResponseData = await queryModule(productId);
                const { data } = response;
                commit('saveModules', data || 0);

                return true;
            } catch (error) {
                return false;
            }
        },
        async fetchSuites({ commit }, productId) {
            try {
                const response: ResponseData = await querySuite(productId);
                const { data } = response;
                commit('saveSuites', data || 0);

                return true;
            } catch (error) {
                return false;
            }
        },
        async fetchTasks({ commit }, productId) {
            try {
                const response: ResponseData = await queryTask(productId);
                const { data } = response;
                commit('saveTasks', data || 0);

                return true;
            } catch (error) {
                return false;
            }
        },
    }
}

export default StoreModel;
