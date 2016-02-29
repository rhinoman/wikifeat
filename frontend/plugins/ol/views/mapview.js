define(['jquery',
	'underscore',
        'marionette',
        'ol'
], function($,_,Marionette,ol){

	return Marionette.ItemView.extend({
		id: 'ol-map-view',
		events: {},

		template: _.template('<div id="mapContainer" class="map"></div>'),

		initialize: function(){
			console.log("Initializing OL Map view");
			this.domReady = $.Deferred();
		},
	         
                // OpenLayers is very DOM dependent.  We MUST make sure Marionette has actually drawn
		// The freaking DOM elements (not just 'rendered' -- thanks for the confusion Marionette docs)
		// to the screen/page before we try to render the map.
		onBeforeRender: function(){
			var self = this;
			var timeoutId = setInterval(function(){
				if(this.$("#mapContainer").get()){
					self.domReady.resolve(this.$("#mapContainer").get());
					clearTimeout(timeoutId);
				}
			}, 100);
		},

		onRender: function(){
			var self = this;
			self.domReady.done(function(){
				self.drawMap();
			});
		},

		drawMap: function(){
			var map = new ol.Map({
				target: 'mapContainer',
				layers: [
					new ol.layer.Tile({
						source: new ol.source.MapQuest({layer: 'osm'})
					})
				],
				view: new ol.View({
					center: ol.proj.transform([37.41, 8.82], 'EPSG:4326', 'EPSG:3857'),
					zoom: 4
				})
			});
		}
		
	});

});
