import {createRouter, createWebHashHistory} from 'vue-router'
import LoginView from '../views/LoginView.vue'
import ChatsView from '../views/ChatsView.vue'

const router = createRouter({
	history: createWebHashHistory(import.meta.env.BASE_URL),
	routes: [
		{path: '/login', component: LoginView},
		{path: '/chats', component: ChatsView, meta: { requiresAuth: true }},
	]
})

router.beforeEach((to, from, next) => {
	const authToken = sessionStorage.getItem('authToken')
	
	if (to.matched.some(record => record.meta.requiresAuth) && !authToken) {
	  next('/login')
	} else {
	  next()
	}
  })
  
export default router
