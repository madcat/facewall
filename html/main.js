


var grid = {}
grid.controller = function(){
  var c = this
  c.state = m.prop("all")
  c.guests = m.prop([])
  c.historyWinners = m.prop([])
  c.getGuests = function(){
    m.request({method:"GET", url:"/guest"}).then(function(data){
      c.guests(data)
      setTimeout(c.pullState, 1000)
    })
  }

  c.pullState =function(){
    m.request({method:"GET", url:"/shuffle/state", deserialize:function(s){
      if (c.state()!=s){
        console.log("changing state to "+s)
        c.state(s)
        if (s.indexOf("shuffling")==0) {
          c.startShuffling()
        } else if (s.indexOf("history")==0){
          m.request({method:"GET", url:"/shuffle/"+s}).then(function(data){
            c.historyWinners(data)
            m.request({method:"GET", url:"/guest"}).then(function(data){
              c.guests(data)
            })
          })
        }
      }
      setTimeout(c.pullState, 1000)
    }})
  }
  c.startShuffling = function(){
    // animate turn off mask on available cells
    var candidates = c.guests().filter(function(g){
      if (g.Prize == "") return true
      return false
    })
    var setRandomGid = function(){
      c.shufflingGid = candidates[Math.floor(Math.random()*candidates.length)].Gid
      m.redraw()
      if (c.state().indexOf("shuffling")==0) setTimeout(setRandomGid, 1000/16)
    }
    setRandomGid()
  }

  c.getGuests()
}
grid.view = function(ctrl){
  return m(".container", [
    m(".grid", [
      ctrl.guests().map(function(g){
        if (g.Prize==""){
          return m(".cell-box", [
            m(".cell.cell-on", {id:"g"+g.Gid}),
            m("span", g.Code),
            m(".mask", {class:ctrl.state().indexOf("shuffling")==0 ? (g.Gid==ctrl.shufflingGid)?"hidden":"" : "hidden"})
            ])
        } else {
          return m(".cell-box", [
            m(".cell.cell-off", {id:"g"+g.Gid}),
            m(".mask", {class:ctrl.state().indexOf("shuffling")==0?"":"hidden"})
            ])
        }
      })
    ]),
    m(".history", {class:ctrl.state().indexOf("history")==0?"":"hidden"}, [
      m(".winner-box", [
        ctrl.historyWinners().map(function(g){
          return m(".winner", g.Code)
        })
      ])
    ])
  ])
}

m.mount(document.getElementById("content"), grid);