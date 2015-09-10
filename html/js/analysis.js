var server = "http://localhost:8081";
var pollRate_ms = "1000";
var processingRate_ms = "250";
var StatusStr = ["Seeded", "Active", "Stable", "Dead"];

var analyses = {};

function getIdStr(id) {
    var idStr = id.toString(16).replace(new RegExp("[/+=]", 'g'), "");
    return idStr.substring(0, idStr.length-1);
}

function createAnalysis(data) {
    var idStr = getIdStr(data.Id);

    // TODO: EACH UPDATE SHOULD HAVE ITS OWN ID

    analyses[idStr] = {
        id: data.Id,
        idAsStr: idStr,
        poller : null,
        generations : [],
        updateQueue : [],
        AddToQueue : function(data) {
                        // console.log("  AddToQueue()", data);
                        // Add each update to the queue
                        for (var i = 0; i < data.Updates.length; i++) {
                            this.updateQueue.push(data.Updates[i])
                        }
                        setTimeout($.proxy(this.Processor, this), processingRate_ms);
                    },
        Processor : function() {
                    // console.log("Process()", this);

                    var update = this.updateQueue.shift();

                    $("#generation-"+this.idAsStr).html(update.Generation);

                    var idPrefix = "#cell-"+this.idAsStr+"-";
                    for (var i = update.Changes.length-1; i >= 0; i--) {
                        var changed = update.Changes[i];

                        switch (changed.Change) {
                        case 0: // Born
                            $(idPrefix+changed.Y+"x"+changed.X).addClass("analysisBoardCellAlive");
                            break;
                        case 1: // Died
                            $(idPrefix+changed.Y+"x"+changed.X).removeClass("analysisBoardCellAlive");
                            break;
                        }
                    }

                    this.generations.push(update)

                    // Keep processing
                    if (this.updateQueue.length > 0) {
                        setTimeout($.proxy(this.Processor, this), processingRate_ms);
                    }
                },
        Start : function() {
                    this.poller = setInterval(function() { pollAnalysisRequest(this.id, this.generations.length+this.updateQueue.length) },
                                                         pollRate_ms);
                    controlAnalysisRequest(this.id, 0);
                },
        Stop : function() {
                    clearInterval(this.poller);
                    controlAnalysisRequest(this.id, 1);
                }
        /*
        */

    };

    // Create the dom entity
    // TODO: consider, perhaps, a map with each cell element for quicker updating?

    var board = $("<span></span>");
    for (var i = 0; i < data.Dims.Height; i++) {
        var row = $("<div></div>");
        for(var j = 0; j < data.Dims.Width; j++) {
            row.append($("<span></span>")
                    .attr("id", "cell-"+idStr+"-"+i+"x"+j)
                    .addClass("analysisBoardCell")
                    )
        }
        board.append(row);
    }

    $("#analyses").html(
        $("<div></div>").attr("id", "analysis-"+idStr)

        // ID
        .append($("<div></div>").addClass("analysisField")
                .append("<span>ID: </span>")
                .append($("<span></span>").addClass("analysisId").text(idStr)))


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
                .append("<span>Generation: </span>")
                .append($("<span></span>").text("0")
                    .attr("id", "generation-"+idStr)
                    .addClass("analysisGeneration")))

        // Control
        .append($("<div></div>")
                .append($("<span></span>").addClass("analysisControl")
                                        .click(function() { startAnalysis(idStr) })
                                        // .click($.proxy(this.Start, this))
                                        .text("Start"))
                .append($("<span></span>").addClass("analysisControl")
                                        .click(function() { stopAnalysis(idStr) })
                                        // .click($.proxy(this.Stop, this))
                                        .text("Stop"))
                )

        // Create the board
        .append($("<div></div>").attr("class", "analysisBoard").attr("id", "board-"+idStr)
                .html(board))
    );
}

function startAnalysis(idStr) {
    analyses[idStr].poller = setInterval(function() { pollAnalysisRequest(analyses[idStr].id,
                                                     analyses[idStr].generations.length+analyses[idStr].updateQueue.length) },
                                         pollRate_ms);
    controlAnalysisRequest(analyses[idStr].id, 0);
}

function stopAnalysis(idStr) {
    clearInterval(analyses[idStr].poller);
    controlAnalysisRequest(analyses[idStr].id, 1);
}

//////////////////// REQUESTORS ////////////////////

function createNewAnalysisRequest() {
    $.post(server+"/analyze", JSON.stringify({"Dims":{"Height": 100, "Width": 200}})) // FIXME
    .done(function( data ) {
        createAnalysis(data);
        pollAnalysisRequest(data.Id, 0);
    });
}

function pollAnalysisRequest(key, startingGen) {
    $.post(server+"/poll", JSON.stringify({"Id": key, "StartingGeneration": startingGen}))
    .done(function( data ) {
        var idStr = getIdStr(data.Id);
        if (idStr in analyses) {
            analyses[idStr].AddToQueue(data);
        } else {
            console.log("Got update for unknown analysis")
            // createAnalysis(data);
        }
    });
}

function controlAnalysisRequest(key, order) {
    $.post(server+"/control", JSON.stringify({"Id":  key, "Order": order}));
}
