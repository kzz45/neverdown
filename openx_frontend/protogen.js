"use strict";

var esprima = require("esprima");
var estraverse = require("estraverse");
var escodegen = require("escodegen");
var fs = require("fs");

function _interopDefaultLegacy(e) {
  return e && typeof e === "object" && "default" in e ? e : { default: e };
}

var esprima__default = /*#__PURE__*/ _interopDefaultLegacy(esprima);
var estraverse__default = /*#__PURE__*/ _interopDefaultLegacy(estraverse);
var escodegen__default = /*#__PURE__*/ _interopDefaultLegacy(escodegen);
var fs__default = /*#__PURE__*/ _interopDefaultLegacy(fs);

fs__default["default"].readFile(
  "./proto/temp.js",
  { encoding: "utf-8" },
  (err, fr) => {
    if (err) {
      console.log(err);
    } else {
      const ast = esprima__default["default"].parseModule(fr, {
        raw: true,
        tokens: true,
        range: true,
        comment: true,
      });
      estraverse__default["default"].traverse(ast, {
        enter: function (node) {
          let addTag = Math.random().toString(36).substr(2);
          if (node.type === "SwitchCase") {
            const nodeConsequents = node.consequent;
            for (let consequent of nodeConsequents) {
              if (consequent.type === "VariableDeclaration") {
                for (let declaration of consequent.declarations) {
                  if (
                    declaration.type === "VariableDeclarator" &&
                    declaration.id.name === "end2"
                  ) {
                    declaration.id.name = "endx" + "_" + addTag;
                  }
                }
              }
              if (consequent.type === "WhileStatement") {
                const whileTest = consequent.test;
                if (whileTest.right.name === "end2") {
                  whileTest.right.name = "endx" + "_" + addTag;
                  addTag++;
                }
              }
            }
          }
        },
      });
      escodegen__default["default"].attachComments(
        ast,
        ast.comments,
        ast.tokens
      );
      const transformCode = escodegen__default["default"].generate(ast, {
        comment: true,
      });
      fs__default["default"].writeFile(
        "./proto/proto.js",
        transformCode,
        (err) => {
          if (err) throw err;
          fs__default["default"].unlink("./proto/temp.js", (err) => {
            if (err) throw err;
            console.log("删除temp.js");
          });
          console.log("proto.js写入成功");
        }
      );
    }
  }
);
