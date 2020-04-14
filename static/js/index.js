//map set
var latlngArr = []
var m = null
var mapURL = 'http://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png'
var shpfileLayer = null
//--map set

var layer;
var leftTreeWdith = 360;
var cols = [];
var joinObj = {};
var addList = []
var shapename = ""

$(function () {
    initMap()
    resetSize();
    // getCols();
    renderForm()
    getFileList();
    $("#exports").click(function () {
        if (shapename == "") {
            layer.msg("请先选择文件")
        } else {
            window.location.href = "/exports";
        }

    });
    $("#joinAdd").click(function () {
        layer.open({
            type: 1,
            skin: 'layui-layer-rim', //加上边框
            area: ['800px', '800px'], //宽高
            content: $("#editDIV"),
            success: function (layero, index) {
                console.log(layero, index);
                bindItemEvent()
                queryExcelList()
                initShpSelect()
                initExcelHead()
            }
        });
    });

    $("#saveShp").click(function () {
        saveShp()
    });

    $("#convBT").click(function () {
        console.log("convBT...")
        var filename = $("#leftShpList").find("dd.layui-this").find("a").data("filename")
        console.log(filename)
        layer.confirm('GCJ02坐标系==>WGS84坐标系？', {
            btn: ['确认', '取消'] //按钮
        }, function () {
            var encoding = $("#selectEncoding").val()
            $.ajax({
                url: "http://localhost:8080/converCoordinate",
                type: "get",
                data: {filename: filename, encoding: encoding},
                dataType: "json",
                success: function (data) {
                    if (data.code != 0) {
                        layer.msg('转换成功', {icon: 1});
                    }
                }
            })
        }, function () {
            // layer.msg('已经取消', {icon: 1});
        });
    });

});

function resetSize() {
    var allHeight = $(document).height()
    var editHeaderHeight = $(".edit-header").height()
    $("#maincenter").height((allHeight - editHeaderHeight) + "px")
    $("#leftList").height((allHeight - editHeaderHeight) + "px")
}

function getCols(filename) {
    var encoding = $("#selectEncoding").val()
    $.ajax({
        url: "http://localhost:8080/getCols",
        type: "get",
        data: {filename: filename, encoding: encoding},
        dataType: "json",
        success: function (data) {
            console.log(data.data)
            cols = getTableCols(data.data)
            initTable(filename)
        }
    })
}

function getTableCols(datas) {
    var results = []
    for (var index in datas) {
        results.push({field: datas[index]["Name"], width: 80, title: datas[index]["Name"]})
    }
    return results;
}

function initLayer() {
    layui.use('element', function () {
        var element = layui.element; //导航的hover效果、二级菜单等功能，需要依赖element模块

        //监听导航点击
        element.on('nav(headNavLayFilter)', function (elem) {
            //console.log(elem)
            var type = elem.data("type")
            switch (type) {
                case "opendir":
                    // CarryOut()
                    break;
                case "openfile":
                    break;
                case "exit":
                    break;
                default:
                    break;
            }
        });
        element.on('nav(leftLayFilter)', function (elem) {
            var type = elem.data("type")
            if (type == "shape") {
                var filename = elem.data("filename")
                getCols(filename)
                shapename = filename
            }
            console.log(elem.data("filename"))
            clearShapeFileLayer()
            loadShapeMap()
            // layer.msg(elem.text()+"11");
        });
    });
}


function saveShp(e) {
    console.log("saveShp....")
    var selectExcel = $("#selectExcel").val()
    var addListParam = addList.filter(function (item) {
        return item != undefined
    })
    if (addListParam.length == 0) {
        layer.msg("新增字段绑定不能为空！")
        return false
    }

    $.ajax({
        url: "http://localhost:8080/saveShp",
        type: "post",
        contentType: 'application/json;charset=utf-8',
        data: JSON.stringify({filename: selectExcel, joinObj: joinObj, addList: addListParam}),
        // data: JSON.stringify({joinObj: "joinObj", addList: "addList"}),
        dataType: "json",
        success: function (data) {
            if (data.code == 200) {
                layer.msg("保存成功")
            } else {
                layer.msg("保存失败")
            }
        }
    })
}

function initTable(filename) {
    layui.use('table', function () {
        var table = layui.table;

        table.render({
            elem: '#shpList',
            url: '/shpList',
            where: {filename: filename},
            cols: [cols],
            // cols: [[
            //     {field: 'id', width: 80, title: 'ID', sort: true}
            //     , {field: 'username', width: 80, title: '用户名'}
            //     , {field: 'sex', width: 80, title: '性别', sort: true}
            //     , {field: 'city', width: 80, title: '城市'}
            //     , {field: 'sign', title: '签名', minWidth: 150}
            //     , {field: 'score', width: 80, title: '评分', sort: true}
            //     , {field: 'classify', width: 80, title: '职业'}
            // ]],
            page: true
        });
    });
}

function clickAddNewFiled(e) {
    var newFiledText = $.trim($("#newFiledText").val())
    var bindFiledText = $.trim($("#selectBindExcel").val().trim())
    var num = addList.length;
    var row = addRow({num: num, newField: newFiledText, bindField: bindFiledText})
    $("#addTbody").append(row)
    addList.push({newFiledText: newFiledText, bindFiledText: bindFiledText})
}

function addRow(obj) {
    var rowTemplate = `
     <tr>
        <td>${obj.num}</td>
        <td>${obj.newField}</td>
        <td>${obj.bindField}</td>
        <td><a href="javascript:;" data-index = ${obj.num} onclick="deleteTableRow(this)">删除</a>|</td>
    </tr>
    `;
    return rowTemplate;
}


function bindItemEvent() {
    $(".bind-item").unbind("mouseover").bind("mouseover", function () {
        $(this).children("img").show()
    }).unbind("mouseout").on("mouseout", function () {
        $(this).children("img").hide()
    })
    $(".bind-item img").unbind("click").bind("click", function () {
        // var index = $(this).data("index");
        // routerList.splice(parseInt(index), 1, '0');
        // removeLayer(routerMarkerList[index])
        var shpName = $(this).data("shp")
        delete joinObj[shpName]
        $(this).parent().remove();
    });
}

function queryExcelList() {
    $.ajax({
        url: "http://localhost:8080/getExcel",
        type: "get",
        data: {},
        dataType: "json",
        success: function (data) {
            console.log(data.data)
            var result = ""
            var option = `<option value=""></option>`;
            for (var i = 0; i < data.data.length; i++) {
                result = result + `<option value="${data.data[i]["file"]}">${data.data[i]["file"]}</option>`;
            }
            result = "<option value=''></option>" + result
            $("#selectExcel").html(result)
            renderForm()
        }
    })
}

function getExcelHead(filename) {
    $.ajax({
        url: "http://localhost:8080/getExcelHead",
        type: "get",
        data: {filename: filename},
        dataType: "json",
        success: function (data) {
            console.log(data.data)
            var result = ""
            var option = `<option value=""></option>`;
            for (var obj in data.data) {
                result = result + `<option value="${obj}">${obj}</option>`;
            }
            $("#selectExcelHead").html(result)
            $("#selectBindExcel").html(result)
            renderForm()
        }
    })
}

function initExcelHead() {
    var selectExcelFile = $("#selectExcel").val()
    if (selectExcelFile != null && selectExcelFile != undefined && selectExcelFile != "") {
        getExcelHead(selectExcelFile)
    }

}

//重新渲染表单
function renderForm() {
    layui.use('form', function () {
        var form = layui.form;//高版本建议把括号去掉，有的低版本，需要加()
        form.render();

        form.on('select(selectExcel)', function (data) {
            console.log(data.elem); //得到select原始DOM对象
            console.log(data.value); //得到被选中的值
            console.log(data.othis); //得到美化后的DOM对象
            initExcelHead()
        });
        form.on('select(selectBindExcel)', function (data) {
            console.log(data.elem); //得到select原始DOM对象
            console.log(data.value); //得到被选中的值
            console.log(data.othis); //得到美化后的DOM对象
            initExcelHead()
        });

        form.on('select(selectEncoding)', function (data) {
            if (shapename == "") {
                layer.msg("请先选择文件")
                return;
            } else {
                getCols(shapename)
            }
        });

    });
}

function initShpSelect() {
    var result = ""
    for (var i = 0; i < cols.length; i++) {
        result = result + `<option value="${cols[i]["field"]}">${cols[i]["field"]}</option>`;
    }
    $("#selectShp").html(result)
    renderForm()
}

function addBindJoin() {
    var shpName = $("#selectShp").val()
    var excelHeadName = $("#selectExcelHead").val()
    if (joinObj[shpName] == null || joinObj[shpName] == undefined) {
        joinObj[shpName] = excelHeadName
        $("#joinListArea").append(`<span class="bind-item">${shpName}=${excelHeadName}<img src="../static/img/close.png" data-shp="${shpName}" alt="" style=""></span>`)
        bindItemEvent()
    } else {
        layer.msg(shpName + "已经绑定！")
    }
}

function deleteTableRow(e) {
    console.log(e)
    var index = $(e).data("index");
    $(e).parent().parent().remove();
    // addList.splice(index, 1);
    delete addList[index]
}

function getFileList() {
    var dirPath = "文件列表"
    $.ajax({
        url: "http://localhost:8080/getFileList",
        type: "post",
        contentType: 'application/json;charset=utf-8',
        // data: {dirPath: dirPath},
        dataType: "json",
        success: function (data) {
            console.log(data)

            var _lis = `<a href="javascript:;">${dirPath}</a><dl class="layui-nav-child">`;


            if (data.code == 200) {
                for (var index in data.data) {
                    console.log(index)
                    var filePath = data.data[index]
                    var fileNames = filePath.split("\\")
                    var fileName = fileNames[fileNames.length - 1]
                    console.log(fileName)
                    _lis = _lis + `<dd><a href="javascript:;" data-type="shape" title="${fileName}" data-filename="${fileName}">${fileName}</a></dd>`;
                }
                _lis = _lis + "</dl>";
                $("#leftShpList").html(_lis)
                initLayer()
            } else {
                layer.msg("获取文件列表有问题")
            }

        }
    })
}


function CarryOut() {
    var inputObj = document.createElement('input')
    inputObj.setAttribute('id', '_ef');
    inputObj.setAttribute('type', 'file');
    inputObj.setAttribute("style", 'visibility:hidden');
    // inputObj.setAttribute("webkitdirectory","");
    document.body.appendChild(inputObj);
    inputObj.click();
    inputObj.value;
    console.log(inputObj.value)
    inputObj.onchange = function (e) {
        console.log("inputObj  change ...", inputObj.value, e)
        var dirPath = getDirPathByFile(inputObj.value)
        getFileList(dirPath)
    }
}

function getDirPathByFile(filepath) {
    var dirPath = ""
    var arrs = filepath.split("\\");
    for (var i = 0; i < arrs.length - 1; i++) {
        dirPath = dirPath + arrs[i] + "\\";
    }
    return dirPath;
}


layui.use('layer', function () {
    layer = layui.layer;
});

function initMap() {
    m = L.map('map').setView([42.5463, 86.0553], 6);
    var watercolor = L.tileLayer(mapURL, {
        attribution: 'Map tiles by <a href="http://stamen.com">Stamen Design</a>, <a href="http://creativecommons.org/licenses/by/3.0">CC BY 3.0</a> &mdash; Map data &copy; <a href="http://openstreetmap.org">OpenStreetMap</a> contributors, <a href="http://creativecommons.org/licenses/by-sa/2.0/">CC-BY-SA</a>'
    }).addTo(m);
    loadShapeMap()

}

function loadShapeMap() {
    console.log("shapename=====", shapename)
    latlngArr = []
    if (shapename != '') {
        var lastIndex = shapename.lastIndexOf(".")
        var shapeNameTemp = shapename.slice(0, lastIndex)
        shpfileLayer = new L.Shapefile('../../shape/' + shapeNameTemp, {
            onEachFeature: function (feature, layer) {
                if (feature.properties) {
                    layer.bindPopup(Object.keys(feature.properties).map(function (k) {
                        return k + ": " + feature.properties[k];
                    }).join("<br />"), {
                        maxHeight: 2000
                    });
                }
                if (feature.geometry.type == "Point") {
                    latlngArr.push(layer._latlng)
                } else {
                    latlngArr = latlngArr.concat(layer._latlngs)
                }


            }
        });
        shpfileLayer.addTo(m);
        shpfileLayer.once("data:loaded", function (data) {
            console.log("finished loaded shapefile");
            //m.setView()
            bounds = L.latLngBounds(latlngArr)
            m.fitBounds(bounds)
            latlngArr = []
        });
    }
}

function clearShapeFileLayer() {
    if (shpfileLayer != null) {
        m.removeLayer(shpfileLayer)
    }
}