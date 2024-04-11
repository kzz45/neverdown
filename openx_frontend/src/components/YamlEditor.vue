<template>
  <div>
    <div class="yaml-editor CodeMirror-wrap">
      <textarea id="yaml" ref="textarea" />
    </div>
  </div>
</template>

<script>
import { markRaw } from "vue";
import _CodeMirror from "codemirror";
import * as jsyaml from "js-yaml";
import "codemirror/addon/lint/lint.css";
import "codemirror/lib/codemirror.css";
import "codemirror/theme/monokai.css";
import "codemirror/mode/yaml/yaml";
import "codemirror/addon/lint/lint";
import "codemirror/addon/lint/yaml-lint";
import "codemirror/mode/php/php";

const CodeMirror = window.CodeMirror || _CodeMirror;
window.jsyaml = jsyaml; // 引入js-yaml为codemirror提高语法检查核心支持

export default {
  name: "YamlEditor",
  props: {
    config: {
      type: String,
      defult: "",
    },
  },
  data() {
    return {
      yamlEditor: false,
    };
  },
  watch: {
    config: {
      handler(val) {
        setTimeout(() => {
          this.yamlEditor.setValue(String(val));
        }, 0);
      },
      immediate: true,
    },
  },
  mounted() {
    this.initMode();
    setTimeout(() => {
      this.initData(String(this.config));
    }, 1);
  },
  methods: {
    initData(val) {
      this.setValue(val);
    },
    getValue() {
      return this.yamlEditor.getValue();
    },
    setValue(value) {
      this.yamlEditor.setValue(value);
      this.yamlEditor.refresh();
    },
    initMode() {
      (this.yamlEditor = markRaw(
        CodeMirror.fromTextArea(document.getElementById("yaml"), {
          mode: "text/x-yaml", // 语法model
          theme: "monokai",
          indentWithTabs: true,
          smartIndent: true,
          lineNumbers: true,
          matchBrackets: true,
          cursorHeight: 1,
          lineWrapping: true,
          autoRefresh: true,
          autofocus: true,
          hintOptions: {
            completeSingle: false,
          },
        })
      )),
        this.yamlEditor.on("change", (cm) => {
          this.$emit("changed", cm.getValue());
        });
    },
  },
};
</script>

<style>
.yaml-editor {
  height: 100%;
  position: relative;
}
.CodeMirror-gutter-wrapper {
  left: -29px !important;
}
.CodeMirror-line {
  z-index: 5 !important;
}
.CodeMirror-cursors {
  z-index: 6 !important;
}
.CodeMirror {
  height: calc(100vh - 360px);
  min-height: 300px;
  text-align: left;
}
.CodeMirror-scroll {
  min-height: 300px;
}
.cm-s-rubyblue span.cm-string {
  color: #f08047;
}
.CodeMirror-gutters {
  z-index: -1 !important;
}
</style>
