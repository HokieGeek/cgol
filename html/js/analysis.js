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

// var cellWidth = 3;
// var cellHeight = 3;
var cellSpacing = 1;
var cellAliveColor = '#4863a0'; // FIXME: hmmm
var cellDeadColor = '#e0e0e0';

function generateRandomSeed(id, boardSize, cellSize, coverage) {
    // console.log("generateRandomSeed()", id, cellWidth, cellHeight, coverage);
    var seed = [];

    var cellWidth = cellSize.width;
    var cellHeight = cellSize.height;

    var adjWidth = cellSpacing + cellWidth;
    var adjHeight = cellSpacing + cellHeight;

    var board = document.getElementById(id);
    var ctx = board.getContext('2d');
    ctx.clearRect(0, 0, board.width, board.height);

    ctx.save();

    var aliveVal = 100-coverage;
    for (var row = boardSize.Height-1; row >= 0; row--) {
      for (var col = boardSize.Width-1; col >= 0; col--) {
        if ((Math.random() * 100) > aliveVal) {
          ctx.save();
          ctx.fillStyle = cellAliveColor;
          // FIXME: yeah, this is no good
          var x = col * adjWidth;
          var y = row * adjHeight;
          // console.log(x,y);
          ctx.translate(x,y);
          // ctx.translate(Math.round(col * adjWidth), Math.round(row * adjHeight));
          ctx.fillRect(0, 0, cellWidth, cellHeight);
          ctx.restore();

          seed.push({"X": col, "Y": row});
        }
      }
    }
    ctx.restore();

    return seed;
}
function generateRandomSeed2(id, boardSize, cellSize, coverage) {
    // console.log("generateRandomSeed()", id, cellWidth, cellHeight, coverage);
    var seed = [];

    var cellWidth = cellSize.Width;
    var cellHeight = cellSize.Height;

    var adjWidth = cellSpacing + cellWidth;
    var adjHeight = cellSpacing + cellHeight;

    var board = document.getElementById(id);
    var numPerRow = board.width / (cellWidth + cellSpacing);
    var numPerCol = board.height / (cellHeight + cellSpacing);

    var ctx = board.getContext('2d');
    ctx.clearRect(0, 0, board.width, board.height);

    ctx.save();

    var aliveVal = 100-coverage;
    for (var row = numPerCol-1; row >= 0; row--) {
      for (var col = numPerRow-1; col >= 0; col--) {
        if ((Math.random() * 100) > aliveVal) {
          var x = Math.round(col * adjWidth);
          var y = Math.round(row * adjHeight);

          ctx.save();
          ctx.fillStyle = cellAliveColor;
          ctx.translate(x, y);
          ctx.fillRect(0, 0, cellWidth, cellHeight);
          ctx.restore();

          seed.push({"X": x, "Y": y});
        }
      }
    }
    ctx.restore();

    return seed;
}

function updateFromInputs(key, width, height) {
    // console.log("updateFromInputs():", key, width, height);
    var saneWidth = Math.floor(width);
    var saneHeight = Math.floor(height);

    var id = "board-"+key;
    var analysis = analyses[key];
    if (width > 0 && height > 0) {
      var canvas = document.getElementById(id);
      canvas.width = saneWidth;
      canvas.height = saneHeight;

      analysis.dimensions.Width = Math.ceil(saneWidth / analysis.elements.cellSize.width);
      analysis.dimensions.Height = Math.ceil(saneHeight / analysis.elements.cellSize.height);
    }

    //cellWidth = parseInt(document.getElementById('cellSize').value);
    analysis.elements.cellSize.width = parseInt($('#cellSize-'+key).val());
    analysis.elements.cellSize.height = analysis.elements.cellSize.width;
    var coverage = parseInt($('#cellDensity-'+key).val());
    analysis.seed = generateRandomSeed(id, analysis.dimensions, analysis.elements.cellSize, coverage);
    // console.log("updateFromInputs():", analysis.seed);
}

function initBoard(key, padre) {
    var analysis = analyses[key];
    var boardWidth = Math.floor(analysis.dimensions.Width * analysis.elements.cellSize.width)
    var boardHeight = Math.floor(analysis.dimensions.Height * analysis.elements.cellSize.height);
    var board = $("<canvas></canvas>").attr("id", "board-"+key)
                                      .addClass("analysisBoard ui-widget-content")
                                      .width(boardWidth)
                                      .height(boardHeight)

    padre.append(
        $("<span></span>")
            .append(board)
            .append($("<br/>"))
            .append($("<span></span>").attr("id", "boardControls-"+key).addClass("newBoardControls")
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
                        $("#boardControls-"+key).remove();
                        createNewAnalysisRequest({"Dims": analyses[key].dimensions, "Pattern": 0, "Seed": analyses[key].seed},
                                                    function( data ) {
                                                        lifeIdtoAnalysisKey[data.Id] = key;
                                                        analyses[key].id = data.Id;
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

    updateFromInputs(key, boardWidth, boardHeight);

    return board;
}

////////////////////  ANALYSIS ////////////////////

// Create an analysis blah blah blah
// Fire the request and include a closure for the return id
// From closure, continue filling out the analysis and add to the map

var lifeIdtoAnalysisKey = {};

function analysisKey() {
    var i = 0;
    return function() {
        return ++i;
    }
}
var getNextKey = analysisKey();

var blah = 0;
function createAnalysis() {
    var key = getNextKey();
    analyses[key] = {
        id: null, // TODO: how about this be the lifed id but this client can keep track of them with its own ID     <<-------
        poller : null,
        processed : 0,
        running : false,
        updateQueue : [],
        seed : [],
        dimensions : {Width: 120, Height: 80}, //{ Width: 500, Height: 300 },
        elements : {
            cellSize : { width: 3, height: 3 },
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
                    // console.log("Process()", key, this);
                    var update = this.updateQueue.shift();

                    if (update != undefined) { // TODO: Why is update sometimes undefined
                        console.log("Processing: ", update.Living.length, update);
                        this.elements.currentGeneration.text(update.Generation);

                        var cellWidth = this.elements.cellSize.width;
                        var cellHeight = this.elements.cellSize.height;
                        var board = this.elements.board;

                        var adjWidth = cellSpacing + cellWidth;
                        var adjHeight = cellSpacing + cellHeight;

                        var numPerRow = board.width / (cellWidth + cellSpacing);
                        var numPerCol = board.height / (cellHeight + cellSpacing);

                        var ctx = board[0].getContext('2d');
                        ctx.clearRect(0, 0, board.width, board.height);

                        // FIXME: go back to doing changed....
                        ctx.save();
                        /*
                        for (var i = update.Living.length-1; i >= 0; i--) {
                              // var x = col * adjWidth;
                              // var y = row * adjHeight;
                              var x = update.Living[i].X * adjWidth;
                              var y = update.Living[i].Y * adjHeight;
                              // console.log(x,y);

                              ctx.save();
                              ctx.fillStyle = "#ff0000";
                              ctx.translate(x, y);
                              // ctx.translate(update.Living.X, update.Living.Y);
                              ctx.fillRect(0, 0, cellWidth, cellHeight);
                              ctx.restore();
                        }
                        */
                        for (var i = update.Changes.length-1; i >= 0; i--) {
                              // var x = col * adjWidth;
                              // var y = row * adjHeight;
                              var x = update.Changes[i].X * adjWidth;
                              var y = update.Changes[i].Y * adjHeight;
                              // console.log(x,y);

                              ctx.save();
                              switch (update.Changes[i].Change) {
                                  case 0: // Born
                                        ctx.fillStyle = "#ff0000";
                                        break;
                                  case 1: // Dead
                                        ctx.fillStyle = cellDeadColor;
                                        break;
                              }
                              ctx.translate(x, y);
                              // ctx.translate(update.Living.X, update.Living.Y);
                              ctx.fillRect(0, 0, cellWidth, cellHeight);
                              ctx.restore();
                        }
                        ctx.restore();

                        this.processed++;

                        // Keep processing
                        // if (this.updateQueue.length > 0) {
                        if (blah <= 3 && this.updateQueue.length > 0 && this.updateQueue.length < updateQueueLimit && this.running) {
                            blah++;
                            setTimeout($.proxy(this.Processor, this), processingRate_ms);
                        }
                    } else {
                        console.log("Update is undefined");
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
        .append($("<div style='height: 40px'></div>") // TODO: WTF
                .append($("<span></span>").addClass("analysisControl")
                                        .click(function() { startAnalysis(key); })
                                        // .click(function() {
                                        //             this.poller = setInterval(function() { pollAnalysisRequest(this.id,
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

    analyses[key].elements.board = initBoard(key, $("#analysis-"+key));
}

function startAnalysis(key) {
    var analysis = analyses[key];
    analysis.poller = setInterval(function() { pollAnalysisRequest(analysis.id,
                                                                   analysis.processed + analysis.updateQueue.length + 1,
                                                                   maxPollGenerations) },
                                  pollRate_ms);
    controlAnalysisRequest(analysis.id, 0);
    analysis.running = true;
}

function stopAnalysis(key) {
    var analysis = analyses[key];
    clearInterval(analysis.poller);
    controlAnalysisRequest(analysis.id, 1);
    analysis.running = false;
}

//////////////////// REQUESTORS ////////////////////

function createNewAnalysisRequest(req, callback) {
    console.log("createNewAnalysisRequest(): ", req);
    $.post(server+"/analyze", JSON.stringify(req))
        .done(callback)
        .fail(function(err) {
            console.log("Got a post error:", err.status, err.responseText);
        });
}

function pollAnalysisRequest(analysisId, startingGen, maxGen) {
    console.log("pollAnalysisRequest()", analysisId, startingGen, maxGen);
    // console.log("analyses: ", analyses);
    $.post(server+"/poll", JSON.stringify({"Id": analysisId, "StartingGeneration": startingGen, "NumMaxGenerations": maxGen}))
    .done(function( data ) {
        var id = lifeIdtoAnalysisKey[data.Id];
        if (id in analyses) {
            analyses[id].AddToQueue(data);
        } else {
            console.log("Got update for unknown analysis: ", data.Id)
        }
    });
}

function controlAnalysisRequest(key, order) {
    $.post(server+"/control", JSON.stringify({"Id":  key, "Order": order}))
        .fail(function(err) {
            console.log("Got a post error:", err.status, err.responseText);
        });
}
