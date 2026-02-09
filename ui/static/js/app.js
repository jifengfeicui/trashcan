// Vue应用主入口
const { createApp } = Vue;

const app = createApp({
    data() {
        return {
            showUploadForm: false,
            uploadMethod: 'click', // 'click' 或 'manual'
            formData: {
                latitude: null,
                longitude: null,
                address: '',
                description: '',
                image: null
            },
            imagePreview: null,
            uploading: false,
            nearbyTrashCans: [],
            showImageModal: false,
            modalImageUrl: '',
            userLocation: null
        };
    },
    mounted() {
        // 监听图片显示事件
        window.addEventListener('showImage', (event) => {
            this.showImage(event.detail.imageUrl);
        });
        // 初始化地图
        this.initMap();
        // 获取用户位置
        this.getUserLocation();
    },
    methods: {
        initMap() {
            // 地图初始化在map.js中
            if (window.mapManager) {
                window.mapManager.init();
            }
        },
        getUserLocation() {
            if (navigator.geolocation) {
                navigator.geolocation.getCurrentPosition(
                    (position) => {
                        this.userLocation = {
                            lat: position.coords.latitude,
                            lng: position.coords.longitude
                        };
                        // 通知地图管理器更新用户位置
                        if (window.mapManager) {
                            window.mapManager.setUserLocation(this.userLocation.lat, this.userLocation.lng);
                        }
                    },
                    (error) => {
                        console.error('获取位置失败:', error);
                        // 默认使用北京坐标
                        this.userLocation = {
                            lat: 39.9042,
                            lng: 116.4074
                        };
                        if (window.mapManager) {
                            window.mapManager.setUserLocation(this.userLocation.lat, this.userLocation.lng);
                        }
                    }
                );
            } else {
                console.error('浏览器不支持地理定位');
                // 默认使用北京坐标
                this.userLocation = {
                    lat: 39.9042,
                    lng: 116.4074
                };
                if (window.mapManager) {
                    window.mapManager.setUserLocation(this.userLocation.lat, this.userLocation.lng);
                }
            }
        },
        searchNearby() {
            if (!this.userLocation) {
                alert('请先获取您的位置');
                return;
            }

            const radius = 5; // 搜索半径5公里
            const limit = 10; // 最多返回10个

            fetch(`/api/trashcans/nearby?lat=${this.userLocation.lat}&lng=${this.userLocation.lng}&radius=${radius}&limit=${limit}`)
                .then(response => response.json())
                .then(data => {
                    if (data.code === 2000) {
                        this.nearbyTrashCans = data.data || [];
                        // 在地图上显示标记
                        if (window.mapManager) {
                            window.mapManager.showTrashCans(this.nearbyTrashCans);
                        }
                    } else {
                        alert('搜索失败: ' + data.msg);
                    }
                })
                .catch(error => {
                    console.error('搜索错误:', error);
                    alert('搜索失败，请稍后重试');
                });
        },
        focusMarker(id) {
            // 聚焦到指定的标记
            if (window.mapManager) {
                window.mapManager.focusMarker(id);
            }
        },
        onUploadMethodChange() {
            if (this.uploadMethod === 'click') {
                // 启用地图点击选择
                if (window.mapManager) {
                    window.mapManager.enableClickToSelect((lat, lng) => {
                        this.formData.latitude = lat;
                        this.formData.longitude = lng;
                    });
                }
            } else {
                // 禁用地图点击选择
                if (window.mapManager) {
                    window.mapManager.disableClickToSelect();
                }
            }
        },
        handleFileSelect(event) {
            const file = event.target.files[0];
            if (file) {
                this.formData.image = file;
                // 预览图片
                const reader = new FileReader();
                reader.onload = (e) => {
                    this.imagePreview = e.target.result;
                };
                reader.readAsDataURL(file);
            }
        },
        resetForm() {
            this.formData = {
                latitude: null,
                longitude: null,
                address: '',
                description: '',
                image: null
            };
            this.imagePreview = null;
            if (window.mapManager) {
                window.mapManager.disableClickToSelect();
            }
        },
        submitTrashCan() {
            if (!this.formData.latitude || !this.formData.longitude) {
                alert('请选择或输入位置坐标');
                return;
            }

            if (!this.formData.image) {
                alert('请选择图片');
                return;
            }

            this.uploading = true;

            const formData = new FormData();
            formData.append('latitude', this.formData.latitude);
            formData.append('longitude', this.formData.longitude);
            formData.append('address', this.formData.address);
            formData.append('description', this.formData.description);
            formData.append('image', this.formData.image);

            fetch('/api/trashcans', {
                method: 'POST',
                body: formData
            })
                .then(response => response.json())
                .then(data => {
                    this.uploading = false;
                    if (data.code === 2000) {
                        alert('上传成功！');
                        this.resetForm();
                        this.showUploadForm = false;
                        // 刷新附近垃圾桶列表
                        this.searchNearby();
                    } else {
                        alert('上传失败: ' + data.msg);
                    }
                })
                .catch(error => {
                    this.uploading = false;
                    console.error('上传错误:', error);
                    alert('上传失败，请稍后重试');
                });
        },
        showImage(imageUrl) {
            this.modalImageUrl = imageUrl;
            this.showImageModal = true;
        }
    }
});

// 挂载应用
const vm = app.mount('#app');

// 暴露给全局，供map.js调用
window.app = vm;

