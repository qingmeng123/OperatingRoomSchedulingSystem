import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: () => import('../views/Login.vue')
    },
    {
      path: '/login',
      name: 'login',
      component: () => import('../views/Login.vue')
    },
    {
      path:'/index',
      component:()=>import('../views/Main.vue'),
      children:[
        {
          path:'',
          name:'index',
          component: () => import('../components/Index.vue'),
          meta: {requiresAuth: false}
        },
        {
          path:'peopleManager',
          name:'peopleManager',
          component: () => import('../components/PeopleManager.vue'),
          meta: {requiresAuth: true, roles: ['1']}
        },
        {
          path:'postManager',
          name:'postManager',
          component: () => import('../components/PostManager.vue'),
          meta: {requiresAuth: true}
        },
        {
          path:'operatingRoomManager',
          name:'operatingRoomManager',
          component: () => import('../components/OperatingRoomManager.vue'),
          meta: {requiresAuth: true}
        },
        {
          path: 'schedulingManager',
          name: 'schedulingManager',
          component: () => import('../components/SchedulingManager.vue'),
          meta: {requiresAuth: true, roles: ['1']}
        },
        {
          path: 'dutyQuery',
          name: 'dutyQuery',
          component: () => import('../components/DutyQuery.vue')
        }
      ]
    }
  ]
})

export default router
