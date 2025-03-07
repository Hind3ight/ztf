import { Mutation, Action } from 'vuex';
import { StoreModuleType } from "@/utils/store";
import { ResponseData } from '@/utils/request';
import {queryProject} from "@/services/project";
import {setCache} from "@/utils/localCache";
import settings from '@/config/settings';
import {saveConfig} from "@/services/config";

export interface ProjectData {
  projects: any[]
  currProject: any
  currConfig: any
  scriptTree: any[]
}

export interface ModuleType extends StoreModuleType<ProjectData> {
  state: ProjectData;
  mutations: {
    saveProjects: Mutation<ProjectData>;
  };
  actions: {
    fetchProject: Action<ProjectData, ProjectData>;
    saveConfig: Action<ProjectData, ProjectData>;
  };
}

const initState: ProjectData = {
  projects: [],
  currProject: {},
  currConfig: {},
  scriptTree: [],
}

const StoreModel: ModuleType = {
  namespaced: true,
  name: 'project',
  state: {
    ...initState
  },
  mutations: {
    saveProjects(state, payload) {
      console.log('payload', payload)

      setCache(settings.currProject, payload.currProject.path);

      state.projects = [...payload.projects];
      state.currProject = payload.currProject;
      state.currConfig = payload.currConfig;
      state.scriptTree = [payload.scriptTree];
    },
  },
  actions: {
    async fetchProject({ commit }, currProjectPath) {
      try {
        const response: ResponseData = await queryProject(currProjectPath);
        const { data } = response;
        commit('saveProjects', data || {});

        return true;
      } catch (error) {
        return false;
      }
    },
    async saveConfig({ commit }, config) {
      const resp: ResponseData = await saveConfig(config);
      console.log('&&&')
      const { data } = resp;
      commit('saveProjects', data);

      return resp;
    },
  }
}

export default StoreModel;
  