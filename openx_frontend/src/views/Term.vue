<template>
  <div id="indexContainer" style="height: 100%; background: #002833">
    <div id="terminal" ref="terminalroot" />
  </div>
</template>
<script>
import "xterm/css/xterm.css";
import "xterm/lib/xterm.js";
import { Terminal } from "xterm";
import { FitAddon } from "xterm-addon-fit";
import { ref } from "vue";
import axios from "axios";
export default {
  name: "Shell",
  data() {
    return {
      order: "",
      term: "", // 保存terminal实例
      rows: 40,
      cols: 100,
      timer: null,
      socket: null,
    };
  },
  setup() {
    const terminalroot = ref(null);
    return {
      terminalroot,
    };
  },
  created() {
    this.wsShell();
    // this.initTerm()
  },

  mounted() {
    //  111111
    this.rows = Math.floor(window.innerHeight / 16 - 6);
    this.cols = Math.floor(window.innerWidth / 14);
    const _this = this;
    const term = new Terminal({
      rendererType: "canvas", // 渲染类型
      rows: parseInt(_this.rows), // 行数
      cols: parseInt(_this.cols), // 不指定行数，自动回车后光标从下一行开始
      convertEol: true, // 启用时，光标将设置为下一行的开头
      // scrollback: 50, //终端中的回滚量
      disableStdin: false, // 是否应禁用输入。
      cursorStyle: "underline", // 光标样式
      cursorBlink: true, // 光标闪烁
      theme: {
        foreground: "#7e9192", // 字体
        background: "#002833", // 背景色
        cursor: "help", // 设置光标
        lineHeight: 16,
      },
    });

    // 换行并输入起始符“$”
    term.prompt = () => {
      term.write("\r");
    };
    term.prompt();
    term.open(_this.terminalroot);
    document
      .getElementsByClassName("xterm-cursor-layer")[0]
      .setAttribute("contenteditable", true);
    var fitAddon = new FitAddon();
    // const attachAddon = new AttachAddon(this.socket)
    term.loadAddon(fitAddon);
    // term.loadAddon(attachAddon)
    fitAddon.fit();

    term.focus();
    // 创建terminal实例
    window.addEventListener("resize", resizeScreen);

    _this.term = term;

    // 内容全屏显示
    function resizeScreen() {
      // 不传size

      try {
        fitAddon.fit();

        // // 窗口大小改变时触发xterm的resize方法，向后端发送行列数，格式由后端决定
        // // 这里不使用size默认参数，因为改变窗口大小只会改变size中的列数而不能改变行数，所以这里不使用size.clos,而直接使用获取我们根据窗口大小计算出来的行列数
        // term.onResize(() => {
        //   _this.send({ type: 'resize', cols: term.cols, rows: term.rows, input: '' })
        // })
      } catch (e) {
        console.log("e", e.message);
      }
    }

    function runFakeTerminal(_this) {
      if (term._initialized) {
        return;
      }
      // 初始化
      term._initialized = true;

      term.prompt();

      // / **
      //     *添加事件监听器，用于按下键时的事件。事件值包含
      //     *将在data事件以及DOM事件中发送的字符串
      //     *触发了它。
      //     * @返回一个IDisposable停止监听。
      //  * /
      //   / ** 更新：xterm 4.x（新增）
      //  *为数据事件触发时添加事件侦听器。发生这种情况
      //  *用户输入或粘贴到终端时的示例。事件值
      //  *是`string`结果的结果，在典型的设置中，应该通过
      //  *到支持pty。
      //  * @返回一个IDisposable停止监听。
      //  * /
      // 支持输入与粘贴方法
      term.onData(function (key) {
        const order = {
          // type: "resize", cols: term.cols, rows: term.rows, input: ''
          input: key,
          type: "input",
          cols: term.cols,
          rows: term.rows,
        };
        _this.send(JSON.stringify(order));
        // 为解决窗体resize方法才会向后端发送列数和行数，所以页面加载时也要触发此方法
        _this.send(
          JSON.stringify({
            input: "",
            type: "resize",
            cols: term.cols,
            rows: term.rows,
          })
        );
      });
      // term.attachCustomKeyEventHandler(e => {
      //   if (e.ctrlKey && e.keyCode === 67) {
      //     e.preventDefault();
      //     document.execCommand('copy');
      //     return false;
      //   }
      //   return true;
      // })
    }
    runFakeTerminal(_this);
  },

  methods: {
    initTerm() {},
    async wsShell() {
      // const token = this.$router.history.current.query.token
      const token = localStorage.getItem("termToken");
      // localStorage.removeItem('termToken')
      let type = "";
      const route = this.$route;
      if (route.query.type === "bash") {
        type = `/ssh/namespace/${route.query.namespace}/pod/${route.query.podname}/shell/${route.query.containername}/shell/token/`;
      } else {
        type = "/log/";
      }
      const cig = await axios.get("config/config.json");
      let env = cig.data;
      const url = "wss://" + String(env.VITE_BASE_URL) + type + token;
      console.log("url===", url);
      this.init(url);
    },
    initXtermPing() {
      var _this = this;
      _this.timer = setInterval(function () {
        const order = { type: "ping", input: "", rows: 0, cols: 0 };
        _this.send(JSON.stringify(order));
      }, 1000);
    },

    init(url) {
      clearInterval(this.timer);
      // 实例化socket
      this.socket = new WebSocket(url);
      // 监听socket连接
      this.socket.onopen = this.open;
      // 监听socket错误信息
      this.socket.onerror = this.error;
      // 监听socket close
      this.socket.onclose = this.close;
      // 监听socket消息
      this.socket.onmessage = this.getMessage;
      // 发送socket消息
      this.socket.onsend = this.send;
    },
    open: function () {
      console.log("socket连接成功");
      this.initXtermPing();
    },
    error: function () {
      console.log("连接错误");
    },
    close: function () {
      if (this.socket) {
        this.socket.close();
      }
      // if (this.term) {
      //   this.term.dispose(document.getElementById('terminal'))
      // }
    },
    getMessage: function (event) {
      let reader = new FileReader();
      let _self = this;
      reader.onload = function (e) {
        _self.term.write(e.target.result);
      };
      reader.readAsText(event.data);
    },
    send: function (order) {
      // this.socket.send(order)
      if (this.socket.readyState === 1) {
        this.socket.send(order);
      }
    },
  },
};
</script>
