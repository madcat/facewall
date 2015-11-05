### api

  GET /guest

  [
    {"1234abcd", "2", "", "group1", "http://x/2.png"},
    {"2345bcde", "3", "", "group2", "http://x/3.png"}
    ...
  ]

  GET /win

  [
    {"1", "1234abcd", "五等奖"},
    {"1", "3456cdef", "五等奖"},
    ...
    {"2", "dddd4444", "四等奖"},
    ...
  ]


#### facewall page

facewall page is stateless. server has state.

  GET /state - facewall javascript periodically pull for its state.

  grid - should GET /guest and update grid when changed to this state, once
  shuffling/{prize} - start shuffle animation when changed to this state, once
  overlay/{order} - should GET /prize and update both grid and overlay display when changed to this state, once

#### ipad remote app

  GET /guest - get all guests

  GET /config - prize configuration

  [
    {"五等奖": 80},
    {"四等奖": 40},
    ...
  ]

  POST /shuffle/start/{prize}

  POST /shuffle/end/ - end last shuffle, then GET /prize to update views, then enable shuffle start in 10 seconds

  POST /overlay/{order} - change server state

### report web

  GET /win - get winning data join with guest data, order by order

#### server

  server start - state="grid"

  /shuffle/start/{prize} - state="shuffling/{prize}"

  /shuffle/end/ - run and write result to win table; state="overlay/{order}"

