var server = "http://localhost:8081";
var pollRate_ms = "500";
var processingRate_ms = "250";
var maxPollGenerations = 15;
var updateQueueLimit = 100;
var StatusStr = ["Seeded", "Active", "Stable", "Dead"];

var analyses = {};

function getIdStr(id) {
    var idStr = id.toString(16).replace(new RegExp("[/+=]", 'g'), "");
    return idStr.substring(0, idStr.length-1);
}

function createAnalysis(data) {
    var idStr = getIdStr(data.Id);

    analyses[idStr] = {
        id: data.Id,
        idAsStr: idStr,
        poller : null,
        processed : 0,
        running : false,
        // generations : [],
        updateQueue : [],
        elements : {
            currentGeneration : null,
            cells : {},
        },
        AddToQueue : function(data) {
                        // console.log("  AddToQueue()", data);
                        // Add each update to the queue
                        for (var i = 0; i < data.Updates.length; i++) {
                            this.updateQueue.push(data.Updates[i]);
                        }

                        if(this.updateQueue.length < updateQueueLimit) {
                        // if(this.updateQueue.length < updateQueueLimit && this.running) {
                            setTimeout($.proxy(this.Processor, this), processingRate_ms);
                        }
                    },
        Processor : function() {
                    // console.log("Process()", this);
                    var update = this.updateQueue.shift();

                    if (update != undefined) { // TODO: Why is update sometimes undefined?
                        // $("#generation-"+this.idAsStr).text(update.Generation);
                        this.elements.currentGeneration.text(update.Generation);

                        // console.log("CELLS: ", this.elements.cells);
                        // var idPrefix = "#cell-"+this.idAsStr+"-";
                        for (var i = update.Changes.length-1; i >= 0; i--) {
                            var changed = update.Changes[i];

                            switch (changed.Change) {
                            case 0: // Born
                                // console.log("BIRTHING: ", changed.Y, changed.X, this.elements.cells[changed.Y][changed.X]);
                                this.elements.cells[changed.Y][changed.X].addClass("analysisBoardCellAlive");
                                // $(idPrefix+changed.Y+"x"+changed.X).addClass("analysisBoardCellAlive");
                                break;
                            case 1: // Died
                                // console.log("KILLING: ", changed.Y, changed.X);
                                this.elements.cells[changed.Y][changed.X].removeClass("analysisBoardCellAlive");
                                // $(idPrefix+changed.Y+"x"+changed.X).removeClass("analysisBoardCellAlive");
                                break;
                            }
                        }

                        this.processed++;

                        // Keep processing
                        // if (this.updateQueue.length > 0) {
                        if (this.updateQueue.length > 0 && this.updateQueue.length < updateQueueLimit && this.running) {
                            setTimeout($.proxy(this.Processor, this), processingRate_ms);
                        }
                    }
                },
        Start : function() {
                    this.poller = setInterval(function() { pollAnalysisRequest(this.id,
                                                                               this.processed + this.updateQueue.length + 1,
                                                                               maxPollGenerations) },
                                                           pollRate_ms);
                    controlAnalysisRequest(this.id, 0);
                    this.running = true;
                },
        Stop : function() {
                    clearInterval(this.poller);
                    controlAnalysisRequest(this.id, 1);
                    this.running = false;
                }

    };

    // Create the dom entity
    var cells = analyses[idStr].elements.cells;
    var board = $("<span></span>");
    for (var y = 0; y < data.Dims.Height; y++) {
        var row = $("<div></div>");
        for(var x = 0; x < data.Dims.Width; x++) {
            var cell = $("<span></span>")
                    // .attr("id", "cell-"+idStr+"-"+y+"x"+x)
                    .addClass("analysisBoardCell");
            // analyses[idStr].elements.cells[y][x] = cell;
            if (!(y in cells)) {
                cells[y] = {};
            }
            cells[y][x] = cell;
            // console.log("CREATING: ", y, x, cells[y][x]);
            row.append(cell);
        }
        board.append(row);
    }
    // console.log("CELLS: ", cells);

    analyses[idStr].elements.currentGeneration = $("<span></span>").text("0")
                                                                   // .attr("id", "generation-"+idStr)
                                                                   .addClass("analysisGeneration");

    $("#analyses").append(
        $("<div></div>").attr("id", "analysis-"+idStr)

        // ID
        /*
        .append($("<div></div>").addClass("analysisField")
                .append("<span>ID: </span>")
                .append($("<span></span>").addClass("analysisId").text(idStr)))
        */


        // Rule
        // .append($("<div></div>")
        //         .append("<span>Rule: </span>")
        //         // .append($("<span></span>").text(data.Id.toString(16)))
        //         .addClass("analysisField").addClass("analysisRule"))

        // Neighbors
        /*
        .append($("<div></div>").addClass("analysisField")
                .append("<span>Neighbors: </span>")
                .append($("<span></span>").addClass("analysisNeighbors").text("TODO")))

        // Status
        .append($("<div></div>").addClass("analysisField")
                .append("<span>Status: </span>")
                .append($("<span></span>").text("Unknown")
                    .attr("id", "status-"+idStr)
                    .addClass("analysisStatus")))

        */

        // Generation
        .append($("<div></div>").addClass("analysisField")
        .append($("<span>Generation: </span>"))
                .append(analyses[idStr].elements.currentGeneration))
                // .append($("<span></span>").text("0")
                    // .attr("id", "generation-"+idStr)
                    // .addClass("analysisGeneration")))

        // Control
        .append($("<div style='height: 40px'></div>")
                /*
                .append($("<span></span>").addClass("analysisControl")
                                        // .click(function() { stopAnalysis(idStr) })
                                        // .click($.proxy(analyses[idStr].Stop, analyses[idStr]))
                                        .text("⟲")) // TODO
                */
                .append($("<span></span>").addClass("analysisControl")
                                        .click(function() { startAnalysis(idStr) })
                                        // .click($.proxy(analyses[idStr].Start, analyses[idStr]))
                                        .text("▶"))
                .append($("<span></span>").addClass("analysisControl")
                                        .click(function() { stopAnalysis(idStr) })
                                        // .click($.proxy(analyses[idStr].Stop, analyses[idStr]))
                                        .text("⬛"))
                )

        // Create the board
        .append($("<div></div>").attr("class", "analysisBoard").attr("id", "board-"+idStr)
                .html(board))
    );
}

function startAnalysis(idStr) {
    var analysis = analyses[idStr];
    analysis.poller = setInterval(function() { pollAnalysisRequest(analysis.id,
                                                                   analysis.processed + analysis.updateQueue.length + 1,
                                                                   maxPollGenerations) },
                                  pollRate_ms);
    controlAnalysisRequest(analysis.id, 0);
    analysis.running = true;
}

function stopAnalysis(idStr) {
    var analysis = analyses[idStr];
    clearInterval(analysis.poller);
    controlAnalysisRequest(analysis.id, 1);
    analysis.running = false;
}

//////////////////// BOARD CREATOR ////////////////////

function resizeCanvas(w, h) {
    var canvas = document.getElementById("board");
    canvas.width = w;
    canvas.height = h;
}

function addBoardCreator() {
    $("body").append($("<div></div>").attr("id", "boardCreator")
                       .resize("resizeCanvas(parseInt(this.style.width), parseInt(this.style.height))")
                       .append($("<div></div>").attr("id", "segrip").addClass("ui-resizable-handle").addClass("ui-resizable-se"))
                       .append($("<div></div>").attr("id", "egrip").addClass("ui-resizable-handle").addClass("ui-resizable-e"))
                       .append($("<div></div>").attr("id", "sgrip").addClass("ui-resizable-handle").addClass("ui-resizable-s"))
                       .append($("<canvas></canvas>").attr("id", "board"))
                      );
    $('#boardCreator').resizable({
          handles: {
            'se': '#segrip',
            'e': '#egrip',
            's': '#sgrip',
          }
        });
}

var cellWidth = 3;
var cellHeight = 3;
var spacing = 1;

function createBoard(id, width, height) {
  if (width > 0 && height > 0) {
    // TODO: figure out real width and height based on cellWidth and cellHeight + spacing
    var canvas = document.getElementById(id);
    canvas.width = width;
    canvas.height = height;
  }
}

function updateBoard(id, cellWidth, cellHeight, spacing, coverage) {
  var adjWidth = spacing+cellWidth;
  var adjHeight = spacing+cellHeight;

  var canvas = document.getElementById(id);
  var ctx = canvas.getContext('2d');
  ctx.clearRect(0, 0, canvas.width, canvas.height);

  var numPerRow = canvas.width/(cellWidth+spacing);
  var numPerCol = canvas.height/(cellHeight+spacing);

  ctx.save();

  var aliveVal = 100-coverage;
  for (var row = 0; row < numPerCol; row++) {
    for (var col = 0; col < numPerRow; col++) {
      if ((Math.random() * 100) > aliveVal) {
        ctx.save();
        ctx.fillStyle = '#4863a0';
        ctx.translate(col*adjWidth, row*adjHeight);
        ctx.fillRect(0,0,cellWidth,cellHeight);
        ctx.restore();
      }
    }
  }
  ctx.restore();
}

function createBoardFromInputs() {
  var width = parseInt(document.getElementById('boardWidth').value);
  var height = parseInt(document.getElementById('boardHeight').value);
  createBoard('canvas', width, height);
}

function updateFromInputs() {
  createBoardFromInputs();

  cellWidth = parseInt(document.getElementById('cellSize').value);
  cellHeight = cellWidth;
  var coverage = parseInt(document.getElementById('coverage').value);
  updateBoard('canvas', cellWidth, cellHeight, spacing, coverage);
}

//////////////////// REQUESTORS ////////////////////

function createNewAnalysisRequest() {
    // $.post(server+"/analyze", JSON.stringify({"Dims":{"Height": 200, "Width": 300}, "Pattern": 0})) // FIXME
    $.post(server+"/analyze", JSON.stringify({"Dims":{"Height": 50, "Width": 100}, "Pattern": 0})) // FIXME
    .done(function( data ) {
        createAnalysis(data);
        pollAnalysisRequest(data.Id, 0, maxPollGenerations);
    });
}

function pollAnalysisRequest(key, startingGen, maxGen) {
    // console.log("pollAnalysisRequest()", startingGen, maxGen);
    $.post(server+"/poll", JSON.stringify({"Id": key, "StartingGeneration": startingGen, "NumMaxGenerations": maxGen}))
    .done(function( data ) {
        var idStr = getIdStr(data.Id);
        if (idStr in analyses) {
            analyses[idStr].AddToQueue(data);
        } else {
            console.log("Got update for unknown analysis")
        }
    });
}

function controlAnalysisRequest(key, order) {
    $.post(server+"/control", JSON.stringify({"Id":  key, "Order": order}));
}
