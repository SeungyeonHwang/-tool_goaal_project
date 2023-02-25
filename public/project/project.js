(function ($) {

    const projectItems = document.querySelectorAll('.project-item');
    projectItems.forEach(item => {
        const projectDescription = item.querySelector('.project-description');
        console.log(projectDescription)
        item.addEventListener('click', () => {
            projectDescription.style.display = 'block';
        });
    });


    $("#color-options .color-option").click(function () {
        $(".color-option").removeClass("active");
        $(this).addClass("active");
        $("#project-color").val($(this).data("color"));
    });
})(jQuery);