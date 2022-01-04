import { RoutesDataItem } from "@/utils/routes";
import BlankLayout from '@/layouts/BlankLayout.vue';

const IndexLayoutRoutes: Array<RoutesDataItem> = [
  {
    icon: 'script',
    title: 'index-layout.menu.script',
    path: '/script',
    redirect: '/script/list',
    component: BlankLayout,
    children: [
      {
        title: 'index-layout.menu.script.list',
        path: 'list',
        component: () => import('@/views/script/list/index.vue'),
        hidden: true,
      },
      {
        title: 'index-layout.menu.script.view',
        path: 'view/:id',
        component: () => import('@/views/script/view/index.vue'),
        hidden: true,
      },
    ],
  },

  {
    icon: 'execution',
    title: 'index-layout.menu.execution',
    path: '/execution',
    redirect: '/execution/history',
    component: BlankLayout,
    children: [
      {
        title: 'index-layout.menu.execution.history',
        path: 'history',
        component: () => import('@/views/execution/history/index.vue'),
        hidden: true,
      },
      {
        title: 'index-layout.menu.execution.result',
        path: 'result/:id',
        component: () => import('@/views/execution/result/index.vue'),
        hidden: true,
      },
    ],
  },

];

export default IndexLayoutRoutes;