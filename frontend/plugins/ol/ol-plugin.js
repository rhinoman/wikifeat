require.config({
	baseUrl: "/app/plugin/OpenLayers/resource/",
	paths: {
	},
	packages: []
});

var OpenLayers = {

    start: function(started){
        require(['ol-app'], function(OLApp){
            OLApp.start();
            OpenLayers.getContentView = function(el, contentId){
                return OLApp.getContentView(el, contentId);
            };
            OpenLayers.getInsertLabel = function(){
                return OLApp.getInsertLabel();
            };
            OpenLayers.getInsertView = function(options){
                return OLApp.getInsertView(options);
            };
            started.resolve();
	    });
    }

};
