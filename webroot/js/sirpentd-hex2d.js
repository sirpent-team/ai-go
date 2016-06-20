function SirpentHex2DGame(game_id, canvas_id) {
  this.canvas = document.getElementById(canvas_id)
  this.context = this.canvas.getContext("2d")

  this.width = this.canvas.width
  this.height = this.canvas.height

  this.hexagon_rings = 30
  this.hexagons_across = this.hexagon_rings * 2 + 1

  // The 0.1 discount is to leave a tiny bit of a border.
  this.radius = Math.min(
    this.width / (this.hexagons_across * Math.sqrt(3)),
    this.height / (this.hexagons_across * 1.5 + 0.5)
  ) - 0.1

  this.drawHexagons()

  var ws = new WebSocket("ws://localhost:8080/worlds/live.json")
  ws.onmessage = function(e) {
    this.clear()
    this.drawHexagons()
    var game_state = JSON.parse(event.data)
    for (player_id in game_state.Plays) {
      var player_state = game_state.Plays[player_id]
      this.drawSnake(player_state.CurrentSnake)
    }
    this.drawHexagon(game_state.Food, "rgb(200, 0, 0)", "rgb(120, 0, 0)")
  }.bind(this)
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
        this.drawHexagon({"Q": cube["X"], "R": cube["Z"]}, "rgb(150,150,150)", null)
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
  return {"X": x, "Y": y}
}

SirpentHex2DGame.prototype.drawHexagon = function (axial, strokeColor, fillColor) {
  var canvas_x = this.width / 2 + this.radius * Math.sqrt(3) * (axial.Q + axial.R/2)
  var canvas_y = this.height / 2 + this.radius * 1.5 * axial.R

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

SirpentHex2DGame.prototype.drawSnake = function (snake) {
  var i
  for (i = 0; i < snake["Segments"].length; i++) {
    var color = (i == 0) ? "rgb(0, 120, 0)" : "rgb(0, 200, 0)"
    this.drawHexagon(snake["Segments"][i], "rgb(0, 120, 0)", color)
  }
}

SirpentHex2DGame.prototype.clear = function () {
  this.context.clearRect(0, 0, this.width, this.height)
}
