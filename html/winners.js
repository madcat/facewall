function winnerRow(ctrl, w){
  return m("tr", [
      m("td.col-sm-1", w.Step),
      m("td.col-sm-1", w.Prize),
      m("td.col-sm-1", w.Code),
      m("td.col-sm-1", w.Tag)
    ])
}

var list = {}
list.controller = function(){
  var c = this
  c.winners = m.prop([])
  c.getWinners = function(){
    m.request({method:"GET", url:"/winner"}).then(function(data){
      c.winners(data)
    })
  }
  c.resetWinners = function(){
    bootbox.confirm("这将清楚所有获奖数据，你确定要这么做?", function(result) {
      if(result){
        bootbox.confirm("确定??", function(result2) {
          if(result2){
            m.request({method:"POST", url:"/shuffle/reset"}).then(function(){
              c.getWinners()
            })
          }
        })
      }
    });
  }
  c.getWinners()
}
list.view = function(ctrl){
  return m(".container-fluid", [
    m(".row", [m("button.btn-danger", {onclick:ctrl.resetWinners}, "重置获奖记录")]),
    m(".row", [
        m("table.table-condensed.table.users", [
          m("thead", [
            m("tr", [
              m("th", "第#次"),
              m("th", "获奖"),
              m("th", "代码"),
              m("th", "分组")
            ])
          ]),
          m("tbody", [
            ctrl.winners().map(function(w){
                return winnerRow(ctrl, w)
            })
          ])
        ])
      ]),
    m(".row", [m(".footer", "共计"+ctrl.winners().length+"人")]),
    ])
}

m.mount(document.getElementById("content"), list);