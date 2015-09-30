require.config({
	baseUrl: "/app/plugin/OpenLayers/resource/",
	paths: {
	},
	packages: [],
});

/*var OpenLayers = {

    start: function(){
        require(['ol-app'], function(OLApp){
            OLApp.start();
        });
    },
    getContentView: function(el, contentId){
        return require(['ol-app'], function(OLApp){
	    return OLApp.getContentView(el, contentId);
	});
    }
};*/

var OpenLayers = {

    start: function(){
        require(['ol-app'], function(OLApp){
	    OLApp.start();
	    OpenLayers.getContentView = function(el, contentId){
                return OLApp.getContentView(el, contentId);
	    }
	});
    }

};
