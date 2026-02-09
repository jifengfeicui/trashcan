// 地图管理器
window.mapManager = {
    map: null,
    markers: [],
    infoWindows: [],
    userMarker: null,
    clickToSelectCallback: null,
    clickToSelectEnabled: false,

    init() {
        // 默认使用北京坐标，等待用户位置更新
        this.map = new AMap.Map('mapContainer', {
            zoom: 13,
            center: [116.4074, 39.9042], // 北京天安门
            viewMode: '3D'
        });

        // 地图点击事件
        this.map.on('click', (e) => {
            if (this.clickToSelectEnabled && this.clickToSelectCallback) {
                this.clickToSelectCallback(e.lnglat.lat, e.lnglat.lng);
                // 在地图上显示临时标记
                this.addTempMarker(e.lnglat.lng, e.lnglat.lat);
            }
        });

        // 添加地图控件
        this.map.addControl(new AMap.Scale());
        this.map.addControl(new AMap.ToolBar());
    },

    setUserLocation(lat, lng) {
        // 设置地图中心到用户位置
        this.map.setCenter([lng, lat]);
        this.map.setZoom(15);

        // 添加用户位置标记
        if (this.userMarker) {
            this.userMarker.setPosition([lng, lat]);
        } else {
            this.userMarker = new AMap.Marker({
                position: [lng, lat],
                title: '我的位置',
                icon: new AMap.Icon({
                    size: new AMap.Size(40, 40),
                    image: 'https://webapi.amap.com/theme/v1.3/markers/n/mid.png',
                    imageOffset: new AMap.Pixel(0, 0),
                    imageSize: new AMap.Size(40, 40)
                })
            });
            this.userMarker.setMap(this.map);
        }
    },

    showTrashCans(trashCans) {
        // 清除之前的标记
        this.clearMarkers();

        // 为每个垃圾桶添加标记
        trashCans.forEach(item => {
            this.addTrashCanMarker(item);
        });

        // 如果有点位，调整地图视野
        if (trashCans.length > 0) {
            const bounds = new AMap.Bounds();
            trashCans.forEach(item => {
                bounds.extend([item.longitude, item.latitude]);
            });
            // 包含用户位置
            if (this.userMarker) {
                const pos = this.userMarker.getPosition();
                bounds.extend([pos.lng, pos.lat]);
            }
            this.map.setBounds(bounds);
        }
    },

    addTrashCanMarker(trashCan) {
        // 创建自定义图标（垃圾桶图标）
        const icon = new AMap.Icon({
            size: new AMap.Size(32, 32),
            image: 'https://webapi.amap.com/theme/v1.3/markers/n/mid.png',
            imageOffset: new AMap.Pixel(0, 0),
            imageSize: new AMap.Size(32, 32)
        });

        const marker = new AMap.Marker({
            position: [trashCan.longitude, trashCan.latitude],
            title: trashCan.address || '垃圾桶',
            icon: icon,
            map: this.map
        });
        
        // 存储垃圾桶ID，用于后续查找
        marker.trashCanId = trashCan.id;

        // 创建信息窗口内容
        let infoContent = `
            <div class="info-window">
                <h4>${trashCan.address || '垃圾桶位置'}</h4>
                ${trashCan.description ? `<p>${trashCan.description}</p>` : ''}
                <p><strong>距离:</strong> ${trashCan.distance.toFixed(2)} 公里</p>
                ${trashCan.image_url ? `
                    <div class="info-image">
                        <img src="${trashCan.image_url}" alt="垃圾桶图片" 
                             onclick="window.mapManager.showImageModal('${trashCan.image_url}')"
                             style="max-width: 200px; cursor: pointer; border-radius: 4px;">
                    </div>
                ` : ''}
                <div class="info-actions">
                    <button onclick="window.mapManager.navigateTo(${trashCan.longitude}, ${trashCan.latitude})" 
                            class="nav-btn">导航到此处</button>
                </div>
            </div>
        `;

        const infoWindow = new AMap.InfoWindow({
            content: infoContent,
            offset: new AMap.Pixel(0, -30)
        });

        // 点击标记显示信息窗口
        marker.on('click', () => {
            // 关闭其他信息窗口
            this.infoWindows.forEach(iw => iw.close());
            infoWindow.open(this.map, marker.getPosition());
        });

        this.markers.push(marker);
        this.infoWindows.push(infoWindow);
    },

    addTempMarker(lng, lat) {
        // 添加临时标记（用于位置选择）
        const tempMarker = new AMap.Marker({
            position: [lng, lat],
            title: '选择的位置',
            icon: new AMap.Icon({
                size: new AMap.Size(24, 24),
                image: 'https://webapi.amap.com/theme/v1.3/markers/n/mid.png',
                imageOffset: new AMap.Pixel(0, 0),
                imageSize: new AMap.Size(24, 24)
            }),
            map: this.map
        });

        // 3秒后移除临时标记
        setTimeout(() => {
            tempMarker.setMap(null);
        }, 3000);
    },

    clearMarkers() {
        // 清除所有垃圾桶标记
        this.markers.forEach(marker => {
            marker.setMap(null);
        });
        this.infoWindows.forEach(iw => {
            iw.close();
        });
        this.markers = [];
        this.infoWindows = [];
    },

    focusMarker(id) {
        // 聚焦到指定ID的标记
        // 查找匹配的标记（通过存储的ID）
        let targetMarker = null;
        let targetInfoWindow = null;
        
        for (let i = 0; i < this.markers.length; i++) {
            const marker = this.markers[i];
            // 通过标记的title或其他属性匹配ID
            // 这里我们需要在创建标记时存储ID
            if (marker.trashCanId === id) {
                targetMarker = marker;
                targetInfoWindow = this.infoWindows[i];
                break;
            }
        }

        if (targetMarker) {
            this.map.setCenter(targetMarker.getPosition());
            this.map.setZoom(16);
            // 关闭其他信息窗口
            this.infoWindows.forEach(iw => iw.close());
            // 打开对应的信息窗口
            if (targetInfoWindow) {
                targetInfoWindow.open(this.map, targetMarker.getPosition());
            }
        }
    },

    enableClickToSelect(callback) {
        this.clickToSelectCallback = callback;
        this.clickToSelectEnabled = true;
        // 改变鼠标样式
        this.map.setDefaultCursor('crosshair');
    },

    disableClickToSelect() {
        this.clickToSelectEnabled = false;
        this.clickToSelectCallback = null;
        this.map.setDefaultCursor('default');
    },

    navigateTo(lng, lat) {
        // 使用高德地图导航
        const url = `https://uri.amap.com/navigation?to=${lng},${lat}&mode=car&policy=1&src=mypage&callnative=1`;
        window.open(url, '_blank');
    },

    showImageModal(imageUrl) {
        // 触发Vue组件的图片模态框
        // 由于Vue实例可能还未完全初始化，使用事件机制
        window.dispatchEvent(new CustomEvent('showImage', { detail: { imageUrl } }));
    }
};

