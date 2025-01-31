<script>
export default {
  data: function(){
    return {
      username:"",
      errorMsg: null,
      loading: false,
      userId: null,
      messsage:""
    }
  },
  methods: {
    async login(){
      this.loading = true
      if (this.username.length < 3){
        console.error("Username must be at least 3 characters long");
        return;
      }
      try{
        console.log("Logging in with username: " + this.username);
        const response = await this.$axios.post("/session",{
          username: this.username
        })
        console.log("Response: ", response.data)
        this.userId = response.data.id
        this.saveToSessionStorage();
        console.log("Token saved correctly")
      }catch(e){
        console.error("Error while logging in:", e);
        console.error("Error message:", e.message);
      }
      setTimeout(() => {
        this.loading = false;
      }, 1000);
    },

    saveToSessionStorage(){
      const bearerToken = `Bearer ${this.userId}`;
      sessionStorage.setItem('authToken', bearerToken);
      sessionStorage.setItem('username', this.username)
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
</div>
</template>

<style>
.animated-bg {
  background: linear-gradient(-45deg, #23a6d5, #23d5ab, #23a6d5, #23d5ab);
  background-size: 400% 400%;
  animation: gradient 15s ease infinite;
}

@keyframes gradient {
  0% {
    background-position: 0% 50%;
  }
  50% {
    background-position: 100% 50%;
  }
  100% {
    background-position: 0% 50%;
  }
}

.card {
  border-radius: 30px;
}

.btn.loading {
  position: relative;
  pointer-events: none;
  color: transparent !important;
}

.btn.loading::after {
  content: "";
  position: absolute;
  top: 50%;
  left: 50%;
  width: 16px;
  height: 16px;
  margin: -8px 0 0 -8px;
  border-radius: 50%;
  border: 2px solid #fff;
  border-top-color: #000;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  0% {
    transform: rotate(0deg);
  }
  100% {
    transform: rotate(360deg);
  }
}
</style>