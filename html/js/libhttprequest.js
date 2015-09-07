var libHTTPRequest_Requests = [];
var libHTTPRequest_Constants = {
    MSXML_VERS : ["MSXML2.XMLHttp.6.0","MSXML2.XMLHttp.5.0",
                  "MSXML2.XMLHttp.4.0","MSXML2.XMLHttp.3.0", 
                  "MSXML2.XMLHttp","Microsoft.XMLHttp"],
    useActiveX : (typeof ActiveXObject != "undefined"),
    useDom : (document.implementation && document.implementation.createDocument),
    useXmlHttp : (typeof XMLHttpRequest != "undefined")
};
function libHTTPRequest_DataManager() {
    //this.Scripts = [ "http://"+window.location.hostname+"/sss/XDRequest.php" ];
    this.Scripts = [ "http://hokiegeek.net/sss/XDRequest.php" ];
    this.StoreLocation = "http://hokiegeek.net/sss/gen";
    //this.StoreLocation = "http://"+window.location.hostname+"/sss/gen";
};

function libHTTPRequest () {
    this.URL = null;
    this.Options = null;
    this.Callback = null;

    this.Request = null;
    this.Store = false;
    this.StoreFileName = "libHTTPRequest"; 
    this.DataManager = new libHTTPRequest_DataManager();
    this.DataScript = null;

    this._load();
};
libHTTPRequest.prototype._load = function() {
    var dm = new libHTTPRequest_DataManager();
    // Ping each script until you find an active one
    for (var i = 0; i < dm.Scripts.length; i++) {
        if (this.ping(dm.Scripts[i])) {
            this.DataScript = dm.Scripts[i];
            break;
        }
    }
}
libHTTPRequest.prototype._XmlHttpReq = function() {
    if (libHTTPRequest_Constants.useXmlHttp) {
        return new XMLHttpRequest();
    } else if (libHTTPRequest_Constants.useActiveX) {
        for (var i = 0; i < libHTTPRequest_Constants.MSXML_VERS.length; i++) {
            try {
                return new ActiveXObject(libHTTPRequest_Constants.MSXML_VERS[i]);
            } catch (oError) {
                // Do nada
            }
        }
    }
    throw new Error("XMLHttp object could not be created.");
}

libHTTPRequest._onreadystatechange = function(reqIdx){
    var request = libHTTPRequest_Requests[reqIdx];
    //console.log("_onreadystatechange: request = ", request);
    var r = request.Request;
    var cb = request.Callback;
    var bk = request.Store;
    var bk_name = request.StoreFileName;
      if (r != null && r.readyState == "4") {
        //console.log("libHTTPRequest._onreadystatechange("+reqIdx+"): ", r.status);
    
        switch(r.status) {
          case 200: // OK
        case 201: // Created
        case 202: // Accepted
        case 203: // Non-Auth Info
        case 204: // No Content
        case 205: // Reset Content
        case 206: // Partial Content
        case 207: // Multi-Status
              var ret = null;
              if (r.responseXML && r.responseXML.length > 0) {
                  ret = r.responseXML;
                //console.log("GOT AN XML ANSWER: ", r.responseXML);
              } else {
                ret = r.responseText;
                //console.log("GOT A JSON ANSWER: ", r.responseText);
              }
              //if (ret != null && cb != null && ret != "Token invalid") {
              if (ret != null && cb != null) {
                //console.log("_onreadystatechange: ", request.URL, ret, cb);
                  cb(ret);
            } else if (bk) {
                libHTTPRequest._loadFromBackup(reqIdx);
            }
            break;
          case 304: // NOT MODIFIED
          case 400: // BAD REQUEST
          case 401: // NOT AUTHORIZED
          case 403: // FORBIDDEN
          case 404: // NOT FOUND
          case 500: // INTERNAL SERVER ERROR
          case 502: // BAD GATEWAY
          case 503: // SERVICE UNAVAILABLE
          default: 
              //console.log("CRAP!!");
              if (bk)
                libHTTPRequest._loadFromBackup(reqIdx);
            break;
        }
      }
}
libHTTPRequest._loadFromBackup = function (reqIdx) {
    var request = libHTTPRequest_Requests[reqIdx];
    var dm = new libHTTPRequest_DataManager();

    /*this.StoreLocation = "http://"+window.location.hostname+"/sss/gen";*/
    this.StoreLocation = "http://hokiegeek.net/sss/gen";
    var file = dm.StoreLocation+"/"+request.StoreFileName+".json";
    //console.log("libHTTPRequest._loadFromBackup("+reqIdx+"): ", file);
    var r = request.getXmlRequest("GET", file, reqIdx);
    request.Request = r;
      r.send(null);
}

/** PUBLIC METHODS **/
libHTTPRequest.prototype.ping = function(url) {
    var r = this._XmlHttpReq();
    var alive = false;
      r.open("GET", url, false);
    try {
        r.send(null);
        if (r != null && r.readyState == "4") {
            if (r.status >= 200 && r.status <= 207)
                alive = true;
        }
    } catch (e) {
    }
    return alive;
}

libHTTPRequest.prototype.getXmlRequest = function(type, u, reqIdx) {
    var r = this._XmlHttpReq();
      r.onreadystatechange = function() { libHTTPRequest._onreadystatechange(reqIdx); };
      r.open(type, u, true);
      return r;
}

libHTTPRequest.prototype.getData = function(u, o, cb, bk) {
    //console.log("libHTTPRequest.getData("+u+", "+o+", cb, bk)");
    if (u == undefined) return;

    // Store 
    this.URL = u;
    this.Options = o;
    this.Callback = cb;
    this.Store = bk;
    if (bk) {
        var modified_url = u;
        if (u.indexOf("?") != -1) modified_url = u.substring(0, u.indexOf("?"));
        modified_url = modified_url.replace(/https?:\/\//g, "")
                                   .replace(/(\/|\.|\-|=|&)/g, "_");
        this.StoreFileName += "__"+modified_url;
        o += "&store="+this.StoreFileName;
    } 

    // Make the request
    libHTTPRequest_Requests.push(this);
    var reqIdx = libHTTPRequest_Requests.length-1;
    //console.log("libHTTPRequest.getData(", u, o, "cb,", bk,"): ", reqIdx);
    //console.log(">>>> this.DataScript = ", this.DataScript);
    var r = this.getXmlRequest("GET", this.DataScript+"?u="+u+"&"+o, reqIdx);
    this.Request = r;
      r.send(null); 
}

libHTTPRequest.prototype.sendData = function(u, p) {
    this.URL = u;
    this.Options = p;

    var params = "";
    for (var i = 0; i < p.length; i++) {
        if (i > 0) params += "&";
        params += p[i];
    }

    // Make the request
    libHTTPRequest_Requests.push(this);
    var reqIdx = libHTTPRequest_Requests.length-1;
    //console.log("libHTTPRequest.sendData(", u, ",", p, "): ", reqIdx);
    var r = this.getXmlRequest("POST", this.DataScript+"?u="+u, reqIdx);
    r.setRequestHeader("Content-type", "application/x-www-form-urlencoded");
    //r.setRequestHeader("Access-Control-Allow-Origin", "*");
    this.Request = r;
      r.send(params); 
}

