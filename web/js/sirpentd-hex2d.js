function SirpentHex2DGame(game_id, canvas_id) {
  this.canvas = document.getElementById(canvas_id)
  this.context = this.canvas.getContext("2d")

  this.width = this.canvas.width
  this.height = this.canvas.height
  this.grid_width = 25 // @TODO
  this.grid_height = 25 // @TODO

  // 1.5radius per x plus 0.5radius to close
  // 1.5rx + 0.5r = width
  // (1.5x + 0.5)r = width
  // r = width / (1.5x + 0.5)

  // sqrt(3)rx + 0.5sqrt(3)r = height
  // sqrt(3)r(x + 0.5) = height
  // r = height / (x + 0.5)sqrt(3)
  // x = height / sqrt(3)r - 0.5

  this.radius = Math.floor(Math.min((this.width - 2) / (1.5 * this.grid_width + 0.5), (this.height - 2) / (Math.sqrt(3) * (0.5 + this.grid_height))))
  this.polygon_width = 2 * this.radius
  this.polygon_height = Math.sqrt(3) * this.radius

  this.drawHexagons()

  var ws = new WebSocket("ws://localhost:8080/worlds/live")
  ws.onmessage = function(e) {
    this.clear()
    this.drawHexagons()
    this.drawSnake(JSON.parse(event.data))
  }.bind(this)
  //ws.send(data);
}

SirpentHex2DGame.prototype.drawHexagons = function () {
  var row, column
  for (row = 0; row < this.grid_height; row++) {
    for (column = 0; column < this.grid_width; column++) {
      this.drawHexagon(column, row, "rgb(150,150,150)", null)//, "rgb(0,0,0)")
    }
  }
}

SirpentHex2DGame.prototype.outlineHexagon = function (x, y) {
  this.context.beginPath()
  this.context.moveTo(x - this.radius, y)
  this.context.lineTo(x - this.radius * 0.5, y + 0.5 * this.polygon_height)
  this.context.lineTo(x + this.radius * 0.5, y + 0.5 * this.polygon_height)
  this.context.lineTo(x + this.radius, y)
  this.context.lineTo(x + this.radius * 0.5, y - 0.5 * this.polygon_height)
  this.context.lineTo(x - this.radius * 0.5, y - 0.5 * this.polygon_height)
  this.context.lineTo(x - this.radius, y)
}

SirpentHex2DGame.prototype.drawHexagon = function (column, row, strokeColor, fillColor) {
  if (column < 0 || row < 0 || column >= this.grid_width || row >= this.grid_height) {
    return
  }
  var canvas_x = 1 + this.polygon_width / 2 + column * 0.75 * this.polygon_width
  var canvas_y = 1 + this.polygon_height / ((column % 2) + 1) + row * this.polygon_height
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

SirpentHex2DGame.prototype.drawSnake = function (snake) {
  console.log(snake)
  var i
  for (i = 0; i < snake["Length"]; i++) {
    var color = (i == 0) ? "rgb(0, 120, 0)" : "rgb(0, 255, 0)"
    this.drawHexagon(snake["Segments"][i]["Position"]["X"], snake["Segments"][i]["Position"]["Y"], "rgb(0, 120, 0)", color)
  }
}

SirpentHex2DGame.prototype.clear = function () {
  this.context.clearRect(0, 0, this.width, this.height)
}
