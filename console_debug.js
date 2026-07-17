(function() {
    'use strict';

    console.log('🎮 调试控制台已加载');

    // 创建控制台 UI
    function createConsoleUI() {
        const container = document.createElement('div');
        container.id = 'game-console-container';
        container.style.cssText = `
            position: fixed;
            bottom: 0;
            left: 0;
            right: 0;
            height: 250px;
            background: rgba(0, 0, 0, 0.92);
            border-top: 2px solid #0f0;
            z-index: 99999;
            display: flex;
            flex-direction: column;
            font-family: 'Courier New', monospace;
            transition: transform 0.3s ease;
            transform: translateY(100%);
        `;

        // 标题栏
        const header = document.createElement('div');
        header.style.cssText = `
            display: flex;
            justify-content: space-between;
            align-items: center;
            padding: 8px 15px;
            background: rgba(0, 0, 0, 0.5);
            border-bottom: 1px solid #333;
            flex-shrink: 0;
            cursor: pointer;
            user-select: none;
        `;

        const title = document.createElement('span');
        title.style.cssText = `
            color: #0f0;
            font-size: 14px;
            font-weight: bold;
        `;
        title.textContent = '🟢 控制台';

        const controls = document.createElement('div');

        const clearBtn = document.createElement('button');
        clearBtn.textContent = '清空';
        clearBtn.style.cssText = `
            margin-right: 8px;
            padding: 4px 12px;
            background: #333;
            color: #fff;
            border: 1px solid #555;
            border-radius: 3px;
            cursor: pointer;
            font-size: 12px;
        `;
        clearBtn.onclick = (e) => {
            e.stopPropagation();
            const content = document.getElementById('console-content');
            if (content) content.innerHTML = '';
        };

        const closeBtn = document.createElement('button');
        closeBtn.textContent = '隐藏';
        closeBtn.style.cssText = `
            padding: 4px 12px;
            background: #666;
            color: #fff;
            border: none;
            border-radius: 3px;
            cursor: pointer;
            font-size: 12px;
        `;
        closeBtn.onclick = (e) => {
            e.stopPropagation();
            toggleConsole(false);
        };

        controls.appendChild(clearBtn);
        controls.appendChild(closeBtn);
        header.appendChild(title);
        header.appendChild(controls);

        const content = document.createElement('div');
        content.id = 'console-content';
        content.style.cssText = `
            flex: 1;
            overflow-y: auto;
            padding: 8px 12px;
            color: #0f0;
            font-size: 12px;
            line-height: 1.6;
            white-space: pre-wrap;
            word-break: break-all;
        `;

        container.appendChild(header);
        container.appendChild(content);
        document.body.appendChild(container);

        header.addEventListener('click', (e) => {
            if (e.target.tagName === 'BUTTON') return;
            const isVisible = container.style.transform !== 'translateY(100%)';
            toggleConsole(!isVisible);
        });

        return container;
    }

    function toggleConsole(show) {
        const container = document.getElementById('game-console-container');
        if (!container) return;

        if (show === undefined) {
            show = container.style.transform === 'translateY(100%)';
        }

        container.style.transform = show ? 'translateY(0)' : 'translateY(100%)';
    }

    function addLog(type, args) {
        const content = document.getElementById('console-content');
        if (!content) return;

        const entry = document.createElement('div');
        const time = new Date().toLocaleTimeString();

        const colors = {
            'LOG': '#0f0',
            'ERROR': '#f44',
            'WARN': '#ff0',
            'INFO': '#4af',
            'DEBUG': '#a4f'
        };

        const color = colors[type] || '#fff';

        const messages = Array.from(args).map(arg => {
            if (arg === null) return 'null';
            if (arg === undefined) return 'undefined';
            if (typeof arg === 'object') {
                try {
                    return JSON.stringify(arg, null, 2);
                } catch (e) {
                    return String(arg);
                }
            }
            return String(arg);
        }).join(' ');

        entry.style.cssText = `
            color: ${color};
            padding: 2px 0;
            border-bottom: 1px solid rgba(255,255,255,0.05);
        `;

        entry.textContent = `[${time}] [${type}] ${messages}`;
        content.appendChild(entry);
        content.scrollTop = content.scrollHeight;

        while (content.children.length > 500) {
            content.removeChild(content.firstChild);
        }
    }

    function overrideConsole() {
        const original = {
            log: console.log,
            error: console.error,
            warn: console.warn,
            info: console.info,
            debug: console.debug
        };

        console.log = function(...args) {
            original.log.apply(console, args);
            addLog('LOG', args);
        };

        console.error = function(...args) {
            original.error.apply(console, args);
            addLog('ERROR', args);
        };

        console.warn = function(...args) {
            original.warn.apply(console, args);
            addLog('WARN', args);
        };

        console.info = function(...args) {
            original.info.apply(console, args);
            addLog('INFO', args);
        };

        console.debug = function(...args) {
            original.debug.apply(console, args);
            addLog('DEBUG', args);
        };
    }

    function catchGlobalErrors() {
        window.addEventListener('error', (e) => {
            console.error('🔥 全局错误:', e.message, 'at', e.filename, 'line', e.lineno);
            return false;
        });

        window.addEventListener('unhandledrejection', (e) => {
            console.error('⚠️ 未处理的 Promise 错误:', e.reason);
        });
    }

    // 初始化
    function init() {
        createConsoleUI();
        overrideConsole();
        catchGlobalErrors();

        // 默认显示
        setTimeout(() => toggleConsole(true), 300);

        // 快捷键 Ctrl+Shift+C
        document.addEventListener('keydown', (e) => {
            if (e.ctrlKey && e.shiftKey && (e.key === 'C' || e.key === 'c')) {
                e.preventDefault();
                toggleConsole();
            }
        });

        console.log('✅ 控制台已启用 | 快捷键: Ctrl+Shift+C');
    }

    if (document.readyState === 'loading') {
        document.addEventListener('DOMContentLoaded', init);
    } else {
        init();
    }
})();