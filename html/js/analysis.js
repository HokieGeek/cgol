//  function add
var analyses = {}

function getIdStr(id) {
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
                    .html("&nbsp"))
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

    analyses[data.Id] = data;
}

function updateBoard(data) {
    console.log("updateBoard()", data);

    var idStr = getIdStr(data.Id);
    for (var i = data.Living.length-1; i >= 0; i--) {
        var living = data.Living[i];
        console.log("Living["+i+"]: ", living);
        console.log("   ID:", "#cell-"+idStr+"-"+living.Y+"x"+living.X)
        $("#cell-"+idStr+"-"+living.Y+"x"+living.X).addClass("analysisBoardCellAlive");
    }
    // TODO
}

function updateAnalysis(data) {
    console.log("updateAnalysis()", data);

    var idStr = getIdStr(data.Id);

    var update = data.Updates[0];
    $("#status-"+idStr).text(update.Status);
    $("#generaton-"+idStr).text(update.Generation);

    // TODO: Loop through the updates
    updateBoard(update);
}

function newAnalysisData(data) {
    if (data.Id in analyses) {
        updateAnalysis(data);
    } else {
        createAnalysis(data);
    }
}

function pollAnalyses() {
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
    $.post( "http://localhost:8081/create", 
            JSON.stringify({"Dims":{"Height": 30, "Width": 200}}))
  .done(function( data ) {
      createAnalysis(data);
      pollAnalyses();
  });
}

