function getBaseUrl() {
    var ishttps = 'https:' == document.location.protocol ? true : false;
    var url = window.location.host;
    if (ishttps) {
        url = 'https://' + url;
    } else {
        url = 'http://' + url;
    }
    return url;
}
function getWsBaseUrl() {
    var ishttps = 'https:' == document.location.protocol ? true : false;
    var url = window.location.host;
    if (ishttps) {
        url = 'wss://' + url;
    } else {
        url = 'ws://' + url;
    }
    return url;
}
//除去html标签
function replaceHtml(str){
    if (!str) {
        return "";
    }
    return str.replace(/<[^>]*>/g, '');
}
//浏览器桌面通知
function notify(title, options, callback, requestPermission) {

    // 先检查浏览器是否支持
    if (!window.Notification) {
        console.log("浏览器不支持notify");
        return;
    }
    options.body=replaceHtml(options.body);
    console.log("浏览器notify权限:", Notification.permission);
    // 检查用户曾经是否同意接受通知
    if (Notification.permission === 'granted') {
        var notification = new Notification(title, options); // 显示通知
        if (notification && callback) {
            notification.onclick = function(event) {
                callback(notification, event);
            }
            setTimeout(function () {
                notification.close();
            },3000);
        }
    } else if (requestPermission === true) {
        Notification.requestPermission().then(function(permission){
            console.log("请求浏览器notify权限:", permission);
            if (permission === 'granted') {
                var notification = new Notification(title, options); // 显示通知
                if (notification && callback) {
                    notification.onclick = function (event) {
                        callback(notification, event);
                    }
                    setTimeout(function () {
                        notification.close();
                    }, 3000);
                }
            } else if (permission === 'default') {
                console.log('用户关闭授权 可以再次请求授权');
            } else {
                console.log('用户拒绝授权 不能显示通知');
            }
        });
    }

}
var titleTimer=0;
var titleNum=0;
var originTitle = document.title;
function flashTitle() {
    if(titleTimer!=0){
        return;
    }
    titleTimer = setInterval(function(){
        titleNum++;
        if (titleNum == 3) {
            titleNum = 1;
        }
        if (titleNum == 1) {
            document.title = '【】' + originTitle;
        }
        if (titleNum == 2) {
            document.title = '【新訊息】' + originTitle;
        }
    }, 500);

}
function clearFlashTitle() {
    clearInterval(titleTimer);
    document.title = originTitle;
}
function replaceContent (content,baseUrl) {// 转义聊天内容中的特殊字符
    if(typeof baseUrl=="undefined"){
        baseUrl="";
    }
    content = (content || '')
        .replace(/img\[(.*?)\]/g, function (face) {  // 转义图片
            var src = face.replace(/^img\[/g, '').replace(/\]/g, '');;
            return '<img onclick="bigPic(src,true)" src="' +baseUrl+ src + '" style="max-width: 150px"/></div>';
        })
        .replace(/\n/g, '<br>'); // 转义换行
    return content;
}
function formatFileSize(fileSize) {
    if (fileSize < 1024) {
        return fileSize + 'B';
    } else if (fileSize < (1024*1024)) {
        var temp = fileSize / 1024;
        temp = temp.toFixed(2);
        return temp + 'KB';
    } else if (fileSize < (1024*1024*1024)) {
        var temp = fileSize / (1024*1024);
        temp = temp.toFixed(2);
        return temp + 'MB';
    } else {
        var temp = fileSize / (1024*1024*1024);
        temp = temp.toFixed(2);
        return temp + 'GB';
    }
}
function bigPic(src,isVisitor){
    if (isVisitor) {
        window.open(src);
        return;
    }
}
function filter (obj){
    var imgType = ["image/jpeg","image/png","image/jpg","image/gif"];
    var filetypes = imgType;
    var isnext = false;
    for (var i = 0; i < filetypes.length; i++) {
        if (filetypes[i] == obj.type) {
            return true;
        }
    }
    return false;
}
function sleep(time) {
    var startTime = new Date().getTime() + parseInt(time, 10);
    while(new Date().getTime() < startTime) {}
}
function checkLang(){
    var langs=["tw","cn","en"];
    var lang=getQuery("lang");
    if(lang!=""&&langs.indexOf(lang) >= 0 ){
        return lang;
    }
    return "tw";
}
function getQuery(key) {
    var query = window.location.search.substring(1);
    var key_values = query.split("&");
    var params = {};
    key_values.map(function (key_val){
        var key_val_arr = key_val.split("=");
        params[key_val_arr[0]] = key_val_arr[1];
    });
    if(typeof params[key]!="undefined"){
        return params[key];
    }
    return "";
}
function utf8ToB64(str) {
    return window.btoa(unescape(encodeURIComponent(str)));
}
function b64ToUtf8(str) {
    return decodeURIComponent(escape(window.atob(str)));
}
//播放声音
function alertSound(id,src){
    var b = document.getElementById(id);
    if(src!=""){
        b.src=src;
    }
    var p = b.play();
    p && p.then(function(){}).catch(function(e){});
}
//日期格式化
function formatDate(dateString, format = 'yyyy-MM-dd HH:mm:ss') {
    const date = new Date(dateString);

    const year = date.getFullYear();
    const month = String(date.getMonth() + 1).padStart(2, '0');
    const day = String(date.getDate()).padStart(2, '0');
    const hour = String(date.getHours()).padStart(2, '0');
    const minute = String(date.getMinutes()).padStart(2, '0');
    const second = String(date.getSeconds()).padStart(2, '0');

    const formattedDate = format
        .replace(/yyyy/g, year)
        .replace(/MM/g, month)
        .replace(/dd/g, day)
        .replace(/HH/g, hour)
        .replace(/mm/g, minute)
        .replace(/ss/g, second);

    return formattedDate;
}
function copyText(text) {
    var target = document.createElement('input')
    target.value = text
    document.body.appendChild(target)
    target.select()
    document.execCommand("copy");
    document.body.removeChild(target);
    return true;
}
;
