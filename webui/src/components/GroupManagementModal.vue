<script>
export default {
  props: {
    show: {
      type: Boolean,
      required: true
    },
    conversation: {
      type: Object,
      required: true
    }
  },
  
  data() {
    return {
      activeTab: 'members',
      newGroupName: '',
      searchQuery: '',
      searchResults: [],
      groupMembers: [],
      isLoading: false,
      isSubmitting: false,
      selectedPhoto: null,
      photoPreview: null,
      error: null,
      success: null,
      showToast: false,
      toastMessage: '',
      toastType: 'success'
    }
  },
  
  watch: {
    show(newVal) {
      if (newVal) {
        this.resetForm();
        this.newGroupName = this.conversation.name || '';
        this.fetchGroupMembers();
      }
    }
  },
  
  methods: {
    async fetchGroupMembers() {
      this.isLoading = true;
      this.groupMembers = [];
      
      try {
        const response = await this.$axios.get(`/conversations/${this.conversation.id}/members/`);
        this.groupMembers = response.data.map(member => ({
          id: member.id,
          username: member.username,
          profilePhoto: member.profilePhoto
        }));
      } catch (error) {
        console.error('Error fetching group members:', error);
        this.showNotification('Failed to load group members', 'error');
      } finally {
        this.isLoading = false;
      }
    },
    async searchUsers() {
      if (!this.searchQuery.trim()) {
        this.searchResults = [];
        return;
      }
      
      try {
        const response = await this.$axios.get(`/users/?username=${this.searchQuery}`);
        
        this.searchResults = response.data.filter(user => 
          !this.groupMembers.some(member => member.id === user.id)
        );
      } catch (error) {
        console.error('Error searching users:', error);
        this.showNotification('Failed to search users', 'error');
      }
    },
    
    async addUserToGroup(user) {
      this.isSubmitting = true;
      this.error = null;
      
      try {
        await this.$axios.post(`/conversations/${this.conversation.id}/members/`, {
          username: user.username
        });
        
    
        this.groupMembers.push(user);
        
      
        this.searchQuery = '';
        this.searchResults = [];
        
        this.showNotification(`${user.username} added to the group`, 'success');
        this.$emit('group-updated');
      } catch (error) {
        console.error('Error adding user to group:', error);
        
        if (error.response && error.response.status === 403) {
          this.error = 'To add a user to the group, you must have already started a conversation with them.';
        } else {
          this.error = 'Failed to add user to the group';
        }
        
        this.showNotification(this.error, 'error');
      } finally {
        this.isSubmitting = false;
      }
    },
    
    async removeUserFromGroup(user) {
      if (!confirm(`Are you sure you want to remove ${user.username} from the group?`)) {
        return;
      }
      
      this.isSubmitting = true;
      
      try {
        await this.$axios.delete(`/conversations/${this.conversation.id}/members/${user.id}`);
        

        this.groupMembers = this.groupMembers.filter(member => member.id !== user.id);
        
        this.showNotification(`${user.username} removed from the group`, 'success');
        this.$emit('group-updated');
      } catch (error) {
        console.error('Error removing user from group:', error);
        this.showNotification('Failed to remove user from the group', 'error');
      } finally {
        this.isSubmitting = false;
      }
    },
    
    async updateGroupName() {
      if (!this.newGroupName.trim() || this.newGroupName === this.conversation.name) {
        return;
      }
      
      if (this.newGroupName.length < 3 || this.newGroupName.length > 16) {
        this.error = 'Group name must be between 3 and 16 characters';
        this.showNotification(this.error, 'error');
        return;
      }
      
      this.isSubmitting = true;
      this.error = null;
      
      try {
        await this.$axios.put(`/conversations/${this.conversation.id}/name`, {
          groupName: this.newGroupName
        });
        
        this.showNotification('Group name updated successfully', 'success');
        this.$emit('group-updated');
      } catch (error) {
        console.error('Error updating group name:', error);
        this.error = 'Failed to update group name';
        this.showNotification(this.error, 'error');
      } finally {
        this.isSubmitting = false;
      }
    },
    
    handleFileSelect(event) {
      const file = event.target.files[0];
      if (!file) return;
      
      if (!file.type.match('image.*')) {
        this.error = 'Please select an image file';
        this.showNotification(this.error, 'error');
        return;
      }
      
      if (file.size > 5 * 1024 * 1024) {
        this.error = 'File size exceeds 5MB limit';
        this.showNotification(this.error, 'error');
        return;
      }
      
      this.selectedPhoto = file;
      this.createPhotoPreview(file);
    },
    
    createPhotoPreview(file) {
      const reader = new FileReader();
      reader.readAsDataURL(file);
      reader.onload = () => {
        this.photoPreview = reader.result;
      };
    },
    
    async updateGroupPhoto() {
      if (!this.selectedPhoto) {
        return;
      }
      
      this.isSubmitting = true;
      this.error = null;
      
      try {
        const base64String = await this.getBase64FromFile(this.selectedPhoto);
        
        await this.$axios.put(`/conversations/${this.conversation.id}/photo`, {
          groupPhoto: base64String
        });
        
        this.showNotification('Group photo updated successfully', 'success');
        this.resetPhotoSelection();
        this.$emit('group-updated');
      } catch (error) {
        console.error('Error updating group photo:', error);
        this.error = 'Failed to update group photo';
        this.showNotification(this.error, 'error');
      } finally {
        this.isSubmitting = false;
      }
    },
    
    async getBase64FromFile(file) {
      return new Promise((resolve, reject) => {
        const reader = new FileReader();
        reader.readAsDataURL(file);
        reader.onload = () => {
          const base64String = reader.result.split(',')[1];
          resolve(base64String);
        };
        reader.onerror = error => reject(error);
      });
    },
    
    resetPhotoSelection() {
      this.selectedPhoto = null;
      this.photoPreview = null;
      
      const fileInput = this.$refs.fileInput;
      if (fileInput) {
        fileInput.value = '';
      }
    },
    
    resetForm() {
      this.activeTab = 'members';
      this.newGroupName = '';
      this.searchQuery = '';
      this.searchResults = [];
      this.error = null;
      this.success = null;
      this.resetPhotoSelection();
    },
    
    closeModal() {
      this.$emit('close');
    },
    
    setActiveTab(tab) {
      this.activeTab = tab;
      this.error = null;
      this.success = null;
    },
    
    getInitials(name) {
      if (!name) return '';
      return name.charAt(0).toUpperCase();
    },
    
    getAvatarColor(name) {
      if (!name) return '#202c33';
      
      let hash = 0;
      for (let i = 0; i < name.length; i++) {
        hash = name.charCodeAt(i) + ((hash << 5) - hash);
      }
      
      const h = hash % 360;
      const s = 65 + (hash % 20);
      const l = 45 + (hash % 10);
      
      return `hsl(${h}, ${s}%, ${l}%)`;
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
  <div v-if="show" class="modal-overlay" @click="closeModal">
    <div class="modal-dialog group-management-modal" @click.stop>
      <div class="modal-content">
      
        <div class="modal-header">
          <div class="d-flex align-items-center">
            <div class="group-avatar me-3">
              <img 
                v-if="conversation && conversation.conversationPhoto" 
                :src="'data:image/jpeg;base64,' + conversation.conversationPhoto"
                alt="Group photo"
                class="group-avatar-img"
              />
              <div v-else class="group-avatar-placeholder">
                {{ getInitials(conversation ? conversation.name : '') }}
              </div>
            </div>
            <div>
              <h5 class="modal-title mb-1">Manage Group</h5>
              <div class="modal-subtitle" v-if="conversation">{{ conversation.name }}</div>
            </div>
          </div>
          <button @click="closeModal" class="close-button">
            <i class="fas fa-times"></i>
          </button>
        </div>
       
        <div class="modal-tabs">
          <button 
            class="tab-button" 
            :class="{ active: activeTab === 'members' }"
            @click="setActiveTab('members')"
          >
            <i class="fas fa-users"></i>
            <span>Members</span>
          </button>
          <button 
            class="tab-button" 
            :class="{ active: activeTab === 'photo' }"
            @click="setActiveTab('photo')"
          >
            <i class="fas fa-camera"></i>
            <span>Photo</span>
          </button>
          <button 
            class="tab-button" 
            :class="{ active: activeTab === 'name' }"
            @click="setActiveTab('name')"
          >
            <i class="fas fa-edit"></i>
            <span>Name</span>
          </button>
        </div>
        
       
        <div class="modal-body">
       
          <div v-if="error" class="error-message mb-3">
            <i class="fas fa-exclamation-circle me-2"></i>
            {{ error }}
          </div>
          
    
          <div v-if="success" class="success-message mb-3">
            <i class="fas fa-check-circle me-2"></i>
            {{ success }}
          </div>
          
      
          <div v-if="activeTab === 'members'" class="tab-content">
            <div class="input-group mb-3">
              <span class="input-group-text">
                <i class="fas fa-search"></i>
              </span>
              <input 
                v-model="searchQuery"
                type="text"
                class="form-control"
                placeholder="Search users to add..."
                @input="searchUsers"
                :disabled="isSubmitting"
              >
            </div>
            
           
            <div v-if="searchResults.length > 0" class="search-results-container mb-4">
              <h6 class="section-title">
                <i class="fas fa-user-plus me-2"></i>
                Add members
              </h6>
              <div 
                v-for="user in searchResults" 
                :key="'search-' + user.id"
                class="user-item"
                @click="addUserToGroup(user)"
              >
                <div class="d-flex align-items-center">
                  <div class="avatar-container" :style="{ backgroundColor: getAvatarColor(user.username) }">
                    <img 
                      v-if="user.profilePhoto" 
                      :src="'data:image/jpeg;base64,' + user.profilePhoto"
                      class="avatar-image"
                      alt="Profile photo"
                    >
                    <div v-else class="avatar-text">
                      {{ getInitials(user.username) }}
                    </div>
                  </div>
                  <div class="user-info ms-2 flex-grow-1">
                    <div class="username">{{ user.username }}</div>
                  </div>
                  <button class="action-button add-button">
                    <i class="fas fa-plus"></i>
                  </button>
                </div>
              </div>
            </div>
            
         
            <div class="current-members-container">
              <h6 class="section-title">
                <i class="fas fa-users me-2"></i>
                Group members ({{ groupMembers.length }})
              </h6>
              
              <div v-if="isLoading" class="text-center p-3">
                <div class="spinner"></div>
              </div>
              
              <div v-else-if="groupMembers.length === 0" class="empty-state">
                <i class="fas fa-users fa-2x mb-2"></i>
                <p>No members in this group yet</p>
              </div>
              
              <div v-else class="members-list">
                <div 
                  v-for="member in groupMembers" 
                  :key="'member-' + member.id"
                  class="user-item"
                >
                  <div class="d-flex align-items-center">
                    <div class="avatar-container" :style="{ backgroundColor: getAvatarColor(member.username) }">
                      <img 
                        v-if="member.profilePhoto" 
                        :src="'data:image/jpeg;base64,' + member.profilePhoto"
                        class="avatar-image"
                        alt="Profile photo"
                      >
                      <div v-else class="avatar-text">
                        {{ getInitials(member.username) }}
                      </div>
                    </div>
                    <div class="user-info ms-2 flex-grow-1">
                      <div class="username">{{ member.username }}</div>
                    </div>
                    <button 
                      v-if="member.username !== sessionStorage.getItem('username')"
                      class="action-button remove-button"
                      @click.stop="removeUserFromGroup(member)"
                      :disabled="isSubmitting"
                    >
                      <i class="fas fa-times"></i>
                    </button>
                  </div>
                </div>
              </div>
            </div>
          </div>
          
         
          <div v-if="activeTab === 'photo'" class="tab-content">
            <div class="current-photo-container mb-4">
              <h6 class="section-title mb-3">Current Group Photo</h6>
              <div class="current-photo">
                <img 
                  v-if="conversation && conversation.conversationPhoto" 
                  :src="'data:image/jpeg;base64,' + conversation.conversationPhoto"
                  alt="Current group photo"
                  class="photo-preview"
                />
                <div v-else class="photo-placeholder">
                  <i class="fas fa-users fa-3x"></i>
                </div>
              </div>
            </div>
            
            <div class="new-photo-container">
              <h6 class="section-title mb-3">Change Group Photo</h6>
              
              <div v-if="photoPreview" class="new-photo-preview mb-3">
                <img :src="photoPreview" alt="New photo preview" class="photo-preview" />
                <button 
                  class="remove-photo-button"
                  @click="resetPhotoSelection"
                  :disabled="isSubmitting"
                >
                  <i class="fas fa-times-circle"></i>
                </button>
              </div>
              
              <div class="d-flex flex-wrap gap-2">
                <button 
                  class="btn-select-photo" 
                  @click="$refs.fileInput.click()"
                  :disabled="isSubmitting"
                >
                  <i class="fas fa-camera me-2"></i>
                  Select Photo
                </button>
                
                <input 
                  ref="fileInput" 
                  type="file" 
                  accept="image/*" 
                  class="d-none" 
                  @change="handleFileSelect" 
                />
                
                <button 
                  v-if="photoPreview"
                  class="btn-update-photo" 
                  @click="updateGroupPhoto"
                  :disabled="isSubmitting"
                >
                  <i v-if="!isSubmitting" class="fas fa-save me-2"></i>
                  <span v-else class="spinner-border spinner-border-sm me-2"></span>
                  Update Photo
                </button>
              </div>
            </div>
          </div>
          
     
          <div v-if="activeTab === 'name'" class="tab-content">
            <h6 class="section-title mb-3">Change Group Name</h6>
            
            <div class="mb-3">
              <label for="groupName" class="form-label">Group Name</label>
              <input 
                id="groupName"
                v-model="newGroupName"
                type="text"
                class="form-control"
                placeholder="Enter group name"
                :minlength="3"
                :maxlength="16"
                :disabled="isSubmitting"
              >
              <small class="form-text text-info">
                Group name must be between 3 and 16 characters
              </small>
            </div>
            
            <button 
              class="btn-update-name" 
              @click="updateGroupName"
              :disabled="!newGroupName.trim() || newGroupName === conversation.name || isSubmitting || newGroupName.length < 3 || newGroupName.length > 16"
            >
              <i v-if="!isSubmitting" class="fas fa-save me-2"></i>
              <span v-else class="spinner-border spinner-border-sm me-2"></span>
              Update Name
            </button>
          </div>
        </div>
        
     
        <div class="modal-footer">
          <button @click="closeModal" class="btn-close-modal">
            Close
          </button>
        </div>
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

<style scoped>
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.7);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1050;
  backdrop-filter: blur(4px);
}

.group-management-modal {
  width: 500px;
  max-width: 92%;
  max-height: 90vh;
  display: flex;
  flex-direction: column;
  animation: modalFadeIn 0.3s ease;
}

@keyframes modalFadeIn {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.modal-content {
  border-radius: 16px;
  background-color: #111b21;
  border: 1px solid #2a3942;
  box-shadow: 0 10px 25px rgba(0, 0, 0, 0.5);
  overflow: hidden;
  display: flex;
  flex-direction: column;
  max-height: 90vh;
}

.modal-header {
  padding: 16px 20px;
  background-color: #202c33;
  border-bottom: 1px solid #2a3942;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.modal-title {
  font-size: 1.2rem;
  font-weight: 600;
  margin: 0;
  color: #e9edef;
}

.modal-subtitle {
  font-size: 0.9rem;
  color: #8696a0;
}

.close-button {
  background: transparent;
  border: none;
  color: #8696a0;
  font-size: 1.1rem;
  cursor: pointer;
  transition: color 0.2s;
  padding: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  border-radius: 50%;
}

.close-button:hover {
  color: #e9edef;
  background-color: #2a3942;
}

.group-avatar {
  width: 48px;
  height: 48px;
  border-radius: 50%;
  overflow: hidden;
  background-color: #2a3942;
  display: flex;
  align-items: center;
  justify-content: center;
}

.group-avatar-img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.group-avatar-placeholder {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.5rem;
  font-weight: 600;
  color: white;
  background-color: #00a884;
}

.modal-tabs {
  display: flex;
  background-color: #202c33;
  border-bottom: 1px solid #2a3942;
}

.tab-button {
  flex: 1;
  padding: 12px;
  background: transparent;
  border: none;
  color: #8696a0;
  font-size: 0.9rem;
  cursor: pointer;
  display: flex;
  flex-direction: column;
  align-items: center;
  transition: all 0.2s ease;
  border-bottom: 3px solid transparent;
}

.tab-button i {
  font-size: 1.2rem;
  margin-bottom: 5px;
}

.tab-button:hover {
  background-color: #2a3942;
  color: #e9edef;
}

.tab-button.active {
  color: #00a884;
  border-bottom-color: #00a884;
}

.modal-body {
  padding: 20px;
  overflow-y: auto;
  max-height: calc(90vh - 200px);
}

.tab-content {
  animation: fadeIn 0.3s ease;
}

@keyframes fadeIn {
  from {
    opacity: 0;
  }
  to {
    opacity: 1;
  }
}

.error-message {
  background-color: rgba(231, 76, 60, 0.15);
  color: #e74c3c;
  padding: 12px;
  border-radius: 8px;
  font-size: 0.9rem;
  border-left: 3px solid #e74c3c;
}

.success-message {
  background-color: rgba(46, 204, 113, 0.15);
  color: #2ecc71;
  padding: 12px;
  border-radius: 8px;
  font-size: 0.9rem;
  border-left: 3px solid #2ecc71;
}

.section-title {
  font-size: 0.95rem;
  font-weight: 600;
  color: #e9edef;
  margin-bottom: 10px;
  display: flex;
  align-items: center;
}

.form-control {
  background-color: #2a3942;
  border: 1px solid #2a3942;
  color: white;
}

.form-control:focus {
  background-color: #36434a;
  border-color: #00a884;
  box-shadow: 0 0 0 2px rgba(0, 168, 132, 0.25);
  color: white;
}

.form-control::placeholder {
  color: #8696a0;
}

.input-group-text {
  background-color: #2a3942;
  border: 1px solid #2a3942;
  color: #8696a0;
}

.form-text {
  color: #8696a0;
  font-size: 0.8rem;
}

.form-text.text-info {
  color: #3498db;
}

.search-results-container, .current-members-container {
  background-color: #202c33;
  border-radius: 10px;
  padding: 15px;
}

.user-item {
  padding: 10px;
  border-radius: 8px;
  margin-bottom: 8px;
  background-color: #2a3942;
  transition: background-color 0.2s ease;
}

.user-item:hover {
  background-color: #36434a;
}

.user-item:last-child {
  margin-bottom: 0;
}

.avatar-container {
  width: 40px;
  height: 40px;
  min-width: 40px;
  border-radius: 50%;
  overflow: hidden;
  display: flex;
  align-items: center;
  justify-content: center;
}

.avatar-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.avatar-text {
  color: white;
  font-size: 1.1rem;
  font-weight: 600;
  text-align: center;
}

.user-info {
  overflow: hidden;
}

.username {
  color: #e9edef;
  font-size: 0.95rem;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.action-button {
  min-width: 32px;
  height: 32px;
  border-radius: 50%;
  border: none;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s ease;
  background-color: transparent;
  color: #e9edef;
}

.add-button {
  color: #00a884;
}

.add-button:hover {
  background-color: rgba(0, 168, 132, 0.2);
  transform: scale(1.1);
}

.remove-button {
  color: #e74c3c;
}

.remove-button:hover {
  background-color: rgba(231, 76, 60, 0.2);
  transform: scale(1.1);
}

.action-button:disabled {
  opacity: 0.5;
  cursor: not-allowed;
  transform: none;
}

.empty-state {
  text-align: center;
  padding: 30px 0;
  color: #8696a0;
}

.spinner {
  width: 30px;
  height: 30px;
  border: 3px solid rgba(0, 168, 132, 0.3);
  border-radius: 50%;
  border-top-color: #00a884;
  animation: spin 1s linear infinite;
  margin: 0 auto;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.current-photo-container, .new-photo-container {
  background-color: #202c33;
  border-radius: 10px;
  padding: 15px;
  margin-bottom: 15px;
}

.current-photo, .new-photo-preview {
  display: flex;
  justify-content: center;
  position: relative;
}

.photo-preview {
  max-width: 200px;
  max-height: 200px;
  border-radius: 10px;
  object-fit: cover;
}

.photo-placeholder {
  width: 200px;
  height: 200px;
  background-color: #2a3942;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #8696a0;
}

.remove-photo-button {
  position: absolute;
  top: -10px;
  right: -10px;
  width: 30px;
  height: 30px;
  border-radius: 50%;
  background-color: #e74c3c;
  border: none;
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  box-shadow: 0 2px 5px rgba(0, 0, 0, 0.3);
  transition: transform 0.2s ease;
}

.remove-photo-button:hover {
  transform: scale(1.1);
}

.remove-photo-button:disabled {
  background-color: #7f8c8d;
  cursor: not-allowed;
  transform: none;
}

.btn-select-photo, .btn-update-photo, .btn-update-name {
  padding: 10px 16px;
  border-radius: 8px;
  font-size: 0.95rem;
  font-weight: 500;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s ease;
}

.btn-select-photo {
  background-color: transparent;
  border: 1px solid #2a3942;
  color: #e9edef;
}

.btn-select-photo:hover:not(:disabled) {
  background-color: #2a3942;
}

.btn-update-photo, .btn-update-name {
  background-color: #00a884;
  border: none;
  color: white;
}

.btn-update-photo:hover:not(:disabled), .btn-update-name:hover:not(:disabled) {
  background-color: #008f6f;
}

.btn-select-photo:disabled, .btn-update-photo:disabled, .btn-update-name:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.gap-2 {
  gap: 8px;
}

.modal-footer {
  padding: 16px 20px;
  background-color: #202c33;
  border-top: 1px solid #2a3942;
  display: flex;
  justify-content: flex-end;
}

.btn-close-modal {
  padding: 8px 20px;
  background-color: transparent;
  border: 1px solid #2a3942;
  border-radius: 8px;
  color: #e9edef;
  font-size: 0.95rem;
  cursor: pointer;
  transition: background-color 0.2s ease;
}

.btn-close-modal:hover {
  background-color: #2a3942;
}

.toast-container {
  bottom: 20px;
  right: 20px;
  z-index: 9999;
}

.custom-toast {
  min-width: 300px;
  background: white;
  padding: 15px 20px;
  border-radius: 12px;
  box-shadow: 0 5px 15px rgba(0, 0, 0, 0.15);
  position: relative;
  overflow: hidden;
}

.toast-content {
  display: flex;
  align-items: center;
}

.toast-icon {
  font-size: 24px;
  margin-right: 12px;
}

.toast-message {
  font-size: 14px;
  font-weight: 500;
}

.toast-progress {
  position: absolute;
  bottom: 0;
  left: 0;
  height: 3px;
  width: 100%;
  animation: progress 3s linear forwards;
}

.success {
  border-left: 4px solid #2ecc71;
}

.error {
  border-left: 4px solid #e74c3c;
}

.success .toast-icon {
  color: #2ecc71;
}

.error .toast-icon {
  color: #e74c3c;
}

.toast-progress.success {
  background: #2ecc71;
}

.toast-progress.error {
  background: #e74c3c;
}

@keyframes progress {
  from { width: 100%; }
  to { width: 0%; }
}

.toast-enter-active {
  animation: slideIn 0.3s ease-out;
}

.toast-leave-active {
  animation: slideOut 0.3s ease-in;
}

@keyframes slideIn {
  from {
    transform: translateX(100%);
    opacity: 0;
  }
  to {
    transform: translateX(0);
    opacity: 1;
  }
}

@keyframes slideOut {
  from {
    transform: translateX(0);
    opacity: 1;
  }
  to {
    transform: translateX(100%);
    opacity: 0;
  }
}

/* Scrollbar styling */
.modal-body::-webkit-scrollbar {
  width: 6px;
}

.modal-body::-webkit-scrollbar-track {
  background: transparent;
}

.modal-body::-webkit-scrollbar-thumb {
  background-color: #374045;
  border-radius: 6px;
}
</style>