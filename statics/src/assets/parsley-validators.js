const
    $ = require('jquery'),
    bluebird = require('bluebird');

bluebird.config({
    longStackTraces: false,
    warnings: false,
});

require('parsleyjs');

const locale = $('html').prop('lang');
window.Parsley.addMessages(locale, {
    notInUse: "This value is already in use.",
});
window.Parsley.setLocale(locale);

window.Parsley.addValidator('notInUse', {
    validateString: function(value, requirement) {
        if (typeof requirement === 'string') requirement = requirement.split(',');
        if (typeof requirement === 'number') requirement = [requirement.toString()];
        return $.inArray( value, requirement ) === -1;
    }
});
