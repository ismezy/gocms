/**
 * Created by 扬 on 2016/9/19.
 */
"use strict"
var app = require("../main")
app.config(function($routeProvider){
    $routeProvider.when('/tile',{template:require("./admin-tile.html")})
})
app.factory('Tiles', ['$resource', function($resource) {
    return $resource('/admin/tile/:Index', null, {
        update: { method:'PUT' }
    })
}])
app.controller('tileController',function($scope,$http,Tiles,Upload){
    $scope.tile = {
        Items:[]
    }
    $scope.add = function () {
        var tile = new Tiles({
            Title       :     "磁贴" + $scope.tile.Items.length,
            Memo        :     "描述",
            ImgPath     :       "",
            Url         :     "/"
        })
        tile.$save().then(function(result){
            $scope.tile.Items.push(tile)
        },function (error) {
        })
    }
    $scope.upload = function (file,index) {
        Upload.upload({
            url:'/admin/image/upload',
            file:file
        }).then(function (ret) {
            $scope.tile.Items[index].ImgPath = ret.data.path
            console.log($scope.tile.Items[index].ImgPath)
        })
    }
    $scope.$watch("tile.Items",function(val,old){
        for(var i = 0; i < val.length || i < old.length; i++){
            var ni = val[i]
            var oi = old[i]
            if(!ni || !oi) return
            console.log(ni,oi)
            if(ni.Title != oi.Title || ni.Memo != oi.Memo || ni.ImgPath != oi.ImgPath || ni.Url != oi.Url){
                ni.$update().then(function(e,e1){
                    console.log(e,e1)
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
    //$http.get("/admin/tile/list").success(function(res){
     //   $scope.tile = res
    //})
})
