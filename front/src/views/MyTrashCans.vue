<template>
  <div class="my-trashcans-container">
    <!-- 加载状态 -->
    <div v-if="loading" class="loading">
      <p>加载中...</p>
    </div>

    <!-- 列表为空 -->
    <div v-else-if="!loading && trashCans.length === 0" class="empty-state">
      <p>您还没有上传任何垃圾桶位置</p>
      <router-link to="/upload" class="btn btn-primary">去上传</router-link>
    </div>

    <!-- 列表展示 -->
    <div v-else class="trashcans-list">
      <div 
        v-for="item in trashCans" 
        :key="item.id" 
        class="trashcan-card"
      >
        <div class="card-image" v-if="item.image_url">
          <img :src="item.image_url" alt="垃圾桶图片" />
        </div>
        <div class="card-content">
          <h3 class="card-title">{{ item.address || '未设置地址' }}</h3>
          <p class="card-description" v-if="item.description">{{ item.description }}</p>
          <div class="card-meta">
            <span class="meta-item">坐标: {{ item.latitude.toFixed(6) }}, {{ item.longitude.toFixed(6) }}</span>
            <span class="meta-item">创建时间: {{ item.created_at }}</span>
          </div>
          <div class="card-actions">
            <button @click="openEditDialog(item)" class="btn btn-secondary btn-sm">编辑</button>
            <button @click="handleDelete(item.id)" class="btn btn-danger btn-sm">删除</button>
          </div>
        </div>
      </div>
    </div>

    <!-- 分页组件 -->
    <div v-if="totalPages > 1" class="pagination">
      <button 
        @click="goToPage(currentPage - 1)" 
        :disabled="currentPage === 1"
        class="page-btn"
      >
        上一页
      </button>
      <span class="page-info">
        第 {{ currentPage }} / {{ totalPages }} 页（共 {{ total }} 条）
      </span>
      <button 
        @click="goToPage(currentPage + 1)" 
        :disabled="currentPage === totalPages"
        class="page-btn"
      >
        下一页
      </button>
    </div>

    <!-- 编辑对话框 -->
    <div v-if="editDialogVisible" class="dialog-overlay" @click="closeEditDialog">
      <div class="dialog-content" @click.stop>
        <div class="dialog-header">
          <h3>编辑垃圾桶信息</h3>
          <button @click="closeEditDialog" class="dialog-close">×</button>
        </div>
        <form @submit.prevent="handleUpdate" class="edit-form">
          <div class="form-group">
            <label>地址：</label>
            <input 
              v-model="editForm.address" 
              type="text" 
              placeholder="请输入地址"
              class="input"
            />
          </div>
          <div class="form-group">
            <label>描述：</label>
            <textarea 
              v-model="editForm.description" 
              rows="3"
              placeholder="请输入描述"
              class="textarea"
            ></textarea>
          </div>
          <div class="form-group">
            <label>图片：</label>
            <div class="upload-area" @click="triggerFileInput">
              <input 
                ref="fileInput"
                type="file" 
                accept="image/*"
                @change="handleFileChange"
                style="display: none"
              />
              <div v-if="!imagePreview && currentEditItem?.image_url" class="current-image">
                <img :src="currentEditItem.image_url" alt="当前图片" />
                <p>点击更换图片</p>
              </div>
              <div v-else-if="imagePreview" class="image-preview">
                <img :src="imagePreview" alt="预览图片" />
                <button type="button" @click.stop="removeImage" class="remove-btn">×</button>
              </div>
              <div v-else class="upload-placeholder">
                <p>点击选择图片</p>
                <p class="hint-text">支持 JPG、PNG 等格式</p>
              </div>
            </div>
          </div>
          <div class="form-actions">
            <button type="submit" :disabled="updating" class="btn btn-primary">
              {{ updating ? '更新中...' : '保存' }}
            </button>
            <button type="button" @click="closeEditDialog" class="btn btn-secondary">
              取消
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { getUserTrashCans, updateTrashCan, deleteTrashCan } from '@/api/trashcan'

const loading = ref(false)
const trashCans = ref([])
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)
const totalPages = ref(0)

// 编辑对话框相关
const editDialogVisible = ref(false)
const currentEditItem = ref(null)
const updating = ref(false)
const fileInput = ref(null)
const imagePreview = ref(null)
const selectedFile = ref(null)

const editForm = reactive({
  address: '',
  description: ''
})

// 加载数据
const loadData = async () => {
  loading.value = true
  try {
    const response = await getUserTrashCans(currentPage.value, pageSize.value)
    if (response.code === 2000 && response.data) {
      trashCans.value = response.data.list || []
      total.value = response.data.total || 0
      totalPages.value = response.data.total_pages || 0
      currentPage.value = response.data.page || 1
    } else {
      alert(response.msg || '加载失败')
    }
  } catch (error) {
    console.error('加载数据失败:', error)
    alert('加载失败: ' + error.message)
  } finally {
    loading.value = false
  }
}

// 翻页
const goToPage = (page) => {
  if (page >= 1 && page <= totalPages.value) {
    currentPage.value = page
    loadData()
  }
}

// 打开编辑对话框
const openEditDialog = (item) => {
  currentEditItem.value = item
  editForm.address = item.address || ''
  editForm.description = item.description || ''
  imagePreview.value = null
  selectedFile.value = null
  editDialogVisible.value = true
}

// 关闭编辑对话框
const closeEditDialog = () => {
  editDialogVisible.value = false
  currentEditItem.value = null
  editForm.address = ''
  editForm.description = ''
  imagePreview.value = null
  selectedFile.value = null
}

// 触发文件选择
const triggerFileInput = () => {
  fileInput.value?.click()
}

// 处理文件选择
const handleFileChange = (event) => {
  const file = event.target.files[0]
  if (file) {
    selectedFile.value = file
    const reader = new FileReader()
    reader.onload = (e) => {
      imagePreview.value = e.target.result
    }
    reader.readAsDataURL(file)
  }
}

// 移除图片
const removeImage = () => {
  imagePreview.value = null
  selectedFile.value = null
  if (fileInput.value) {
    fileInput.value.value = ''
  }
}

// 更新垃圾桶信息
const handleUpdate = async () => {
  if (!currentEditItem.value) return

  updating.value = true
  try {
    const formData = new FormData()
    formData.append('address', editForm.address)
    formData.append('description', editForm.description)
    
    // 如果选择了新图片，添加到表单
    if (selectedFile.value) {
      formData.append('image', selectedFile.value)
    }

    const response = await updateTrashCan(currentEditItem.value.id, formData)
    if (response.code === 2000) {
      alert('更新成功')
      closeEditDialog()
      loadData() // 重新加载数据
    } else {
      alert(response.msg || '更新失败')
    }
  } catch (error) {
    console.error('更新失败:', error)
    alert('更新失败: ' + error.message)
  } finally {
    updating.value = false
  }
}

// 删除垃圾桶
const handleDelete = async (id) => {
  if (!confirm('确定要删除这个垃圾桶位置吗？此操作不可恢复。')) {
    return
  }

  try {
    const response = await deleteTrashCan(id)
    if (response.code === 2000) {
      alert('删除成功')
      loadData() // 重新加载数据
    } else {
      alert(response.msg || '删除失败')
    }
  } catch (error) {
    console.error('删除失败:', error)
    alert('删除失败: ' + error.message)
  }
}

onMounted(() => {
  loadData()
})
</script>

<style scoped>
.my-trashcans-container {
  width: 100%;
}

.loading, .empty-state {
  text-align: center;
  padding: 60px 20px;
  color: var(--text-secondary);
}

.empty-state p {
  margin-bottom: 20px;
  font-size: 16px;
}

.trashcans-list {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 20px;
  margin-bottom: 30px;
}

.trashcan-card {
  background: var(--bg-primary);
  border: 1px solid var(--border-color);
  border-radius: 8px;
  overflow: hidden;
  transition: var(--transition-base);
}

.trashcan-card:hover {
  box-shadow: var(--shadow-md);
  border-color: var(--border-color-hover);
}

.card-image {
  width: 100%;
  height: 200px;
  overflow: hidden;
  background: var(--bg-secondary);
}

.card-image img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.card-content {
  padding: 15px;
}

.card-title {
  margin: 0 0 10px 0;
  font-size: 18px;
  color: var(--text-primary);
  font-weight: 600;
}

.card-description {
  margin: 0 0 10px 0;
  color: var(--text-secondary);
  font-size: 14px;
  line-height: 1.5;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.card-meta {
  display: flex;
  flex-direction: column;
  gap: 5px;
  margin-bottom: 15px;
  font-size: 12px;
  color: var(--text-secondary);
}

.meta-item {
  display: block;
}

.card-actions {
  display: flex;
  gap: 10px;
}

.btn {
  padding: 8px 16px;
  border: none;
  border-radius: 4px;
  font-size: 14px;
  cursor: pointer;
  transition: var(--transition-base);
  font-weight: 500;
}

.btn-sm {
  padding: 6px 12px;
  font-size: 13px;
}

.btn-primary {
  background: var(--color-primary);
  color: white;
}

.btn-primary:hover:not(:disabled) {
  background: var(--color-primary-light);
}

.btn-secondary {
  background: var(--bg-secondary);
  color: var(--color-primary);
  border: 1px solid var(--color-primary);
}

.btn-secondary:hover {
  background: var(--bg-tertiary);
}

.btn-danger {
  background: #dc3545;
  color: white;
}

.btn-danger:hover {
  background: #c82333;
}

.btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

/* 分页组件 */
.pagination {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 20px;
  padding: 20px 0;
  border-top: 1px solid var(--border-color);
}

.page-btn {
  padding: 8px 16px;
  border: 1px solid var(--border-color);
  background: var(--bg-secondary);
  color: var(--text-primary);
  border-radius: 4px;
  cursor: pointer;
  transition: var(--transition-base);
}

.page-btn:hover:not(:disabled) {
  background: var(--bg-tertiary);
  border-color: var(--border-color-hover);
}

.page-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.page-info {
  color: var(--text-secondary);
  font-size: 14px;
}

/* 编辑对话框 */
.dialog-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
  padding: 20px;
}

.dialog-content {
  background: var(--bg-primary);
  border-radius: 8px;
  width: 100%;
  max-width: 600px;
  max-height: 90vh;
  overflow-y: auto;
  box-shadow: var(--shadow-lg);
}

.dialog-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px;
  border-bottom: 1px solid var(--border-color);
}

.dialog-header h3 {
  margin: 0;
  font-size: 20px;
  color: var(--text-primary);
}

.dialog-close {
  background: none;
  border: none;
  font-size: 28px;
  color: var(--text-secondary);
  cursor: pointer;
  padding: 0;
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 4px;
  transition: var(--transition-base);
}

.dialog-close:hover {
  background: var(--bg-tertiary);
  color: var(--text-primary);
}

.edit-form {
  padding: 20px;
}

.form-group {
  margin-bottom: 20px;
}

.form-group label {
  display: block;
  margin-bottom: 8px;
  color: var(--text-secondary);
  font-size: 14px;
  font-weight: 500;
}

.input, .textarea {
  width: 100%;
  padding: 12px;
  border: 1px solid var(--border-color);
  border-radius: 4px;
  font-size: 14px;
  box-sizing: border-box;
  background: var(--bg-secondary);
  color: var(--text-primary);
  transition: var(--transition-base);
  font-family: inherit;
}

.textarea {
  resize: vertical;
  min-height: 80px;
}

.input:focus, .textarea:focus {
  outline: none;
  border-color: var(--border-color-focus);
  box-shadow: 0 0 0 2px rgba(139, 111, 71, 0.1);
}

.upload-area {
  border: 2px dashed var(--border-color);
  border-radius: 4px;
  padding: 20px;
  text-align: center;
  cursor: pointer;
  transition: var(--transition-base);
  background: var(--bg-secondary);
}

.upload-area:hover {
  border-color: var(--color-primary);
  background: var(--bg-tertiary);
}

.current-image {
  position: relative;
}

.current-image img {
  max-width: 100%;
  max-height: 200px;
  border-radius: 4px;
  margin-bottom: 10px;
}

.image-preview {
  position: relative;
  display: inline-block;
}

.image-preview img {
  max-width: 100%;
  max-height: 200px;
  border-radius: 4px;
}

.remove-btn {
  position: absolute;
  top: -10px;
  right: -10px;
  width: 28px;
  height: 28px;
  border-radius: 50%;
  background: #dc3545;
  color: white;
  border: none;
  cursor: pointer;
  font-size: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
  line-height: 1;
}

.upload-placeholder p {
  margin: 5px 0;
  color: var(--text-secondary);
}

.hint-text {
  font-size: 12px;
}

.form-actions {
  display: flex;
  gap: 10px;
  justify-content: flex-end;
  margin-top: 30px;
}

/* 移动端适配 */
@media (max-width: 768px) {
  .trashcans-list {
    grid-template-columns: 1fr;
    gap: 15px;
  }

  .card-image {
    height: 150px;
  }

  .pagination {
    flex-direction: column;
    gap: 10px;
  }

  .page-info {
    order: -1;
  }

  .dialog-content {
    max-width: 100%;
    margin: 10px;
  }

  .form-actions {
    flex-direction: column;
  }

  .form-actions .btn {
    width: 100%;
  }
}
</style>

