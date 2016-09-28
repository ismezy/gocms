/**
 * Created by 扬 on 2016/9/23.
 */
"use strict"
var app = require("../main")
app.config(["$routeProvider",function($routeProvider){
    $routeProvider.when('/menu',{template:require("./menu.html")})
}])
app.factory('Menus', ['$resource', function($resource) {
    return $resource('/admin/menu/:id', null, {
        update: { method:'PUT' }
    })
}])
app.controller('menuController',["$scope","$http","Menus",function($scope,$http,Menus) {
    $scope.menus = []
    Menus.query(function(menus){
        $scope.menus = menus
    })
    $scope.addMenu = function () {
        var menu = new Menus({
            Title: "新菜单" + $scope.menus.length,
            Path: "/",
            Code: "M" + $scope.menus.length,
            Index: $scope.menus.length
        })
        menu.$save().then(function(result){
            $scope.inserted = result
            $scope.menus.push(result)
        })
        return false;
    }
    $scope.saveMenu = function (data,index) {
        var old = $scope.menus[index]
        var nm = new Menus({
            Id: old.Id,
            Path: data.path,
            Title: data.title,
            Index: index,
            Code: data.code
        })
        return nm.$update()
    }
    $scope.removeMenu = function (index) {
        var old = $scope.menus[index]
        old.$delete({id:old.Id}).then(function(){
            $scope.menus.splice(index,1)
        })
    }
}])