(function ($) {
    $("#color-options .color-option").click(function() {
        $(".color-option").removeClass("active");
        $(this).addClass("active");
        $("#project-color").val($(this).data("color"));
    });
})(jQuery);