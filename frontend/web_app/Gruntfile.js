module.exports = function(grunt){
  //Project configuration
  
  grunt.initConfig({
    bower: {
      target: {
        rjsConfig: 'app/scripts/main.js'
      }
    },
    connect: {
      server:{
        options: {
          port: 8000,
          base: 'app',
          keepalive: true
        }
      }
    },
    sass: {
      dist: {
        options: {
          style: 'expanded'
        },
        files: {
          'app/styles/main.css': 'app/styles/main.scss'
        }
      }
    }
  });

  grunt.loadNpmTasks('grunt-bower-requirejs');
  grunt.loadNpmTasks('grunt-contrib-connect');
  grunt.loadNpmTasks('grunt-contrib-sass');
  grunt.registerTask('default', ['bower']);
};
