<script>
export default {
  data: function(){
    return {
      username:"",
      loading: false,
      userId: null,
      showToast: false,
      toastMessage: '',
      toastType: 'success'
    }
  },
  
  created() {
    const authToken = sessionStorage.getItem('authToken')
    if (authToken) {
      this.$router.push('/chats')
    } else {
      sessionStorage.clear()
    }
  },
  
  methods: {
    async login(){
      this.loading = true
      this.showToast = false
      try{
        console.log("Logging in with username: " + this.username);
        const response = await this.$axios.post("/session",{
          username: this.username
        })
        console.log("Response: ", response.data)
        this.userId = response.data.id
        this.saveToSessionStorage();
        console.log("Token saved correctly")
        this.showNotification("Login successful!", "success");
        setTimeout(() => {
          this.$router.push('/chats');
        }, 1000);
      }catch(e){
        console.error("Error while logging in:", e);
        console.error("Error message:", e.message);
        this.showNotification("Login failed: " + e.message, "error");
      }
      setTimeout(() => {
        this.loading = false;
      }, 1000);
    },
    saveToSessionStorage(){
      const bearerToken = `Bearer ${this.userId}`;
      sessionStorage.setItem('authToken', bearerToken);
      sessionStorage.setItem('username', this.username)
    },
    showNotification(message, type) {
      this.toastMessage = message;
      this.toastType = type;
      this.showToast = true;
      setTimeout(() => { 
        this.showToast = false; 
      }, 3000);
    }
  }
}
</script>

<template>
  <div class="container-fluid animated-bg vh-100 d-flex align-items-center justify-content-center">
    <div class="card shadow-lg" style="width: 450px;">
      <div class="card-body p-4">
        <h2 class="card-title text-center mb-4">Login</h2>
        <form @submit.prevent="login">
          <div class="mb-4">
            <label for="username" class="form-label fw-bold">Username:</label>
            <input
              type="text"
              class="form-control form-control-lg"
              id="username"
              v-model="username"
              required
              minlength="3"
              maxlength="16"
              placeholder="Enter username"
            />
          </div>
          <div class="d-grid">
            <button type="submit" class="btn btn-outline-dark btn-lg">
              <span v-if="!loading">Login</span>
              <span v-else class="spinner-border spinner-border-sm"></span>
            </button>
          </div>
        </form>
      </div>
    </div>

    <transition name="toast">
      <div v-if="showToast" class="toast-container position-fixed p-3">
        <div class="custom-toast" :class="toastType">
          <div class="toast-content">
            <i v-if="toastType === 'success'" class="fas fa-check-circle toast-icon"></i>
            <i v-else class="fas fa-exclamation-circle toast-icon"></i>
            <div class="toast-message">{{ toastMessage }}</div>
          </div>
          <div class="toast-progress" :class="toastType"></div>
        </div>
      </div>
    </transition>
  </div>
</template>

<style>
@import '../assets/login.css';
</style>