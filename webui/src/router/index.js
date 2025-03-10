import {createRouter, createWebHashHistory} from 'vue-router'
import LoginView from '../views/LoginView.vue'
import ChatsView from '../views/ChatsView.vue'
import ProfileView from '../views/ProfileView.vue'

const router = createRouter({
	history: createWebHashHistory(import.meta.env.BASE_URL),
	routes: [
		{path: '/login', component: LoginView, meta: { requiresGuest: true }},
		{path: '/chats', component: ChatsView, meta: { requiresAuth: true }},
		{path: '/me', component: ProfileView, meta: { requiresAuth: true }},
	]
})

router.beforeEach((to, from, next) => {
	const authToken = sessionStorage.getItem('authToken')
	
	if (to.matched.some(record => record.meta.requiresAuth) && !authToken) {
	  next('/login')
	}

	else if (to.matched.some(record => record.meta.requiresGuest) && authToken) {
	  next('/chats')
	}
	else {
	  next()
	}
})
  
export default router