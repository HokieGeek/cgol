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

//////////////////// BOARD CREATOR ////////////////////

var cellWidth = 3;
var cellHeight = 3;
var cellSpacing = 1;
var cellAliveColor = '#4863a0'; // FIXME: hmmm

function generateRandomSeed(id, cellWidth, cellHeight, spacing, coverage) {
    var seed = [];

  var adjWidth = cellSpacing+cellWidth;
  var adjHeight = cellSpacing+cellHeight;

  var board = document.getElementById(id);
  var numPerRow = board.width/(cellWidth+cellSpacing);
  var numPerCol = board.height/(cellHeight+cellSpacing);

  var ctx = board.getContext('2d');
  ctx.clearRect(0, 0, board.width, board.height);

  ctx.save();

  var aliveVal = 100-coverage;
  for (var row = numPerCol-1; row >= 0; row--) {
    for (var col = numPerRow-1; col >= 0; col--) {
      if ((Math.random() * 100) > aliveVal) {
        var x = col*adjWidth;
        var y = row*adjHeight;

        ctx.save();
        ctx.fillStyle = cellAliveColor;
        ctx.translate(x, y);
        ctx.fillRect(0, 0, cellWidth, cellHeight);
        ctx.restore();

        seed.push({"X":x, "Y": y});
      }
    }
  }
  ctx.restore();

  return seed;
}

function updateFromInputs(key, width, height) {
  var id = "board-"+key;
  if (width > 0 && height > 0) {
    // TODO: figure out real width and height based on cellWidth and cellHeight + spacing
    var canvas = document.getElementById(id);
    canvas.width = width;
    canvas.height = height;
  }

  //cellWidth = parseInt(document.getElementById('cellSize').value);
  cellWidth = parseInt($('#cellSize-'+key).val());
  cellHeight = cellWidth;
  var coverage = parseInt($('#cellDensity-'+key).val());
  analyses[key].seed = generateRandomSeed(id, cellWidth, cellHeight, cellSpacing, coverage);
}

function initBoard(key, padre) {
    var board = $("<canvas></canvas>").attr("id", "board-"+key).addClass("analysisBoard ui-widget-content");

    padre.append(
        $("<span></span>")
            .append(board)
            .append($("<br/>"))
            // TODO: add the next "row" of controls to their own span that does vertical-align: center
            .append($("<span></span>").addClass("newBoardControls")
                .append($("<span></span>").text("Cell size: "))
                .append($("<input></input>").attr("type", "range")
                                            .attr("id", "cellSize-"+key)
                                            .attr("min", "1").attr("max", "5")
                                            .attr("value", "2")
                                            .change(function() { updateFromInputs(key, -1, -1); })
                                            .addClass("analysisBoardCellSizeSelector")
                       )
                .append($("<span></span>").text("Cell density: "))
                .append($("<input></input>").attr("type", "range")
                                                 .attr("id", "cellDensity-"+key)
                                                 .attr("min", "1").attr("max", "100")
                                                 .attr("value", "60")
                                                 .change(function() { updateFromInputs(key, -1, -1); }))
                .append($("<button></button>").text("Create")
                    .click(function() {
                        $("#board-"+key).resizable('destroy');
                        createNewAnalysisRequestNEW({"Dims": analyses[key].dimensions, "Pattern": 0, "Seed": analyses[key].seed},
                                                    function( data ) {
                                                        console.log("Returned analysis id: ", data.Id);
                                                        lifeIdtoAnalysisId[data.Id] = key;
                                                        analyses[key].id = data.Id;
                                                        console.log("The map: ", lifeIdtoAnalysisId);
                                                        pollAnalysisRequest(data.Id, 0, maxPollGenerations);
                                                    });
                    })
                )
            )
    );

    board.resizable({
      helper: "analysisBoard-resizable-helper",
      stop: function( event, ui ) { updateFromInputs(key, ui.size.width, ui.size.height); }
    });

    updateFromInputs(key);

    return board;
}

//////////////////// NEW ANALYSIS ////////////////////

// Create an analysis blah blah blah
// Fire the request and include a closure for the return id
// From closure, continue filling out the analysis and add to the map

var lifeIdtoAnalysisId = {};

function createAnalysisNEW() {
    var key = 42;
    analyses[key] = {
        id: null, // TODO: how about this be the lifed id but this client can keep track of them with its own ID     <<-------
        idAsStr: null,
        poller : null,
        processed : 0,
        running : false,
        updateQueue : [],
        seed : [],
        dimensions : { Width: 300, Height: 200 },
        elements : {
            currentGeneration : null,
            board : null,
        },
        AddToQueue : function(data) {
                        console.log("  AddToQueue()", data);
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

                        // TODO: update the canvas board

                        // for (var i = update.Changes.length-1; i >= 0; i--) {
                        //     var changed = update.Changes[i];

                        //     switch (changed.Change) {
                        //     case 0: // Born
                        //         this.elements.cells[changed.Y][changed.X].addClass("analysisBoardCellAlive");
                        //         break;
                        //     case 1: // Died
                        //         this.elements.cells[changed.Y][changed.X].removeClass("analysisBoardCellAlive");
                        //         break;
                        //     }
                        // }

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

    // Create the generation object
    analyses[key].elements.currentGeneration = $("<span></span>").text("0").addClass("analysisGeneration");

    // Create a div for this analysis and attach it to the primary div
    $("#analyses").append(
        $("<div></div>").attr("id", "analysis-"+key)

        // Generation field
        .append($("<div></div>").addClass("analysisField")
        .append($("<span>Generation: </span>")).append(analyses[key].elements.currentGeneration))

        // Control
        .append($("<div style='height: 40px'></div>")
                .append($("<span></span>").addClass("analysisControl")
                                        .click(function() { startAnalysis(key); })
                                        // .click(function() {
                                        //             this.poller = setInterval(function() { pollAnalysisRequest(key,
                                        //                                                                        this.processed + this.updateQueue.length + 1,
                                        //                                                                        maxPollGenerations) },
                                        //                                                    pollRate_ms);
                                        //             controlAnalysisRequest(this.id, 0);
                                        //             this.running = true;
                                        //         })
                                        .text("▶"))
                .append($("<span></span>").addClass("analysisControl")
                                        .click(function() { stopAnalysis(key); })
                                        // .click(function() {
                                        //         clearInterval(this.poller);
                                        //         controlAnalysisRequest(key, 1);
                                        //         this.running = false;
                                        //     })
                                        .text("⬛"))
                ) // Control

    );

    initBoard(key, $("#analysis-"+key));
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

//////////////////// REQUESTORS ////////////////////

function createNewAnalysisRequestNEW(req, callback) {
    console.log("createNewAnalysisRequestNEW(): ", req);
    $.post(server+"/analyze", JSON.stringify(req)).done(callback);
}

function pollAnalysisRequest(analysisId, startingGen, maxGen) {
    console.log("pollAnalysisRequest()", analysisId, startingGen, maxGen);
    console.log("analyses: ", analyses);
    $.post(server+"/poll", JSON.stringify({"Id": analysisId, "StartingGeneration": startingGen, "NumMaxGenerations": maxGen}))
    .done(function( data ) {
        var id = lifeIdtoAnalysisId[data.Id];
        if (id in analyses) {
            analyses[id].AddToQueue(data);
        } else {
            console.log("Got update for unknown analysis: ", data.Id)
        }
    });
}

function controlAnalysisRequest(key, order) {
    $.post(server+"/control", JSON.stringify({"Id":  key, "Order": order}));
}
