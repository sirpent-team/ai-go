function SirpentHex2DGame(game_id, canvas_id) {
  this.canvas = document.getElementById(canvas_id)
  this.context = this.canvas.getContext("2d")

  this.width = this.canvas.width
  this.height = this.canvas.height

  /*this.hexagon_rings = 30
  this.hexagons_across = this.hexagon_rings * 2 + 1

  // The 0.1 discount is to leave a tiny bit of a border.
  this.radius = Math.min(
    this.width / (this.hexagons_across * Math.sqrt(3)),
    this.height / (this.hexagons_across * 1.5 + 0.5)
  ) - 0.1*/

  //this.drawHexagons()

  var grid = false
  var ws = new WebSocket("ws://localhost:8080/worlds/live.json")
  ws.onmessage = function(e) {
    var game_state = JSON.parse(event.data)

    if (!grid) {
      grid = game_state
      this.hexagon_rings = grid.Rings
      this.hexagons_across = this.hexagon_rings * 2 + 1
      this.radius = Math.min(
        this.width / (this.hexagons_across * Math.sqrt(3)),
        this.height / (this.hexagons_across * 1.5 + 0.5)
      ) - 0.1
      this.clear()
      this.drawHexagons()
      return
    }

    this.clear()
    this.drawHexagons()
    for (player_id in game_state.Plays) {
      var player_state = game_state.Plays[player_id]
      this.drawPlayerSnake(player_state)
    }
    this.drawHexagon(game_state.Food, "rgb(200, 0, 0)", "rgb(120, 0, 0)")
  }.bind(this)
  ws.onclose = function(e) {
    console.log("onclose!")
    setTimeout(function() {
      window.location.reload()
    }, 2500)
  }
  //ws.send(data)
}

SirpentHex2DGame.prototype.drawHexagons = function () {
  var cube = {"X": 0, "Y": 0, "Z": 0}

  for (cube["X"] = -this.hexagon_rings; cube["X"] <= this.hexagon_rings; cube["X"]++) {
    for (cube["Y"] = -this.hexagon_rings; cube["Y"] <= this.hexagon_rings; cube["Y"]++) {
      for (cube["Z"] = -this.hexagon_rings; cube["Z"] <= this.hexagon_rings; cube["Z"]++) {
        if (cube["X"] + cube["Y"] + cube["Z"] != 0) {
          continue
        }
        this.drawHexagon(cube, "rgb(150,150,150)", null)
      }
    }
  }
}

SirpentHex2DGame.prototype.outlineHexagon = function (x, y) {
  this.context.beginPath()

  var i
  for (i = 0; i < 6; i++) {
    var hc = this.hexCorner({"X": x, "Y": y}, this.radius, i)
    if (i == 0) {
      this.context.moveTo(hc.X, hc.Y)
    } else {
      this.context.lineTo(hc.X, hc.Y)
    }
  }
}

// hex corner, http://www.redblobgames.com/grids/hexagons/#coordinates
SirpentHex2DGame.prototype.hexCorner = function (center, radius, i) {
  var angle_deg = 60 * i   + 30
  var angle_rad = Math.PI / 180 * angle_deg
  var x = center["X"] + radius * Math.cos(angle_rad)
  var y = center["Y"] + radius * Math.sin(angle_rad)
  //console.log({"X": x, "Y": y})
  return {"X": x, "Y": y}
}

SirpentHex2DGame.prototype.drawHexagon = function (hex_vector, strokeColor, fillColor) {
  var canvas_x = this.width / 2 + this.radius * Math.sqrt(3) * (hex_vector.X + hex_vector.Z/2)
  var canvas_y = this.height / 2 + this.radius * 1.5 * hex_vector.Z

  this.outlineHexagon(canvas_x, canvas_y)

  var tmp

  if (fillColor) {
    tmp = this.context.fillStyle
    this.context.fillStyle = fillColor
    this.context.fill()
    this.context.fillStyle = tmp
  }

  this.context.closePath()

  if (strokeColor) {
    tmp = this.context.strokeStyle
    this.context.strokeStyle = strokeColor
    this.context.stroke()
    this.context.strokeStyle = tmp
  }
}

SirpentHex2DGame.prototype.axialToCube = function (axial) {
  return {"X": axial["Q"], "Z": axial["R"], "Y": -axial["Q"] -axial["R"]}
}

SirpentHex2DGame.prototype.drawPlayerSnake = function (player) {
  var i
  for (i = 0; i < player["Snake"].length; i++) {
    var r = (player["Alive"]) ? 0 : 255
    var color = (i == 0) ? "rgb(" + r + ", 120, 0)" : "rgb(" + r + ", 200, 0)"
    this.drawHexagon(player["Snake"][i], "rgb(" + r + ", 120, 0)", color)
  }
}

SirpentHex2DGame.prototype.clear = function () {
  this.context.clearRect(0, 0, this.width, this.height)
}
