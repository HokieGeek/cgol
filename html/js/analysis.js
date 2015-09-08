//  function add
var analyses = {}

function getIdStr(id) {
    console.log("getIdStr(): ", id);
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
                    // .html("&nbsp")
                    )
        }
        span.append(row);
    }
    return span;
}

function createAnalysis(data) {
    console.log("createAnalysis()", data)

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
        .append($("<div></div>").addClass("analysisField")
                .append("<span>Neighbors: </span>")
                .append($("<span></span>").addClass("analysisNeighbors").text("TODO")))

        // Status
        .append($("<div></div>").addClass("analysisField")
                .append("<span>Status: </span>")
                .append($("<span></span>").text("Unknown")
                    .attr("id", "status-"+idStr)
                    .addClass("analysisStatus")))

        // Generation
        .append($("<div></div>").addClass("analysisField")
                .append("<span>Generation: </span>")
                .append($("<span></span>").text("0")
                    .attr("id", "generation-"+idStr)
                    .addClass("analysisGeneration")))

        // Living cells
        .append($("<div></div>").attr("class", "analysisBoard").attr("id", "board-"+idStr)
                .html(createBoard(data)))
    );

    analyses[data.Id] = [];
}

function updateBoard(idStr, data) {
    console.log("updateBoard()", data);

    // var idStr = getIdStr(data.Id);
    for (var i = data.Living.length-1; i >= 0; i--) {
        var changed = data.Changes[i];
        // console.log("Changes["+i+"]: ", changed);
        // console.log("   ID:", "#cell-"+idStr+"-"+living.Y+"x"+living.X)

        switch (changed.Change) {
        case 0: // Born
            $("#cell-"+idStr+"-"+changed.Y+"x"+changed.X).addClass("analysisBoardCellAlive");
            break;
        case 1: // Died
            $("#cell-"+idStr+"-"+changed.Y+"x"+changed.X).removeClass("analysisBoardCellAlive");
            break;
        }
    }
}

var StatusStr = ["Seeded", "Active", "Stable", "Dead"]

function processAnalysisUpdate(aId, gen) {
    console.log("processAnalysisUpdate()", aId, gen);

    console.log("analyses", analyses);

    var id = parseInt(aId, 16);
    console.log("    id: ", id)
    var update = analyses[id][gen];
    console.log("    update = ", update)

    var idStr = getIdStr(update.Id);

    $("#status-"+idStr).text(StatusStr[update.Status]);
    $("#generaton-"+idStr).text(update.Generation);

    updateBoard(idStr, update);
}

function updateAnalysis(data) {
    console.log("updateAnalysis()", data, data.Updates.length);

    for (var i = 0; i < data.Updates.length; i++) {
        console.log("updateAnalysis(): ", i);
        analyses[data.Id].push(data.Updates);
        // scheduleUpdateProcessing(data.Id, i, (i * 1000));
        // setTimeout(function() { eval("processAnalysisUpdate("+data.Id+", "+i+")"); }, (i * 1000));
        // console.log("WANT TO CALL: setTimeout(function() { processAnalysisUpdate('"+data.Id+"', "+i+"); }, "+(i * 1000)+");")
        eval("setTimeout(function() { processAnalysisUpdate('"+data.Id+"', "+i+"); }, "+(i * 1000)+");")
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
    console.log("pollAnalyses()")
    for (var key in analyses) {
        console.log("pollAnalyses(): ", key)
        $.post( "http://localhost:8081/poll", 
            JSON.stringify({"Id":key}))
    .done(function( data ) {
        newAnalysisData(data);
    });
    }
}

function createNewAnalysis() {
    $.post( "http://localhost:8081/analyze", 
            JSON.stringify({"Dims":{"Height": 100, "Width": 200}}))
  .done(function( data ) {
      createAnalysis(data);
      setTimeout(pollAnalyses, 2000) // setInterval
  });
}

