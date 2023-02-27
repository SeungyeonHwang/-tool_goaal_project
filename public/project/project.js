(function ($) {
    'use strict';

    $(function () {
        var projectListItem = $('.project-list');
        var projectName = $('#project-name');
        var projectCode = $('#project-code');
        var projectDescription = $('#project-description');
        var projectColor = $('#project-color');
        var projectPriority = $('#project-priority');

        $('#search-project').keyup(function() {
            var searchValue = $(this).val().toLowerCase();
            console.log(searchValue)
            $('.project-item').each(function() {
                var itemName = $(this).find('.project-name').text().toLowerCase();
                var itemCode = $(this).find('.project-code').text().toLowerCase();
                if (itemName.indexOf(searchValue) >= 0 || itemCode.indexOf(searchValue) >= 0) {
                    $(this).show();
                } else {
                    $(this).hide();
                }
            });
        });

        $('#project-create-btn').on("click", function () {
            var name = projectName.val();
            var code = projectCode.val();
            var description = projectDescription.val();
            var selectedColor = $('.color-option.active').data("color");
            var priority = projectPriority.val();

            if (!name || !code) {
                alert("Please fill in project name and code.");
                return;
            }

            $.post("/projects", {
                name: name,
                code: code,
                description: description,
                color: selectedColor,
                priority: priority
            });
            projectName.val("");
            projectCode.val("");
            projectDescription.val("");
        });

        var addItem = function (item) {
            var color = item.color || "#A9A9A9";
            var priorityText = ""; // 우선순위에 따른 한자 표시를 위한 문자열 변수
            switch(item.priority) {
                case "high":
                    priorityText = "上";
                    break;
                case "mid":
                    priorityText = "中";
                    break;
                case "low":
                    priorityText = "下";
                    break;
                default:
                    priorityText = "";
                    break;
            }
            var listItemHtml =
            "<li class='project-item'>" +
                "<div class='project-color' style='background-color: " + color + ";'></div>" +
                "<div class='project-name'>" + item.name + "</div>" +
                "<div class='project-code'>&nbsp;(" + item.code + ")</div>" +
                "<div class='project-priority'>" + priorityText + "</div>" +
            "</li>";
            projectListItem.append(listItemHtml);
        };
        
        $.get('/projects', function (items) {
            items.forEach(e => {
                addItem(e);
            });
        });

        $("#color-options .color-option").click(function () {
            $(".color-option").removeClass("active");
            $(this).addClass("active");
            projectColor.val($(this).data("color"));
        });
    });
})(jQuery);