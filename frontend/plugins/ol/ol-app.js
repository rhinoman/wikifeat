define(['jquery',
	'underscore',
	'marionette',
	'views/menuview',
	'views/mapview'
], function($,_,Marionette,OLMenuView,OLMapView){

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
	
	return OLApp;
});
