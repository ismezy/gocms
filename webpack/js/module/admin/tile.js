/**
 * Created by 扬 on 2016/9/19.
 */
"use strict"
var app = require("../main")
app.config(["$routeProvider",function($routeProvider){
    $routeProvider.when('/tile',{template:require("./tile.html")})
}])
app.factory('Tiles', ['$resource', function($resource) {
    return $resource('/admin/tile/item/:index', null, {
        update: { method:'PUT' }
    })
}])
app.controller('tileController',["$scope","$http","Tiles","Upload",function($scope,$http,Tiles,Upload){
    $scope.tile = {
        Items:[]
    }
    $scope.addEnable = true
    $scope.add = function () {
        $scope.addEnable = false
        var tile = new Tiles({
            Title       :     "磁贴" + $scope.tile.Items.length,
            Memo        :     "描述",
            ImgPath     :       "",
            Url         :     "/"
        })
        tile.$save().then(function(result){
            $scope.tile.Items.push(tile)
            $scope.addEnable = true
        },function (error) {
            $scope.addEnable = true
        })
    }
    $scope.upload = function (file,index) {
        Upload.upload({
            url:'/admin/image/upload',
            file:file
        }).then(function (ret) {
            $scope.tile.Items[index].ImgPath = ret.data.path
        })
    }
    $scope.$watch("tile.Items",function(val,old){
        if(val.length != old.length) return
        for(var i = 0; i < val.length; i++){
            var ni = val[i]
            var oi = old[i]
            if(!ni || !oi) return
            if(ni.Title != oi.Title || ni.Memo != oi.Memo || ni.ImgPath != oi.ImgPath || ni.Url != oi.Url){
                ni.$update().then(function(e,e1){
                },function () {
                    $scope.tile.Itmes[i].Title = oi.Title
                    $scope.tile.Itmes[i].Memo = oi.Memo
                    $scope.tile.Itmes[i].ImgPath = oi.ImgPath
                    $scope.tile.Itmes[i].Url = oi.Url
                })
            }
        }
    },true)
    Tiles.query(function(tiles,header){
        $scope.tile.Items = tiles
    })
    $scope.removeTile = function (index){
        Tiles.delete({index:index},{},function(){
            $scope.tile.Items.splice(index,1)
        })
    }
}])
