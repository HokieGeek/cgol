//  function add
var analyses = {}

function createBoard(data) {
    var span = $("<span></span>");
    for (var i = 0; i < data.Dims.Height; i++) {
        var row = $("<div></div>");
        for(var j = 0; j < data.Dims.Width; j++) {
            row.append($("<span></span>")
                    .attr("id", "cell-"+data.Id+"-"+i+"x"+j)
                    .addClass("analysisBoardCell")
                    .html("&nbsp"))
        }
        span.append(row);
    }
    return span;
}

function createAnalysis(data) {
    console.log("createAnalysis()", data)
    $("#analyses").html(
        $("<div></div>").attr("id", "analysis-"+data.Id.toString(16))

        // ID
        .append($("<div></div>")
                .append("<span>ID: </span>")
                .append($("<span></span>").addClass("analysisId").text(data.Id.toString(16))))
                .addClass("analysisField")

        // Rule
        // .append($("<div></div>")
        //         .append("<span>Rule: </span>")
        //         // .append($("<span></span>").text(data.Id.toString(16)))
        //         .addClass("analysisField").addClass("analysisRule"))

        // Neighbors
        .append($("<div></div>")
                .append("<span>Neighbors: </span>")
                .append($("<span></span>").addClass("analysisNeighbors").text("TODO")))
                .addClass("analysisField")

        // Status
        .append($("<div></div>")
                .append("<span>Status: </span>")
                .append($("<span></span>").addClass("analysisStatus")))
                .addClass("analysisField")

        // Generation
        .append($("<div></div>")
                .append("<span>Generation: </span>")
                .append($("<span></span>").addClass("analysisGeneration")))
                .addClass("analysisField")

        // Living cells
        .append($("<div></div>").attr("class", "analysisBoard").html(createBoard(data)))
    )

    analyses[data.Id] = data
}

function updateBoard(data) {
    // TODO
}

function updateAnalysis(data) {
    console.log("updateAnalysis()", data)
    // TODO: Loop through the updates
}

function newAnalysisData(data) {
    if (data.Id in analyses) {
        updateAnalysis(data)
    } else {
        createAnalysis(data)
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
      pollAnalyses()
  });
}

