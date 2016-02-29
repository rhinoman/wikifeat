define([
	'views/menuview',
	'views/mapview',
    'views/insertview'
], function(OLMenuView,OLMapView,OLInsertView){

	var OLApp = new Backbone.Marionette.Application();

	OLApp.on("start", function() {
		console.log("The OpenLayers plugin has started.");
	});

	OLApp.getContentView = function(el, contentId){
		console.log("Showing OpenLayers content");
		return new OLMapView({el: el});
	};

	OLApp.getInsertLabel = function(){
        return "<span class='glyphicon glyphicon-globe'></span>&nbsp;OpenLayers";
	};

    OLApp.getInsertView = function(options){
        return new OLInsertView(options);
    };
	
	return OLApp;
});
