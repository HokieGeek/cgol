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
    // console.log("generateRandomSeed()", id, boardSize, cellSize, coverage);
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
          // FIXME: yeah, this is no good
          var x = col * adjWidth;
          var y = row * adjHeight;

          ctx.fillStyle = cellAliveColor;
          // if ((col % 2) == 0) ctx.fillStyle = '#00ff00';
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

    var id = "board-"+key;
    var analysis = analyses[key];

    // Determine new cell size
    analysis.elements.cellSize.width = parseInt($('#cellSize-'+key).val());
    analysis.elements.cellSize.height = analysis.elements.cellSize.width;
    // console.log("cellSize = ", analysis.elements.cellSize.width);

    var canvas = document.getElementById(id);
    if (width > 0 && height > 0) {
      var displayCellWidth = analysis.elements.cellSize.width + cellSpacing;
      var displayCellHeight = analysis.elements.cellSize.height + cellSpacing;
      canvas.width = Math.floor(width) + (displayCellWidth - (Math.floor(width) % displayCellWidth));
      canvas.height = Math.floor(height) + (displayCellHeight - (Math.floor(height) % displayCellHeight));
      // var adjHeight = Math.floor(height);
      // console.log("ADJ % w,j: ", adjWidth % (analysis.elements.cellSize.width + cellSpacing), adjHeight % (analysis.elements.cellSize.height + cellSpacing));
      // console.log("ADJ % w,j: ", adjWidth % (analysis.elements.cellSize.width + cellSpacing), adjHeight % (analysis.elements.cellSize.height + cellSpacing));
      // canvas.width = adjWidth;
      // canvas.height = adjHeight;
      // canvas.width = adjWidth + cellSpacing;
      // canvas.height = adjHeight + cellSpacing;
    }
    analysis.dimensions.Width = Math.floor(canvas.width / (analysis.elements.cellSize.width + cellSpacing));
    analysis.dimensions.Height = Math.floor(canvas.height / (analysis.elements.cellSize.height + cellSpacing));
    // console.log("updateFromInputs(): board w,h = ", analysis.dimensions.Width, analysis.dimensions.Height);

    // Determine cell coverage and rebuild seed
    var coverage = parseInt($('#cellDensity-'+key).val());
    analysis.seed = generateRandomSeed(id, analysis.dimensions, analysis.elements.cellSize, coverage);
    // console.log("updateFromInputs():", analysis.seed);
}

function initBoard(key, padre) {
    var analysis = analyses[key];
    var boardWidth = Math.floor(analysis.dimensions.Width * (analysis.elements.cellSize.width + cellSpacing))
    var boardHeight = Math.floor(analysis.dimensions.Height * (analysis.elements.cellSize.height + cellSpacing));
    // console.log("initBoard():   cell = ", analysis.elements.cellSize.width);
    // console.log("initBoard():  width = ", analysis.dimensions.Width, boardWidth);
    // console.log("initBoard(): height = ", analysis.dimensions.Height, boardHeight);
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
                                            .attr("value", analysis.elements.cellSize.width)
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
      // resize: function(event, ui) { console.log("resizing w,h: ", ui.size.width, ui.size.height); },
      stop: function(event, ui) { updateFromInputs(key, ui.size.width, ui.size.height); }
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
        processing: false,
        running : false,
        updateQueue : [],
        seed : [],
        dimensions : {Width: 120, Height: 80}, //{ Width: 500, Height: 300 },
        // dimensions : {Width: 10, Height: 10}, //{ Width: 500, Height: 300 },
        elements : {
            cellSize : { width: 3, height: 3 },
            currentGeneration : null,
            board : null,
        },
        AddToQueue : function(data) {
                        // console.log("  AddToQueue()", data);
                        // console.log("  AddtoQueue(): ", this.updateQueue.length);
                        // Add each update to the queue
                        for (var i = 0; i < data.Updates.length; i++) {
                            this.updateQueue.push(data.Updates[i]);
                        }

                        if (!this.processing && this.updateQueue.length > 0 && this.updateQueue.length < updateQueueLimit) {
                        // if(this.updateQueue.length < updateQueueLimit && this.running) {
                            this.processing = true;
                            // console.log(">>>> WILL BE PROCESSING <<<<");
                            // TODO: only do this timeout if one doesn't already exist
                            setTimeout($.proxy(this.Processor, this), processingRate_ms);
                        }
                    },
        Processor : function() {
                    // console.log("Process()", key, this);
                    // console.log("Process(): updateQueue = ", this.updateQueue.length);
                    var update = this.updateQueue.shift();

                    if (update != undefined) { // TODO: Why is update sometimes undefined
                        // this.processing = true;
                        // console.log("Processing: ", update.Living.length, update);
                        this.elements.currentGeneration.text(update.Generation);

                        var cellWidth = this.elements.cellSize.width;
                        var cellHeight = this.elements.cellSize.height;
                        var board = this.elements.board;

                        var adjWidth = cellSpacing + cellWidth;
                        var adjHeight = cellSpacing + cellHeight;

                        var numPerRow = board.width / (cellWidth + cellSpacing);
                        var numPerCol = board.height / (cellHeight + cellSpacing);

                        var ctx = board[0].getContext('2d');
                        ctx.save()

                        ctx.setTransform(1, 0, 0, 1, 0, 0);
                        ctx.clearRect(0, 0, ctx.canvas.width, ctx.canvas.height);

                        ctx.fillStyle = cellAliveColor;
                        for (var i = update.Living.length-1; i >= 0; i--) {
                              var x = update.Living[i].X * adjWidth;
                              var y = update.Living[i].Y * adjHeight;
                              // console.log(x,y);

                              ctx.save();
                              ctx.translate(x, y);
                              ctx.fillRect(0, 0, cellWidth, cellHeight);
                              ctx.restore();
                        }
                        ctx.restore();

                        this.processed++;

                        // Keep processing
                        if (this.updateQueue.length > 0 && this.running) {
                            // console.log(".... Will process again ....");
                            setTimeout($.proxy(this.Processor, this), processingRate_ms);
                        } else {
                            this.processing = false;
                            // console.log(">>>> NO MORE PROCESSING <<<< ", this.updateQueue.length, this.running);
                        }
                    } else {
                        // console.log("Update is undefined: num_updates = " + this.updateQueue.length);
                        this.processing = false;
                    }
                },
        Start : function() {
                    this.poller = setInterval(function() {  var startingGen = this.processed + this.updateQueue.length + 1;
                                                            // console.log("### startingGen = ", this.processed, this.updateQueue.length, 1);
                                                            pollAnalysisRequest(this.id,
                                                                                startingGen,
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
                                        .click(function() {
                                            // console.log("Running? ", this.running, this);
                                            if (analyses[key].running) {
                                                stopAnalysis(key);
                                                this.innerHTML = '▶';
                                            } else {
                                                startAnalysis(key);
                                                this.innerHTML = '▮▮';
                                            }
                                        })
                                        // .click(function() {
                                        //             this.poller = setInterval(function() { pollAnalysisRequest(this.id,
                                        //                                                                        this.processed + this.updateQueue.length + 1,
                                        //                                                                        maxPollGenerations) },
                                        //                                                    pollRate_ms);
                                        //             controlAnalysisRequest(this.id, 0);
                                        //             this.running = true;
                                        //         })
                                        .text("▶"))
                /*
                .append($("<span></span>").addClass("analysisControl")
                                        .click(function() { stopAnalysis(key); })
                                        // .click(function() {
                                        //         clearInterval(this.poller);
                                        //         controlAnalysisRequest(key, 1);
                                        //         this.running = false;
                                        //     })
                                        .text("▮▮"))
                                        // .text("⬛"))
                */
                ) // Control

    );

    analyses[key].elements.board = initBoard(key, $("#analysis-"+key));
}

function startAnalysis(key) {
    var analysis = analyses[key];
    analysis.poller = setInterval(function() {  var startingGen = analysis.processed + analysis.updateQueue.length + 1;
                                                // console.log("))) startingGen = ", analysis.processed, analysis.updateQueue.length, 1);
                                                pollAnalysisRequest(analysis.id,
                                                                    startingGen,
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
    // console.log("createNewAnalysisRequest(): ", req);
    $.post(server+"/analyze", JSON.stringify(req))
        .done(callback)
        .fail(function(err) {
            console.log("Got a post error:", err.status, err.responseText);
        });
}

function pollAnalysisRequest(analysisId, startingGen, maxGen) {
    // console.log("pollAnalysisRequest()", analysisId, startingGen, maxGen);
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
