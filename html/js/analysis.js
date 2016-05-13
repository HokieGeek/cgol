var server = "http://localhost:8081";
var pollRate_ms = "250";
var processingRate_ms = "125";
var maxPollGenerations = 25;
var updateQueueLimit = 100;
var StatusStr = ["Seeded", "Active", "Stable", "Dead"];

var analyses = {};

function getIdStr(id) {
    var idStr = id.toString(16).replace(new RegExp("[/+=]", 'g'), "");
    return idStr.substring(0, idStr.length-1);
}

//////////////////// BOARD CREATOR ////////////////////

var cellSpacing = 1;
var cellAliveColor = '#4863a0'; // FIXME: hmmm
var cellDeadColor = '#e0e0e0';

function generateRandomSeed(id, boardSize, cellSize, coverage) {
    // console.log("generateRandomSeed()", id, boardSize, cellSize, coverage);
    var seed = [];

    var adjWidth = cellSize.width + cellSpacing;
    var adjHeight = cellSize.height + cellSpacing;

    var board = document.getElementById(id);
    var ctx = board.getContext('2d');
    ctx.save();

    ctx.setTransform(1, 0, 0, 1, 0, 0);
    ctx.clearRect(0, 0, board.width, board.height);

    ctx.fillStyle = cellAliveColor;

    var aliveVal = 100-coverage;
    for (var row = boardSize.Height-1; row >= 0; row--) {
      for (var col = boardSize.Width-1; col >= 0; col--) {
        if ((Math.random() * 100) > aliveVal) {
          var x = col * adjWidth;
          var y = row * adjHeight;

          ctx.fillRect(x, y, cellSize.width, cellSize.height);

          seed.push({"X": col, "Y": row});
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

    var canvas = document.getElementById(id);
    if (width > 0 && height > 0) {
      var displayCellWidth = analysis.elements.cellSize.width + cellSpacing;
      var displayCellHeight = analysis.elements.cellSize.height + cellSpacing;
      canvas.width = Math.floor(width) + (displayCellWidth - (Math.floor(width) % displayCellWidth));
      canvas.height = Math.floor(height) + (displayCellHeight - (Math.floor(height) % displayCellHeight));
    }
    analysis.dimensions.Width = Math.floor(canvas.width / (analysis.elements.cellSize.width + cellSpacing));
    analysis.dimensions.Height = Math.floor(canvas.height / (analysis.elements.cellSize.height + cellSpacing));

    // Determine cell coverage and rebuild seed
    var coverage = parseInt($('#cellDensity-'+key).val());
    analysis.seed = generateRandomSeed(id, analysis.dimensions, analysis.elements.cellSize, coverage);
}

function initBoard(key, padre) {
    var analysis = analyses[key];
    var boardWidth = Math.floor(analysis.dimensions.Width * (analysis.elements.cellSize.width + cellSpacing))
    var boardHeight = Math.floor(analysis.dimensions.Height * (analysis.elements.cellSize.height + cellSpacing));
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
                        $("#analysis-"+key+" .analysisControl").css('visibility', 'visible');
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
        dimensions : { Width: 120, Height: 80 },
        elements : {
            cellSize : { width: 3, height: 3 },
            currentGeneration : null,
            board : null,
        },
        AddToQueue : function(data) {
                        // console.log("  AddToQueue()", data);
                        // Add each update to the queue
                        for (var i = 0; i < data.Updates.length; i++) {
                            this.updateQueue.push(data.Updates[i]);
                        }

                        if (!this.processing && this.updateQueue.length > 0 && this.updateQueue.length < updateQueueLimit) {
                            this.processing = true;
                            // TODO: only do this timeout if one doesn't already exist
                            setTimeout($.proxy(this.Processor, this), processingRate_ms);
                        }
                    },
        Processor : function() {
                    // console.log("Process()", key, this);
                    var update = this.updateQueue.shift();

                    if (update != undefined) {
                        // console.log("Processing: ", update.Living.length, update);
                        this.elements.currentGeneration.text(update.Generation);

                        var cellWidth = this.elements.cellSize.width;
                        var cellHeight = this.elements.cellSize.height;

                        var adjWidth = cellWidth + cellSpacing;
                        var adjHeight = cellHeight + cellSpacing;

                        var board = this.elements.board;
                        var ctx = board[0].getContext('2d');
                        ctx.save()

                        ctx.setTransform(1, 0, 0, 1, 0, 0);
                        ctx.clearRect(0, 0, ctx.canvas.width, ctx.canvas.height);

                        ctx.fillStyle = cellAliveColor;
                        for (var i = update.Living.length-1; i >= 0; i--) {
                              var x = update.Living[i].X * adjWidth;
                              var y = update.Living[i].Y * adjHeight;

                              ctx.fillRect(x, y, cellWidth, cellHeight);
                        }
                        ctx.restore();

                        this.processed++;

                        // Keep processing
                        if (this.updateQueue.length > 0 && this.running) {
                            setTimeout($.proxy(this.Processor, this), processingRate_ms);
                        } else {
                            this.processing = false;
                        }
                    } else {
                        this.processing = false;
                    }
                },
        Start : function() {
                    var me = this;
                    this.poller = setInterval(function() {  var startingGen = me.processed + me.updateQueue.length + 1;
                                                                // console.log("))) startingGen = ", this.processed, this.updateQueue.length, 1);
                                                                pollAnalysisRequest(me.id,
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
                                            var analysis = analyses[key];
                                            if (analysis.running) {
                                                analysis.Stop();
                                                this.innerHTML = '▶';
                                            } else {
                                                analysis.Start();
                                                this.innerHTML = '▮▮';
                                            }
                                        })
                                        .text("▶"))
                ) // Control

    );

    analyses[key].elements.board = initBoard(key, $("#analysis-"+key));
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
            if (err.status != 0) {
                console.log("Got a post error:", err.status, err.responseText);
            }
        });
}
