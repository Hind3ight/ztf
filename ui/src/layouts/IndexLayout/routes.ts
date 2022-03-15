import { RoutesDataItem } from "@/utils/routes";
import BlankLayout from '@/layouts/BlankLayout.vue';

const IndexLayoutRoutes: Array<RoutesDataItem> = [
  {
    icon: 'script',
    title: 'script',
    path: '/script',
    redirect: '/script/list',
    component: BlankLayout,
    children: [
      {
        title: 'script.list',
        path: 'list',
        component: () => import('@/views/script/index/main.vue'),
        hidden: true,
      },
    ],
  },

  {
    icon: 'execution',
    title: 'execution',
    path: '/exec',
    redirect: '/exec/history',
    component: BlankLayout,
    children: [
      {
        title: 'execution.history',
        path: 'history',
        component: () => import('@/views/exec/history/index.vue'),
        hidden: true,
      },
      {
        title: 'execution.result.func',
        path: 'history/func/:seq',
        component: () => import('@/views/exec/history/result-func.vue'),
        hidden: true,
      },
      {
        title: 'execution.result.unit',
        path: 'history/unit/:seq',
        component: () => import('@/views/exec/history/result-unit.vue'),
        hidden: true,
      },

      {
        title: 'execution',
        path: 'run',
        component: BlankLayout,
        hidden: true,
        children: [
          {
            title: 'execution.execCase',
            path: 'case/:seq/:scope',
            component: () => import('@/views/exec/exec/case.vue'),
            hidden: true,
          },
          {
            title: 'execution.execModule',
            path: 'module/:productId/:moduleId/:seq/:scope',
            component: () => import('@/views/exec/exec/module.vue'),
            hidden: true,
          },
          {
            title: 'execution.execSuite',
            path: 'suite/:productId/:suiteId/:seq/:scope',
            component: () => import('@/views/exec/exec/suite.vue'),
            hidden: true,
          },
          {
            title: 'execution.execTask',
            path: 'task/:productId/:taskId/:seq/:scope',
            component: () => import('@/views/exec/exec/task.vue'),
            hidden: true,
          },
          {
            title: 'execution.execUnit',
            path: 'unit',
            component: () => import('@/views/exec/exec/unit.vue'),
            hidden: true,
          },
        ]
      },
    ],
  },

  {
    icon: 'config',
    title: 'zentao_config',
    path: '/config',
    component: () => import('@/views/config/index.vue'),
  },

  {
    icon: 'sync',
    title: 'sync',
    path: '/sync',
    component: () => import('@/views/sync/index.vue'),
  },

  {
    icon: '',
    title: 'empty',
    path: '/site',
    redirect: '/site/list',
    component: BlankLayout,
    hidden: true,
    children: [
      {
        title: 'zentao_site',
        path: 'list',
        component: () => import('@/views/site/index.vue'),
        hidden: true,
      },
      {
        title: 'edit_site',
        path: 'edit/:id',
        component: () => import('@/views/site/edit.vue'),
        hidden: true,
      },
    ],
  },

];

export default IndexLayoutRoutes;