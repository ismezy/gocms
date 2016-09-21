/**
 * Created by 扬 on 2016/9/19.
 */
"use strict"

var angular = require("angular")
require("angular-route")
require("angular-animate")
require("angular-resource")
require("angular-xeditable")
require("ng-file-upload")
require("../../node_modules/angular-xeditable/dist/css/xeditable.css")
var app = angular.module("adminApp",["ngRoute","xeditable","ngResource","ngFileUpload","ngAnimate"])
app.run(function(editableOptions) {
    editableOptions.theme = 'bs3'
})
module.exports = app
