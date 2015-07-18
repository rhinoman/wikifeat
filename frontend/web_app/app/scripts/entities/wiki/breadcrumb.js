/** This file is part of WikiFeat.
 *
 *     WikiFeat is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 2 of the License, or
 * (at your option) any later version.
 *
 *     WikiFeat is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 *     You should have received a copy of the GNU General Public License
 * along with WikiFeat.  If not, see <http://www.gnu.org/licenses/>.
 */
/**
 * Created by jcadam on 3/1/15.
 */

define([
    'jquery',
    'underscore',
    'backbone'
], function($,_,Backbone){

    return Backbone.Model.extend({

        defaults: {
            name: "Untitled",
            pageId: "",
            wikiId: "",
            wikiSlug: "",
            pageSlug: "",
            parent: ""
        }

    });

});
