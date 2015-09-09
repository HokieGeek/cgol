//  function add
var analyses = {};
var analysesIdMap = {};
var updateQueues = {}

function getIdStr(id) {
    // console.log("getIdStr(): ", id);
    var idStr = id.toString(16).replace(new RegExp("[/+=]", 'g'), "");
    return idStr.substring(0, idStr.length-1);
}

function createBoard(data) {
    var idStr = getIdStr(data.Id);

    var span = $("<span></span>");
    for (var i = 0; i < data.Dims.Height; i++) {
        var row = $("<div></div>");
        for(var j = 0; j < data.Dims.Width; j++) {
            row.append($("<span></span>")
                    .attr("id", "cell-"+idStr+"-"+i+"x"+j)
                    .addClass("analysisBoardCell")
                    )
        }
        span.append(row);
    }
    return span;
}

function createAnalysis(data) {
    // console.log("createAnalysis()", data)

    var idStr = getIdStr(data.Id);

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
                                        .click(function() { controlAnalysis(data.Id,0) })
                                        .text("Start"))
                .append($("<span></span>").addClass("analysisControl")
                                        .click(function() { controlAnalysis(data.Id,1) })
                                        .text("Stop"))
                )

        // Living cells
        .append($("<div></div>").attr("class", "analysisBoard").attr("id", "board-"+idStr)
                .html(createBoard(data)))
    );

    analysesIdMap[idStr] = data.Id;
    analyses[data.Id] = [];

    updateQueues[data.Id] = { "timeout" : null, "updates" : [] };
}

function updateBoard(idStr, data) {
    // console.log("updateBoard()", data);

    var idPrefix = "#cell-"+idStr+"-";
    for (var i = data.Changes.length-1; i >= 0; i--) {
        var changed = data.Changes[i];

        switch (changed.Change) {
        case 0: // Born
            $(idPrefix+changed.Y+"x"+changed.X).addClass("analysisBoardCellAlive");
            break;
        case 1: // Died
            $(idPrefix+changed.Y+"x"+changed.X).removeClass("analysisBoardCellAlive");
            break;
        }
    }
}

/*
function getLivingMap(living) {
    var map = {};

    for (var i = living.length-1; i >= 0; i--) {
        var loc = living[i]
        if (!(loc.Y in map)) {
            map[loc.Y] = {};
        }
        if (!(loc.X in map[loc.Y])) {
            map[loc.Y][loc.X] = true;
        }

    }

    return map;
}

function updateBoard2(idStr, data) {
    console.log("   updateBoard2()", idStr, data.Living);
    livingMap = getLivingMap(data.Living);

    // var idStr = getIdStr(data.Id);
    var board = $("<span></span>");
    for (var y = 0; y < data.Dims.Height; y++) {
        var row = $("<div></div>");
        for(var x = 0; x < data.Dims.Width; x++) {
            var cell = $("<span></span>")
                    .attr("id", "cell-"+idStr+"-"+y+"x"+x)
                    .addClass("analysisBoardCell")
            if (y in livingMap && x in livingMap[y]) {
                cell.addClass("analysisBoardCellAlive")
            }
            row.append(cell);
        }
        board.append(row);
    }
    return board;
}
*/

var StatusStr = ["Seeded", "Active", "Stable", "Dead"]

function processAnalysisUpdate(idStr, gen) {
    console.log("   processAnalysisUpdate()", idStr, gen);

    var id = analysesIdMap[idStr];
    var update = analyses[id][gen];
    // console.log("   update = ", update);

    // $("#status-"+idStr).text(StatusStr[update.Status]);
    $("#generation-"+idStr).html(update.Generation);

    updateBoard(idStr, update);
    // console.log("Adding created row:", $("#board-"+idStr));
    // console.log("New board:", newBoard);
    // var newBoard = updateBoard2(idStr, update);
    // $("#board-"+idStr).html(newBoard);
}

function updateAnalysis(data) {
    // console.log("  updateAnalysis()", data);
    for (var i = 0; i < data.Updates.length; i++) {
        var idStr = getIdStr(data.Id);

        analyses[data.Id].push(data.Updates[i]);

        // setTimeout(function() { processAnalysisUpdate(idStr, i); }, (i * 500));
        setTimeout(function() { eval("processAnalysisUpdate('"+idStr+"', "+i+")"); }, (i * 500));
        // eval("setTimeout(function() { processAnalysisUpdate('"+idStr+"', "+i+"); }, "+(i * 500)+");")
    }
}

function newAnalysisData(data) {
    if (data.Id in analyses) {
        updateAnalysis(data);
    } else {
        createAnalysis(data);
    }
}

function pollAnalyses() {
    // console.log("pollAnalyses()")
    for (var key in analyses) {
        console.log("pollAnalyses(): ", key)
        $.post( "http://localhost:8081/poll", 
            JSON.stringify({"Id": key, "StartingGeneration": analyses[key].length}))
    .done(function( data ) {
        newAnalysisData(data);
    });
    }
}

function pollAnalysis(key, startingGen) {
        console.log("pollAnalysis(): ", key, startingGen)
        $.post( "http://localhost:8081/poll", 
            JSON.stringify({"Id": key, "StartingGeneration": startingGen}))
    .done(function( data ) {
        newAnalysisData(data);
    });
}

function createNewAnalysis() {
    $.post( "http://localhost:8081/analyze", 
            JSON.stringify({"Dims":{"Height": 100, "Width": 200}}))
  .done(function( data ) {
      createAnalysis(data);
      // setInterval(pollAnalyses, 1500) // setInterval
  });
}

function controlAnalysis(key, order) {
    $.post( "http://localhost:8081/control", 
            JSON.stringify({"Id":  key, "Order": order}))
    switch order {
    case 0: // Start
        // TODO: create the polling call updateQueues[key].poller = setInterval(function() { pollAnalysis(key, analyses[key].length) }, 1500);
        break;
    case 1: // Stop
        // TODO: clearInterval(updateQueues[key].poller)
        break;
    }
}
