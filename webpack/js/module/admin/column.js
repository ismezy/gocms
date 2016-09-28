/**
 * Created by 扬 on 2016/9/21.
 */
"use strict"

var app = require("../main")
app.config(["$routeProvider",function($routeProvider){
    $routeProvider.when('/column',{template:require("./column.html")})
}])
app.factory('Column', ['$resource', function($resource) {
    return $resource('/admin/column/:Index', null, {
        update: { method:'PUT' }
    })
}])

app.controller('columnController',["$scope","$http","Column",function($scope,$http,Column){

}])
