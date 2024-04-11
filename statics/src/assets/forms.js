import {getCache} from "@/assets/cache";

const
    $ = require('jquery');

export let $forms;

window.ParsleyConfig = {
    classHandler: function (f) {
        return f.$element.closest('.error-handler');
    },
    errorClass: "is-invalid",
    errorsContainer: function (f) {
        return f.$element.closest('.error-handler');
    },
    errorsWrapper: "<ul class='list-unstyled help-block is-invalid col-lg-12'></ul>",
    excluded: "input[type=button], input[type=submit], input[type=reset], input[type=hidden], [disabled], :hidden:not(.modal *)",
    trigger: "change keyup select"
};
require('parsleyjs');
require('./parsley-validators');

export function initParsley($forms) {
    $forms.each((_, form) => {
        const $form = $(form);
        $form.off('submit.Parsley');
        const parsley = $form.parsley();
        parsley.on('field:validated', function () {
            $form.find(':submit').attr('disabled', !parsley.isValid());
        });
        if (!$form.hasClass('wait')) {
            parsley.validate();
        } else {
            $form.find(':submit').attr('disabled', true);
        }
        form.addEventListener('submit', e => {
            e.preventDefault();
        });
    });
}

initParsley($('form[data-parsley-validate]'));

export function resetParsley($forms) {
    $forms.each((_, form) => {
        $(form).parsley().destroy();
    });
    initParsley($forms);
}

export function showMessage(isShow, status, _show_time) {
    const ajax_alert = $(".ajax-alert");
    const ajax_msg = $("#ajax-msg");
    if (isShow) {
        ajax_alert.removeClass('d-none').css('display', '');
        let show_time = 3000;
        if (!status) {
            ajax_alert.addClass('alert-success');
        } else {
            show_time = 5000;
            ajax_alert.addClass('alert-' + status);
        }
        if (_show_time) {
            show_time = _show_time;
        }
        setTimeout(function () {
            ajax_msg.text('');
            ajax_alert.addClass('d-none').removeClass('alert-danger alert-warning alert-info');
        }, show_time);
    } else {
        ajax_msg.text('');
        ajax_alert.addClass('d-none');
    }
}

export function initModal() {
    $('.ajax-alert').on('close.bs.alert', function (e) {
        e.preventDefault();
        showMessage(false);
    });

    $('input, select').on('keyup keydown keypress change paste', function () {
        showMessage(false);
    });
}

export function checkLogin($this) {
    if (getCache('session', 'hasAuth') !== 'true') {
        $this.$parent.logout();
        if ($this.$route.name !== 'Login') {
            $this.$router.push('/login');
        }
    }
}
